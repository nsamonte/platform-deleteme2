package migrations

import (
	"fmt"

	v0alpha "github.com/arundo/golang-crud-skeleton/crudv0alpha"
	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

type skeletons v0alpha.Skeleton

type envs struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
}

// MigrateSkeletons executes skeletons table migrations
func MigrateSkeletons(db *gorm.DB) error {
	var e envs
	if err := envconfig.Process("", &e); err != nil {
		log.Fatal(err)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: fmt.Sprintf("2019.109.0000-skeletons-%s", e.ServiceName),
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&skeletons{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("skeletons").Error
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return err
	}
	log.Printf("Skeleton table migrations completed")
	return nil
}
