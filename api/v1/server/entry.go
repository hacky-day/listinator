package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"

	"github.com/shaardie/listinator/database"
)

func (s server) entryList() echo.HandlerFunc {
	type input struct {
		ListID string `query:"ListID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		if i.ListID == "" {
			return echo.ErrBadRequest.SetInternal(errors.New("missing ListID"))
		}

		es := []database.Entry{}
		if err := s.db.Where("list_id = ?", i.ListID).Order("updated_at asc").Order("bought asc").Find(&es).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to get entries from database, %w", err))
		}
		return c.JSON(http.StatusOK, es)
	}
}

func (s server) entryCreate() echo.HandlerFunc {
	type input struct {
		ID              uuid.UUID  `json:"ID"`
		Name            string     `json:"Name"`
		Number          string     `json:"Number"`
		Bought          bool       `json:"Bought"`
		TypeID          uuid.UUID  `json:"TypeID"`
		ListID          uuid.UUID  `json:"ListID"`
		ClientUpdatedAt *time.Time `json:"ClientUpdatedAt"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		clientUpdatedAt := i.ClientUpdatedAt
		if clientUpdatedAt == nil {
			now := time.Now()
			clientUpdatedAt = &now
		}

		e := database.Entry{
			Model:           database.Model{ID: i.ID},
			Name:            i.Name,
			Number:          i.Number,
			Bought:          i.Bought,
			TypeID:          i.TypeID,
			ListID:          i.ListID,
			ClientUpdatedAt: clientUpdatedAt,
		}
		// The offline write outbox replays creates on reconnect, and a
		// dropped response after a successful insert must not turn into a
		// duplicate row. Client-supplied IDs make replays idempotent: a
		// conflict on the primary key means this exact create already
		// landed, so skip the insert instead of erroring or duplicating.
		result := s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&e)
		if result.Error != nil {
			return echo.ErrInternalServerError.SetInternal(result.Error)
		}
		if result.RowsAffected == 0 {
			if err := s.db.First(&e, i.ID).Error; err != nil {
				return echo.ErrInternalServerError.SetInternal(err)
			}
			return c.JSON(http.StatusOK, e)
		}

		s.entryPubSub.Publish(e.ListID, entryEvent{
			Action: "create",
			Entry:  e,
		})
		return c.JSON(http.StatusCreated, e)
	}
}

func (s server) entryGet() echo.HandlerFunc {
	type input struct {
		ID uuid.UUID `param:"ID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		var e database.Entry
		if err := s.db.First(&e, i.ID).Error; err != nil {
			return echo.NotFoundHandler(c)
		}
		return c.JSON(http.StatusOK, e)
	}
}

func (s server) entryUpdate() echo.HandlerFunc {
	type input struct {
		ID              uuid.UUID  `param:"ID"`
		Name            string     `json:"Name"`
		Number          string     `json:"Number"`
		Bought          bool       `json:"Bought"`
		TypeID          uuid.UUID  `json:"TypeID"`
		ListID          uuid.UUID  `json:"ListID"`
		ClientUpdatedAt *time.Time `json:"ClientUpdatedAt"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		var e database.Entry
		if err := s.db.First(&e, i.ID).Error; err != nil {
			return echo.NotFoundHandler(c)
		}

		clientUpdatedAt := i.ClientUpdatedAt
		if clientUpdatedAt == nil {
			now := time.Now()
			clientUpdatedAt = &now
		}
		// Edit-time last-write-wins: a queued offline write that reaches the
		// server after a newer edit already landed must not clobber it just
		// because it arrived later. A nil e.ClientUpdatedAt (row predates
		// this column, or was never touched by an offline-aware client) has
		// nothing to compare against, so it always accepts the write.
		if e.ClientUpdatedAt != nil && clientUpdatedAt.Before(*e.ClientUpdatedAt) {
			return c.JSON(http.StatusOK, e)
		}

		e.Name = i.Name
		e.Number = i.Number
		e.Bought = i.Bought
		e.TypeID = i.TypeID
		e.ListID = i.ListID
		e.ClientUpdatedAt = clientUpdatedAt

		if err := s.db.Save(&e).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		s.entryPubSub.Publish(e.ListID, entryEvent{
			Action: "update",
			Entry:  e,
		})
		return c.JSON(http.StatusOK, e)
	}
}

