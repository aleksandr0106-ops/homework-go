package storage

import "homework-go/homework-7/internal/item"

type MemoryStorage struct {
	items []item.Item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{items: []item.Item{}}
}

func (s *MemoryStorage) Save(items []item.Item) error {
	s.items = make([]item.Item, len(items))
	for i, it := range items {
		s.items[i] = it.Clone()
	}
	return nil
}

func (s *MemoryStorage) Load() ([]item.Item, error) {
	copied := make([]item.Item, len(s.items))
	for i, it := range s.items {
		copied[i] = it.Clone()
	}
	return copied, nil
}
