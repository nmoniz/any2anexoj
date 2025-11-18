package internal

import (
	"bytes"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestTableWriter_Write(t *testing.T) {
	tNow := time.Now()

	tests := []struct {
		name            string
		items           []ReportItem
		wantTotalSpent  decimal.Decimal
		wantTotalEarned decimal.Decimal
		wantTotalTaxes  decimal.Decimal
		wantTotalFees   decimal.Decimal
	}{
		{
			name: "empty",
		},
		{
			name: "single item positive",
			items: []ReportItem{
				{
					BuyValue:      decimal.NewFromFloat(100.0),
					SellValue:     decimal.NewFromFloat(200.0),
					SellTimestamp: tNow,
					Taxes:         decimal.NewFromFloat(2.5),
					Fees:          decimal.NewFromFloat(2.5),
				},
			},
			wantTotalSpent:  decimal.NewFromFloat(100.0),
			wantTotalEarned: decimal.NewFromFloat(200.0),
			wantTotalTaxes:  decimal.NewFromFloat(2.5),
			wantTotalFees:   decimal.NewFromFloat(2.5),
		},
		{
			name: "single item negative",
			items: []ReportItem{
				{
					BuyValue:      decimal.NewFromFloat(200.0),
					SellValue:     decimal.NewFromFloat(150.0),
					SellTimestamp: tNow,
					Taxes:         decimal.NewFromFloat(2.5),
					Fees:          decimal.NewFromFloat(2.5),
				},
			},
			wantTotalSpent:  decimal.NewFromFloat(200.0),
			wantTotalEarned: decimal.NewFromFloat(150.0),
			wantTotalTaxes:  decimal.NewFromFloat(2.5),
			wantTotalFees:   decimal.NewFromFloat(2.5),
		},
		{
			name: "multiple items",
			items: []ReportItem{
				{
					Symbol:        "US1912161007",
					BuyValue:      decimal.NewFromFloat(100.0),
					SellValue:     decimal.NewFromFloat(200.0),
					SellTimestamp: tNow,
					Taxes:         decimal.NewFromFloat(2.5),
					Fees:          decimal.NewFromFloat(2.5),
				},
				{
					Symbol:        "US1912161007",
					BuyValue:      decimal.NewFromFloat(200.0),
					SellValue:     decimal.NewFromFloat(150.0),
					SellTimestamp: tNow.Add(1),
					Taxes:         decimal.NewFromFloat(2.5),
					Fees:          decimal.NewFromFloat(2.5),
				},
			},
			wantTotalSpent:  decimal.NewFromFloat(300.0),
			wantTotalEarned: decimal.NewFromFloat(350.0),
			wantTotalTaxes:  decimal.NewFromFloat(5.0),
			wantTotalFees:   decimal.NewFromFloat(5.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tw := NewTableWriter(buf)

			for _, item := range tt.items {
				err := tw.Write(t.Context(), item)
				if err != nil {
					t.Fatalf("unexpected error on write: %v", err)
				}
			}

			if tw.table.Length() != len(tt.items) {
				t.Fatalf("want %d items in table but got %d", len(tt.items), tw.table.Length())
			}

			if !tw.totalSpent.Equal(tt.wantTotalSpent) {
				t.Errorf("want totalSpent to be %v but got %v", tt.wantTotalSpent, tw.totalSpent)
			}

			if !tw.totalEarned.Equal(tt.wantTotalEarned) {
				t.Errorf("want totalEarned to be %v but got %v", tt.wantTotalEarned, tw.totalEarned)
			}

			if !tw.totalTaxes.Equal(tt.wantTotalTaxes) {
				t.Errorf("want totalTaxes to be %v but got %v", tt.wantTotalTaxes, tw.totalTaxes)
			}

			if !tw.totalFees.Equal(tt.wantTotalFees) {
				t.Errorf("want totalFees to be %v but got %v", tt.wantTotalFees, tw.totalFees)
			}
		})
	}
}
