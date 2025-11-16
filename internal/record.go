package internal

import (
	"time"

	"github.com/shopspring/decimal"
)

type Record interface {
	Symbol() string
	Side() Side
	Price() decimal.Decimal
	Quantity() decimal.Decimal
	Timestamp() time.Time
	Fees() decimal.Decimal
	Taxes() decimal.Decimal
}
