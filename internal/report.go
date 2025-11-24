package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/shopspring/decimal"
)

type Record interface {
	Symbol() string
	Nature() Nature
	BrokerCountry() int64
	AssetCountry() int64
	Side() Side
	Price() decimal.Decimal
	Quantity() decimal.Decimal
	Timestamp() time.Time
	Fees() decimal.Decimal
	Taxes() decimal.Decimal
}

type RecordReader interface {
	// ReadRecord should return Records until an error is found.
	ReadRecord(context.Context) (Record, error)
}

type ReportItem struct {
	Symbol        string
	Nature        Nature
	BrokerCountry int64
	AssetCountry  int64
	BuyValue      decimal.Decimal
	BuyTimestamp  time.Time
	SellValue     decimal.Decimal
	SellTimestamp time.Time
	Fees          decimal.Decimal
	Taxes         decimal.Decimal
}

func (ri ReportItem) RealisedPnL() decimal.Decimal {
	return ri.SellValue.Sub(ri.BuyValue)
}

type ReportWriter interface {
	// ReportWriter writes report items
	Write(context.Context, ReportItem) error
}

func BuildReport(ctx context.Context, reader RecordReader, writer ReportWriter) error {
	buys := make(map[string]*FillerQueue)

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
				buyQueue = new(FillerQueue)
				buys[rec.Symbol()] = buyQueue
			}

			err = processRecord(ctx, buyQueue, rec, writer)
			if err != nil {
				return fmt.Errorf("processing record: %w", err)
			}
		}
	}
}

func processRecord(ctx context.Context, q *FillerQueue, rec Record, writer ReportWriter) error {
	switch rec.Side() {
	case SideBuy:
		q.Push(NewFiller(rec))

	case SideSell:
		unmatchedQty := rec.Quantity()

		for unmatchedQty.IsPositive() {
			buy, ok := q.Peek()
			if !ok {
				return ErrInsufficientBoughtVolume
			}

			matchedQty, filled := buy.Fill(unmatchedQty)

			if filled {
				_, ok := q.Pop()
				if !ok {
					return fmt.Errorf("pop empty filler queue")
				}
			}

			unmatchedQty = unmatchedQty.Sub(matchedQty)

			buyValue := matchedQty.Mul(buy.Price())
			sellValue := matchedQty.Mul(rec.Price())

			err := writer.Write(ctx, ReportItem{
				Symbol:        rec.Symbol(),
				BrokerCountry: rec.BrokerCountry(),
				AssetCountry:  rec.AssetCountry(),
				BuyValue:      buyValue,
				BuyTimestamp:  buy.Timestamp(),
				SellValue:     sellValue,
				SellTimestamp: rec.Timestamp(),
				Fees:          buy.Fees().Add(rec.Fees()),
				Taxes:         buy.Taxes().Add(rec.Taxes()),
				Nature:        buy.Nature(),
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
