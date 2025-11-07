package main

import (
	"container/list"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
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

	r := NewRecordReader(f)

	assets := make(map[string]*list.List)
	for {
		record, err := r.ReadRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read statement record: %w", err)
		}

		switch record.Direction() {
		case DirectionBuy:
			lst, ok := assets[record.Symbol()]
			if !ok {
				lst = list.New()
				assets[record.Symbol()] = lst
			}
			lst.PushBack(record)

		case DirectionSell:
			lst, ok := assets[record.Symbol()]
			if !ok {
				return ErrSellWithoutBuy
			}

			first := lst.Front()
			if first == nil {
				return ErrSellWithoutBuy
			}

		default:
			return fmt.Errorf("unknown direction: %s", record.Direction())
		}

	}

	slog.Info("Finish processing statement", slog.Any("assets_count", len(assets)))

	return nil
}

var ErrSellWithoutBuy = fmt.Errorf("found sell without bought volume")

type Record struct {
	symbol    string
	direction Direction
	quantity  *big.Float
	price     *big.Float
}

func (r Record) Symbol() string {
	return r.symbol
}

func (r Record) Direction() Direction {
	return r.direction
}

func (r Record) Quantity() *big.Float {
	return r.quantity
}

func (r Record) Price() *big.Float {
	return r.price
}

type RecordReader struct {
	reader *csv.Reader
}

func NewRecordReader(r io.Reader) *RecordReader {
	return &RecordReader{
		reader: csv.NewReader(r),
	}
}

func (rr RecordReader) ReadRecord() (Record, error) {
	for {
		raw, err := rr.reader.Read()
		if err != nil {
			return Record{}, fmt.Errorf("read record: %w", err)
		}

		var dir Direction
		switch strings.ToLower(raw[0]) {
		case "market buy":
			dir = DirectionBuy
		case "market sell":
			dir = DirectionSell
		case "action", "stock split open", "stock split close":
			continue
		default:
			return Record{}, fmt.Errorf("unhandled record: %s", raw[0])
		}
		qant, _, err := big.ParseFloat(raw[6], 10, 20, big.ToZero)
		if err != nil {
			return Record{}, fmt.Errorf("parse quantity: %w", err)
		}

		price, _, err := big.ParseFloat(raw[7], 10, 20, big.ToZero)
		if err != nil {
			return Record{}, fmt.Errorf("parse price: %w", err)
		}

		return Record{
			symbol:    raw[2],
			direction: dir,
			quantity:  qant,
			price:     price,
		}, nil
	}
}

type Direction uint

const (
	DirectionUnknown Direction = 0
	DirectionBuy               = 1
	DirectionSell              = 2
)

func (d Direction) String() string {
	switch d {
	case 1:
		return "buy"
	case 2:
		return "sell"
	default:
		return "unknown"
	}
}

func (d Direction) IsBuy() bool {
	return d == DirectionBuy
}

func (d Direction) IsSell() bool {
	return d == DirectionSell
}
