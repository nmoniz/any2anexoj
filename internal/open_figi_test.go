package internal_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
)

func TestOpenFIGI_SecurityTypeByISIN(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		response *http.Response
		isin     string
		want     string
		wantErr  bool
	}{
		{
			name: "all good",
			response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`[{"data":[{"figi":"BBG000BJJR23","name":"AIRBUS SE","ticker":"EADSF","exchCode":"US","compositeFIGI":"BBG000BJJR23","securityType":"Common Stock","marketSector":"Equity","shareClassFIGI":"BBG001S8TFZ6","securityType2":"Common Stock","securityDescription":"EADSF"},{"figi":"BBG000BJJXJ2","name":"AIRBUS SE","ticker":"EADSF","exchCode":"PQ","compositeFIGI":"BBG000BJJR23","securityType":"Common Stock","marketSector":"Equity","shareClassFIGI":"BBG001S8TFZ6","securityType2":"Common Stock","securityDescription":"EADSF"}]}]`)),
			},
			isin: "NL0000235190",
			want: "Common Stock",
		},
		{
			name: "bas status code",
			response: &http.Response{
				Status:     http.StatusText(http.StatusTooManyRequests),
				StatusCode: http.StatusTooManyRequests,
			},
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "empty top-level",
			response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`[]`)),
			},
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "empty data elements",
			response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`[{"data":[]}]`)),
			},
			isin:    "NL0000235190",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return tt.response, nil
			})

			of := internal.NewOpenFIGI(c)

			got, gotErr := of.SecurityTypeByISIN(context.Background(), tt.isin)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("want success but failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("want error but none")
			}

			if tt.want != got {
				t.Fatalf("want security type to be %s but got %s", tt.want, got)
			}
		})
	}
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(t testing.TB, fn RoundTripFunc) *http.Client {
	t.Helper()

	return &http.Client{
		Timeout:   time.Second,
		Transport: fn,
	}
}
