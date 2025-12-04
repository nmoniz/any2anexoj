package internal_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/shopspring/decimal"
)

func TestAggregatorWriter_Write(t *testing.T) {
	tests := []struct {
		name       string
		items      []internal.ReportItem
		wantEarned decimal.Decimal
		wantSpent  decimal.Decimal
		wantFees   decimal.Decimal
		wantTaxes  decimal.Decimal
	}{
		{
			name: "single write updates all totals",
			items: []internal.ReportItem{
				{
					Symbol:        "AAPL",
					BuyValue:      decimal.NewFromFloat(100.50),
					SellValue:     decimal.NewFromFloat(150.75),
					Fees:          decimal.NewFromFloat(2.50),
					Taxes:         decimal.NewFromFloat(5.25),
					BuyTimestamp:  time.Now(),
					SellTimestamp: time.Now(),
				},
			},
			wantEarned: decimal.NewFromFloat(150.75),
			wantSpent:  decimal.NewFromFloat(100.50),
			wantFees:   decimal.NewFromFloat(2.50),
			wantTaxes:  decimal.NewFromFloat(5.25),
		},
		{
			name: "multiple writes accumulate totals",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.NewFromFloat(100.00),
					SellValue: decimal.NewFromFloat(120.00),
					Fees:      decimal.NewFromFloat(1.00),
					Taxes:     decimal.NewFromFloat(2.00),
				},
				{
					BuyValue:  decimal.NewFromFloat(200.00),
					SellValue: decimal.NewFromFloat(250.00),
					Fees:      decimal.NewFromFloat(3.00),
					Taxes:     decimal.NewFromFloat(4.00),
				},
				{
					BuyValue:  decimal.NewFromFloat(50.00),
					SellValue: decimal.NewFromFloat(55.00),
					Fees:      decimal.NewFromFloat(0.50),
					Taxes:     decimal.NewFromFloat(1.50),
				},
			},
			wantEarned: decimal.NewFromFloat(425.00),
			wantSpent:  decimal.NewFromFloat(350.00),
			wantFees:   decimal.NewFromFloat(4.50),
			wantTaxes:  decimal.NewFromFloat(7.50),
		},
		{
			name:       "empty writer returns zero totals",
			items:      []internal.ReportItem{},
			wantEarned: decimal.Zero,
			wantSpent:  decimal.Zero,
			wantFees:   decimal.Zero,
			wantTaxes:  decimal.Zero,
		},
		{
			name: "handles zero values",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.Zero,
					SellValue: decimal.Zero,
					Fees:      decimal.Zero,
					Taxes:     decimal.Zero,
				},
			},
			wantEarned: decimal.Zero,
			wantSpent:  decimal.Zero,
			wantFees:   decimal.Zero,
			wantTaxes:  decimal.Zero,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aw := &internal.AggregatorWriter{}
			ctx := context.Background()

			for _, item := range tt.items {
				if err := aw.Write(ctx, item); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			assertDecimalEqual(t, "TotalEarned", tt.wantEarned, aw.TotalEarned())
			assertDecimalEqual(t, "TotalSpent", tt.wantSpent, aw.TotalSpent())
			assertDecimalEqual(t, "TotalFees", tt.wantFees, aw.TotalFees())
			assertDecimalEqual(t, "TotalTaxes", tt.wantTaxes, aw.TotalTaxes())
		})
	}
}

