package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shaardie/listinator/database"
)

// listCreate returns a handler function that creates a new empty list.
// The list is created with just a UUID and timestamps - no additional data is required.
//
// Returns:
//   - 201 Created: JSON representation of the newly created list with its UUID
//   - 500 Internal Server Error: If database operation fails
//
// Example response:
//   {
//     "ID": "123e4567-e89b-12d3-a456-426614174000",
//     "CreatedAt": "2023-01-01T12:00:00Z",
//     "UpdatedAt": "2023-01-01T12:00:00Z",
//     "DeletedAt": null,
//     "Entries": []
//   }
func (s server) listCreate() echo.HandlerFunc {
	return func(c echo.Context) error {
		l := database.List{}
		if err := s.db.Create(&l).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		return c.JSON(http.StatusCreated, l)
	}
}
