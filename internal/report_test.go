package internal_test

import (
	"context"
	"io"
	"math/big"
	"testing"
	"time"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
	"git.naterciomoniz.net/applications/broker2anexoj/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestReporter_Run(t *testing.T) {
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
	writer.EXPECT().Write(gomock.Any(), gomock.Eq(internal.ReportItem{
		BuyValue:      new(big.Float).SetFloat64(200.0),
		BuyTimestamp:  now,
		SellValue:     new(big.Float).SetFloat64(250.0),
		SellTimestamp: now.Add(1),
		Fees:          new(big.Float),
		Taxes:         new(big.Float),
	})).Times(1)

	gotErr := internal.BuildReport(t.Context(), reader, writer)
	if gotErr != nil {
		t.Fatalf("got unexpected err: %v", gotErr)
	}
}

func mockRecord(ctrl *gomock.Controller, price, quantity float64, side internal.Side, ts time.Time) *mocks.MockRecord {
	rec := mocks.NewMockRecord(ctrl)
	rec.EXPECT().Price().Return(big.NewFloat(price)).AnyTimes()
	rec.EXPECT().Quantity().Return(big.NewFloat(quantity)).AnyTimes()
	rec.EXPECT().Side().Return(side).AnyTimes()
	rec.EXPECT().Symbol().Return("TEST").AnyTimes()
	rec.EXPECT().Timestamp().Return(ts).AnyTimes()
	rec.EXPECT().Fees().Return(new(big.Float)).AnyTimes()
	rec.EXPECT().Taxes().Return(new(big.Float)).AnyTimes()
	return rec
}
