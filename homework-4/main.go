package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Item interface {
	Name() string
	Weight() float64
}

type Weapon struct {
	name       string
	weight     float64
	Damage     int
	Durability int
}

func (w *Weapon) Name() string    { return w.name }
func (w *Weapon) Weight() float64 { return w.weight }
func (w *Weapon) Use() string {
	if w.Durability <= 0 {
		return fmt.Sprintf("Оружие %s сломано!", w.name)
	}
	w.Durability--
	return fmt.Sprintf("Взмах оружием: %s (Урон: %d, Осталось прочности: %d)", w.name, w.Damage, w.Durability)
}

type Armor struct {
	name    string
	weight  float64
	Defense int
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) Use() string     { return fmt.Sprintf("Надеваем броню: %s", a.name) }

type Potion struct {
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string    { return p.name }
func (p *Potion) Weight() float64 { return p.weight }
func (p *Potion) Use() string {
	if p.Charges <= 0 {
		return fmt.Sprintf("Зелье %s пустое!", p.name)
	}
	p.Charges--
	return fmt.Sprintf("Выпиваем зелье: %s (%s). Осталось зарядов: %d", p.name, p.Effect, p.Charges)
}

var _ Item = (*Weapon)(nil)
var _ Item = (*Armor)(nil)
var _ Item = (*Potion)(nil)

func DescribeItem(i Item) string {
	if i == nil {
		return "Предмет отсутствует"
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "Предмет отсутствует"
	}

	return fmt.Sprintf("%s (вес: %.2f)", i.Name(), i.Weight())
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](items []T, transform func(T) R) []R {
	result := make([]R, 0, len(items))
	for _, item := range items {
		result = append(result, transform(item))
	}
	return result
}

func Find[T any](items []T, condition func(T) bool) (T, bool) {
	for _, item := range items {
		if condition(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

type Inventory struct {
	Items []Item
}

func (inv *Inventory) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) TotalWeight() float64 {
	var total float64
	for _, item := range inv.Items {
		total += item.Weight()
	}
	return total
}

func (inv *Inventory) Weapons() []*Weapon {
	weaponItems := Filter(inv.Items, func(i Item) bool {
		_, ok := i.(*Weapon)
		return ok
	})
	return Map(weaponItems, func(i Item) *Weapon {
		return i.(*Weapon)
	})
}

func (inv *Inventory) BrokenItems() []Item {
	return Filter(inv.Items, func(i Item) bool {
		switch v := i.(type) {
		case *Weapon:
			return v.Durability <= 0
		case *Potion:
			return v.Charges <= 0
		}
		return false
	})
}

func (inv *Inventory) ItemNames() []string {
	return Map(inv.Items, func(i Item) string {
		return i.Name()
	})
}

func (inv *Inventory) FindItemByName(name string) (Item, bool) {
	return Find(inv.Items, func(i Item) bool {
		return i.Name() == name
	})
}

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader)
}

func (w *Weapon) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Weapon|%s|%d|%d|%.2f\n", w.name, w.Damage, w.Durability, w.weight)
}

func (w *Weapon) Deserialize(reader io.Reader) {
	data, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) == 5 {
		w.name = parts[1]
		fmt.Sscanf(parts[2], "%d", &w.Damage)
		fmt.Sscanf(parts[3], "%d", &w.Durability)
		fmt.Sscanf(parts[4], "%f", &w.weight)
	}
}

func (a *Armor) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Armor|%s|%d|%.2f\n", a.name, a.Defense, a.weight)
}

func (a *Armor) Deserialize(reader io.Reader) {
	date, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(date)), "|")
	if len(parts) == 4 {
		a.name = parts[1]
		fmt.Sscanf(parts[2], "%d", &a.Defense)
		fmt.Sscanf(parts[3], "%f", &a.weight)
	}
}

func (p *Potion) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Potion|%s|%s|%d|%.2f\n", p.name, p.Effect, p.Charges, p.weight)
}

