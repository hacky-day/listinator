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
	wg.Add(len(subs))
	for k, sub := range subs {
		go func(k uuid.UUID, ch chan T) {
			select {
			case ch <- t:
			case <-time.After(100 * time.Millisecond):
				slog.Error("failed to send to sub", "sub", k)
			}
			wg.Done()
		}(k, sub)
	}
	wg.Wait()
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
	close(subs[id])
	delete(subs, id)
}
