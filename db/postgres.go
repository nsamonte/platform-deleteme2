package db

import (
	"fmt"
	"log"
	"time"

	"github.com/arundo/golang-crud-skeleton/migrations"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
)

type envs struct {
	MaxIdleConns int `envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	MaxOpenConns int `envconfig:"POSTGRES_MAX_OPEN_CONNS"`
}

// GetDBConnection returns postgres db instance
func GetDBConnection(connString string) (*gorm.DB, error) {

	var e envs
	if err := envconfig.Process("", &e); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect postgres (%s). %v", connString, err)
	}
	db.LogMode(true)
	db.DB().SetConnMaxLifetime(time.Minute)
	db.DB().SetMaxIdleConns(e.MaxIdleConns)
	db.DB().SetMaxOpenConns(e.MaxOpenConns)
	return db, err
}

// MigrateSkeletons migrates old user
func MigrateSkeletons(db *gorm.DB) error {
	err := migrations.MigrateSkeletons(db)
	if err != nil {
		return fmt.Errorf("Failed to migrate skeletons. %v", err)
	}
	return nil
}
