package cli

import (
	"bytes"
	"errors"
	"homework-go/homework-7/internal/inventory"
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

	if err := app.Run([]string{"use", "-id", "999"}); err == nil {
		t.Error("Ожидали ошибку при поиске несуществующего ID")
	}
}
