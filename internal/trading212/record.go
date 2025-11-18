package trading212

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/biter777/countries"
	"github.com/nmoniz/any2anexoj/internal"
	"github.com/shopspring/decimal"
)

type Record struct {
	symbol    string
	side      internal.Side
	quantity  decimal.Decimal
	price     decimal.Decimal
	timestamp time.Time
	fees      decimal.Decimal
	taxes     decimal.Decimal
}

func (r Record) Symbol() string {
	return r.symbol
}

func (r Record) BrokerCountry() int64 {
	return int64(Country)
}

func (r Record) AssetCountry() int64 {
	return int64(countries.ByName(r.Symbol()[:2]).Info().Code)
}

func (r Record) Side() internal.Side {
	return r.side
}

func (r Record) Quantity() decimal.Decimal {
	return r.quantity
}

func (r Record) Price() decimal.Decimal {
	return r.price
}

func (r Record) Timestamp() time.Time {
	return r.timestamp
}

func (r Record) Fees() decimal.Decimal {
	return r.fees
}

func (r Record) Taxes() decimal.Decimal {
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

func (rr RecordReader) ReadRecord(_ context.Context) (internal.Record, error) {
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
			taxes:     stampDutyTax.Add(frenchTxTax),
		}, nil
	}
}

// parseFloat attempts to parse a string using a standard precision and rounding mode.
// Using this function helps avoid issues around converting values due to sligh parameter changes.
func parseDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

// parseOptinalDecimal behaves the same as parseDecimal but returns 0 when len(s) is 0 instead of
// error.
// Using this function helps avoid issues around converting values due to sligh parameter changes.
func parseOptinalDecimal(s string) (decimal.Decimal, error) {
	if len(s) == 0 {
		return decimal.Decimal{}, nil
	}

	return parseDecimal(s)
}
