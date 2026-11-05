package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrBroken    = errors.New("item is broken")
	ErrExhausted = errors.New("item is exhausted")
	ErrUsage     = errors.New("usage error")
)

type Item interface {
	Name() string
	Weight() float64
	Use() (string, error)
}

type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("line %d: parse error %q: %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error { return e.Err }

type Weapon struct {
	ID         int
	name       string
	weight     float64
	Damage     int
	Durability int
}

func (w *Weapon) Name() string { return w.name }

func (w *Weapon) Weight() float64 { return w.weight }

func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", ErrBroken
	}
	w.Durability--
	return fmt.Sprintf("Взмах оружием: %s", w.name), nil
}

type Armor struct {
	ID         int
	name       string
	weight     float64
	Defense    int
	Durability int
}

func (a *Armor) Name() string { return a.name }

func (a *Armor) Weight() float64 { return a.weight }

func (a *Armor) Use() (string, error) {
	if a.Durability <= 0 {
		return "", ErrBroken
	}
	a.Durability--
	return fmt.Sprintf("Надеваем броню: %s (Осталось прочности: %d)", a.name, a.Durability), nil
}

type Potion struct {
	ID      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string { return p.name }

func (p *Potion) Weight() float64 { return p.weight }

func (p *Potion) Use() (string, error) {
	if p.Charges <= 0 {
		return "", ErrExhausted
	}
	p.Charges--
	return fmt.Sprintf("Выпиваем зелье: %s (%s)", p.name, p.Effect), nil
}

type Inventory struct {
	Items  []Item
	NextID int
}

func (inv *Inventory) AddItem(item Item) int {
	inv.NextID++
	id := inv.NextID
	switch v := item.(type) {
	case *Weapon:
		v.ID = id
	case *Armor:
		v.ID = id
	case *Potion:
		v.ID = id
	}
	inv.Items = append(inv.Items, item)
	return id
}

func (inv *Inventory) DeleteItem(id int) bool {
	for i, item := range inv.Items {
		var itemID int
		switch v := item.(type) {
		case *Weapon:
			itemID = v.ID
		case *Armor:
			itemID = v.ID
		case *Potion:
			itemID = v.ID
		}
		if itemID == id {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
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
		case *Armor:
			return v.Durability <= 0
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

type DangerousItem struct{ name string }

func (d *DangerousItem) Name() string { return d.name }

func (d *DangerousItem) Weight() float64 { return 100 }

func (d *DangerousItem) Use() (string, error) {
	panic("Взорвалось при использовании!")
}

func UseDangerousItem(item Item) (message string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			switch value := recovered.(type) {
			case error:
				err = fmt.Errorf("recovered panic: %w", value)
			default:
				err = fmt.Errorf("recovered panic: %v", value)
			}
		}
	}()
	return item.Use()
}

func (w *Weapon) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "Weapon|%d|%s|%d|%d|%.2f\n", w.ID, w.name, w.Damage, w.Durability, w.weight)
	return err
}

func (w *Weapon) Deserialize(reader io.Reader, lineNum int) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: "", Err: err}
	}

	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 6 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid weapon format")}
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	w.ID = id
	w.name = parts[2]

	damage, err := strconv.Atoi(parts[3])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	w.Damage = damage

	durability, err := strconv.Atoi(parts[4])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	w.Durability = durability

	weight, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	w.weight = weight

	return nil
}

func (a *Armor) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "Armor|%d|%s|%d|%d|%.2f\n", a.ID, a.name, a.Defense, a.Durability, a.weight)
	return err
}

func (a *Armor) Deserialize(reader io.Reader, lineNum int) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: "", Err: err}
	}

	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 6 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid armor format")}
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	a.ID = id
	a.name = parts[2]

	defense, err := strconv.Atoi(parts[3])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	a.Defense = defense

	durability, err := strconv.Atoi(parts[4])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}
	a.Durability = durability

	weight, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}
	a.weight = weight

	return nil
}

