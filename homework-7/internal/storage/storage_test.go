package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"homework-go/homework-7/internal/item"
)

func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage()
	originalWeapon := item.NewWeapon(1, "Меч", 5.0, 10, 5)
	items := []item.Item{originalWeapon}

	if err := store.Save(items); err != nil {
		t.Fatalf("Ошибка Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil || len(loaded) != 1 {
		t.Fatalf("Ошибка Load: %v, len=%d", err, len(loaded))
	}

	originalWeapon.Durability = 1
	if loaded[0].(*item.Weapon).Durability != 5 {
		t.Errorf("Копия в памяти изменилась вместе с оригиналом! Ожидали 5, получили %d", loaded[0].(*item.Weapon).Durability)
	}

	loaded[0].(*item.Weapon).Durability = 3
	if originalWeapon.Durability != 1 {
		t.Errorf("Оригинал изменился вместе с копией! Ожидали 1, получили %d", originalWeapon.Durability)
	}
}

func TestFileStorage_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_inventory.txt")
	store := NewFileStorage(filePath)

	weapon := item.NewWeapon(1, "Меч", 5.5, 10, 4)
	armor := item.NewArmor(2, "Щит", 3.0, 7, 3)
	potion := item.NewPotion(3, "Зелье", 0.5, "HP", 2)
	items := []item.Item{weapon, armor, potion}

	if err := store.Save(items); err != nil {
		t.Fatalf("Ошибка сохранения: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Ошибка загрузки: %v", err)
	}
	if len(loaded) != 3 {
		t.Fatalf("Ожидали 3 предмета, получили %d", len(loaded))
	}

	lw := loaded[0].(*item.Weapon)
	if lw.ID() != 1 || lw.Name() != "Меч" || lw.Weight() != 5.5 || lw.Damage != 10 || lw.Durability != 4 {
		t.Errorf("Оружие загрузилось с ошибками: %+v", lw)
	}

	la := loaded[1].(*item.Armor)
	if la.ID() != 2 || la.Name() != "Щит" || la.Weight() != 3.0 || la.Defense != 7 || la.Durability != 3 {
		t.Errorf("Броня загрузилась с ошибками: %+v", la)
	}

	lp := loaded[2].(*item.Potion)
	if lp.ID() != 3 || lp.Name() != "Зелье" || lp.Weight() != 0.5 || lp.Effect != "HP" || lp.Charges != 2 {
		t.Errorf("Зелье загрузилось с ошибками: %+v", lp)
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

func TestFileStorage_CorruptedAndDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "corrupted.txt")

	content := `Weapon|1|Меч|10|5|5.0
	Weapon|1|Меч|10|5|5.0
	Weapon|2|Щит|abc|5|3.0
	Unknown|3|Зелье|10|5|1.0
	Armor|4|Броня|5|`

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	store := NewFileStorage(filePath)
	loaded, err := store.Load()

	if len(loaded) != 1 {
		t.Errorf("Ожидали 1 валидный предмет, загрузилось %d", len(loaded))
	}

	if err == nil {
		t.Fatal("Ожидали ошибку парсинга, получили nil")
	}

	var parseErr *ParseItemError
	if !errors.As(err, &parseErr) {
		t.Fatalf("Ожидали ParseItemError, получили %T", err)
	}

	if !strings.Contains(err.Error(), "line 2") || !strings.Contains(err.Error(), "duplicate ID 1") {
		t.Errorf("Не нашли ошибку дубликата в цепочке ошибок: %v", err)
	}

	if !strings.Contains(err.Error(), "line 3") {
		t.Errorf("Не нашли ошибку парсинга в строке 3: %v", err)
	}
}

func TestFileStorage_SaveError(t *testing.T) {
	store := NewFileStorage("non_existent_dir/inventory.txt")

	err := store.Save([]item.Item{item.NewWeapon(1, "Меч", 1.0, 10, 5)})
	if err == nil {
		t.Error("Ожидали ошибку записи, но Save вернул nil")
	}
}
