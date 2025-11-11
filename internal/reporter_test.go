package internal_test

import (
	"io"
	"math/big"
	"testing"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
	"git.naterciomoniz.net/applications/broker2anexoj/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestReporter_Run(t *testing.T) {
	ctrl := gomock.NewController(t)

	rec := mocks.NewMockRecord(ctrl)
	rec.EXPECT().Price().Return(big.NewFloat(1.25)).AnyTimes()
	rec.EXPECT().Quantity().Return(big.NewFloat(10)).AnyTimes()
	rec.EXPECT().Side().Return(internal.SideBuy).AnyTimes()
	rec.EXPECT().Symbol().Return("TEST").AnyTimes()

	reader := mocks.NewMockRecordReader(ctrl)
	records := []internal.Record{
		rec,
		rec,
	}
	reader.EXPECT().ReadRecord().DoAndReturn(func() (internal.Record, error) {
		if len(records) > 0 {
			r := records[0]
			records = records[1:]
			return r, nil
		} else {
			return nil, io.EOF
		}
	}).AnyTimes()

	reporter := internal.NewReporter(reader)
	gotErr := reporter.Run()
	if gotErr != nil {
		t.Fatalf("got unexpected err: %v", gotErr)
	}
}
