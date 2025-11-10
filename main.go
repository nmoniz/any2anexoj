package main

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
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

	r := trading212.NewRecordReader(f)

	assets := make(map[string]*list.List)
	for {
		record, err := r.ReadRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read statement record: %w", err)
		}

		switch record.Side() {
		case internal.SideBuy:
			lst, ok := assets[record.Symbol()]
			if !ok {
				lst = list.New()
				assets[record.Symbol()] = lst
			}
			lst.PushBack(record)

		case internal.SideSell:
			lst, ok := assets[record.Symbol()]
			if !ok {
				return ErrSellWithoutBuy
			}

			unmatchedQty := new(big.Float).Copy(record.Quantity())
			zero := new(big.Float)

			for unmatchedQty.Cmp(zero) > 0 {
				front := lst.Front()
				if front == nil {
					return ErrSellWithoutBuy
				}

				next, ok := front.Value.(internal.Record)
				if !ok {
					return fmt.Errorf("unexpected record type: %T", front)
				}

				var matchedQty *big.Float
				if next.Quantity().Cmp(unmatchedQty) > 0 {
					matchedQty = unmatchedQty
					next.Quantity().Sub(next.Quantity(), unmatchedQty)
				} else {
					matchedQty = next.Quantity()
					lst.Remove(front)
				}

				unmatchedQty.Sub(unmatchedQty, matchedQty)

				sellValue := new(big.Float).Mul(matchedQty, record.Price())
				buyValue := new(big.Float).Mul(matchedQty, next.Price())
				realisedPnL := new(big.Float).Sub(sellValue, buyValue)
				slog.Info("Realised PnL",
					slog.Any("Symbol", record.Symbol()),
					slog.Any("PnL", realisedPnL))
			}

		default:
			return fmt.Errorf("unknown side: %s", record.Side())
		}
	}

	slog.Info("Finish processing statement", slog.Any("assets_count", len(assets)))

	return nil
}

var ErrSellWithoutBuy = fmt.Errorf("found sell without bought volume")
