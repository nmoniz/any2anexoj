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
		slog.Error("fatal error", slog.Any("err", err))
	}
}

func run() error {
	f, err := os.Open("test.csv")
	if err != nil {
		return fmt.Errorf("open statement: %w", err)
	}

	reader := trading212.NewRecordReader(f)

	reporter := internal.NewReporter(reader)

	err = reporter.Run()
	if err != nil {
		return err
	}

	slog.Info("Finish processing statement")

	return nil
}

var ErrSellWithoutBuy = fmt.Errorf("found sell without bought volume")
