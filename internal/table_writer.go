package internal

import (
	"context"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/shopspring/decimal"
)

// TableWriter writes a simple, human readable, table row to the provided io.Writer for each
// ReportItem received.
type TableWriter struct {
	table  table.Writer
	output io.Writer

	totalEarned decimal.Decimal
	totalSpent  decimal.Decimal
	totalFees   decimal.Decimal
	totalTaxes  decimal.Decimal
}

func NewTableWriter(w io.Writer) *TableWriter {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetAutoIndex(true)
	t.SetStyle(table.StyleLight)

	t.AppendHeader(table.Row{"", "", "Realisation", "Realisation", "Realisation", "Realisation", "Acquisition", "Acquisition", "Acquisition", "Acquisition", "", "", ""}, table.RowConfig{AutoMerge: true})
	t.AppendHeader(table.Row{"Source Country", "Code", "Year", "Month", "Day", "Value", "Year", "Month", "Day", "Value", "Expenses", "Paid Taxes", "Counter Country"})

	return &TableWriter{
		table:  t,
		output: w,
	}
}

func (tw *TableWriter) Write(_ context.Context, ri ReportItem) error {
	tw.totalEarned = tw.totalEarned.Add(ri.SellValue)
	tw.totalSpent = tw.totalSpent.Add(ri.BuyValue)
	tw.totalFees = tw.totalFees.Add(ri.Fees)
	tw.totalTaxes = tw.totalTaxes.Add(ri.Taxes)

	tw.table.AppendRow(table.Row{ri.AssetCountry, ri.Nature, ri.SellTimestamp.Year(), int(ri.SellTimestamp.Month()), ri.SellTimestamp.Day(), ri.SellValue, ri.BuyTimestamp.Year(), ri.BuyTimestamp.Month(), ri.BuyTimestamp.Day(), ri.BuyValue, ri.Fees, ri.Taxes, ri.BrokerCountry})

	return nil
}

func (tw *TableWriter) Render() {
	tw.table.AppendFooter(table.Row{"SUM", "SUM", "SUM", "SUM", "SUM", tw.totalEarned, "", "", "", tw.totalSpent, tw.totalFees, tw.totalTaxes}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignRight})
	tw.table.Render()
}
