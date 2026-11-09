package main

import (
	"errors"
	"flag"
	"fmt"
	"homework-go/homework-7/internal/cli"
	"homework-go/homework-7/internal/inventory"
	"homework-go/homework-7/internal/storage"
	"os"
)

func main() {
	useMemory := flag.Bool("memory", false, "Использовать хранилище в памяти")
	flag.Parse()

	var store inventory.Storage
	if *useMemory {
		store = storage.NewMemoryStorage()
	} else {
		store = storage.NewFileStorage("inventory.txt")
	}

	invService := inventory.NewService(store)

	app := cli.NewApp(invService, os.Stdout, os.Stderr)

	err := app.Run(flag.Args())
	if err != nil && !errors.Is(err, flag.ErrHelp) {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	os.Exit(cli.ExitCode(err))
}
