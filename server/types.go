package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shaardie/listinator/database"
)

// typeList returns a handler function that retrieves all available item types.
// Types are used to categorize list entries (e.g., fruit, vegetable, etc.)
// and include both a name identifier and a visual icon.
//
// Returns:
//   - 200 OK: JSON array of all types, ordered alphabetically by name
//   - 500 Internal Server Error: If database query fails
//
// Example response:
//   [
//     {"Name": "bread", "Icon": "🥖"},
//     {"Name": "dairy", "Icon": "🧀"},
//     {"Name": "fruit", "Icon": "🍎"}
//   ]
func (s server) typeList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ts := []database.Type{}
		if err := s.db.Order("name").Find(&ts).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get types from database, %w", err))
		}
		return c.JSON(http.StatusOK, ts)
	}
}
