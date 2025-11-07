package main

import (
	"container/list"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
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

	r := csv.NewReader(f)

	assets := make(map[string]*list.List)
	for {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read statement record: %w", err)
		}

		switch strings.ToLower(record[0]) {
		case "market buy":
			lst, ok := assets[record[2]]
			if !ok {
				lst = list.New()
				assets[record[2]] = lst
			}
			lst.PushBack(record[12])

		case "market sell":
			lst, ok := assets[record[2]]
			if !ok {
				return ErrSellWithoutBuy
			}

			first := lst.Front()
			if first == nil {
				return ErrSellWithoutBuy
			}

			slog.Info("Realised PnL", slog.Any("record", record))

		case "action", "stock split open", "stock split close":
			// ignored

		default:
			return fmt.Errorf("unhandled record: %s", record[0])
		}

	}

	slog.Info("Finish processing statement", slog.Any("assets_count", len(assets)))

	return nil
}

var ErrSellWithoutBuy = fmt.Errorf("found sell without bought volume")
