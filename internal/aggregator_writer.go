package internal

import (
	"context"
	"iter"
	"sync"

	"github.com/shopspring/decimal"
)

// AggregatorWriter tracks ReportItem totals.
type AggregatorWriter struct {
	mu sync.RWMutex

	items []ReportItem

	totalEarned decimal.Decimal
	totalSpent  decimal.Decimal
	totalFees   decimal.Decimal
	totalTaxes  decimal.Decimal
}

func NewAggregatorWriter() *AggregatorWriter {
	return &AggregatorWriter{}
}

func (aw *AggregatorWriter) Write(_ context.Context, ri ReportItem) error {
	aw.mu.Lock()
	defer aw.mu.Unlock()

	aw.items = append(aw.items, ri)

	aw.totalEarned = aw.totalEarned.Add(ri.SellValue.Round(2))
	aw.totalSpent = aw.totalSpent.Add(ri.BuyValue.Round(2))
	aw.totalFees = aw.totalFees.Add(ri.Fees.Round(2))
	aw.totalTaxes = aw.totalTaxes.Add(ri.Taxes.Round(2))

	return nil
}

func (aw *AggregatorWriter) Iter() iter.Seq[ReportItem] {
	aw.mu.RLock()
	itemsCopy := make([]ReportItem, len(aw.items))
	copy(itemsCopy, aw.items)
	aw.mu.RUnlock()

	return func(yield func(ReportItem) bool) {
		for _, ri := range itemsCopy {
			if !yield(ri) {
				return
			}
		}
	}
}

func (aw *AggregatorWriter) TotalEarned() decimal.Decimal {
	aw.mu.RLock()
	defer aw.mu.RUnlock()
	return aw.totalEarned
}

func (aw *AggregatorWriter) TotalSpent() decimal.Decimal {
	aw.mu.RLock()
	defer aw.mu.RUnlock()
	return aw.totalSpent
}

func (aw *AggregatorWriter) TotalFees() decimal.Decimal {
	aw.mu.RLock()
	defer aw.mu.RUnlock()
	return aw.totalFees
}

func (aw *AggregatorWriter) TotalTaxes() decimal.Decimal {
	aw.mu.RLock()
	defer aw.mu.RUnlock()
	return aw.totalTaxes
}
