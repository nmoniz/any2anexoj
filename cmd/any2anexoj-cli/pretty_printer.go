package main

import (
	"fmt"
	"io"

	"github.com/biter777/countries"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/nmoniz/any2anexoj/internal"
)

// PrettyPrinter writes a simple, human readable, table row to the provided io.Writer for each
// ReportItem received.
type PrettyPrinter struct {
	table      table.Writer
	output     io.Writer
	translator Translator
}

type Translator interface {
	Translate(key string, count int, values map[string]any) string
}

func NewPrettyPrinter(w io.Writer, tr Translator) *PrettyPrinter {
	tw := table.NewWriter()
	tw.SetOutputMirror(w)
	tw.SetAutoIndex(true)
	tw.SetStyle(table.StyleLight)
	tw.SetColumnConfigs([]table.ColumnConfig{
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

	return &PrettyPrinter{
		table:      tw,
		output:     w,
		translator: tr,
	}
}

func (pp *PrettyPrinter) Render(aw *internal.AggregatorWriter) {
	realizationTxt := pp.translator.Translate("realization", 1, nil)
	acquisitionTxt := pp.translator.Translate("acquisition", 1, nil)
	yearTxt := pp.translator.Translate("year", 1, nil)
	monthTxt := pp.translator.Translate("month", 1, nil)
	dayTxt := pp.translator.Translate("day", 1, nil)
	valorTxt := pp.translator.Translate("value", 1, nil)

	pp.table.AppendHeader(table.Row{"", "", realizationTxt, realizationTxt, realizationTxt, realizationTxt, acquisitionTxt, acquisitionTxt, acquisitionTxt, acquisitionTxt, "", "", ""}, table.RowConfig{AutoMerge: true})
	pp.table.AppendHeader(table.Row{
		pp.translator.Translate("source_country", 1, nil), pp.translator.Translate("code", 1, nil),
		yearTxt, monthTxt, dayTxt, valorTxt,
		yearTxt, monthTxt, dayTxt, valorTxt,
		pp.translator.Translate("expenses", 2, nil), pp.translator.Translate("foreign_tax_paid", 1, nil), pp.translator.Translate("counter_country", 1, nil),
	})

	for ri := range aw.Iter() {
		pp.table.AppendRow(table.Row{
			ri.AssetCountry, ri.Nature,
			ri.SellTimestamp.Year(), int(ri.SellTimestamp.Month()), ri.SellTimestamp.Day(), ri.SellValue.StringFixed(2),
			ri.BuyTimestamp.Year(), int(ri.BuyTimestamp.Month()), ri.BuyTimestamp.Day(), ri.BuyValue.StringFixed(2),
			ri.Fees.StringFixed(2), ri.Taxes.StringFixed(2), ri.BrokerCountry,
		})
	}

	pp.table.AppendFooter(table.Row{"SUM", "SUM", "SUM", "SUM", "SUM", aw.TotalEarned(), "", "", "", aw.TotalSpent(), aw.TotalFees(), aw.TotalTaxes()}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignRight})
	pp.table.Render()
}

func colEuros(n int) table.ColumnConfig {
	return table.ColumnConfig{
		Number:      n,
		Align:       text.AlignRight,
		AlignFooter: text.AlignRight,
		AlignHeader: text.AlignRight,
		WidthMin:    12,
		WidthMax:    15,
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
		WidthMax:    12,
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
