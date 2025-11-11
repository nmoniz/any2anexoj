package internal

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"sync"
)

type RecordReader interface {
	// ReadRecord should return Records until an error is found.
	ReadRecord() (Record, error)
}

// Reporter consumes each record to produce ReportItem.
type Reporter struct {
	reader RecordReader
}

func NewReporter(rr RecordReader) *Reporter {
	return &Reporter{
		reader: rr,
	}
}

func (r *Reporter) Run() error {
	forewarders := make(map[string]chan Record)

	aggregator := make(chan processResult)
	defer close(aggregator)

	go func() {
		for result := range aggregator {
			fmt.Printf("%v\n", result)
		}
	}()

	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	for {
		rec, err := r.reader.ReadRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		router, ok := forewarders[rec.Symbol()]
		if !ok {
			router = make(chan Record, 1)
			defer close(router)

			wg.Go(func() {
				processRecords(router, aggregator)
			})

			forewarders[rec.Symbol()] = router
		}

		router <- rec
	}
}

func processRecords(records <-chan Record, results chan<- processResult) {
	var q RecordQueue

	for rec := range records {
		switch rec.Side() {
		case SideBuy:
			q.Push(rec)

		case SideSell:
			unmatchedQty := new(big.Float).Copy(rec.Quantity())
			zero := new(big.Float)

			for unmatchedQty.Cmp(zero) > 0 {
				buy, ok := q.Pop()
				if !ok {
					results <- processResult{
						err: ErrSellWithoutBuy,
					}
					return
				}

				var matchedQty *big.Float
				if buy.Quantity().Cmp(unmatchedQty) > 0 {
					matchedQty = unmatchedQty
					buy.Quantity().Sub(buy.Quantity(), unmatchedQty)
				} else {
					matchedQty = buy.Quantity()
				}

				unmatchedQty.Sub(unmatchedQty, matchedQty)

				sellValue := new(big.Float).Mul(matchedQty, rec.Price())
				buyValue := new(big.Float).Mul(matchedQty, buy.Price())
				realisedPnL := new(big.Float).Sub(sellValue, buyValue)
				slog.Info("Realised PnL",
					slog.Any("Symbol", rec.Symbol()),
					slog.Any("PnL", realisedPnL),
					slog.Any("Timestamp", rec.Timestamp()))

				results <- processResult{
					item: ReportItem{},
				}
			}

		default:
			results <- processResult{
				err: fmt.Errorf("unknown side: %v", rec.Side()),
			}
			return
		}
	}
}

type processResult struct {
	item ReportItem
	err  error
}

type ReportItem struct{}

var ErrSellWithoutBuy = fmt.Errorf("found sell without bought volume")
