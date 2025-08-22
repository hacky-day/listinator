package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shaardie/listinator/database"
)

func (s server) typeList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ts := []database.Type{}
		if err := s.db.Order("name").Find(&ts).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get types from database, %w", err))
		}
		return c.JSON(http.StatusOK, ts)
	}
}

func (s server) typeCreate() echo.HandlerFunc {
	type input struct {
		Name     string `json:"Name"`
		Icon     string `json:"Icon"`
		Priority int    `json:"Priority"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		t := database.Type{
			Name:     i.Name,
			Icon:     i.Icon,
			Priority: i.Priority,
		}
		if err := s.db.Create(&t).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		return c.JSON(http.StatusCreated, t)
	}
}

func (s server) typeUpdate() echo.HandlerFunc {
	type input struct {
		Name     string `param:"Name"`
		Icon     string `json:"Icon"`
		Priority int    `json:"Priority"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		t := database.Type{
			Name: i.Name,
		}
		if err := s.db.First(&t).Error; err != nil {
			// TODO: handle different errors
			return echo.NotFoundHandler(c)
		}
		t.Icon = i.Icon
		t.Priority = i.Priority

		if err := s.db.Save(&t).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		return c.JSON(http.StatusOK, t)
	}
}

func (s server) typeDelete() echo.HandlerFunc {
	type input struct {
		Name string `param:"Name"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		t := database.Type{
			Name: i.Name,
		}
		if err := s.db.Delete(&t).Error; err != nil {
			// TODO: handle different errors
			return echo.NotFoundHandler(c)
		}
		return c.JSON(http.StatusOK, t)
	}
}
