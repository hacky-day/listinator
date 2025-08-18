// Package main is the entry point for the Listinator application.
// Listinator is a Vue.js-based web application for managing shopping lists and to-do lists.
// It provides a RESTful API backend written in Go and a Vue.js frontend with TypeScript.
package main

import (
	"embed"
	"os"
	"path"

	"github.com/shaardie/listinator/database"
	"github.com/shaardie/listinator/server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// frontendFS embeds the built Vue.js frontend assets into the Go binary.
// This allows for single-binary deployment without external file dependencies.
//
//go:embed frontend/*
var frontendFS embed.FS

// main is the application entry point.
// It initializes the database, sets up the HTTP server with API routes,
// and starts serving both the API and the embedded frontend.
//
// Environment variables:
//   - LISTINATOR_DATABASE_DIR: Required. Directory where the SQLite database file will be stored.
//
// The application serves:
//   - API endpoints at /api/v1/* for list and entry management
//   - Static frontend assets at / (embedded Vue.js application)
//
// Default port: 8080
func main() {
	// Get database directory from environment variable
	p := os.Getenv("LISTINATOR_DATABASE_DIR")
	dbPath := path.Join(p, "listinator.db")

	// Initialize database connection and run migrations
	db, err := database.Init(dbPath)
	if err != nil {
		panic(err)
	}

	// Create Echo HTTP server instance
	e := echo.New()
	e.Use(middleware.Logger()) // Log all HTTP requests

	// Development vs Production frontend serving strategy:
	// In development: serve from filesystem to allow hot reload
	// In production: serve embedded frontend from binary
	if _, err := os.Stat("frontend"); os.IsNotExist(err) {
		// Production: serve embedded frontend assets
		e.StaticFS("/", echo.MustSubFS(frontendFS, "frontend"))
	} else {
		// Development: serve from filesystem
		e.Static("/", "frontend")
	}

	// Initialize API server and register all routes
	s := server.New(db)
	s.SetupRoutes(e)

	// Start HTTP server on port 8080
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
