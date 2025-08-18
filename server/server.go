// Package server provides HTTP handlers and API routes for the Listinator application.
// It implements RESTful API endpoints for managing lists, entries, and types.
package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// server holds the database connection and provides HTTP handler methods.
// All API endpoints are methods on this struct to maintain access to the database.
type server struct {
	db *gorm.DB // Database connection for data operations
}

// New creates a new server instance with the provided database connection.
// The server will use this database connection for all data operations.
//
// Parameters:
//   - db: A configured GORM database instance
//
// Returns:
//   - server: A new server instance ready to handle HTTP requests
func New(db *gorm.DB) server {
	return server{
		db: db,
	}
}

// SetupRoutes configures all API routes on the provided Echo instance.
// This method registers all available endpoints for the Listinator API v1.
//
// API Routes:
//   - GET /api/v1/entries - List all entries for a specific list
//   - POST /api/v1/entries - Create a new entry
//   - PUT /api/v1/entries/:id - Update an existing entry
//   - DELETE /api/v1/entries/:id - Delete an entry
//   - POST /api/v1/lists - Create a new list
//   - GET /api/v1/types - Get all available item types
//
// Parameters:
//   - e: The Echo instance to register routes on
func (s server) SetupRoutes(e *echo.Echo) {
	// Entry management endpoints
	e.GET("/api/v1/entries", s.entryList())
	e.POST("/api/v1/entries", s.entryCreate())
	e.PUT("/api/v1/entries/:id", s.entryUpdate())
	e.DELETE("/api/v1/entries/:id", s.entryDelete())

	// List management endpoints
	e.POST("/api/v1/lists", s.listCreate())

	// Type/category endpoints
	e.GET("/api/v1/types", s.typeList())
}
