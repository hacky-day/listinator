package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type migration struct {
	version uint
	up      func(db *gorm.DB) error
}

var migrations = []migration{
	{
		version: 1,
		up: func(db *gorm.DB) error {
			return nil
		},
	},
}

const currentVersion = 1

func migrate(ctx context.Context, db *gorm.DB) error {
	if !db.Migrator().HasTable("version") { // no version table, so it's a fresh one, just automigrate
		if err := db.AutoMigrate(&Entry{}, &Type{}, &List{}, &User{}, &Version{}); err != nil {
			return fmt.Errorf("unable to run auto migration, %w", err)
		}
		db.Logger.Info(ctx, "Successfully initialized database")
	} else { // Run all migration step up to the current version
		version := Version{}

		if err := db.Order("id desc").First(&version).Error; err != nil {
			return fmt.Errorf("unable to get latest version from database, %w", err)
		}
		for _, m := range migrations {
			if version.ID < m.version {
				if err := m.up(db); err != nil {
					return fmt.Errorf("unable to migrate to version %v, %w", m.version, err)
				}
				db.Logger.Info(context.TODO(), "migrated database to version %v", m.version)
			}
		}
	}
	// save current version in database
	db.Save(&Version{
		Model: gorm.Model{
			ID: currentVersion,
		},
	})
	db.Logger.Info(ctx, "Database version %v", currentVersion)
	return nil
}
