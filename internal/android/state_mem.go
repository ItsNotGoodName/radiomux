package android

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/ItsNotGoodName/radiomux/internal"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/pkg/diff"
	"github.com/rs/zerolog/log"
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

func NewStateMemStore(statePubSub StatePubSub, bus core.Bus, playerStore core.PlayerStore) *StateMemStore {
	s := &StateMemStore{
		statePubSub: statePubSub,
		playerStore: playerStore,
		statesMu:    sync.Mutex{},
		states:      []State{},
	}

	_, err := s.sync(context.TODO())
	if err != nil {
		panic(err)
	}

	bus.OnPlayerCreated(s.onPlayerCreated)
	bus.OnPlayerDeleted(s.onPlayerDeleted)

	return s
}

type StateMemStore struct {
	statePubSub StatePubSub
	playerStore core.PlayerStore

	statesMu sync.Mutex
	states   []State
}

func (s *StateMemStore) onPlayerCreated(ctx context.Context, evt core.EventPlayerCreated) error {
	s.statesMu.Lock()
	_, found := s.get(evt.ID)
	if found {
		s.statesMu.Unlock()
	} else {
		state := NewState(evt.ID)
		s.states = append(s.states, state)
		s.statesMu.Unlock()

		s.statePubSub.Broadcast(evt.ID, diff.ChangedAll)
	}

	return nil
}

func (s *StateMemStore) onPlayerDeleted(ctx context.Context, evt core.EventPlayerDeleted) error {
	s.statesMu.Lock()
	s.states = slices.DeleteFunc(s.states, func(s State) bool { return s.ID == evt.ID })
	s.statesMu.Unlock()

	s.statePubSub.Broadcast(evt.ID, diff.ChangedAll)

	return nil
}

func (s *StateMemStore) get(id int64) (int, bool) {
	return slices.BinarySearchFunc(s.states, State{ID: id}, func(s1, s2 State) int {
		return cmp.Compare(s1.ID, s2.ID)
	})
}

func (s *StateMemStore) sync(ctx context.Context) (bool, error) {
	ids, err := core.PlayerIDS(ctx, s.playerStore)
	if err != nil {
		return false, err
	}
	slices.SortFunc(ids, func(a, b int64) int {
		return cmp.Compare(a, b)
	})

	// Create merge new states with old states
	var createdStateIDS []int64
	var states []State
	for _, id := range ids {
		_, found := slices.BinarySearchFunc(s.states, State{ID: id}, func(s1, s2 State) int {
			return cmp.Compare(s1.ID, s2.ID)
		})
		if !found {
			state := NewState(id)
			states = append(states, state)
			createdStateIDS = append(createdStateIDS, id)
		}
	}

	// Check which states were deleted
	var deletedStateIDS []int64
	for _, old := range s.states {
		_, found := slices.BinarySearchFunc(states, State{ID: old.ID}, func(s1, s2 State) int {
			return cmp.Compare(s1.ID, s2.ID)
		})
		if !found {
			deletedStateIDS = append(deletedStateIDS, old.ID)
		}
	}

	s.states = states

	return len(createdStateIDS)+len(deletedStateIDS) > 0, nil
}

func (s *StateMemStore) List() []State {
	s.statesMu.Lock()
	states := slices.Clone(s.states)
	s.statesMu.Unlock()

	return states
}

func (s *StateMemStore) Get(id int64) (State, error) {
	s.statesMu.Lock()
	index, found := s.get(id)
	if !found {
		s.statesMu.Unlock()
		return State{}, internal.ErrNotFound
	}
	s.statesMu.Unlock()

	return s.states[index], nil
}

func (s *StateMemStore) Update(id int64, fn func(state State, changed diff.Changed) (State, diff.Changed)) error {
	s.statesMu.Lock()
	index, found := s.get(id)
	if !found {
		s.statesMu.Unlock()
		return internal.ErrNotFound
	}
	state := s.states[index]

	state, changed := fn(state, diff.ChangedNone)
	s.states[index] = state
	s.statesMu.Unlock()

	s.statePubSub.Broadcast(id, changed)

	return nil
}

var _ StateStore = (*StateMemStore)(nil)