func (p *Potion) Deserialize(reader io.Reader) {
	data, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) == 5 {
		p.name = parts[1]
		p.Effect = parts[2]
		fmt.Sscanf(parts[3], "%d", &p.Charges)
		fmt.Sscanf(parts[4], "%f", &p.weight)
	}
}

func (inv *Inventory) Save(w io.Writer) {
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			storable.Serialize(w)
		}
	}
}

func (inv *Inventory) Load(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Weapon|") {
			w := &Weapon{}
			w.Deserialize(strings.NewReader(line))
			inv.AddItem(w)
		} else if strings.HasPrefix(line, "Armor|") {
			a := &Armor{}
			a.Deserialize(strings.NewReader(line))
			inv.AddItem(a)
		} else if strings.HasPrefix(line, "Potion|") {
			p := &Potion{}
			p.Deserialize(strings.NewReader(line))
			inv.AddItem(p)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения инвентаря:", err)
	}
}

func main() {
	inv := &Inventory{}

	sword := &Weapon{name: "Меч", weight: 5.0, Damage: 10, Durability: 2}
	shield := &Armor{name: "Щит", weight: 4.5, Defense: 5}
	shieldА := &Armor{name: "Щит грубой силы", weight: 5.5, Defense: 7}
	potion := &Potion{name: "Лечебное", weight: 0.5, Effect: "+50 HP", Charges: 1}
	brokenBow := &Weapon{name: "Сломанный лук", weight: 1.0, Damage: 5, Durability: 0}

	inv.AddItem(sword)
	inv.AddItem(shield)
	inv.AddItem(shieldА)

	inv.AddItem(potion)
	inv.AddItem(brokenBow)

	fmt.Println("1. Использование")
	fmt.Println(sword.Use())
	fmt.Println(sword.Use())
	fmt.Println(sword.Use())

	fmt.Println(potion.Use())
	fmt.Println(potion.Use())

	fmt.Println("\n 2. Описание")
	fmt.Println(DescribeItem(sword))
	fmt.Println(DescribeItem(nil))
	var nilWeapon *Weapon
	fmt.Println(DescribeItem(nilWeapon))

	fmt.Println("\n 3. Оружие")
	for _, w := range inv.Weapons() {
		fmt.Printf("%s (Урон: %d)\n", w.Name(), w.Damage)
	}

	fmt.Println("\n 4. Сломанные предметы")
	for _, b := range inv.BrokenItems() {
		fmt.Println(b.Name())
	}

	fmt.Println("\n 5. Названия")
	names := inv.ItemNames()
	for _, name := range names {
		fmt.Println(name)
	}

	fmt.Println("\n 6. Поиск по имени")
	found, ok := inv.FindItemByName("Щит")
	if ok {
		fmt.Printf("Найдено: %s\n", found.Name())
	} else {
		fmt.Println("Не найдено")
	}

	fmt.Println("\n 7. Бонус: Сохранение и загрузка")
	file, _ := os.Create("inventory.txt")
	inv.Save(file)
	file.Close()

	newInv := &Inventory{}
	loadFile, _ := os.Open("inventory.txt")
	newInv.Load(loadFile)
	loadFile.Close()

	fmt.Println("Загруженные предметы:")
	for _, item := range newInv.Items {
		switch v := item.(type) {
		case *Weapon:
			fmt.Printf("- Оружие: %s (Урон: %d, Прочность: %d, Вес: %.2f)\n",
				v.Name(), v.Damage, v.Durability, v.Weight())
		case *Armor:
			fmt.Printf("- Броня: %s (Защита: %d, Вес: %.2f)\n",
				v.Name(), v.Defense, v.Weight())
		case *Potion:
			fmt.Printf("- Зелье: %s (Эффект: %s, Заряды: %d, Вес: %.2f)\n",
				v.Name(), v.Effect, v.Charges, v.Weight())
		default:
			fmt.Printf("- Предмет: %s (Вес: %.2f)\n", item.Name(), item.Weight())
		}
	}
}
