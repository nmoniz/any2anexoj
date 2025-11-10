package trading212

import (
	"bytes"
	"io"
	"math/big"
	"testing"
	"time"

	"git.naterciomoniz.net/applications/broker2anexoj/internal"
)

func TestRecordReader_ReadRecord(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    Record
		wantErr bool
	}{
		{
			name:    "empty reader",
			r:       bytes.NewBufferString(""),
			want:    Record{},
			wantErr: true,
		},
		{
			name: "well formed buy",
			r:    bytes.NewBufferString(`Market buy,2025-07-03 10:44:29,SYM123456ABXY,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,7.3690000000,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want: Record{
				symbol:    "SYM123456ABXY",
				direction: internal.DirectionBuy,
				quantity:  ShouldParseDecimal(t, "2.4387014200"),
				price:     ShouldParseDecimal(t, "7.3690000000"),
				timestamp: time.Date(2025, 7, 3, 10, 44, 29, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "well formed sell",
			r:    bytes.NewBufferString(`Market sell,2025-08-04 11:45:30,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want: Record{
				symbol:    "IE000GA3D489",
				direction: internal.DirectionSell,
				quantity:  ShouldParseDecimal(t, "2.4387014200"),
				price:     ShouldParseDecimal(t, "7.9999999999"),
				timestamp: time.Date(2025, 8, 4, 11, 45, 30, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "malformed direction",
			r:       bytes.NewBufferString(`Aljksdaf Balsjdkf,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "empty direction",
			r:       bytes.NewBufferString(`,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,0x1234,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "malformed qantity",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,0x1234,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "empty qantity",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "malformed price",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,0b101010,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "empty price",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "malformed timestamp",
			r:       bytes.NewBufferString(`Market sell,2006-01-02T15:04:05Z07:00,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
		{
			name:    "empty timestamp",
			r:       bytes.NewBufferString(`Market sell,,IE000GA3D489,ABXY,"Aspargus Brocoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			want:    Record{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := NewRecordReader(tt.r)
			got, gotErr := rr.ReadRecord()
			if gotErr != nil {
				if !tt.wantErr {
					t.Fatalf("ReadRecord() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatalf("ReadRecord() expected an error")
			}

			if got.symbol != tt.want.symbol {
				t.Fatalf("want symbol %v but got %v", tt.want.symbol, got.symbol)
			}

			if got.direction != tt.want.direction {
				t.Fatalf("want direction %v but got %v", tt.want.direction, got.direction)
			}

			if got.price.Cmp(tt.want.price) != 0 {
				t.Fatalf("want price %v but got %v", tt.want.price, got.price)
			}

			if got.quantity.Cmp(tt.want.quantity) != 0 {
				t.Fatalf("want quantity %v but got %v", tt.want.quantity, got.quantity)
			}

			if !got.timestamp.Equal(tt.want.timestamp) {
				t.Fatalf("want timestamp %v but got %v", tt.want.timestamp, got.timestamp)
			}
		})
	}
}

func ShouldParseDecimal(t testing.TB, sf string) *big.Float {
	t.Helper()

	bf, err := parseDecimal(sf)
	if err != nil {
		t.Fatalf("parsing decimal: %s", sf)
	}
	return bf
}
