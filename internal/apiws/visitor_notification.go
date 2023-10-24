package apiws

import (
	"context"
)

type bufferVisitor struct {
	events chan []byte
}

func newBufferVisitor(count int) *bufferVisitor {
	return &bufferVisitor{
		events: make(chan []byte, count),
	}
}

func (h *bufferVisitor) Push(event []byte) error {
	select {
	case h.events <- event:
	default:
	}

	return nil
}

// HasMore implements visitor.
func (h *bufferVisitor) HasMore() bool {
	return len(h.events) > 0
}

// Visit implements visitor.
func (h *bufferVisitor) Visit(ctx context.Context) ([]byte, error) {
	select {
	case event := <-h.events:
		return event, nil
	default:
		return nil, errVisitorEmpty
	}
}

var _ visitor = (*bufferVisitor)(nil)
