package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes the database connection and performs necessary setup operations.
// It opens a SQLite database, runs migrations, and seeds default data.
//
// Parameters:
//   - dsn: The database connection string/path for the SQLite database file
//
// Returns:
//   - *gorm.DB: A configured GORM database instance
//   - error: Any error that occurred during initialization
func Init(dsn string) (*gorm.DB, error) {
	dialector := sqlite.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database, %w", err)
	}

	// Run database migrations to create/update tables
	err = db.AutoMigrate(&Entry{}, &Type{}, &List{})
	if err != nil {
		return nil, fmt.Errorf("unable to migrate database models, %w", err)
	}

	// Seed database with default item types
	types := []Type{
		{Name: "fruit", Icon: "🍎"},
		{Name: "vegetable", Icon: "🥦"},
		{Name: "drink", Icon: "🍹"},
		{Name: "meat", Icon: "🍖"},
		{Name: "snack", Icon: "🍿"},
		{Name: "dairy", Icon: "🧀"},
		{Name: "bread", Icon: "🥖"},
		{Name: "condiment", Icon: "🧂"},
		{Name: "frozen", Icon: "❄️"},
		{Name: "canned", Icon: "🥫"},
		{Name: "spice", Icon: "🌶️"},
		{Name: "unknown", Icon: "🤷‍♀️"},
	}
	
	// Insert default types (will update existing ones due to GORM's Save behavior)
	for _, t := range types {
		if err := db.Save(&t).Error; err != nil {
			return nil, fmt.Errorf("unable to create default type '%s' in database, %w", t.Name, err)
		}
	}

	return db, nil
}