func TestAggregatorWriter_Rounding(t *testing.T) {
	tests := []struct {
		name       string
		items      []internal.ReportItem
		wantEarned decimal.Decimal
		wantSpent  decimal.Decimal
		wantFees   decimal.Decimal
		wantTaxes  decimal.Decimal
	}{
		{
			name: "rounds to 2 decimal places",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.NewFromFloat(100.123456),
					SellValue: decimal.NewFromFloat(150.987654),
					Fees:      decimal.NewFromFloat(2.555555),
					Taxes:     decimal.NewFromFloat(5.444444),
				},
			},
			wantEarned: decimal.NewFromFloat(150.99),
			wantSpent:  decimal.NewFromFloat(100.12),
			wantFees:   decimal.NewFromFloat(2.56),
			wantTaxes:  decimal.NewFromFloat(5.44),
		},
		{
			name: "rounding accumulates correctly across multiple writes",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.NewFromFloat(10.111),
					SellValue: decimal.NewFromFloat(15.999),
					Fees:      decimal.NewFromFloat(0.555),
					Taxes:     decimal.NewFromFloat(1.445),
				},
				{
					BuyValue:  decimal.NewFromFloat(20.222),
					SellValue: decimal.NewFromFloat(25.001),
					Fees:      decimal.NewFromFloat(0.444),
					Taxes:     decimal.NewFromFloat(0.555),
				},
			},
			// Each write rounds individually, then accumulates
			// First: 10.11 + 20.22 = 30.33
			// Second: 16.00 + 25.00 = 41.00
			// Fees: 0.56 + 0.44 = 1.00
			// Taxes: 1.45 + 0.56 = 2.01
			wantSpent:  decimal.NewFromFloat(30.33),
			wantEarned: decimal.NewFromFloat(41.00),
			wantFees:   decimal.NewFromFloat(1.00),
			wantTaxes:  decimal.NewFromFloat(2.01),
		},
		{
			name: "handles small fractions",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.NewFromFloat(0.001),
					SellValue: decimal.NewFromFloat(0.009),
					Fees:      decimal.NewFromFloat(0.0055),
					Taxes:     decimal.NewFromFloat(0.0045),
				},
			},
			wantSpent:  decimal.NewFromFloat(0.00),
			wantEarned: decimal.NewFromFloat(0.01),
			wantFees:   decimal.NewFromFloat(0.01),
			wantTaxes:  decimal.NewFromFloat(0.00),
		},
		{
			name: "handles large numbers with precision",
			items: []internal.ReportItem{
				{
					BuyValue:  decimal.NewFromFloat(999999.996),
					SellValue: decimal.NewFromFloat(1000000.004),
					Fees:      decimal.NewFromFloat(12345.678),
					Taxes:     decimal.NewFromFloat(54321.123),
				},
			},
			wantSpent:  decimal.NewFromFloat(1000000.00),
			wantEarned: decimal.NewFromFloat(1000000.00),
			wantFees:   decimal.NewFromFloat(12345.68),
			wantTaxes:  decimal.NewFromFloat(54321.12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aw := &internal.AggregatorWriter{}
			ctx := context.Background()

			for _, item := range tt.items {
				if err := aw.Write(ctx, item); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			assertDecimalEqual(t, "TotalEarned", tt.wantEarned, aw.TotalEarned())
			assertDecimalEqual(t, "TotalSpent", tt.wantSpent, aw.TotalSpent())
			assertDecimalEqual(t, "TotalFees", tt.wantFees, aw.TotalFees())
			assertDecimalEqual(t, "TotalTaxes", tt.wantTaxes, aw.TotalTaxes())
		})
	}
}

func TestAggregatorWriter_Items(t *testing.T) {
	aw := &internal.AggregatorWriter{}
	ctx := context.Background()

	for range 10 {
		item := internal.ReportItem{Symbol: "TEST"}
		if err := aw.Write(ctx, item); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	count := 0
	for range aw.Iter() {
		count++
	}

	if count != 5 {
		t.Errorf("expected for loop to stop at 5 items, got %d", count)
	}
}

func TestAggregatorWriter_ThreadSafety(t *testing.T) {
	aw := &internal.AggregatorWriter{}
	ctx := context.Background()

	numGoroutines := 100
	writesPerGoroutine := 100

	var wg sync.WaitGroup
	for range numGoroutines {
		wg.Go(func() {
			for range writesPerGoroutine {
				item := internal.ReportItem{
					BuyValue:  decimal.NewFromFloat(1.00),
					SellValue: decimal.NewFromFloat(2.00),
					Fees:      decimal.NewFromFloat(0.10),
					Taxes:     decimal.NewFromFloat(0.20),
				}
				if err := aw.Write(ctx, item); err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}

	wg.Wait()

	// Verify totals are correct
	wantWrites := numGoroutines * writesPerGoroutine
	wantSpent := decimal.NewFromFloat(float64(wantWrites) * 1.00)
	wantEarned := decimal.NewFromFloat(float64(wantWrites) * 2.00)
	wantFees := decimal.NewFromFloat(float64(wantWrites) * 0.10)
	wantTaxes := decimal.NewFromFloat(float64(wantWrites) * 0.20)

	assertDecimalEqual(t, "TotalSpent", wantSpent, aw.TotalSpent())
	assertDecimalEqual(t, "TotalEarned", wantEarned, aw.TotalEarned())
	assertDecimalEqual(t, "TotalFees", wantFees, aw.TotalFees())
	assertDecimalEqual(t, "TotalTaxes", wantTaxes, aw.TotalTaxes())
}

// Helper function to assert decimal equality
func assertDecimalEqual(t *testing.T, name string, expected, actual decimal.Decimal) {
	t.Helper()

	if !expected.Equal(actual) {
		t.Errorf("want %s to be %s but got %s", name, expected.String(), actual.String())
	}
}
