// Package database provides data models and database operations for the Listinator application.
// It includes models for lists, entries, and types used in the shopping list management system.
package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model is the base model struct that provides common fields for all database entities.
// It includes a UUID primary key and standard GORM timestamps for audit purposes.
type Model struct {
	ID        uuid.UUID `gorm:"primaryKey;type=uuid"` // UUID primary key
	CreatedAt time.Time                               // Automatically set when record is created
	UpdatedAt time.Time                               // Automatically updated when record is modified
	DeletedAt gorm.DeletedAt `gorm:"index"`           // Soft delete timestamp
}

// BeforeCreate is a GORM hook that ensures each model has a UUID before being saved to the database.
// If no UUID is set, it generates a new one automatically.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// List represents a shopping list or to-do list in the system.
// Each list can contain multiple entries and is identified by a UUID.
type List struct {
	Model

	// Entries contains all the items associated with this list
	Entries []Entry
}

// Entry represents a single item within a list (e.g., a shopping list item).
// Each entry has a name, optional number/quantity, completion status, and type classification.
type Entry struct {
	Model

	Name   string // The name/description of the entry item
	Number string // Optional quantity or number (e.g., "2", "500g")

	Bought bool // Whether the item has been purchased/completed

	TypeID string // Foreign key referencing the Type table
	Type   Type   `json:"-"` // Associated type/category (excluded from JSON serialization)

	ListID uuid.UUID // Foreign key referencing the parent List
}

// BeforeCreate is a GORM hook that prepares an Entry for database insertion.
// It ensures the entry has a UUID and sets a default type if none is specified.
func (e *Entry) BeforeCreate(tx *gorm.DB) error {
	if err := e.Model.BeforeCreate(tx); err != nil {
		return err
	}
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	// Set default type if none specified
	if e.TypeID == "" {
		e.TypeID = "unknown"
	}
	return nil
}

// Type represents a category or classification for list entries.
// Types have both a name identifier and an icon for visual representation.
type Type struct {
	Name string `gorm:"primaryKey"` // The type identifier (e.g., "fruit", "vegetable")
	Icon string                     // Unicode emoji or icon representing this type
}
