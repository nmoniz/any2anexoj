package internal

import "testing"

func TestSide_String(t *testing.T) {
	tests := []struct {
		name string
		side Side
		want string
	}{
		{"buy", SideBuy, "buy"},
		{"sell", SideSell, "sell"},
		{"unknown", SideUnknown, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.String(); got != tt.want {
				t.Errorf("want Side.String() to be %v but got %v", tt.want, got)
			}
		})
	}
}

func TestSide_IsBuy(t *testing.T) {
	tests := []struct {
		name string
		side Side
		want bool
	}{
		{"buy", SideBuy, true},
		{"sell", SideSell, false},
		{"unknown", SideUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.IsBuy(); got != tt.want {
				t.Errorf("want Side.IsBuy() to be %v but got %v", tt.want, got)
			}
		})
	}
}

func TestSide_IsSell(t *testing.T) {
	tests := []struct {
		name string
		side Side
		want bool
	}{
		{"buy", SideBuy, false},
		{"sell", SideSell, true},
		{"unknown", SideUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.IsSell(); got != tt.want {
				t.Errorf("want Side.IsSell() to be %v but got %v", tt.want, got)
			}
		})
	}
}
