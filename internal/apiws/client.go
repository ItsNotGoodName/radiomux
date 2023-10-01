package apiws

import (
	"errors"
)

var ErrVisitorEmpty = errors.New("visitor empty")

type Visitors interface {
	Visit() ([]byte, error)
}

type Client struct {
	visitors []Visitors
}

func NewClient(visitors ...Visitors) Client {
	return Client{
		visitors: visitors,
	}
}

func (c Client) Flush() ([]byte, error) {
	for _, v := range c.visitors {
		data, err := v.Visit()
		if err != nil {
			if errors.Is(err, ErrVisitorEmpty) {
				continue
			}
			return nil, err
		}

		return data, nil
	}

	return nil, ErrVisitorEmpty
}
