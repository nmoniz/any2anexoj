package internal_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/nmoniz/any2anexoj/internal/mocks"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
)

func TestBuildReport(t *testing.T) {
	now := time.Now()
	ctrl := gomock.NewController(t)

	reader := mocks.NewMockRecordReader(ctrl)
	records := []internal.Record{
		mockRecord(ctrl, 20.0, 10.0, internal.SideBuy, now),
		mockRecord(ctrl, 25.0, 10.0, internal.SideSell, now.Add(1)),
	}
	reader.EXPECT().ReadRecord(gomock.Any()).DoAndReturn(func(ctx context.Context) (internal.Record, error) {
		if len(records) > 0 {
			r := records[0]
			records = records[1:]
			return r, nil
		} else {
			return nil, io.EOF
		}
	}).Times(3)

	writer := mocks.NewMockReportWriter(ctrl)
	writer.EXPECT().Write(gomock.Any(), eqReportItem(internal.ReportItem{
		BuyValue:      decimal.NewFromFloat(200.0),
		BuyTimestamp:  now,
		SellValue:     decimal.NewFromFloat(250.0),
		SellTimestamp: now.Add(1),
		Fees:          decimal.Decimal{},
		Taxes:         decimal.Decimal{},
	})).Times(1)

	gotErr := internal.BuildReport(t.Context(), reader, writer)
	if gotErr != nil {
		t.Fatalf("got unexpected err: %v", gotErr)
	}
}

func mockRecord(ctrl *gomock.Controller, price, quantity float64, side internal.Side, ts time.Time) *mocks.MockRecord {
	rec := mocks.NewMockRecord(ctrl)
	rec.EXPECT().Price().Return(decimal.NewFromFloat(price)).AnyTimes()
	rec.EXPECT().Quantity().Return(decimal.NewFromFloat(quantity)).AnyTimes()
	rec.EXPECT().Side().Return(side).AnyTimes()
	rec.EXPECT().Symbol().Return("TEST").AnyTimes()
	rec.EXPECT().Timestamp().Return(ts).AnyTimes()
	rec.EXPECT().Fees().Return(decimal.Decimal{}).AnyTimes()
	rec.EXPECT().Taxes().Return(decimal.Decimal{}).AnyTimes()
	return rec
}

func eqReportItem(ri internal.ReportItem) ReportItemMatcher {
	return ReportItemMatcher{
		ReportItem: ri,
	}
}

type ReportItemMatcher struct {
	internal.ReportItem
}

// Matches implements gomock.Matcher.
func (m ReportItemMatcher) Matches(x any) bool {
	if x == nil {
		return false
	}

	switch other := x.(type) {
	case internal.ReportItem:
		return m.BuyValue.Equal(other.BuyValue) &&
			m.BuyTimestamp.Equal(other.BuyTimestamp) &&
			m.SellValue.Equal(other.SellValue) &&
			m.SellTimestamp.Equal(other.SellTimestamp) &&
			m.Fees.Equal(other.Fees) &&
			m.Taxes.Equal(other.Taxes)
	default:
		return false
	}
}

func (m ReportItemMatcher) String() string {
	return fmt.Sprintf("is equivalent to %v", m.ReportItem)
}

var _ gomock.Matcher = (*ReportItemMatcher)(nil)
