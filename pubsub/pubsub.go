// Package pubsub
package pubsub

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PubSub[K comparable, T any] struct {
	buffer int
	m      *sync.Mutex
	store  map[K]map[uuid.UUID]chan T
}

func New[K comparable, T any](buffer int) PubSub[K, T] {
	return PubSub[K, T]{
		m:      &sync.Mutex{},
		store:  map[K]map[uuid.UUID]chan T{},
		buffer: buffer,
	}
}

// getOrInit gets or initializes the store part for the Key k.
// This is for internal use and has therefor the Mutex not locked, so use with care.
func (ps PubSub[K, T]) getOrInit(k K) map[uuid.UUID]chan T {
	subs, ok := ps.store[k]

	// Key is not present, so init it.
	if !ok {
		ps.store[k] = map[uuid.UUID]chan T{}
		subs = ps.store[k]
	}

	return subs
}

func (ps PubSub[K, T]) Publish(k K, t T) error {
	ps.m.Lock()
	defer ps.m.Unlock()

	subs := ps.getOrInit(k)

	// Send to all subs parallel and with timeout
	var wg sync.WaitGroup
	var muStale sync.Mutex
	stale := []uuid.UUID{}
	wg.Add(len(subs))
	for id, sub := range subs {
		go func(id uuid.UUID, ch chan T) {
			defer wg.Done()
			select {
			case ch <- t:
			case <-time.After(100 * time.Millisecond):
				slog.Error("subscriber too slow to keep up, disconnecting it", "sub", id)
				muStale.Lock()
				stale = append(stale, id)
				muStale.Unlock()
			}
		}(id, sub)
	}
	wg.Wait()

	// A subscriber that could not keep up would otherwise silently miss this
	// event forever. Disconnect it instead, so its SSE handler returns, the
	// client's EventSource sees the closed connection and reconnects, and a
	// full resync (see entryList) picks up whatever was missed.
	for _, id := range stale {
		close(subs[id])
		delete(subs, id)
	}
	return nil
}

func (ps PubSub[K, T]) Subscribe(k K) (uuid.UUID, chan T, error) {
	ps.m.Lock()
	defer ps.m.Unlock()

	subs := ps.getOrInit(k)

	// add new sub to subs
	id, err := uuid.NewUUID()
	if err != nil {
		return id, nil, fmt.Errorf("unable to create uuid, %w", err)
	}
	subs[id] = make(chan T, ps.buffer)

	return id, subs[id], nil
}

func (ps PubSub[K, T]) Unsubscribe(k K, id uuid.UUID) {
	ps.m.Lock()
	defer ps.m.Unlock()

	subs, ok := ps.store[k]

	// Key is not present, so there is nothing to unsubscribe
	if !ok {
		return
	}

	// Publish may already have disconnected this subscriber (slow consumer),
	// in which case it is no longer in subs. Closing a nil/missing channel
	// again would panic, so only close it if it is still there.
	ch, ok := subs[id]
	if !ok {
		return
	}
	close(ch)
	delete(subs, id)
}
