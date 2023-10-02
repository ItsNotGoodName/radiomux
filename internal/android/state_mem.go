package android

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func NewStateMemPubSub() *StateMemPubSub {
	return &StateMemPubSub{
		mu:          sync.Mutex{},
		id:          0,
		subscribers: make(map[int32]chan StateChange),
	}
}

type StateMemPubSub struct {
	mu          sync.Mutex
	id          int32
	subscribers map[int32]chan StateChange
}

func (s *StateMemPubSub) Broadcast(id int64, changed diff.Changed) {
	if changed == diff.ChangedNone {
		return
	}

	s.mu.Lock()
	for _, v := range s.subscribers {
		select {
		case v <- StateChange{ID: id, Changed: changed}:
		}
	}
	s.mu.Unlock()
}

func (s *StateMemPubSub) Subscribe() (<-chan StateChange, func()) {
	msg := make(chan StateChange)

	id := atomic.AddInt32(&s.id, 1)

	s.mu.Lock()
	log.Debug().Str("package", "android").Int32("id", id).Msg("Subscribe to state")
	s.subscribers[id] = msg
	s.mu.Unlock()

	return msg, sync.OnceFunc(func() {
		s.mu.Lock()
		log.Debug().Str("package", "android").Int32("id", id).Msg("Unsubscribe from state")
		delete(s.subscribers, id)
		close(msg)
		s.mu.Unlock()
	})
}

var _ StatePubSub = (*StateMemPubSub)(nil)

func NewStateMemStore(statePubSub StatePubSub, bus core.Bus) (*StateMemStore, func()) {
	s := &StateMemStore{
		statePubSub:    statePubSub,
		statesMu:       sync.Mutex{},
		states:         []State{},
		deletedPlayers: make(map[int64]struct{}),
	}

	return s, bus.OnPlayerDeleted(s.onPlayerDeleted)
}

type StateMemStore struct {
	statePubSub StatePubSub

	statesMu sync.Mutex
	states   []State
	// TODO: do a db call instead of keeping track of deleted players
	deletedPlayers map[int64]struct{}
}

func (s *StateMemStore) onPlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) error {
	s.statesMu.Lock()
	s.states = lo.Filter(s.states, func(s State, i int) bool { return s.ID != evt.ID })
	s.deletedPlayers[evt.ID] = struct{}{}
	s.statesMu.Unlock()

	s.statePubSub.Broadcast(evt.ID, diff.ChangedAll)

	return nil
}

func (s *StateMemStore) List() []State {
	s.statesMu.Lock()
	states := make([]State, 0, len(s.states))
	for _, state := range s.states {
		states = append(states, state)
	}
	s.statesMu.Unlock()

	return states
}

func (s *StateMemStore) get(id int64) (State, int, error) {
	if _, found := s.deletedPlayers[id]; found {
		return State{}, 0, internal.ErrNotFound
	}

	for i := range s.states {
		if s.states[i].ID == id {
			return s.states[i], i, nil
		}
	}

	return State{}, -1, nil
}

func (s *StateMemStore) Get(id int64) (State, error) {
	s.statesMu.Lock()
	state, index, err := s.get(id)
	if err != nil {
		s.statesMu.Unlock()
		return State{}, err
	}
	s.statesMu.Unlock()

	if index == -1 {
		return State{}, internal.ErrNotFound
	}

	return state, nil
}

func (s *StateMemStore) Update(id int64, fn func(state State, changed diff.Changed) (State, diff.Changed)) error {
	s.statesMu.Lock()
	state, index, err := s.get(id)
	if err != nil {
		s.statesMu.Unlock()
		return err
	}
	if index == -1 {
		state = NewState(id)
		s.states = append(s.states, state)
		index = len(s.states) - 1
	}

	state, changed := fn(state, diff.ChangedNone)
	s.states[index] = state
	s.statesMu.Unlock()

	s.statePubSub.Broadcast(id, changed)

	return nil
}

var _ StateStore = (*StateMemStore)(nil)
