package internal

import (
	"context"
	"fmt"
	"io"

	"github.com/biter777/countries"
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
	t.SetColumnConfigs([]table.ColumnConfig{
		colCountry(1),
		colOther(2),
		colOther(3),
		colOther(4),
		colOther(5),
		colEuros(6),
		colOther(7),
		colOther(8),
		colOther(9),
		colEuros(10),
		colEuros(11),
		colEuros(12),
		colCountry(13),
	})

	t.AppendHeader(table.Row{"", "", "Realisation", "Realisation", "Realisation", "Realisation", "Acquisition", "Acquisition", "Acquisition", "Acquisition", "", "", ""}, table.RowConfig{AutoMerge: true})
	t.AppendHeader(table.Row{"Source Country", "Code", "Year", "Month", "Day", "Value", "Year", "Month", "Day", "Value", "Expenses", "Paid Taxes", "Counter Country"})

	return &TableWriter{
		table:  t,
		output: w,
	}
}

func (tw *TableWriter) Write(_ context.Context, ri ReportItem) error {
	tw.totalEarned = tw.totalEarned.Add(ri.SellValue.Round(2))
	tw.totalSpent = tw.totalSpent.Add(ri.BuyValue.Round(2))
	tw.totalFees = tw.totalFees.Add(ri.Fees.Round(2))
	tw.totalTaxes = tw.totalTaxes.Add(ri.Taxes.Round(2))

	tw.table.AppendRow(table.Row{
		ri.AssetCountry, ri.Nature,
		ri.SellTimestamp.Year(), int(ri.SellTimestamp.Month()), ri.SellTimestamp.Day(), ri.SellValue.StringFixed(2),
		ri.BuyTimestamp.Year(), int(ri.BuyTimestamp.Month()), ri.BuyTimestamp.Day(), ri.BuyValue.StringFixed(2),
		ri.Fees.StringFixed(2), ri.Taxes.StringFixed(2), ri.BrokerCountry,
	})

	return nil
}

func (tw *TableWriter) Render() {
	tw.table.AppendFooter(table.Row{"SUM", "SUM", "SUM", "SUM", "SUM", tw.totalEarned.StringFixed(2), "", "", "", tw.totalSpent.StringFixed(2), tw.totalFees.StringFixed(2), tw.totalTaxes.StringFixed(2)}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignRight})
	tw.table.Render()
}

func colEuros(n int) table.ColumnConfig {
	return table.ColumnConfig{
		Number:      n,
		Align:       text.AlignRight,
		AlignFooter: text.AlignRight,
		AlignHeader: text.AlignRight,
		WidthMin:    12,
		Transformer: func(val any) string {
			return fmt.Sprintf("%v €", val)
		},
		TransformerFooter: func(val any) string {
			return fmt.Sprintf("%v €", val)
		},
	}
}

func colOther(n int) table.ColumnConfig {
	return table.ColumnConfig{
		Number:      n,
		Align:       text.AlignLeft,
		AlignFooter: text.AlignLeft,
		AlignHeader: text.AlignLeft,
	}
}

func colCountry(n int) table.ColumnConfig {
	return table.ColumnConfig{
		Number:           n,
		Align:            text.AlignLeft,
		AlignFooter:      text.AlignLeft,
		AlignHeader:      text.AlignLeft,
		WidthMax:         24,
		WidthMaxEnforcer: text.Trim,
		Transformer: func(val any) string {
			countryCode := val.(int64)
			return fmt.Sprintf("%v - %s", val, countries.ByNumeric(int(countryCode)).Info().Name)
		},
	}
}
