package trading212

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
	"github.com/shopspring/decimal"
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
			name: "well-formed buy",
			r:    bytes.NewBufferString(`Market buy,2025-07-03 10:44:29,SYM123456ABXY,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.3690000000,USD,1.17995999,,"EUR",15.25,"EUR",0.25,"EUR",0.02,"EUR",,`),
			want: Record{
				symbol:       "SYM123456ABXY",
				side:         internal.SideBuy,
				quantity:     ShouldParseDecimal(t, "2.4387014200"),
				price:        ShouldParseDecimal(t, "7.3690000000"),
				timestamp:    time.Date(2025, 7, 3, 10, 44, 29, 0, time.UTC),
				fees:         ShouldParseDecimal(t, "0.02"),
				taxes:        ShouldParseDecimal(t, "0.25"),
				natureGetter: func() internal.Nature { return internal.NatureG01 },
			},
		},
		{
			name: "well-formed sell",
			r:    bytes.NewBufferString(`Market sell,2025-08-04 11:45:30,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",0.1,"EUR"`),
			want: Record{
				symbol:       "IE000GA3D489",
				side:         internal.SideSell,
				quantity:     ShouldParseDecimal(t, "2.4387014200"),
				price:        ShouldParseDecimal(t, "7.9999999999"),
				timestamp:    time.Date(2025, 8, 4, 11, 45, 30, 0, time.UTC),
				fees:         ShouldParseDecimal(t, "0.02"),
				taxes:        ShouldParseDecimal(t, "0.1"),
				natureGetter: func() internal.Nature { return internal.NatureG01 },
			},
		},
		{
			name:    "malformed side",
			r:       bytes.NewBufferString(`Aljksdaf Balsjdkf,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "empty side",
			r:       bytes.NewBufferString(`,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,0x1234,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "malformed qantity",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,0x1234,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "empty qantity",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "malformed price",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,0b101010,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "empty price",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:39,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "malformed fees",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:30,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,BAD,"EUR",0.1,"EUR"`),
			wantErr: true,
		},
		{
			name:    "malformed taxes",
			r:       bytes.NewBufferString(`Market sell,2025-08-04 11:45:30,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",BAD,"EUR"`),
			wantErr: true,
		},
		{
			name:    "malformed timestamp",
			r:       bytes.NewBufferString(`Market sell,2006-01-02T15:04:05Z07:00,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
		{
			name:    "empty timestamp",
			r:       bytes.NewBufferString(`Market sell,,IE000GA3D489,ABXY,"Aspargus Broccoli",EOF987654321,2.4387014200,7.9999999999,USD,1.17995999,,"EUR",15.25,"EUR",,,0.02,"EUR",,`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := NewRecordReader(tt.r, NewFigiClientSecurityTypeStub(t, "Common Stock"))
			got, gotErr := rr.ReadRecord(t.Context())
			if gotErr != nil {
				if !tt.wantErr {
					t.Fatalf("ReadRecord() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatalf("ReadRecord() expected an error")
			}

			if got.Symbol() != tt.want.symbol {
				t.Fatalf("want symbol %v but got %v", tt.want.symbol, got.Symbol())
			}

			if got.Side() != tt.want.side {
				t.Fatalf("want side %v but got %v", tt.want.side, got.Side())
			}

			if got.Price().Cmp(tt.want.price) != 0 {
				t.Fatalf("want price %v but got %v", tt.want.price, got.Price())
			}

			if got.Quantity().Cmp(tt.want.quantity) != 0 {
				t.Fatalf("want quantity %v but got %v", tt.want.quantity, got.Quantity())
			}

			if !got.Timestamp().Equal(tt.want.timestamp) {
				t.Fatalf("want timestamp %v but got %v", tt.want.timestamp, got.Timestamp())
			}

			if got.Fees().Cmp(tt.want.fees) != 0 {
				t.Fatalf("want fees %v but got %v", tt.want.fees, got.Fees())
			}

			if got.Taxes().Cmp(tt.want.taxes) != 0 {
				t.Fatalf("want taxes %v but got %v", tt.want.taxes, got.Taxes())
			}

			if tt.want.natureGetter != nil && tt.want.Nature() != got.Nature() {
				t.Fatalf("want nature %v but got %v", tt.want.Nature(), got.Nature())
			}
		})
	}
}

func Test_figiNatureGetter(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		of   *internal.OpenFIGI
		want internal.Nature
	}{
		{
			name: "Common Stock translates to G01",
			of:   NewFigiClientSecurityTypeStub(t, "Common Stock"),
			want: internal.NatureG01,
		},
		{
			name: "ETP translates to G20",
			of:   NewFigiClientSecurityTypeStub(t, "ETP"),
			want: internal.NatureG20,
		},
		{
			name: "Other translates to Unknown",
			of:   NewFigiClientSecurityTypeStub(t, "Other"),
			want: internal.NatureUnknown,
		},
		{
			name: "Request fails",
			of:   NewFigiClientErrorStub(t, fmt.Errorf("boom")),
			want: internal.NatureUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := figiNatureGetter(t.Context(), tt.of, "IR123456789")
			got := getter()
			if tt.want != got {
				t.Errorf("want %v but got %v", tt.want, got)
			}
		})
	}
}

func ShouldParseDecimal(t testing.TB, sf string) decimal.Decimal {
	t.Helper()

	bf, err := parseDecimal(sf)
	if err != nil {
		t.Fatalf("parsing decimal: %s", sf)
	}
	return bf
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewFigiClientSecurityTypeStub(t testing.TB, securityType string) *internal.OpenFIGI {
	t.Helper()

	c := &http.Client{
		Timeout: time.Second,
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`[{"data":[{"securityType":%q}]}]`, securityType))),
				Request:    req,
			}, nil
		}),
	}

	return internal.NewOpenFIGI(c)
}

func NewFigiClientErrorStub(t testing.TB, err error) *internal.OpenFIGI {
	t.Helper()

	c := &http.Client{
		Timeout: time.Second,
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, err
		}),
	}

	return internal.NewOpenFIGI(c)
}
