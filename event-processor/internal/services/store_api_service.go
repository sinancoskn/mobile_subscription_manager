package services

import (
	"bytes"
	"context"
	"encoding/json"
	"event-processor/internal/config"
	"fmt"
	"io"
	"net/http"
)

// StoreApiService provides methods to interact with the store API
type StoreApiService struct {
	apiHost string
	client  *http.Client
}

// NewStoreApiService creates a new instance of StoreApiService
func NewStoreApiService(config *config.Config) (*StoreApiService, error) {
	apiHost := config.StoreApiHost
	return &StoreApiService{
		apiHost: apiHost,
		client:  &http.Client{}, // Default HTTP client
	}, nil
}

// ValidateReceipt calls the store API to validate the receipt
func (s *StoreApiService) ValidateReceipt(ctx context.Context, data string) (map[string]interface{}, error) {
	// Construct the endpoint URL
	endpoint := fmt.Sprintf("%s/validate-receipt", s.apiHost)

	requestBody := map[string]interface{}{
		"receipt": data,
	}

	// Marshal the request data into JSON
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("store API call failed with HTTP code %d: %s", resp.StatusCode, string(body))
	}

	// Unmarshal the JSON response
	var decodedResponse map[string]interface{}
	if err := json.Unmarshal(body, &decodedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return decodedResponse, nil
}
