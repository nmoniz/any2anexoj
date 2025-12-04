package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/nmoniz/any2anexoj/internal/trading212"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
)

// TODO: once we support more brokers or exchanges we should make this parameter required and
// remove/change default
var platform = pflag.StringP("platform", "p", "trading212", "one of the supported platforms")

var lang = pflag.StringP("language", "l", language.Portuguese.String(), "2 letter language code")

var readerFactories = map[string]func() internal.RecordReader{
	"trading212": func() internal.RecordReader {
		return trading212.NewRecordReader(os.Stdin, internal.NewOpenFIGI(&http.Client{Timeout: 5 * time.Second}))
	},
}

func main() {
	pflag.Parse()

	if platform == nil || len(*platform) == 0 {
		slog.Error("--platform flag is required")
		os.Exit(1)
	}

	if lang == nil || len(*lang) == 0 {
		slog.Error("--language flag is required")
		os.Exit(1)
	}

	err := run(context.Background(), *platform, *lang)
	if err != nil {
		slog.Error("found a fatal issue", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, platform, lang string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	factory, ok := readerFactories[platform]
	if !ok {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	reader := factory()

	writer := internal.NewAggregatorWriter()

	eg.Go(func() error {
		return internal.BuildReport(ctx, reader, writer)
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	loc, err := NewLocalizer(lang)
	if err != nil {
		return fmt.Errorf("create localizer: %w", err)
	}

	printer := NewPrettyPrinter(os.Stdout, loc)

	printer.Render(writer)

	return nil
}
