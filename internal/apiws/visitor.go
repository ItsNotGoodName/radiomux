package apiws

import (
	"context"
	"errors"
)

var errVisitorEmpty = errors.New("visitor empty")

type visitor interface {
	Visit(ctx context.Context) ([]byte, error)
	HasMore() bool
}

type vistors struct {
	visitors []visitor
}

func newVisitors(visitors ...visitor) vistors {
	return vistors{
		visitors: visitors,
	}
}

func (c vistors) Visit(ctx context.Context) ([]byte, error) {
	for _, v := range c.visitors {
		data, err := v.Visit(ctx)
		if err != nil {
			if errors.Is(err, errVisitorEmpty) {
				continue
			}
			return nil, err
		}

		return data, nil
	}

	return nil, errVisitorEmpty
}

func (c vistors) HasMore() bool {
	for _, v := range c.visitors {
		if v.HasMore() {
			return true
		}
	}
	return false
}
