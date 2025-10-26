package main

import (
	"embed"
	"log/slog"
	"os"
	"path"

	"github.com/gorilla/sessions"
	"github.com/shaardie/listinator/api/v1/server"
	"github.com/shaardie/listinator/database"
	"github.com/shaardie/listinator/logger"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	p := os.Getenv("LISTINATOR_DATABASE_DIR")
	dbPath := path.Join(p, "listinator.db")

	if err := logger.Init(); err != nil {
		panic(err)
	}

	sessionSecret := os.Getenv("LISTINATOR_SESSION_SECRET")
	if sessionSecret == "" {
		panic("session secret missing")
	}

	// init database
	db, err := database.Init(dbPath)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	// API V1
	apiV1 := e.Group("/api/v1")
	sV1 := server.New(db)
	sV1.SetupRoutes(apiV1)

	// Embeded Frontend
	e.StaticFS("/", echo.MustSubFS(frontendFS, "frontend/dist"))

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
