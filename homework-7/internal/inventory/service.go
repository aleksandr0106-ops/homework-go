package inventory

import (
	"fmt"
	"homework-go/homework-7/internal/item"
	"strings"
)

type Storage interface {
	Load() ([]item.Item, error)
	Save(items []item.Item) error
}

type Service struct {
	store  Storage
	items  []item.Item
	nextID int
}

func NewService(store Storage) *Service {
	return &Service{
		store:  store,
		items:  []item.Item{},
		nextID: 1,
	}
}

func (s *Service) Load() error {
	items, err := s.store.Load()
	if err != nil {
		return err
	}

	s.items = items

	maxID := 0
	for _, it := range items {
		if it.ID() > maxID {
			maxID = it.ID()
		}
	}
	s.nextID = maxID + 1
	return nil
}

func (s *Service) Save() error {
	return s.store.Save(s.items)
}

func (s *Service) Items() []item.Item {
	return s.items
}

func (s *Service) Add(name string, weight float64) (int, error) {
	if strings.TrimSpace(name) == "" {
		return 0, fmt.Errorf("имя не может быть пустым")
	}
	if strings.ContainsAny(name, "|\n\r") {
		return 0, fmt.Errorf("имя содержит недопустимые символы")
	}

	id := s.nextID
	it := item.NewWeapon(id, name, weight, 10, 5)
	s.items = append(s.items, it)
	s.nextID++

	return id, nil
}

func (s *Service) Delete(id int) error {
	for idx, it := range s.items {
		if it.ID() == id {
			s.items = append(s.items[:idx], s.items[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("предмет с ID %d не найден", id)
}

func (s *Service) FindByID(id int) (item.Item, error) {
	for _, it := range s.items {
		if it.ID() == id {
			return it, nil
		}
	}
	return nil, fmt.Errorf("предмет с ID %d не найден", id)
}
