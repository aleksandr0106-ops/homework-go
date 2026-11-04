package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
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

func (w *Weapon) Name() string    { return w.name }
func (w *Weapon) Weight() float64 { return w.weight }
func (w *Weapon) Use() (string, error) {
	if w.Durability <= 0 {
		return "", ErrBroken
	}
	w.Durability--
	return fmt.Sprintf("Взмах оружием: %s", w.name), nil
}

type Armor struct {
	ID      int
	name    string
	weight  float64
	Defense int
}

func (a *Armor) Name() string    { return a.name }
func (a *Armor) Weight() float64 { return a.weight }
func (a *Armor) Use() (string, error) {
	return fmt.Sprintf("Надеваем броню: %s", a.name), nil
}

type Potion struct {
	ID      int
	name    string
	weight  float64
	Effect  string
	Charges int
}

func (p *Potion) Name() string    { return p.name }
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

type DangerousItem struct{ name string }

func (d *DangerousItem) Name() string    { return d.name }
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

func (w *Weapon) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Weapon|%d|%s|%d|%d|%.2f\n", w.ID, w.name, w.Damage, w.Durability, w.weight)
}
func (w *Weapon) Deserialize(reader io.Reader, lineNum int) error {
	data, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 6 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid weapon format")}
	}
	w.ID, _ = strconv.Atoi(parts[1])
	w.name = parts[2]
	w.Damage, _ = strconv.Atoi(parts[3])
	w.Durability, _ = strconv.Atoi(parts[4])
	w.weight, _ = strconv.ParseFloat(parts[5], 64)
	return nil
}

func (a *Armor) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Armor|%d|%s|%d|%.2f\n", a.ID, a.name, a.Defense, a.weight)
}
func (a *Armor) Deserialize(reader io.Reader, lineNum int) error {
	data, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 5 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid armor format")}
	}
	a.ID, _ = strconv.Atoi(parts[1])
	a.name = parts[2]
	a.Defense, _ = strconv.Atoi(parts[3])
	a.weight, _ = strconv.ParseFloat(parts[4], 64)
	return nil
}

func (p *Potion) Serialize(writer io.Writer) {
	fmt.Fprintf(writer, "Potion|%d|%s|%s|%d|%.2f\n", p.ID, p.name, p.Effect, p.Charges, p.weight)
}
func (p *Potion) Deserialize(reader io.Reader, lineNum int) error {
	data, _ := io.ReadAll(reader)
	parts := strings.Split(strings.TrimSpace(string(data)), "|")
	if len(parts) != 6 {
		return &ParseItemError{Line: lineNum, Value: string(data), Err: errors.New("invalid potion format")}
	}
	p.ID, _ = strconv.Atoi(parts[1])
	p.name = parts[2]
	p.Effect = parts[3]
	p.Charges, _ = strconv.Atoi(parts[4])
	p.weight, _ = strconv.ParseFloat(parts[5], 64)
	return nil
}

type Storable interface {
	Serialize(w io.Writer)
	Deserialize(r io.Reader, lineNum int) error
}

func (inv *Inventory) Save(w io.Writer) error {
	for _, item := range inv.Items {
		if storable, ok := item.(Storable); ok {
			storable.Serialize(w)
		}
	}
	return nil
}

func (inv *Inventory) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var errs []error
	lineNum := 0
	maxID := 0

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

		inv.Items = append(inv.Items, item)
		switch v := item.(type) {
		case *Weapon:
			if v.ID > maxID {
				maxID = v.ID
			}
		case *Armor:
			if v.ID > maxID {
				maxID = v.ID
			}
		case *Potion:
			if v.ID > maxID {
				maxID = v.ID
			}
		}
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
		loadErr := inv.Load(file)
		file.Close()
		if loadErr != nil {
			fmt.Fprintf(stderr, "Предупреждение при загрузке: %v\n", loadErr)
		}
	}

	switch cmd {
	case "list":
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

	case "add":
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		fs.SetOutput(stderr)
		name := fs.String("name", "", "Имя предмета")
		weight := fs.Float64("weight", 1.0, "Вес предмета")
		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if *name == "" {
			return fmt.Errorf("%w: -name is required", ErrUsage)
		}

		item := &Weapon{name: *name, weight: *weight, Damage: 10, Durability: 5}
		id := inv.AddItem(item)
		fmt.Fprintf(stdout, "Добавлен предмет с ID: %d\n", id)

	case "use":
		fs := flag.NewFlagSet("use", flag.ContinueOnError)
		fs.SetOutput(stderr)
		id := fs.Int("id", 0, "ID предмета")
		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
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
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, msg)

	case "delete":
		fs := flag.NewFlagSet("delete", flag.ContinueOnError)
		fs.SetOutput(stderr)
		id := fs.Int("id", 0, "ID предмета")
		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}

		if inv.DeleteItem(*id) {
			fmt.Fprintln(stdout, "Предмет удален")
		} else {
			return fmt.Errorf("предмет с ID %d не найден", *id)
		}

	default:
		return fmt.Errorf("%w: unknown command %q", ErrUsage, cmd)
	}

	if cmd != "list" {
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
