// Package server provides HTTP handlers and API routes for the Listinator web application.
// 
// This package implements a RESTful API server using the Echo framework that manages
// shopping lists and to-do lists. It provides endpoints for:
//
//   - Creating and managing lists
//   - Adding, updating, and deleting list entries  
//   - Categorizing entries with types/icons
//   - Marking entries as purchased/completed
//
// The server uses GORM as the ORM layer for database operations and supports
// JSON request/response format for all API endpoints.
//
// All API routes are prefixed with /api/v1/ and follow REST conventions.
package server
