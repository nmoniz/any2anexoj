package internal

type Side uint

const (
	SideUnknown Side = iota
	SideBuy
	SideSell
)

func (d Side) String() string {
	switch d {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	default:
		return "unknown"
	}
}

func (d Side) IsBuy() bool {
	return d == SideBuy
}

func (d Side) IsSell() bool {
	return d == SideSell
}
