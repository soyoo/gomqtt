// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"sync"
	"time"
)

/* message */

func newMessage(topic string, payload []byte) *Message {
	return &Message{
		Topic:   topic,
		Payload: payload,
	}
}

/* futureStore */

// a futureStore is used to store active Futures
type futureStore struct {
	sync.Mutex

	store map[uint16]Future
}

// newFutureStore will create a new futureStore
func newFutureStore() *futureStore {
	return &futureStore{
		store: make(map[uint16]Future),
	}
}

// put will save a Future to the store
func (s *futureStore) put(id uint16, future Future) {
	s.Lock()
	defer s.Unlock()

	s.store[id] = future
}

// get will retrieve a Future from the store
func (s *futureStore) get(id uint16) Future {
	s.Lock()
	defer s.Unlock()

	return s.store[id]
}

// del will remove a Future from the store
func (s *futureStore) del(id uint16) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
}

// return a slice with all stored futures
func (s *futureStore) all() []Future {
	s.Lock()
	defer s.Unlock()

	all := make([]Future, len(s.store))

	i := 0
	for _, future := range s.store {
		all[i] = future
		i++
	}

	return all
}

// will wait until all futures have completed and removed or timeout is reached
func (s *futureStore) await(timeout time.Duration) error {
	stop := time.Now().Add(timeout)

	for {
		// get futures
		futures := s.all()

		// return if no futures are left
		if len(futures) == 0 {
			return nil
		}

		// wait for next future to complete
		err := futures[0].Wait(stop.Sub(time.Now()))
		if err != nil {
			return err
		}
	}
}

/* counter */

// a counter keeps track of packet ids
type counter struct {
	sync.Mutex

	id uint16
}

// newCounter will return a new counter
func newCounter() *counter {
	return &counter{}
}

// next will generate the next packet id
func (c *counter) next() uint16 {
	c.Lock()
	defer func() {
		c.id++
		c.Unlock()
	}()

	return c.id
}

// resume will continue to count from the provided id
func (c *counter) resume(id uint16) {
	c.Lock()
	defer c.Unlock()

	c.id = id + 1
}

/* state */

const (
	stateInitialized byte = iota
	stateConnecting
	stateConnected
	stateDisconnecting
	stateDisconnected
)

// a state keeps track of the clients current state
type state struct {
	sync.Mutex

	current byte
}

// create new state
func newState() *state {
	return &state{
		current: stateInitialized,
	}
}

// set will change to the specified state
func (s *state) set(state byte) {
	s.Lock()
	defer s.Unlock()

	s.current = state
}

// get will retrieve the current state
func (s *state) get() byte {
	s.Lock()
	defer s.Unlock()

	return s.current
}

/* tracker */

// a tracker keeps track of keep alive intervals
type tracker struct {
	sync.Mutex

	last    time.Time
	pings   uint8
	timeout time.Duration
}

// returns a new tracker
func newTracker(timeout time.Duration) *tracker {
	return &tracker{
		last:    time.Now(),
		timeout: timeout,
	}
}

// updates the tracker
func (t *tracker) reset() {
	t.Lock()
	defer t.Unlock()

	t.last = time.Now()
}

// returns the current time window
func (t *tracker) window() time.Duration {
	t.Lock()
	defer t.Unlock()

	return t.timeout - time.Since(t.last)
}

// mark ping
func (t *tracker) ping() {
	t.Lock()
	defer t.Unlock()

	t.pings++
}

// mark pong
func (t *tracker) pong() {
	t.Lock()
	defer t.Unlock()

	t.pings--
}

// returns if pings are pending
func (t *tracker) pending() bool {
	t.Lock()
	defer t.Unlock()

	return t.pings > 0
}
