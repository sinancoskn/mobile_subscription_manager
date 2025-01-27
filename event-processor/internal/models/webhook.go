package models

import (
	"time"

	"gorm.io/gorm"
)

// Webhook represents a row in the "webhooks" table
type Webhook struct {
	ID             int64             `gorm:"primaryKey" json:"id"`
	AppID          string            `gorm:"size:255;not null" json:"app_id"` // Identifier for the app
	URL            string            `gorm:"size:1024;not null" json:"url"`   // Callback URL
	DefaultHeaders map[string]string `gorm:"-" json:"default_headers"`        // Headers for the callback (stored as JSON in DB)
	TriggerEvents  []string          `gorm:"-" json:"trigger_events"`         // List of trigger events (stored as JSON in DB)
	TriedCount     int               `gorm:"default:0" json:"tried_count"`    // Number of times the webhook was attempted
	LastAttempt    time.Time         `json:"last_attempt"`                    // Timestamp of the last attempt
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the database table name for the Webhook model
func (Webhook) TableName() string {
	return "webhooks"
}

// WebhookRepository is a repository for managing webhooks
type WebhookRepository struct {
	db *gorm.DB
}

// NewWebhookRepository creates a new instance of WebhookRepository
func NewWebhookRepository(db *gorm.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

// GetWebhooksByEvent fetches all webhooks triggered by a specific event
func (r *WebhookRepository) GetWebhooksByEvent(appID int, event string) ([]Webhook, error) {
	var webhooks []Webhook

	// Query webhooks by app_id and trigger_events
	err := r.db.Where(
		"app_id = ? AND trigger_events @> ?", // Postgres JSONB containment operator @>
		appID,                                // App ID to match
		`["`+event+`"]`,                      // JSON array containing the event
	).Find(&webhooks).Error

	return webhooks, err
}
