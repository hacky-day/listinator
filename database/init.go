package database

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(ctx context.Context, dsn string) (*gorm.DB, error) {
	dialector := sqlite.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &slogLogger{},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to open database, %w", err)
	}

	if err := migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("unable to migrate database models, %w", err)
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

	// types to database
	types := []Type{
		{Name: "fruit", Icon: "🍎", Color: "crimson"},
		{Name: "vegetable", Icon: "🥦", Color: "green"},
		{Name: "drink", Icon: "🍹", Color: "orange"},
		{Name: "meat", Icon: "🍖", Color: "red"},
		{Name: "snack", Icon: "🍿", Color: "yellow"},
		{Name: "dairy", Icon: "🧀", Color: "gold"},
		{Name: "bread", Icon: "🥖", Color: "saddlebrown"},
		{Name: "condiment", Icon: "🧂", Color: "gray"},
		{Name: "frozen", Icon: "❄️", Color: "lightblue"},
		{Name: "canned", Icon: "🥫", Color: "silver"},
		{Name: "spice", Icon: "🌶️", Color: "darkred"},
		{Name: "unknown", Icon: "🤷‍♀️", Color: "black"},
	}
	for _, t := range types {
		if err := db.Save(&t).Error; err != nil {
			return nil, fmt.Errorf("unable to create default type '%s' in database, %w", t.Name, err)
		}
	}

	return db, nil
}
