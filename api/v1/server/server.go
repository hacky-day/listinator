package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shaardie/listinator/pubsub"
	"gorm.io/gorm"
)

type server struct {
	db *gorm.DB

	// Entry
	entryPubSub pubsub.PubSub[uuid.UUID, entryEvent]
}

func New(db *gorm.DB, logger echo.Logger) server {
	return server{
		db:          db,
		entryPubSub: pubsub.New[uuid.UUID, entryEvent](logger, 16),
	}
}

func (s server) SetupRoutes(g *echo.Group) {
	// entries
	g.GET("/entries", s.entryList())
	g.POST("/entries", s.entryCreate())
	g.GET("/entries/:id", s.entryGet())
	g.PUT("/entries/:id", s.entryUpdate())
	g.DELETE("/entries/:id", s.entryDelete())
	g.GET("/entries/events", s.entryGetEvents())

	// lists
	g.POST("/lists", s.listCreate())

	// types
	g.GET("/types", s.typeList())

	// Login, Logout and stuff
	g.GET("/session", s.sessionMiddleware(s.sessionGet()))
	g.POST("/session", s.sessionCreate())
	g.DELETE("/session", s.sessionDelete())
}
