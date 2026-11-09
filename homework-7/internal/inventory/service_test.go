package inventory

import (
	"homework-go/homework-7/internal/item"
	"testing"
)

type mockStorage struct {
	items []item.Item
	err   error
}

func (m *mockStorage) Load() ([]item.Item, error) {
	return m.items, m.err
}

func (m *mockStorage) Save(items []item.Item) error {
	m.items = items
	return m.err
}

func TestService_AddAndFind(t *testing.T) {
	t.Parallel()

	store := &mockStorage{}
	srv := NewService(store)

	id, err := srv.Add("Меч", 5.0)
	if err != nil {
		t.Fatalf("Ожидали успех, получили: %v", err)
	}
	if id != 1 {
		t.Errorf("Ожидали ID 1, получили %d", id)
	}

	found, err := srv.FindByID(1)
	if err != nil {
		t.Fatalf("Ожидали найти предмет, получили: %v", err)
	}
	if found.Name() != "Меч" {
		t.Errorf("Ожидали 'Меч', получили '%s'", found.Name())
	}

	_, err = srv.Add("   ", 5.0)
	if err == nil {
		t.Error("Ожидали ошибку для пустого имени, получили nil")
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	store := &mockStorage{}
	srv := NewService(store)
	srv.Add("Меч", 5.0)

	err := srv.Delete(1)
	if err != nil {
		t.Errorf("Ожидали успех, получили: %v", err)
	}
	if len(srv.Items()) != 0 {
		t.Errorf("Ожидали пустой инвентарь, осталось %d", len(srv.Items()))
	}

	err = srv.Delete(99)
	if err == nil {
		t.Error("Ожидали ошибку при удалении несуществующего ID")
	}
}

func TestService_Validation(t *testing.T) {
	t.Parallel()

	store := &mockStorage{}
	srv := NewService(store)

	_, err := srv.Add("Меч|Супер", 5.0)
	if err == nil {
		t.Error("Ожидали ошибку из-за символа '|', получили nil")
	}
}

func TestService_LoadSave(t *testing.T) {
	t.Parallel()

	store := &mockStorage{
		items: []item.Item{item.NewWeapon(5, "Топор", 8.0, 15, 5)},
	}
	srv := NewService(store)

	err := srv.Load()
	if err != nil {
		t.Fatalf("Ошибка Load: %v", err)
	}
	if len(srv.Items()) != 1 {
		t.Fatalf("Ожидали 1 предмет, получили %d", len(srv.Items()))
	}

	id, _ := srv.Add("Кинжал", 1.0)
	if id != 6 {
		t.Errorf("Ожидали ID 6, получили %d", id)
	}

	err = srv.Save()
	if err != nil {
		t.Errorf("Ошибка Save: %v", err)
	}
}
