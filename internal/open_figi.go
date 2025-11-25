package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/biter777/countries"
	"golang.org/x/time/rate"
)

// OpenFIGI is a small adapter for the openfigi.com api
type OpenFIGI struct {
	client         *http.Client
	mappingLimiter *rate.Limiter

	mu sync.RWMutex
	// TODO: there's no eviction policy at the moment as this is only used by short-lived application
	// which processes a relatively small amount of records. We need to consider using an external
	// cache lib (like golang-lru or go-cache) if this becomes a problem or implement this ourselves.
	securityTypeCache map[string]string
}

func NewOpenFIGI(c *http.Client) *OpenFIGI {
	return &OpenFIGI{
		client:         c,
		mappingLimiter: rate.NewLimiter(rate.Every(time.Minute), 25), // https://www.openfigi.com/api/documentation#rate-limits

		securityTypeCache: make(map[string]string),
	}
}

func (of *OpenFIGI) SecurityTypeByISIN(ctx context.Context, isin string) (string, error) {
	of.mu.RLock()
	if secType, ok := of.securityTypeCache[isin]; ok {
		of.mu.RUnlock()
		return secType, nil
	}
	of.mu.RUnlock()

	of.mu.Lock()
	defer of.mu.Unlock()

	// we check again because there could be more than one concurrent cache miss and we want only one
	// of them to result in an actual request. When the first one releases the lock the following
	// reads will hit the cache.
	if secType, ok := of.securityTypeCache[isin]; ok {
		return secType, nil
	}

	if len(isin) != 12 || countries.ByName(isin[:2]) == countries.Unknown {
		return "", fmt.Errorf("invalid ISIN: %s", isin)
	}

	rawBody, err := json.Marshal([]mappingRequestBody{{
		IDType:  "ID_ISIN",
		IDValue: isin,
	}})
	if err != nil {
		return "", fmt.Errorf("marshal mapping request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openfigi.com/v3/mapping", bytes.NewBuffer(rawBody))
	if err != nil {
		return "", fmt.Errorf("create mapping request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	err = of.mappingLimiter.Wait(ctx)
	if err != nil {
		return "", fmt.Errorf("wait for mapping request capacity: %w", err)
	}

	res, err := of.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("make mapping request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("bad mapping response status code: %s", res.Status)
	}

	var resBody []mappingResponseBody
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resBody) == 0 {
		return "", fmt.Errorf("missing top-level elements")
	}

	if len(resBody[0].Data) == 0 {
		return "", fmt.Errorf("missing data elements")
	}

	// It is not possible that an isin is assign to diferent security types, therefore we can assume
	// all entries have the same securityType value.
	secType := resBody[0].Data[0].SecurityType
	if secType == "" {
		return "", fmt.Errorf("empty security type returned for ISIN: %s", isin)
	}

	of.securityTypeCache[isin] = secType

	return secType, nil
}

type mappingRequestBody struct {
	IDType  string `json:"idType"`
	IDValue string `json:"idValue"`
}

type mappingResponseBody struct {
	Data []struct {
		FIGI         string `json:"figi"`
		SecurityType string `json:"securityType"`
		Ticker       string `json:"ticker"`
	} `json:"data"`
}
