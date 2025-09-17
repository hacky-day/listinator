package server

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shaardie/listinator/core/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

const (
	sessionKey = "listinator_session"
	uuidKey    = "uuid"
	userKey    = "user"
)

func (s server) sessionCreate() echo.HandlerFunc {
	type input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		if i.Name == "" || i.Password == "" {
			return echo.ErrBadRequest
		}

		u := database.User{}
		if err := s.db.Find(&u, "name = ?", i.Name).Error; err != nil {
			return echo.ErrUnauthorized.SetInternal(fmt.Errorf("unable to get user %v from database, %w", i.Name, err))
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(i.Password)); err != nil {
			return echo.ErrUnauthorized.SetInternal(fmt.Errorf("wrong password for user %v, %w", i.Name, err))
		}

		sess, err := session.Get(sessionKey, c)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get session, %w", err))
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 365, // one year
			HttpOnly: true,
			Secure:   true,
		}
		sess.Values[uuidKey] = u.ID.String()
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to save session, %w", err))
		}

		return nil
	}
}

func (s server) sessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get uuid from session cookie
		sess, err := session.Get(sessionKey, c)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get session, %w", err))
		}
		uuidAny, ok := sess.Values[uuidKey]
		if !ok {
			return echo.ErrUnauthorized
		}

		// Check if string
		uuidStr, ok := uuidAny.(string)
		if !ok {
			return echo.ErrUnauthorized.SetInternal(errors.New("uuid in cookie is no string"))
		}

		// Parse UUID
		uuidObj, err := uuid.Parse(uuidStr)
		if err != nil {
			return echo.ErrUnauthorized.SetInternal(fmt.Errorf("unable to parse uuid, %w", err))
		}

		u := database.User{
			Model: database.Model{
				ID: uuidObj,
			},
		}
		if err := s.db.First(&u).Error; err != nil {
			return echo.ErrUnauthorized.SetInternal(fmt.Errorf("unable to get user from database, %w", err))
		}

		c.Set(userKey, &u)
		return next(c)
	}
}

func (s server) sessionGet() echo.HandlerFunc {
	return func(c echo.Context) error {
		userAny := c.Get(userKey)
		if userAny == nil {
			return echo.ErrUnauthorized
		}
		user, ok := userAny.(*database.User)
		if !ok {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("wrong type %T context", userAny))
		}
		return c.JSON(200, user)
	}
}

func (s server) sessionDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(sessionKey, c)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get session, %w", err))
		}
		sess.Options.MaxAge = -1
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to save session, %w", err))
		}
		return nil
	}
}
