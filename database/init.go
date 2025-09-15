package database

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(dsn string) (*gorm.DB, error) {
	dialector := sqlite.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &slogLogger{},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to open database, %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get *sql.DB from *gorm.DB, %w", err)
	}

	if err := migrate(sqlDB); err != nil {
		return nil, fmt.Errorf("migration failed, %w", err)
	}

	// Create admin user with password or update, if already present.
	adminPassword := os.Getenv("LISTINATOR_ADMIN_PASSWORD")
	if adminPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to hash admin password, %w", err)
		}

		// This can probably be done nicer, but I think the race condition is not that important during startup
		admin := User{Name: "admin", PasswordHash: string(hash), IsAdmin: true}
		x := db.Model(&admin).Where("name = ?", "admin").Updates(&admin)
		if x.Error != nil {
			return nil, fmt.Errorf("unable to update admin, %w", err)
		}
		if x.RowsAffected == 0 {
			if err := db.Create(&admin).Error; err != nil {
				return nil, fmt.Errorf("unable to create admin, %w", err)
			}
		}
	}

	return db, nil
}
