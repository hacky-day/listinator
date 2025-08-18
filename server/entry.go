package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shaardie/listinator/database"
)

// entryList returns a handler function that retrieves all entries for a specific list.
// The list ID must be provided as a query parameter.
//
// Query Parameters:
//   - ListID: UUID of the list to retrieve entries for (required)
//
// Returns:
//   - 200 OK: JSON array of entries ordered by update time and completion status
//   - 400 Bad Request: If ListID is missing or invalid
//   - 500 Internal Server Error: If database query fails
func (s server) entryList() echo.HandlerFunc {
	type input struct {
		ListID string `query:"ListID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		if i.ListID == "" {
			return echo.ErrBadRequest.SetInternal(errors.New("missing ListID"))
		}

		es := []database.Entry{}
		// Order by updated_at and bought status to show unbought items first
		if err := s.db.Where("list_id = ?", i.ListID).Order("updated_at asc").Order("bought asc").Find(&es).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get entries from database, %w", err))
		}
		return c.JSON(http.StatusOK, es)
	}
}

// entryCreate returns a handler function that creates a new entry in the database.
//
// Request Body (JSON):
//   - Name: The name/description of the entry item (string)
//   - Number: Optional quantity or number (string)
//   - Bought: Whether the item is already purchased (boolean)
//   - TypeID: The type/category identifier (string)
//   - ListID: UUID of the parent list (UUID)
//
// Returns:
//   - 201 Created: JSON representation of the created entry
//   - 400 Bad Request: If request body is invalid
//   - 500 Internal Server Error: If database operation fails
func (s server) entryCreate() echo.HandlerFunc {
	type input struct {
		Name   string    `json:"Name"`
		Number string    `json:"Number"`
		Bought bool      `json:"Bought"`
		TypeID string    `json:"TypeID"`
		ListID uuid.UUID `json:"ListID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		e := database.Entry{
			Name:   i.Name,
			Number: i.Number,
			Bought: i.Bought,
			TypeID: i.TypeID,
			ListID: i.ListID,
		}
		if err := s.db.Create(&e).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		return c.JSON(http.StatusCreated, e)
	}
}

// entryUpdate returns a handler function that updates an existing entry.
//
// URL Parameters:
//   - id: UUID of the entry to update
//
// Request Body (JSON):
//   - Name: The updated name/description (string)
//   - Number: Updated quantity or number (string)
//   - Bought: Updated purchase status (boolean)
//   - TypeID: Updated type/category identifier (string)
//   - ListID: Updated parent list UUID (UUID)
//
// Returns:
//   - 200 OK: JSON representation of the updated entry
//   - 400 Bad Request: If request parameters or body are invalid
//   - 404 Not Found: If entry with given ID doesn't exist
//   - 500 Internal Server Error: If database operation fails
func (s server) entryUpdate() echo.HandlerFunc {
	type input struct {
		ID     uuid.UUID `param:"ID"`
		Name   string    `json:"Name"`
		Number string    `json:"Number"`
		Bought bool      `json:"Bought"`
		TypeID string    `json:"TypeID"`
		ListID uuid.UUID `json:"ListID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		e := database.Entry{
			Model: database.Model{
				ID: i.ID,
			},
		}
		// First, verify the entry exists
		if err := s.db.First(&e).Error; err != nil {
			// TODO: handle different errors (not found vs other database errors)
			return echo.NotFoundHandler(c)
		}
		
		// Update all fields with new values
		e.Name = i.Name
		e.Number = i.Number
		e.Bought = i.Bought
		e.TypeID = i.TypeID
		e.ListID = i.ListID

		if err := s.db.Save(&e).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		return c.JSON(http.StatusOK, e)
	}
}

// entryDelete returns a handler function that soft-deletes an entry from the database.
//
// URL Parameters:
//   - id: UUID of the entry to delete
//
// Returns:
//   - 200 OK: JSON representation of the deleted entry
//   - 400 Bad Request: If ID parameter is invalid
//   - 404 Not Found: If entry with given ID doesn't exist
//   - 500 Internal Server Error: If database operation fails
func (s server) entryDelete() echo.HandlerFunc {
	type input struct {
		ID uuid.UUID `param:"ID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		e := database.Entry{
			Model: database.Model{
				ID: i.ID,
			},
		}
		// Perform soft delete using GORM's Delete method
		if err := s.db.Delete(&e).Error; err != nil {
			// TODO: handle different errors (not found vs other database errors)
			return echo.NotFoundHandler(c)
		}
		return c.JSON(http.StatusOK, e)
	}
}
