package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID `gorm:"primaryKey;type=uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type User struct {
	Model

	Name         string
	PasswordHash string `json:"-"`
	IsAdmin      bool
}

type List struct {
	Model

	Entries []Entry
}

type Entry struct {
	Model

	Name   string
	Number string

	Bought bool

	TypeID uuid.UUID
	Type   Type `json:"-"`

	ListID uuid.UUID

	// ClientUpdatedAt is when the user actually made the edit, as opposed to
	// UpdatedAt (server-write time). It drives edit-time last-write-wins for
	// offline clients replaying queued writes: a pointer so rows written
	// before this column existed, and payloads from older/non-offline-aware
	// clients, have no value to compare against and always accept the write.
	ClientUpdatedAt *time.Time
}

func (e *Entry) BeforeCreate(tx *gorm.DB) error {
	if err := e.Model.BeforeCreate(tx); err != nil {
		return err
	}
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.TypeID == uuid.Nil {
		e.TypeID = uuid.MustParse("c29ebd85-812e-4cf6-bfc7-c8368eb83334")
	}
	return nil
}

type Type struct {
	Model

	Name string
	// Slug is a stable, language independent identifier used by the frontend
	// to look up a translated name for this type.
	Slug     string
	Color    string
	Priority int
}
