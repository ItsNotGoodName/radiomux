package core

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/ItsNotGoodName/radiomux/internal"
)

func GenerateToken() (string, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func PlayerExists(ctx context.Context, store PlayerStore, id int64) (bool, error) {
	_, err := store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return false, nil
		} else {
			return true, nil
		}
	}

	return true, nil
}

func PlayerIDS(ctx context.Context, store PlayerStore) ([]int64, error) {
	items, err := store.List(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].ID)
	}

	return ids, nil
}
