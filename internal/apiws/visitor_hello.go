package apiws

import (
	"context"

	"github.com/ItsNotGoodName/radiomux/internal/openapi"
)

type helloVisitor struct {
	first bool
}

func newHelloVisitor() *helloVisitor {
	return &helloVisitor{
		first: true,
	}
}

// HasMore implements visitor.
func (h *helloVisitor) HasMore() bool {
	return h.first
}

// Visit implements visitor.
func (h *helloVisitor) Visit(ctx context.Context) ([]byte, error) {
	if !h.first {
		return nil, errVisitorEmpty
	}
	h.first = false

	return openapi.ConvertNotification(openapi.Notification{
		Title:       "WebSocket",
		Description: "Connected to server.",
	})
}

var _ visitor = (*helloVisitor)(nil)
