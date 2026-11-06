package storage

import (
	"bufio"
	"errors"
	"fmt"
	"homework-go/homework-6/internal/item"
	"os"
	"strconv"
	"strings"
)

type ParseItemError struct {
	Line  int
	Value string
	Err   error
}

func (e *ParseItemError) Error() string {
	return fmt.Sprintf("line %d: parse error %q: %v", e.Line, e.Value, e.Err)
}

func (e *ParseItemError) Unwrap() error { return e.Err }

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func (s *FileStorage) Save(items []item.Item) error {
	tmpFile, err := os.CreateTemp(".", "inv_*.tmp")
	if err != nil {
		return fmt.Errorf("создание временного файла: %w", err)
	}

	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	writer := bufio.NewWriter(tmpFile)
	for _, i := range items {
		var line string
		switch v := i.(type) {
		case *item.Weapon:
			line = fmt.Sprintf("Weapon|%d|%s|%d|%d|%.2f\n", v.ID(), v.Name(), v.Damage, v.Durability, v.Weight())
		case *item.Armor:
			line = fmt.Sprintf("Armor|%d|%s|%d|%d|%.2f\n", v.ID(), v.Name(), v.Defense, v.Durability, v.Weight())
		case *item.Potion:
			line = fmt.Sprintf("Potion|%d|%s|%s|%d|%.2f\n", v.ID(), v.Name(), v.Effect, v.Charges, v.Weight())
		}

		if _, err := writer.WriteString(line); err != nil {
			tmpFile.Close()
			return fmt.Errorf("запись в буфер: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("сброс буфера: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("закрытие временного файла: %w", err)
	}

	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("замена файла: %w", err)
	}
	return nil
}

func (s *FileStorage) Load() ([]item.Item, error) {
	file, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []item.Item{}, nil
		}
		return nil, fmt.Errorf("открытие файла: %w", err)
	}
	defer file.Close()

	var items []item.Item
	var errs []error
	lineNum := 0
	seenIDs := make(map[int]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		it, err := parseItem(line, lineNum)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if seenIDs[it.ID()] {
			errs = append(errs, &ParseItemError{Line: lineNum, Value: line, Err: fmt.Errorf("duplicate ID %d", it.ID())})
			continue
		}

		seenIDs[it.ID()] = true
		items = append(items, it)
	}

	if err := scanner.Err(); err != nil {
		errs = append(errs, err)
	}
	return items, errors.Join(errs...)
}

func parseItem(line string, lineNum int) (item.Item, error) {
	parts := strings.Split(strings.TrimSpace(line), "|")
	if len(parts) < 6 {
		return nil, &ParseItemError{Line: lineNum, Value: line, Err: errors.New("invalid format")}
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
	}

	switch parts[0] {
	case "Weapon":
		damage, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}

		dur, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}

		w, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}
		return item.NewWeapon(id, parts[2], w, damage, dur), nil

	case "Armor":
		def, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}

		dur, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}

		w, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}
		return item.NewArmor(id, parts[2], w, def, dur), nil

	case "Potion":
		ch, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}

		w, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, &ParseItemError{Line: lineNum, Value: line, Err: err}
		}
		return item.NewPotion(id, parts[2], w, parts[3], ch), nil
	}

	return nil, &ParseItemError{Line: lineNum, Value: line, Err: errors.New("unknown type")}
}