func (p *Potion) Serialize(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "Potion|%d|%s|%s|%d|%.2f\n", p.ID, p.name, p.Effect, p.Charges, p.weight)
	return err
}

func (p *Potion) Deserialize(reader io.Reader, lineNum int) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: "", Err: err}
	}

	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 6 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid potion format")}
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	p.ID = id
	p.name = parts[2]
	p.Effect = parts[3]

	charges, err := strconv.Atoi(parts[4])
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	p.Charges = charges

	weight, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: err}
	}

	p.weight = weight

	return nil
}

type Storable interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader, lineNum int) error
}

func (inv *Inventory) Save(w io.Writer) error {
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			if err := storable.Serialize(w); err != nil {
				return fmt.Errorf("ошибка сохранения предмета: %w", err)
			}
		}
	}
	return nil
}

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var errs []error
	lineNum := 0
	maxID := 0

	seenIDs := make(map[int]bool)
	for _, item := range inv.Items {
		var existingID int
		switch v := item.(type) {
		case *Weapon:
			existingID = v.ID
		case *Armor:
			existingID = v.ID
		case *Potion:
			existingID = v.ID
		}
		seenIDs[existingID] = true
		if existingID > maxID {
			maxID = existingID
		}
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var item Item
		if strings.HasPrefix(line, "Weapon|") {
			item = &Weapon{}
		} else if strings.HasPrefix(line, "Armor|") {
			item = &Armor{}
		} else if strings.HasPrefix(line, "Potion|") {
			item = &Potion{}
		} else {
			errs = append(errs, &ParseItemError{Line: lineNum, Value: line, Err: errors.New("unknown type")})
			continue
		}

		if storable, ok := item.(Storable); ok {
			if err := storable.Deserialize(strings.NewReader(line), lineNum); err != nil {
				errs = append(errs, err)
				continue
			}
		}

		var currentID int
		switch v := item.(type) {
		case *Weapon:
			currentID = v.ID
		case *Armor:
			currentID = v.ID
		case *Potion:
			currentID = v.ID
		}

		if seenIDs[currentID] {
			errs = append(errs, &ParseItemError{
				Line:  lineNum,
				Value: line,
				Err:   fmt.Errorf("duplicate ID %d", currentID),
			})
			continue
		}

		seenIDs[currentID] = true
		if currentID > maxID {
			maxID = currentID
		}

		inv.Items = append(inv.Items, item)
	}

	inv.NextID = maxID
	if err := scanner.Err(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: command is required: list, add, use, delete", ErrUsage)
	}

	cmd := args[0]
	inv := &Inventory{}

	file, err := os.Open("inventory.txt")
	if err == nil {
		defer file.Close()
		if loadErr := inv.Load(file); loadErr != nil {
			var parseErr *ParseItemError
			if errors.As(loadErr, &parseErr) {
				fmt.Fprintf(stderr, "Внимание: обнаружена поврежденная строка %d: %q\n", parseErr.Line, parseErr.Value)
			}
			return fmt.Errorf("ошибка загрузки инвентаря: %w", loadErr)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}

	var successMsg string

	switch cmd {
	case "list":
		fs := flag.NewFlagSet("list", flag.ContinueOnError)
		fs.SetOutput(stderr)

		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("%w: unexpected arguments: %v", ErrUsage, fs.Args())
		}

		for _, item := range inv.Items {
			var id int
			switch v := item.(type) {
			case *Weapon:
				id = v.ID
			case *Armor:
				id = v.ID
			case *Potion:
				id = v.ID
			}
			fmt.Fprintf(stdout, "[%d] %s (вес: %.2f)\n", id, item.Name(), item.Weight())
		}
		return nil

	case "add":
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		fs.SetOutput(stderr)
		name := fs.String("name", "", "Имя предмета")
		weight := fs.Float64("weight", 1.0, "Вес предмета")

		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("%w: unexpected arguments: %v", ErrUsage, fs.Args())
		}
		if *name == "" {
			return fmt.Errorf("%w: -name is required", ErrUsage)
		}

		if strings.ContainsAny(*name, "|\n\r") {
			return fmt.Errorf("%w: name cannot contain '|' or newline characters", ErrUsage)
		}

		item := &Weapon{name: *name, weight: *weight, Damage: 10, Durability: 5}
		id := inv.AddItem(item)
		successMsg = fmt.Sprintf("Добавлен предмет с ID: %d\n", id)

	case "use":
		fs := flag.NewFlagSet("use", flag.ContinueOnError)
		fs.SetOutput(stderr)
		id := fs.Int("id", 0, "ID предмета")

		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("%w: unexpected arguments: %v", ErrUsage, fs.Args())
		}
		if *id == 0 {
			return fmt.Errorf("%w: -id is required", ErrUsage)
		}

		var foundItem Item
		for _, item := range inv.Items {
			var itemID int
			switch v := item.(type) {
			case *Weapon:
				itemID = v.ID
			case *Armor:
				itemID = v.ID
			case *Potion:
				itemID = v.ID
			}
			if itemID == *id {
				foundItem = item
				break
			}
		}

		if foundItem == nil {
			return fmt.Errorf("предмет с ID %d не найден", *id)
		}

		msg, err := foundItem.Use()
		if errors.Is(err, ErrBroken) {
			return fmt.Errorf("предмет %q сломан и не может быть использован", foundItem.Name())
		} else if errors.Is(err, ErrExhausted) {
			return fmt.Errorf("у предмета %q закончились заряды", foundItem.Name())
		} else if err != nil {
			return err
		}
		successMsg = msg + "\n"

	case "delete":
		fs := flag.NewFlagSet("delete", flag.ContinueOnError)
		fs.SetOutput(stderr)
		id := fs.Int("id", 0, "ID предмета")

		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("%w: unexpected arguments: %v", ErrUsage, fs.Args())
		}

		if *id == 0 {
			return fmt.Errorf("%w: -id is required", ErrUsage)
		}

		if inv.DeleteItem(*id) {
			successMsg = "Предмет удален\n"
		} else {
			return fmt.Errorf("предмет с ID %d не найден", *id)
		}

	default:
		return fmt.Errorf("%w: unknown command %q", ErrUsage, cmd)
	}

	tmpFile, err := os.CreateTemp(".", "inv_*.tmp")
	if err != nil {
		return fmt.Errorf("создание временного файла: %w", err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if err := inv.Save(tmpFile); err != nil {
		tmpFile.Close()
		return fmt.Errorf("сохранение: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("закрытие: %w", err)
	}

	if err := os.Rename(tmpName, "inventory.txt"); err != nil {
		return fmt.Errorf("замена файла: %w", err)
	}

	if successMsg != "" {
		fmt.Fprint(stdout, successMsg)
	}

	return nil
}

func exitCode(err error) int {
	switch {
	case err == nil, errors.Is(err, flag.ErrHelp):
		return 0
	case errors.Is(err, ErrUsage):
		return 2
	default:
		return 1
	}
}

func report(err error, stderr io.Writer) int {
	if err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(stderr, "error:", err)
	}
	return exitCode(err)
}

func main() {
	os.Exit(report(run(os.Args[1:], os.Stdout, os.Stderr), os.Stderr))
}

// 1. Почему ошибка сломанного предмета не должна быть паникой?
//    Паника предназначена для критических ошибок, после которых программа не может продолжать работу.
//    Сломанный меч — это штатная ситуация бизнес-логики. Если паниковать, программа упадет, и игрок потеряет прогресс.
// 2. Какая строка после panic не будет выполнена?
//    Ни одна. Как только срабатывает panic, нормальное выполнение функции мгновенно прекращается.
//    Начинают выполняться только отложенные функции (defer), после чего программа падает (или паника перехватывается).
// 3. Почему recover должен находиться в отложенной функции?
//    Потому что при панике нормальный поток управления останавливается. Единственный код,
//    который гарантированно выполнится после паники — это функции, зарегистрированные через defer.
