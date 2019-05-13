package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/arundo/data-sdk-go/authentication"
	"github.com/arundo/data-sdk-go/ginutils"
	"github.com/arundo/data-sdk-go/status"
	"github.com/arundo/data-sdk-go/utils"
	v0alpha "github.com/arundo/golang-crud-skeleton/crudv0alpha"
	pgdb "github.com/arundo/golang-crud-skeleton/db"

	_ "github.com/joho/godotenv/autoload"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type envs struct {
	ServiceName          string `envconfig:"SERVICE_NAME"`
	EnableIPWhitelist    bool   `envconfig:"ENABLE_IP_WHITELIST"`
	EndpointMetricsPort  int    `envconfig:"ENDPOINT_METRICS_PORT"`
	GormLogMode          string `envconfig:"GORM_LOG_MODE"`
	HeartbeatFrequency   int    `envconfig:"HEARTBEAT_FREQUENCY"`
	IPWhiteList          string `envconfig:"IP_WHITELIST"`
	LogLevel             string `envconfig:"LOG_LEVEL"`
	PostgresHost         string `envconfig:"POSTGRES_HOSTNAME"`
	PostgresUser         string `envconfig:"POSTGRES_USERNAME"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	PostgresPort         string `envconfig:"POSTGRES_PORT"`
	PostgresDatabase     string `envconfig:"POSTGRES_DATABASE_NAME"`
	PostgresSSLMode      string `envconfig:"POSTGRES_SSL_MODE"`
	Port                 string `envconfig:"PORT"`
	InternalAccessAPIKey string `envconfig:"INTERNAL_ACCESS_API_KEY"`
}

func main() {
	var e envs
	if err := envconfig.Process("", &e); err != nil {
		log.Fatal(err)
	}

	utils.SetLogLevel(e.LogLevel)
	utils.StandardFormatLogs()
	status.InitMetrics()
	srv := utils.StartPprof(6060)
	defer func() {
		if err := srv.Shutdown(context.TODO()); err != nil {
			log.WithError(err)
		}
	}()

	dbConnString := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s", e.PostgresHost, e.PostgresPort, e.PostgresDatabase, e.PostgresSSLMode, e.PostgresUser, e.PostgresPassword)
	db, err := pgdb.GetDBConnection(dbConnString)
	if err != nil {
		log.Fatal(err)
	}

	if err = pgdb.MigrateSkeletons(db); err != nil {
		log.Fatal(err)
	}

	logMode, err := strconv.ParseBool(e.GormLogMode)
	if err != nil {
		log.WithError(err).Error("GORM_LOG_MODE error")
	}
	db.LogMode(logMode)
	defer db.Close()

	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		ginutils.FormatResponse(c, http.StatusNotFound, "about:blank", http.StatusText(http.StatusNotFound), "Page not found", utils.WhereAmI())
	})
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))
	r.Use(status.MetricsMiddlewares()...)

	noAuth := r.Group("")
	noAuth.GET("skeletons/health", v0alpha.HealthCheck(db))

	iAuthdv0Alpha := r.Group("v0alpha", gin.BasicAuth(gin.Accounts{
		"internal-access": e.InternalAccessAPIKey,
	}))
	iAuthdv0Alpha.Use(authentication.InternalAuthMiddleware())
	if e.EnableIPWhitelist {
		iAuthdv0Alpha.Use(authentication.WhiteListIPs(e.IPWhiteList))
	}
	iAuthdv0Alpha.POST("/iskeletons/:company", v0alpha.CreateSkeleton(db))
	iAuthdv0Alpha.GET("/iskeletons/:company", v0alpha.GetSkeleton(db))
	iAuthdv0Alpha.PUT("/iskeletons/:company/:id", v0alpha.UpdateSkeleton(db))

	authdv0Alpha := r.Group("v0alpha")
	jwtMidleware := authentication.JwtMiddlewareInit()
	authdv0Alpha.Use(authentication.Auth0Middleware(jwtMidleware))
	authdv0Alpha.POST("/skeletons", v0alpha.CreateSkeleton(db))
	authdv0Alpha.GET("/skeletons", v0alpha.GetSkeleton(db))
	authdv0Alpha.GET("/skeletons/:skeletonId", v0alpha.GetSkeleton(db))
	authdv0Alpha.PUT("/skeletons/:skeletonId", v0alpha.UpdateSkeleton(db))
	authdv0Alpha.DELETE("/skeletons/:skeletonId", v0alpha.DeleteSkeleton(db))

	r.Use(gin.Recovery())

	ctx, cancel := context.WithCancel(context.Background())

	status.StartMetrics(e.EndpointMetricsPort)
	status.StartLoggingMetrics()
	go status.Beat(ctx, e.ServiceName, e.HeartbeatFrequency)
	go startServer(ctx, r, e.Port)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	sig := <-signalChan
	cancel()
	log.Warningf("Exited %s. Received signal: %+v", e.ServiceName, sig)
	status.StopLoggingMetrics()
}

func startServer(ctx context.Context, r *gin.Engine, port string) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(ctx)
		}
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server. %v", err)
	}
}
