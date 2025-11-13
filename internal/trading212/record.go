package trading212

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
)

type Record struct {
	symbol    string
	side      internal.Side
	quantity  *big.Float
	price     *big.Float
	timestamp time.Time
	fees      *big.Float
	taxes     *big.Float
}

func (r Record) Symbol() string {
	return r.symbol
}

func (r Record) Side() internal.Side {
	return r.side
}

func (r Record) Quantity() *big.Float {
	return r.quantity
}

func (r Record) Price() *big.Float {
	return r.price
}

func (r Record) Timestamp() time.Time {
	return r.timestamp
}

func (r Record) Fees() *big.Float {
	return r.fees
}

func (r Record) Taxes() *big.Float {
	return r.taxes
}

type RecordReader struct {
	reader *csv.Reader
}

func NewRecordReader(r io.Reader) *RecordReader {
	return &RecordReader{
		reader: csv.NewReader(r),
	}
}

const (
	MarketBuy  = "market buy"
	MarketSell = "market sell"
	LimitBuy   = "limit buy"
	LimitSell  = "limit sell"
)

func (rr RecordReader) ReadRecord() (internal.Record, error) {
	for {
		raw, err := rr.reader.Read()
		if err != nil {
			return Record{}, fmt.Errorf("read record: %w", err)
		}

		var side internal.Side
		switch strings.ToLower(raw[0]) {
		case MarketBuy, LimitBuy:
			side = internal.SideBuy
		case MarketSell, LimitSell:
			side = internal.SideSell
		case "action", "stock split open", "stock split close":
			continue
		default:
			return Record{}, fmt.Errorf("parse record type: %s", raw[0])
		}

		qant, err := parseDecimal(raw[6])
		if err != nil {
			return Record{}, fmt.Errorf("parse record quantity: %w", err)
		}

		price, err := parseDecimal(raw[7])
		if err != nil {
			return Record{}, fmt.Errorf("parse record price: %w", err)
		}

		ts, err := time.Parse(time.DateTime, raw[1])
		if err != nil {
			return Record{}, fmt.Errorf("parse record timestamp: %w", err)
		}

		convertionFee, err := parseOptinalDecimal(raw[16])
		if err != nil {
			return Record{}, fmt.Errorf("parse record convertion fee: %w", err)
		}

		stampDutyTax, err := parseOptinalDecimal(raw[14])
		if err != nil {
			return Record{}, fmt.Errorf("parse record stamp duty tax: %w", err)
		}

		frenchTxTax, err := parseOptinalDecimal(raw[18])
		if err != nil {
			return Record{}, fmt.Errorf("parse record french transaction tax: %w", err)
		}

		return Record{
			symbol:    raw[2],
			side:      side,
			quantity:  qant,
			price:     price,
			timestamp: ts,
			fees:      convertionFee,
			taxes:     new(big.Float).Add(stampDutyTax, frenchTxTax),
		}, nil
	}
}

// parseFloat attempts to parse a string using a standard precision and rounding mode.
// Using this function helps avoid issues around converting values due to sligh parameter changes.
func parseDecimal(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 128, big.ToZero)
	return f, err
}

// parseOptinalDecimal behaves the same as parseDecimal but returns 0 when len(s) is 0 instead of
// error.
// Using this function helps avoid issues around converting values due to sligh parameter changes.
func parseOptinalDecimal(s string) (*big.Float, error) {
	if len(s) == 0 {
		return new(big.Float), nil
	}

	return parseDecimal(s)
}
