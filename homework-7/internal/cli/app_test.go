package cli

import (
	"bytes"
	"errors"
	"homework-go/homework-7/internal/inventory"
	"homework-go/homework-7/internal/item"
	"homework-go/homework-7/internal/storage"
	"strings"
	"testing"
)

func setupApp() (*App, *bytes.Buffer, *bytes.Buffer) {
	store := storage.NewMemoryStorage()
	inv := inventory.NewService(store)

	var out, errOut bytes.Buffer
	app := NewApp(inv, &out, &errOut)

	return app, &out, &errOut
}

type errorMockStorage struct {
	loadErr   error
	saveErr   error
	loadCalls int
	saveCalls int
}

func (m *errorMockStorage) Load() ([]item.Item, error) {
	m.loadCalls++
	return nil, m.loadErr
}

func (m *errorMockStorage) Save(items []item.Item) error {
	m.saveCalls++
	return m.saveErr
}

func TestApp_UnknownCommand(t *testing.T) {
	app, _, _ := setupApp()

	err := app.Run([]string{"unknown"})

	if !errors.Is(err, ErrUsage) {
		t.Errorf("Ожидали ErrUsage, получили %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Ожидали текст 'unknown command' в ошибке, получили: %v", err)
	}
}

func TestApp_AddAndList(t *testing.T) {
	app, out, _ := setupApp()

	err := app.Run([]string{"add", "-name", "Меч", "-weight", "5.5"})
	if err != nil {
		t.Fatalf("Ошибка добавления: %v", err)
	}
	if !strings.Contains(out.String(), "Добавлен предмет с ID: 1") {
		t.Errorf("Неверный вывод после добавления: %s", out.String())
	}

	out.Reset()

	err = app.Run([]string{"list"})
	if err != nil {
		t.Fatalf("Ошибка вывода списка: %v", err)
	}
	if !strings.Contains(out.String(), "[1] Меч (вес: 5.50)") {
		t.Errorf("Неверный вывод списка: %s", out.String())
	}
}

func TestApp_UseAndDelete(t *testing.T) {
	app, out, _ := setupApp()

	_ = app.Run([]string{"add", "-name", "Меч", "-weight", "5.5"})
	out.Reset()

	err := app.Run([]string{"use", "-id", "1"})
	if err != nil {
		t.Fatalf("Ошибка использования: %v", err)
	}
	if !strings.Contains(out.String(), "Взмах оружием: Меч") {
		t.Errorf("Неверный вывод Use: %s", out.String())
	}
	out.Reset()

	err = app.Run([]string{"delete", "-id", "1"})
	if err != nil {
		t.Fatalf("Ошибка удаления: %v", err)
	}
	if !strings.Contains(out.String(), "Предмет удален") {
		t.Errorf("Неверный вывод Delete: %s", out.String())
	}
}

func TestApp_EdgeCases(t *testing.T) {
	app, _, _ := setupApp()

	if err := app.Run([]string{}); !errors.Is(err, ErrUsage) {
		t.Error("Ожидали ErrUsage при пустых аргументах")
	}

	if err := app.Run([]string{"add", "-weight", "5.0"}); !errors.Is(err, ErrUsage) {
		t.Error("Ожидали ErrUsage при отсутствии -name")
	}

	if err := app.Run([]string{"add", "-unknown"}); !errors.Is(err, ErrUsage) {
		t.Error("Ожидали ErrUsage при неизвестном флаге")
	}

	if err := app.Run([]string{"use", "-id", "999"}); err == nil {
		t.Error("Ожидали ошибку при поиске несуществующего ID")
	}
}

func TestApp_LoadError(t *testing.T) {
	loadErr := errors.New("load failed")
	store := &errorMockStorage{loadErr: loadErr}
	inv := inventory.NewService(store)
	var out, errOut bytes.Buffer
	app := NewApp(inv, &out, &errOut)

	err := app.Run([]string{"add", "-name", "Меч", "-weight", "5"})

	if !errors.Is(err, loadErr) {
		t.Errorf("Ожидали ошибку загрузки, получили %v", err)
	}
	if store.saveCalls != 0 {
		t.Errorf("Save не должен вызываться после ошибки Load! Вызван %d раз", store.saveCalls)
	}
}

func TestApp_SaveError(t *testing.T) {
	saveErr := errors.New("save failed")
	store := &errorMockStorage{saveErr: saveErr}
	inv := inventory.NewService(store)
	var out, errOut bytes.Buffer
	app := NewApp(inv, &out, &errOut)

	err := app.Run([]string{"add", "-name", "Меч", "-weight", "5"})

	if !errors.Is(err, saveErr) {
		t.Errorf("Ожидали ошибку сохранения, получили %v", err)
	}
	if strings.Contains(out.String(), "Добавлен предмет") {
		t.Errorf("Сообщение об успехе не должно печататься, если Save упал! Вывод: %s", out.String())
	}
}
