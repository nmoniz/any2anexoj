package internal_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/nmoniz/any2anexoj/internal"
)

func TestOpenFIGI_SecurityTypeByISIN(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		client  *http.Client
		isin    string
		want    string
		wantErr bool
	}{
		{
			name: "all good",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[{"data":[{"figi":"BBG000BJJR23","name":"AIRBUS SE","ticker":"EADSF","exchCode":"US","compositeFIGI":"BBG000BJJR23","securityType":"Common Stock","marketSector":"Equity","shareClassFIGI":"BBG001S8TFZ6","securityType2":"Common Stock","securityDescription":"EADSF"},{"figi":"BBG000BJJXJ2","name":"AIRBUS SE","ticker":"EADSF","exchCode":"PQ","compositeFIGI":"BBG000BJJR23","securityType":"Common Stock","marketSector":"Equity","shareClassFIGI":"BBG001S8TFZ6","securityType2":"Common Stock","securityDescription":"EADSF"}]}]`)),
				}, nil
			}),
			isin: "NL0000235190",
			want: "Common Stock",
		},
		{
			name: "bad status code",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     http.StatusText(http.StatusTooManyRequests),
					StatusCode: http.StatusTooManyRequests,
				}, nil
			}),
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "bad json",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"bad": "json"}`)),
				}, nil
			}),
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "empty top-level",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[]`)),
				}, nil
			}),
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "empty data elements",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[{"data":[]}]`)),
				}, nil
			}),
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "client error",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("boom")
			}),
			isin:    "NL0000235190",
			wantErr: true,
		},
		{
			name: "empty isin",
			client: NewTestClient(t, func(req *http.Request) (*http.Response, error) {
				t.Fatalf("should not make api request")
				return nil, nil
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			of := internal.NewOpenFIGI(tt.client)

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

func TestOpenFIGI_SecurityTypeByISIN_Cache(t *testing.T) {
	var alreadyCalled bool
	c := NewTestClient(t, func(req *http.Request) (*http.Response, error) {
		if alreadyCalled {
			t.Fatalf("want requests to be cached")
		}

		alreadyCalled = true
		return &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`[{"data":[{"securityType":"Common Stock"}]}]`)),
		}, nil
	})

	of := internal.NewOpenFIGI(c)

	got, gotErr := of.SecurityTypeByISIN(t.Context(), "NL0000235190")
	if gotErr != nil {
		t.Fatalf("want 1st success call but got error: %v", gotErr)
	}

	if got != "Common Stock" {
		t.Fatalf("want 1st securityType to be %q but got %q", "Common Stock", got)
	}

	got, gotErr = of.SecurityTypeByISIN(t.Context(), "NL0000235190")
	if gotErr != nil {
		t.Fatalf("want 2nd success call but got error: %v", gotErr)
	}

	if got != "Common Stock" {
		t.Fatalf("want 2nd securityType to be %q but got %q", "Common Stock", got)
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
