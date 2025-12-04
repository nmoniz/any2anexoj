package main

import "testing"

func TestNewLocalizer(t *testing.T) {
	tests := []struct {
		name string
		lang string
	}{
		{"english", "en"},
		{"portuguese", "pt"},
		{"english with region", "en-US"},
		{"portuguese with region", "pt-BR"},
		{"unknown language falls back to default", "!!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLocalizer(tt.lang)
			if err != nil {
				t.Fatalf("want success call but failed: %v", err)
			}
		})
	}
}
