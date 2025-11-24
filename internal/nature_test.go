package internal_test

import (
	"testing"

	"github.com/nmoniz/any2anexoj/internal"
)

func TestNature_StringUnknow(t *testing.T) {
	tests := []struct {
		name   string
		nature internal.Nature
		want   string
	}{
		{
			name: "return unknown",
			want: "unknown",
		},
		{
			name:   "return unknown",
			nature: internal.NatureG01,
			want:   "G01",
		},
		{
			name:   "return unknown",
			nature: internal.NatureG20,
			want:   "G20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.nature.String()
			if tt.want != got {
				t.Fatalf("want %q but got %q", tt.want, got)
			}
		})
	}
}
