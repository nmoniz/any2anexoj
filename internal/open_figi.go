package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenFIGI struct {
	client http.Client

	cache map[string]string
}

func NewOpenFIGI() *OpenFIGI {
	return &OpenFIGI{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (of *OpenFIGI) Category(ctx context.Context, isin, ticker string) (string, error) {
	rawBody, err := json.Marshal(mappingRequestBody{
		IdType:  "ID_ISIN",
		IdValue: isin,
	})
	if err != nil {
		return "", fmt.Errorf("marshal mapping request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openfigi.com/v3/mapping", bytes.NewBuffer(rawBody))
	if err != nil {
		return "", fmt.Errorf("create mapping request: %w", err)
	}

	res, err := of.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("make mapping request: %w", err)
	}
	defer res.Body.Close()

	var resBody mappingResponseBody
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	return "", nil
}

type mappingRequestBody struct {
	IdType  string `json:"idType"`
	IdValue string `json:"idValue"`
}

type mappingResponseBody struct{}
