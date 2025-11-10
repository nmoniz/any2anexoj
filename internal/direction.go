package internal

type Direction uint

const (
	DirectionUnknown Direction = 0
	DirectionBuy     Direction = 1
	DirectionSell    Direction = 2
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
