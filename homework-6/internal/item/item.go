package item

import (
	"errors"
	"fmt"
)

var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
)

type Item interface {
	ID() int
	SetID(id int)
	Name() string
	Weight() float64
	Use() (string, error)
	IsBroken() bool
}

type Weapon struct {
	id         int
	name       string
	weight     float64
	Damage     int
	Durability int
}

func NewWeapon(id int, name string, weight float64, damage, durability int) *Weapon {
	return &Weapon{id: id, name: name, weight: weight, Damage: damage, Durability: durability}
}

func (w *Weapon) ID() int { return w.id }

func (w *Weapon) SetID(id int) { w.id = id }

func (w *Weapon) Name() string { return w.name }

func (w *Weapon) Weight() float64 { return w.weight }

func (w *Weapon) IsBroken() bool { return w.Durability <= 0 }

func (w *Weapon) Use() (string, error) {
	if w.IsBroken() {
		return "", ErrBroken
	}
	w.Durability--
	return fmt.Sprintf("Взмах оружием: %s", w.name), nil
}

type Armor struct {
	id         int
	name       string
	weight     float64
	Defense    int
	Durability int
}

func NewArmor(id int, name string, weight float64, defense, durability int) *Armor {
	return &Armor{id: id, name: name, weight: weight, Defense: defense, Durability: durability}
}

func (a *Armor) ID() int { return a.id }

func (a *Armor) SetID(id int) { a.id = id }

func (a *Armor) Name() string { return a.name }

func (a *Armor) Weight() float64 { return a.weight }

func (a *Armor) IsBroken() bool { return a.Durability <= 0 }

func (a *Armor) Use() (string, error) {
	if a.IsBroken() {
		return "", ErrBroken
	}
	a.Durability--
	return fmt.Sprintf("Надеваем броню: %s (Осталось прочности: %d)", a.name, a.Durability), nil
}

type Potion struct {
	id      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func NewPotion(id int, name string, weight float64, effect string, charges int) *Potion {
	return &Potion{id: id, name: name, weight: weight, Effect: effect, Charges: charges}
}

func (p *Potion) ID() int { return p.id }

func (p *Potion) SetID(id int) { p.id = id }

func (p *Potion) Name() string { return p.name }

func (p *Potion) Weight() float64 { return p.weight }

func (p *Potion) IsBroken() bool { return p.Charges <= 0 }

func (p *Potion) Use() (string, error) {
	if p.IsBroken() {
		return "", ErrExhausted
	}
	p.Charges--
	return fmt.Sprintf("Выпиваем зелье: %s (%s)", p.name, p.Effect), nil
}
