package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"event-processor/internal/models"
)

type WebhookService struct {
	repo *models.WebhookRepository // Injected repository for database operations
}

// NewWebhookService creates a new instance of WebhookService
func NewWebhookService(repo *models.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

// HandleEvent processes an incoming event and triggers relevant webhooks
func (s *WebhookService) HandleEvent(ctx context.Context, data map[string]interface{}) error {
	log.Printf("Processing event: %v", data)

	// Extract app_id from the event data
	appID, ok := data["app_id"].(float64) // JSON numbers are decoded as float64
	if !ok {
		err := errors.New("invalid or missing app_id in event data")
		log.Printf("%v: %v", err, data)
		return err
	}

	// Extract status from the event data
	status, ok := data["status"].(string)
	if !ok {
		err := errors.New("invalid or missing status in event data")
		log.Printf("%v: %v", err, data)
		return err
	}

	// Fetch webhooks for the given app_id and status
	webhooks, err := s.repo.GetWebhooksByEvent(int(appID), status)
	if err != nil {
		log.Printf("Failed to fetch webhooks for app_id %d and status %s: %v", int(appID), status, err)
		return err
	}

	// Process each webhook
	for _, webhook := range webhooks {
		err := s.invokeWebhook(webhook, data)
		if err != nil {
			log.Printf("Error invoking webhook %s: %v", webhook.URL, err)
			// Optional: Add retry logic or update retry counts in the database
			// _ = s.repo.IncrementTriedCount(webhook.ID)
			// _ = s.repo.UpdateLastAttempt(webhook.ID, time.Now())
			// Continue processing the next webhook
			continue
		}
		log.Printf("Successfully invoked webhook: %s", webhook.URL)
	}

	return nil // Indicate successful processing of the event
}

// invokeWebhook sends the event data to the webhook's URL
func (s *WebhookService) invokeWebhook(webhook models.Webhook, data map[string]interface{}) error {
	// Serialize the data into JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set default headers
	for key, value := range webhook.DefaultHeaders {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client and send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("received non-success status code: " + resp.Status)
	}

	return nil
}
