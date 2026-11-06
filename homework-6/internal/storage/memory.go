package storage

import "homework-go/homework-6/internal/item"

type MemoryStorage struct {
	items []item.Item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{items: []item.Item{}}
}

func (s *MemoryStorage) Save(items []item.Item) error {
	s.items = append([]item.Item(nil), items...)
	return nil
}

func (s *MemoryStorage) Load() ([]item.Item, error) {
	return append([]item.Item(nil), s.items...), nil
}
