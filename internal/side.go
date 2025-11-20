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

// IsBuy returns true if the s == SideBuy
func (d Side) IsBuy() bool {
	return d == SideBuy
}

// IsSell returns true if the s == SideSell
func (d Side) IsSell() bool {
	return d == SideSell
}
