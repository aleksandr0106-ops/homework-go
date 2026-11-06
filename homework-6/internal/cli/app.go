package cli

import (
	"errors"
	"flag"
	"fmt"
	"homework-go/homework-6/internal/inventory"
	"homework-go/homework-6/internal/item"
	"homework-go/homework-6/internal/storage"
	"io"
)

var ErrUsage = errors.New("usage error")

type App struct {
	inv    *inventory.Service
	stdout io.Writer
	stderr io.Writer
}

func NewApp(inv *inventory.Service, stdout, stderr io.Writer) *App {
	return &App{
		inv:    inv,
		stdout: stdout,
		stderr: stderr,
	}
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: command is required: list, add, use, delete", ErrUsage)
	}
	if err := a.inv.Load(); err != nil {
		var parseErr *storage.ParseItemError
		if errors.As(err, &parseErr) {
			fmt.Fprintf(a.stderr, "Внимание: обнаружена поврежденная строка %d: %q\n", parseErr.Line, parseErr.Value)
		}
		return fmt.Errorf("ошибка загрузки инвентаря: %w", err)
	}

	cmd := args[0]
	var successMsg string

	switch cmd {
	case "list":
		fs := flag.NewFlagSet("list", flag.ContinueOnError)
		fs.SetOutput(a.stderr)
		if err := fs.Parse(args[1:]); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("%w: unexpected arguments: %v", ErrUsage, fs.Args())
		}

		for _, it := range a.inv.Items() {
			fmt.Fprintf(a.stdout, "[%d] %s (вес: %.2f)\n", it.ID(), it.Name(), it.Weight())
		}
		return nil

	case "add":
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		fs.SetOutput(a.stderr)
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

		id, err := a.inv.Add(*name, *weight)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
		successMsg = fmt.Sprintf("Добавлен предмет с ID: %d\n", id)

	case "use":
		fs := flag.NewFlagSet("use", flag.ContinueOnError)
		fs.SetOutput(a.stderr)
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

		foundItem, err := a.inv.FindByID(*id)
		if err != nil {
			return err
		}

		msg, err := foundItem.Use()
		if errors.Is(err, item.ErrBroken) {
			return fmt.Errorf("предмет %q сломан и не может быть использован", foundItem.Name())
		} else if errors.Is(err, item.ErrExhausted) {
			return fmt.Errorf("у предмета %q закончились заряды", foundItem.Name())
		} else if err != nil {
			return err
		}
		successMsg = msg + "\n"

	case "delete":
		fs := flag.NewFlagSet("delete", flag.ContinueOnError)
		fs.SetOutput(a.stderr)
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

		if err := a.inv.Delete(*id); err != nil {
			return err
		}
		successMsg = "Предмет удален\n"

	default:
		return fmt.Errorf("%w: unknown command %q", ErrUsage, cmd)
	}

	if err := a.inv.Save(); err != nil {
		return fmt.Errorf("сохранение: %w", err)
	}

	if successMsg != "" {
		fmt.Fprint(a.stdout, successMsg)
	}

	return nil
}
