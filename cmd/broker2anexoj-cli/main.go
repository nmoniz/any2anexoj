package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
	"git.naterciomoniz.net/applications/broker2anexoj/internal/trading212"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run(context.Background())
	if err != nil {
		slog.Error("found a fatal issue", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	f, err := os.Open("test.csv")
	if err != nil {
		return fmt.Errorf("open statement: %w", err)
	}

	reader := trading212.NewRecordReader(f)

	writer := internal.NewStdOutLogger()

	eg.Go(func() error {
		return internal.BuildReport(ctx, reader, writer)
	})

	err = eg.Wait()
	if err != nil {
		return err
	}

	slog.Info("Finish processing statement")

	return nil
}
