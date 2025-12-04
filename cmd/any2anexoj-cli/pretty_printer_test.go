package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/shopspring/decimal"
)

func TestPrettyPrinter_Render(t *testing.T) {
	// Create test data
	aw := internal.NewAggregatorWriter()
	ctx := context.Background()

	// Add some sample report items
	err := aw.Write(ctx, internal.ReportItem{
		Symbol:        "AAPL",
		Nature:        internal.NatureG01,
		BrokerCountry: 826, // United Kingdom
		AssetCountry:  840, // United States
		BuyValue:      decimal.NewFromFloat(100.50),
		BuyTimestamp:  time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		SellValue:     decimal.NewFromFloat(150.75),
		SellTimestamp: time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC),
		Fees:          decimal.NewFromFloat(2.50),
		Taxes:         decimal.NewFromFloat(5.00),
	})
	if err != nil {
		t.Fatalf("failed to write first report item: %v", err)
	}

	err = aw.Write(ctx, internal.ReportItem{
		Symbol:        "GOOGL",
		Nature:        internal.NatureG20,
		BrokerCountry: 826, // United Kingdom
		AssetCountry:  840, // United States
		BuyValue:      decimal.NewFromFloat(200.00),
		BuyTimestamp:  time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC),
		SellValue:     decimal.NewFromFloat(225.50),
		SellTimestamp: time.Date(2023, 9, 5, 0, 0, 0, 0, time.UTC),
		Fees:          decimal.NewFromFloat(3.00),
		Taxes:         decimal.NewFromFloat(7.50),
	})
	if err != nil {
		t.Fatalf("failed to write second report item: %v", err)
	}

	// Create English localizer
	localizer, err := NewLocalizer("en")
	if err != nil {
		t.Fatalf("failed to create localizer: %v", err)
	}

	// Create pretty printer with buffer
	var buf bytes.Buffer
	pp := NewPrettyPrinter(&buf, localizer)

	// Render the table
	pp.Render(aw)

	// Get the output
	got := buf.String()

	// Expected output
	want := `┌───┬────────────────────────────┬───────────────────────────────────┬───────────────────────────────────┬──────────────────────────────────────────────────────────┐
│   │                            │            REALIZATION            │            ACQUISITION            │                                                          │
│   │ SOURCE COUNTRY      │ CODE │ YEAR │ MONTH │ DAY │        VALUE │ YEAR │ MONTH │ DAY │        VALUE │ EXPENSES AND CH │ TAX PAID ABROAD │ COUNTER COUNTRY      │
│   │                     │      │      │       │     │              │      │       │     │              │           ARGES │                 │                      │
├───┼─────────────────────┼──────┼──────┼───────┼─────┼──────────────┼──────┼───────┼─────┼──────────────┼─────────────────┼─────────────────┼──────────────────────┤
│ 1 │ 840 - United States │ G01  │ 2023 │ 6     │ 20  │     150.75 € │ 2023 │ 1     │ 15  │     100.50 € │          2.50 € │          5.00 € │ 826 - United Kingdom │
│ 2 │ 840 - United States │ G20  │ 2023 │ 9     │ 5   │     225.50 € │ 2023 │ 3     │ 10  │     200.00 € │          3.00 € │          7.50 € │ 826 - United Kingdom │
├───┼─────────────────────┴──────┴──────┴───────┴─────┼──────────────┼──────┴───────┴─────┼──────────────┼─────────────────┼─────────────────┼──────────────────────┤
│   │                                             SUM │     376.25 € │                    │      300.5 € │           5.5 € │          12.5 € │                      │
└───┴─────────────────────────────────────────────────┴──────────────┴────────────────────┴──────────────┴─────────────────┴─────────────────┴──────────────────────┘
`

	// Compare output
	if got != want {
		t.Errorf("PrettyPrinter.Render() output doesn't match expected.\n\nGot:\n%s\n\nWant:\n%s", got, want)
	}
}
