package internal

import "testing"

func TestDirection_String(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want string
	}{
		{"buy", DirectionBuy, "buy"},
		{"sell", DirectionSell, "sell"},
		{"unknown", DirectionUnknown, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("want Direction.String() to be %v but got %v", tt.want, got)
			}
		})
	}
}

func TestDirection_IsBuy(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want bool
	}{
		{"buy", DirectionBuy, true},
		{"sell", DirectionSell, false},
		{"unknown", DirectionUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsBuy(); got != tt.want {
				t.Errorf("want Direction.IsBuy() to be %v but got %v", tt.want, got)
			}
		})
	}
}

func TestDirection_IsSell(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want bool
	}{
		{"buy", DirectionBuy, false},
		{"sell", DirectionSell, true},
		{"unknown", DirectionUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsSell(); got != tt.want {
				t.Errorf("want Direction.IsSell() to be %v but got %v", tt.want, got)
			}
		})
	}
}
