package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

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
		Name   string    `json:"Name"`
		Number string    `json:"Number"`
		Bought bool      `json:"Bought"`
		TypeID string    `json:"TypeID"`
		ListID uuid.UUID `json:"ListID"`
	}
	return func(c echo.Context) error {
		var i input
		if err := c.Bind(&i); err != nil {
			return echo.ErrBadRequest.SetInternal(err)
		}

		e := database.Entry{
			Name:   i.Name,
			Number: i.Number,
			Bought: i.Bought,
			TypeID: i.TypeID,
			ListID: i.ListID,
		}
		if err := s.db.Create(&e).Error; err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}

		s.entryPubSub.Publish(e.ListID, entryEvent{
			Action: "create",
			Entry:  e,
		})
		return c.JSON(http.StatusCreated, e)
	}
}

func (s server) entryUpdate() echo.HandlerFunc {
	type input struct {
		ID     uuid.UUID `param:"ID"`
		Name   string    `json:"Name"`
		Number string    `json:"Number"`
		Bought bool      `json:"Bought"`
		TypeID string    `json:"TypeID"`
		ListID uuid.UUID `json:"ListID"`
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
		e.Name = i.Name
		e.Number = i.Number
		e.Bought = i.Bought
		e.TypeID = i.TypeID
		e.ListID = i.ListID

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
	sendEvent := func(c echo.Context, t time.Time, event string, data string) error {
		w := c.Response()
		defer w.Flush()
		if _, err := fmt.Fprintf(w, "id: %v\n", t.UnixMilli()); err != nil {
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
		now := time.Now()
		if err := sendEvent(c, now, "ping", fmt.Sprintf("%v", now)); err != nil {
			return fmt.Errorf("failed to send event, %w", err)
		}
		return nil
	}

	// replayEvents is a helper function which reploys all entries from the database,
	// which where created, updated oder delete after the last sync as events to the client.
	replayEvents := func(c echo.Context, listID uuid.UUID, lastSync time.Time) error {
		es := []database.Entry{}
		if err := s.db.
			Unscoped(). // also get soft-deleted entries
			Where("list_id = ?", listID).
			Where("updated_at > ? OR deleted_at > ?", lastSync, lastSync).
			Order("MAX(updated_at, COALESCE(deleted_at, updated_at)) ASC").
			Find(&es).Error; err != nil {
			return fmt.Errorf("unable to get entries from database, %w", err)
		}
		for _, e := range es {
			marshaledEntry, err := json.Marshal(e)
			if err != nil {
				return fmt.Errorf("failed to marshal JSON, %w", err)
			}
			// Entry was deleted after last sync.
			if e.DeletedAt.Time.After(lastSync) {
				// Entry was created and deleted after the last sync. So we can ignore it.
				if e.CreatedAt.After(lastSync) {
					continue
				}
				sendEvent(c, e.DeletedAt.Time, "delete", string(marshaledEntry))
				continue
			}

			// Entry was created after last sync.
			if e.CreatedAt.After(lastSync) {
				sendEvent(c, e.UpdatedAt, "create", string(marshaledEntry))
				continue
			}

			// Entry was updated, but not created after last sync.
			if e.UpdatedAt.After(lastSync) {
				sendEvent(c, e.UpdatedAt, "update", string(marshaledEntry))
			}
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

		// ff it's a SSE reconnect a Last-Event-ID Header is sent.
		// We use the date in the header to calculate, which database changes the client musst have missed and replay them to the client.
		r := c.Request()
		lastEventID := r.Header.Get("Last-Event-ID")
		if lastEventID != "" {
			UnixMilliLastEventID, err := strconv.ParseInt(lastEventID, 10, 0)
			if err != nil {
				return echo.ErrBadRequest.SetInternal(fmt.Errorf("unable to parse Last-Event-ID Header, %w", err))
			}
			if err := replayEvents(c, i.ListID, time.UnixMilli(UnixMilliLastEventID)); err != nil {
				return echo.ErrInternalServerError.SetInternal(fmt.Errorf("failed to replay events to client, %w", err))
			}
		}

		for {
			select {
			// return, if connection is closed from client
			case <-r.Context().Done():
				return nil
			// send a ping every 5 seconds to keep the SSE connection. This can probably be less than that
			case <-time.After(5 * time.Second):
				if err := ping(c); err != nil {
					return echo.ErrInternalServerError.SetInternal(fmt.Errorf("failed to ping, %w", err))
				}
			// receive data from pubsub channel and send them to the client
			case ee := <-ch:
				entry, err := json.Marshal(ee.Entry)
				if err != nil {
					return echo.ErrInternalServerError.SetInternal(fmt.Errorf("failed to marshal JSON, %w", err))
				}

				// time depends on the type of action
				var t time.Time
				switch ee.Action {
				case "create":
					t = ee.Entry.CreatedAt
				case "delete":
					t = ee.Entry.DeletedAt.Time
				default:
					t = ee.Entry.UpdatedAt
				}

				sendEvent(c, t, ee.Action, string(entry))
			}
		}
	}
}
