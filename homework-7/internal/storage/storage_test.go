package storage

import (
	"homework-go/homework-7/internal/item"
	"os"
	"path/filepath"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()
	items := []item.Item{item.NewWeapon(1, "Меч", 5.0, 10, 5)}

	if err := store.Save(items); err != nil {
		t.Fatalf("Ошибка Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil || len(loaded) != 1 {
		t.Fatalf("Ошибка Load: %v, len=%d", err, len(loaded))
	}
	if loaded[0].Name() != "Меч" {
		t.Errorf("Ожидали 'Меч', получили '%s'", loaded[0].Name())
	}
}

func TestFileStorage_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_inventory.txt")

	store := NewFileStorage(filePath)
	items := []item.Item{
		item.NewWeapon(1, "Меч", 5.0, 10, 5),
		item.NewPotion(2, "Зелье", 0.5, "HP", 3),
	}

	err := store.Save(items)
	if err != nil {
		t.Fatalf("Ошибка сохранения: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Файл не был создан на диске!")
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Ошибка загрузки: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("Ожидали 2 предмета, получили %d", len(loaded))
	}
	if loaded[1].Name() != "Зелье" {
		t.Errorf("Ожидали 'Зелье', получили '%s'", loaded[1].Name())
	}
}

func TestFileStorage_LoadMissingFile(t *testing.T) {
	store := NewFileStorage("missing_file_12345.txt")
	loaded, err := store.Load()

	if err != nil {
		t.Fatalf("Ожидали nil, получили ошибку: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("Ожидали пустой список, получили %d", len(loaded))
	}
}
