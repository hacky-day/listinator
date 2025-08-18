// Package database provides data models and database operations for the Listinator application.
//
// This package defines the core data structures used throughout the application:
//
//   - Model: Base struct with UUID primary keys and GORM timestamps
//   - List: Represents a shopping list or to-do list container
//   - Entry: Individual items within a list (with name, quantity, type, completion status)
//   - Type: Categories for entries with name and icon (fruit, vegetable, etc.)
//
// The package handles database initialization, schema migrations, and seeding
// of default data. It uses SQLite as the database backend with GORM as the ORM.
//
// Database features:
//   - UUID primary keys for all entities
//   - Soft deletes with GORM's DeletedAt field  
//   - Automatic timestamps (CreatedAt, UpdatedAt)
//   - Foreign key relationships between lists and entries
//   - Pre-seeded item types with emoji icons
package database
