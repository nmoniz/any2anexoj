package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/nmoniz/any2anexoj/internal/trading212"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

// TODO: once we support more brokers or exchanges we should make this parameter required and
// remove/change default
var platform = pflag.StringP("platform", "p", "trading212", "one of the supported platforms")

var supportedPlatforms = map[string]func() internal.RecordReader{
	"trading212": func() internal.RecordReader { return trading212.NewRecordReader(os.Stdin) },
}

func main() {
	pflag.Parse()

	if platform == nil || len(*platform) == 0 {
		slog.Error("--platform flag is required")
		os.Exit(1)
	}

	err := run(context.Background(), *platform)
	if err != nil {
		slog.Error("found a fatal issue", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, platform string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	factory, ok := supportedPlatforms[platform]
	if !ok {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	reader := factory()

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