func (s server) entryDelete() echo.HandlerFunc {
	type input struct {
		ID              uuid.UUID  `param:"ID"`
		ClientUpdatedAt *time.Time `json:"ClientUpdatedAt"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		var e database.Entry
		if err := s.db.First(&e, i.ID).Error; err != nil {
			return echo.NotFoundHandler(c)
		}

		clientUpdatedAt := i.ClientUpdatedAt
		if clientUpdatedAt == nil {
			now := time.Now()
			clientUpdatedAt = &now
		}
		// Symmetric to entryUpdate's last-write-wins check: a delete queued
		// offline that reaches the server after a newer edit already landed
		// must not destroy that edit just because it arrived later. Reject
		// with the current row instead of deleting it.
		if e.ClientUpdatedAt != nil && clientUpdatedAt.Before(*e.ClientUpdatedAt) {
			return c.JSON(http.StatusConflict, e)
		}

		if err := s.db.Delete(&e).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to delete entry %v, %w", e, err))
		}
		s.entryPubSub.Publish(e.ListID, entryEvent{
			Action: "delete",
			Entry:  e,
		})
		return c.JSON(http.StatusOK, e)
	}
}

type entryEvent struct {
	Action string
	Entry  database.Entry
}

func (s server) entryGetEvents() echo.HandlerFunc {
	type input struct {
		ListID uuid.UUID `query:"ListID"`
	}

	// sendEvent is a helper function to simplify sending SSE Events
	sendEvent := func(c echo.Context, event string, data string) error {
		w := c.Response()
		defer w.Flush()
		if _, err := fmt.Fprintf(w, "id: %v\n", time.Now().UnixMilli()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "data: %v\n", data); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "event: %v\n", event); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
		return nil
	}

	// ping just send a ping SSE Event with the current time
	ping := func(c echo.Context) error {
		if err := sendEvent(c, "ping", "ping"); err != nil {
			return fmt.Errorf("failed to send event, %w", err)
		}
		return nil
	}

	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		// subscribe to the pubsub channel for this list
		id, ch, err := s.entryPubSub.Subscribe(i.ListID)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to subscribe, %w", err))
		}
		// unsubscribe after return
		defer s.entryPubSub.Unsubscribe(i.ListID, id)

		w := c.Response()
		// set header for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {
			// return, if connection is closed from client
			case <-c.Request().Context().Done():
				return nil
			// send a ping every 5 seconds to keep the SSE connection. This can probably be less than that
			case <-time.After(5 * time.Second):
				if err := ping(c); err != nil {
					return echo.ErrInternalServerError.SetInternal(fmt.Errorf("failed to ping, %w", err))
				}
			// receive data from pubsub channel and send them to the client
			case ee, ok := <-ch:
				// channel was closed because we could not keep up (see
				// pubsub.Publish); end the stream so the client notices and
				// reconnects with a fresh full sync instead of busy-looping
				// on the now permanently-ready closed channel.
				if !ok {
					return echo.ErrInternalServerError.SetInternal(errors.New("disconnected, too slow to keep up with events"))
				}
				entry, err := json.Marshal(ee.Entry)
				if err != nil {
					return echo.ErrInternalServerError.SetInternal(fmt.Errorf("failed to marshal JSON, %w", err))
				}
				if err := sendEvent(c, ee.Action, string(entry)); err != nil {
					return echo.ErrInternalServerError.SetInternal(fmt.Errorf("unable to send event, %w", err))
				}
			}
		}
	}
}
