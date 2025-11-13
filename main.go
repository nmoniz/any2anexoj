package main

import (
	"fmt"
	"log/slog"
	"os"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
	"git.naterciomoniz.net/applications/broker2anexoj/internal/trading212"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("found a fatal issue", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	f, err := os.Open("test.csv")
	if err != nil {
		return fmt.Errorf("open statement: %w", err)
	}

	reader := trading212.NewRecordReader(f)

	writer := internal.NewStdOutLogger()

	err = internal.BuildReport(reader, writer)
	if err != nil {
		return err
	}

	slog.Info("Finish processing statement")

	return nil
}
