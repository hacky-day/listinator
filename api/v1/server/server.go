package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB
}

func New(db *gorm.DB) server {
	return server{
		db: db,
	}
}

func (s server) SetupRoutes(g *echo.Group) {
	// entries
	g.GET("/entries", s.entryList())
	g.POST("/entries", s.entryCreate())
	g.PUT("/entries/:id", s.entryUpdate())
	g.DELETE("/entries/:id", s.entryDelete())

	// lists
	g.POST("/lists", s.listCreate())

	// types
	g.GET("/types", s.typeList())
}
