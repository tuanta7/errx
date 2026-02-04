package main

import (
	"encoding/json"
	"errors"

	"github.com/tuanta7/errx"
	"github.com/tuanta7/errx/predefined"
)

type Repository struct {
	cache *Cache
}

func NewRepository(cache *Cache) *Repository {
	return &Repository{cache: cache}
}

func (r *Repository) GetCounter(key string) (*Counter, error) {
	value, err := r.cache.Get(key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, predefined.ErrRecordNotFound
		}

		return nil, err
	}

	counter := &Counter{}
	err = json.Unmarshal(value, counter)
	if err != nil {
		return nil, errx.New("failed to unmarshal counter", err)
	}

	return counter, nil
}

func (r *Repository) SetCounter(key string, counter *Counter) error {
	data, err := json.Marshal(counter)
	if err != nil {
		return errx.New("failed to marshal counter", err)
	}

	r.cache.Set(key, data)
	return nil
}
