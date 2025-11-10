package internal

import (
	"math/big"
	"time"
)

type Record interface {
	Symbol() string
	Price() *big.Float
	Quantity() *big.Float
	Direction() Direction
	Timestamp() time.Time
}

type RecordReader interface {
	ReadRecord() (Record, error)
}
