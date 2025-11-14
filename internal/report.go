package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"
)

type RecordReader interface {
	// ReadRecord should return Records until an error is found.
	ReadRecord(context.Context) (Record, error)
}

type ReportWriter interface {
	// ReportWriter writes report items
	Write(context.Context, ReportItem) error
}

func BuildReport(ctx context.Context, reader RecordReader, writer ReportWriter) error {
	buys := make(map[string]*RecordQueue)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			rec, err := reader.ReadRecord(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}

			buyQueue, ok := buys[rec.Symbol()]
			if !ok {
				buyQueue = new(RecordQueue)
				buys[rec.Symbol()] = buyQueue
			}

			err = processRecord(ctx, buyQueue, rec, writer)
			if err != nil {
				return fmt.Errorf("processing record: %w", err)
			}
		}
	}
}

func processRecord(ctx context.Context, q *RecordQueue, rec Record, writer ReportWriter) error {
	switch rec.Side() {
	case SideBuy:
		q.Push(rec)

	case SideSell:
		unmatchedQty := new(big.Float).Copy(rec.Quantity())
		zero := new(big.Float)

		for unmatchedQty.Cmp(zero) > 0 {
			buy, ok := q.Peek()
			if !ok {
				return ErrInsufficientBoughtVolume
			}

			var matchedQty *big.Float
			if buy.Quantity().Cmp(unmatchedQty) > 0 {
				matchedQty = unmatchedQty
				buy.Quantity().Sub(buy.Quantity(), unmatchedQty)
			} else {
				matchedQty = buy.Quantity()
				q.Pop()
			}

			unmatchedQty.Sub(unmatchedQty, matchedQty)

			sellValue := new(big.Float).Mul(matchedQty, rec.Price())
			buyValue := new(big.Float).Mul(matchedQty, buy.Price())

			err := writer.Write(ctx, ReportItem{
				BuyValue:      buyValue,
				BuyTimestamp:  buy.Timestamp(),
				SellValue:     sellValue,
				SellTimestamp: rec.Timestamp(),
				Fees:          new(big.Float).Add(buy.Fees(), rec.Fees()),
				Taxes:         new(big.Float).Add(buy.Taxes(), rec.Fees()),
			})
			if err != nil {
				return fmt.Errorf("write report item: %w", err)
			}
		}

	default:
		return fmt.Errorf("unknown side: %v", rec.Side())
	}

	return nil
}

type ReportItem struct {
	BuyValue      *big.Float
	BuyTimestamp  time.Time
	SellValue     *big.Float
	SellTimestamp time.Time
	Fees          *big.Float
	Taxes         *big.Float
}

func (ri ReportItem) RealisedPnL() *big.Float {
	return new(big.Float).Sub(ri.SellValue, ri.BuyValue)
}

var ErrInsufficientBoughtVolume = fmt.Errorf("insufficient bought volume")
