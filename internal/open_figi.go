package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenFIGI struct {
	client *http.Client
}

func NewOpenFIGI(c *http.Client) *OpenFIGI {
	return &OpenFIGI{
		client: c,
	}
}

func (of *OpenFIGI) SecurityTypeByISIN(ctx context.Context, isin string) (string, error) {
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
	return resBody[0].Data[0].SecurityType, nil
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
