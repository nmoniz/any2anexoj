package internal_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/shopspring/decimal"
)

func TestReportLogger_Write(t *testing.T) {
	tNow := time.Now()

	tests := []struct {
		name  string
		items []internal.ReportItem
		want  []string
	}{
		{
			name: "empty",
		},
		{
			name: "single item positive",
			items: []internal.ReportItem{
				{
					BuyValue:      decimal.NewFromFloat(100.0),
					SellValue:     decimal.NewFromFloat(200.0),
					SellTimestamp: tNow,
				},
			},
			want: []string{
				fmt.Sprintf("%6d: realised 100 on %s\n", 1, tNow.Format(time.RFC3339)),
			},
		},
		{
			name: "single item negative",
			items: []internal.ReportItem{
				{
					BuyValue:      decimal.NewFromFloat(200.0),
					SellValue:     decimal.NewFromFloat(150.0),
					SellTimestamp: tNow,
				},
			},
			want: []string{
				fmt.Sprintf("%6d: realised -50 on %s\n", 1, tNow.Format(time.RFC3339)),
			},
		},
		{
			name: "multiple items",
			items: []internal.ReportItem{
				{
					BuyValue:      decimal.NewFromFloat(100.0),
					SellValue:     decimal.NewFromFloat(200.0),
					SellTimestamp: tNow,
				},
				{
					BuyValue:      decimal.NewFromFloat(200.0),
					SellValue:     decimal.NewFromFloat(150.0),
					SellTimestamp: tNow.Add(1),
				},
			},
			want: []string{
				fmt.Sprintf("%6d: realised 100 on %s\n", 1, tNow.Format(time.RFC3339)),
				fmt.Sprintf("%6d: realised -50 on %s\n", 2, tNow.Add(1).Format(time.RFC3339)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rw := internal.NewReportLogger(buf)

			for _, item := range tt.items {
				err := rw.Write(t.Context(), item)
				if err != nil {
					t.Fatalf("unexpected error on write: %v", err)
				}
			}

			for _, wantLine := range tt.want {
				gotLine, err := buf.ReadString(byte('\n'))
				if err != nil {
					t.Fatalf("unexpected error on buffer read: %v", err)
				}
				if wantLine != gotLine {
					t.Fatalf("want line %q but got %q", wantLine, gotLine)
				}
			}
		})
	}
}
