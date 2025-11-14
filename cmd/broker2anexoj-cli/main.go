package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/nmoniz/any2anexoj/internal/trading212"
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

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	reader := trading212.NewRecordReader(os.Stdin)

	writer := internal.NewStdOutLogger()

	eg.Go(func() error {
		return internal.BuildReport(ctx, reader, writer)
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	slog.Info("Finish processing statement")

	return nil
}
