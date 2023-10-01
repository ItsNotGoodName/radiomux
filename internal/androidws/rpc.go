package androidws

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// rpcTimeout is the maximum time to wait for a rpc request.
const rpcTimeout = 5 * time.Second

type rpcMailbox struct {
	id int32

	rpcMu      sync.Mutex
	rpcReplies map[int32]chan<- struct{}
}

func newRPCMailbox() *rpcMailbox {
	return &rpcMailbox{
		rpcMu:      sync.Mutex{},
		rpcReplies: map[int32]chan<- struct{}{},
		id:         0,
	}
}

func (r *rpcMailbox) NextID() int32 {
	return atomic.AddInt32(&r.id, 1)
}

func (r *rpcMailbox) Create(id int32) (<-chan struct{}, func()) {
	replyC := make(chan struct{}, 1)

	r.rpcMu.Lock()
	r.rpcReplies[id] = replyC
	r.rpcMu.Unlock()

	return replyC, func() {
		r.rpcMu.Lock()
		delete(r.rpcReplies, id)
		r.rpcMu.Unlock()
	}
}

func (r *rpcMailbox) Get(id int32) (chan<- struct{}, bool) {
	r.rpcMu.Lock()
	replyC, found := r.rpcReplies[id]
	r.rpcMu.Unlock()

	return replyC, found
}

func rpcWait(ctx context.Context, replyC <-chan struct{}) error {
	t := time.NewTimer(rpcTimeout)
	select {
	case <-t.C:
		return errors.New("player did not reply")
	case <-ctx.Done():
		if !t.Stop() {
			<-t.C
		}
		return ctx.Err()
	case <-replyC:
		if !t.Stop() {
			<-t.C
		}
		return nil
	}
}
