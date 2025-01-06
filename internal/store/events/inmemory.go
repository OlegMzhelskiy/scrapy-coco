package events

import (
	"errors"

	"scraper_nike/internal/models"
)

var ErrNotFound = errors.New("event not found")

type MemoStore struct {
	store map[string]models.Event
}

func NewMemoryStore() *MemoStore {
	s := make(map[string]models.Event)

	return &MemoStore{
		store: s,
	}
}

func (s MemoStore) Save(e models.Event) error {
	s.store[e.Key()] = e

	return nil
}

func (s MemoStore) Get(key string) (models.Event, error) {
	event, ok := s.store[key]
	if !ok {
		return event, ErrNotFound
	}

	return event, nil
}
