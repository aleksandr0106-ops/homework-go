package item

import (
	"errors"
	"testing"
)

func TestWeaponUse(t *testing.T) {
	tests := []struct {
		name        string
		initialDur  int
		expectedErr error
		expectedDur int
	}{
		{"Успешный удар", 2, nil, 1},
		{"Последний удар", 1, nil, 0},
		{"Удар сломанным", 0, ErrBroken, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := NewWeapon(1, "Меч", 5.0, 10, tt.initialDur)
			_, err := w.Use()

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Ожидали ошибку %v, получили %v", tt.expectedErr, err)
			}
			if w.Durability != tt.expectedDur {
				t.Errorf("Ожидали прочность %d, получили %d", tt.expectedDur, w.Durability)
			}
		})
	}
}
func TestPotionUse(t *testing.T) {
	p := NewPotion(1, "Зелье", 0.5, "HP", 1)

	_, err := p.Use()
	if err != nil {
		t.Fatalf("Ожидали успех, получили ошибку: %v", err)
	}
	if p.Charges != 0 {
		t.Errorf("Ожидали 0 зарядов, осталось: %d", p.Charges)
	}

	_, err = p.Use()
	if !errors.Is(err, ErrExhausted) {
		t.Errorf("Ожидали ошибку ErrExhausted, получили %v", err)
	}
}
func TestClone(t *testing.T) {
	original := NewWeapon(1, "Меч", 5.0, 10, 5)
	clone := original.Clone().(*Weapon)

	original.Durability = 1

	if clone.Durability != 5 {
		t.Errorf("Клон изменился вслед за оригиналом! Ожидали 5, получили %d", clone.Durability)
	}
}
func TestArmorUse(t *testing.T) {
	a := NewArmor(1, "Щит", 5.0, 10, 2)

	_, err := a.Use()
	if err != nil || a.Durability != 1 {
		t.Errorf("Ожидали успех и прочность 1, получили err=%v, dur=%d", err, a.Durability)
	}

	_, err = a.Use()
	if err != nil || a.Durability != 0 {
		t.Errorf("Ожидали успех и прочность 0, получили err=%v, dur=%d", err, a.Durability)
	}

	_, err = a.Use()
	if !errors.Is(err, ErrBroken) {
		t.Errorf("Ожидали ErrBroken, получили %v", err)
	}
}
func TestItemProperties(t *testing.T) {
	w := NewWeapon(1, "Меч", 5.5, 10, 5)

	if w.ID() != 1 || w.Name() != "Меч" || w.Weight() != 5.5 {
		t.Error("Геттеры возвращают неверные значения")
	}

	w.SetID(42)
	if w.ID() != 42 {
		t.Errorf("SetID не сработал, ожидали 42, получили %d", w.ID())
	}
}
