package models

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID        int64     `gorm:"primaryKey" json:"id"`              // Primary key
	UID       string    `gorm:"type:uuid;not null" json:"uid"`     // Unique user identifier
	AppID     int       `gorm:"not null" json:"app_id"`            // App ID
	Receipt   string    `gorm:"size:255;not null" json:"receipt"`  // Receipt string
	Status    string    `gorm:"size:20;not null" json:"status"`    // Subscription status
	ExpireAt  time.Time `gorm:"type:timestamptz" json:"expire_at"` // Expiration timestamp
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`  // Creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`  // Update timestamp
}

// SubscriptionRepository manages subscription-related database operations
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// GetCountForProcessing fetches the count of records that need processing
func (r *SubscriptionRepository) GetCountForProcessing() (int64, error) {
	var count int64
	err := r.db.Model(&Subscription{}).
		Where("expire_at <= NOW() AND status != ?", "canceled").
		Order("id ASC").
		Count(&count).Error
	return count, err
}

// UpdateSubscriptionStatus updates the status of a subscription
func (r *SubscriptionRepository) UpdateSubscriptionStatus(subscriptionID int64, status string) error {
	return r.db.Model(&Subscription{}).
		Where("id = ?", subscriptionID).
		Update("status", status).Error
}

func (r *SubscriptionRepository) FetchSubscriptionsForBatch(start int64, end int64) ([]Subscription, error) {
	var subscriptions []Subscription

	// Fetch subscriptions based on the given range and processing criteria
	err := r.db.Model(&Subscription{}).
		Where("expire_at <= NOW() AND status != ?", "canceled").
		Order("id ASC").
		Offset(int(start - 1)).      // Subtract 1 because `start` is inclusive
		Limit(int(end - start + 1)). // The range from start to end (inclusive)
		Find(&subscriptions).Error

	return subscriptions, err
}

func (r *SubscriptionRepository) BulkUpdateSubscriptions(subscriptions []Subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}

	// Build a bulk update query
	query := "UPDATE subscriptions SET status = CASE id "
	idList := []interface{}{}

	// Construct CASE statements
	for _, sub := range subscriptions {
		query += "WHEN ? THEN ? "
		idList = append(idList, sub.ID, sub.Status)
	}

	query += "END, expire_at = CASE id "

	// Add CASE statements for expire_at
	for _, sub := range subscriptions {
		query += "WHEN ? THEN ? "
		idList = append(idList, sub.ID, sub.ExpireAt)
	}

	query += "END, updated_at = ? WHERE id IN ("
	idList = append(idList, time.Now().UTC())

	// Add IDs for the WHERE clause
	for i, sub := range subscriptions {
		if i > 0 {
			query += ", "
		}
		query += "?"
		idList = append(idList, sub.ID)
	}
	query += ");"

	// Execute the raw query
	return r.db.Exec(query, idList...).Error
}

// GenerateMockSubscriptions generates mock data for subscriptions
func (r *SubscriptionRepository) GenerateMockSubscriptions(total int) error {
	const chunkSize = 1000                               // Insert 10,000 rows per batch
	statuses := []string{"active", "pending", "expired"} // Example statuses

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < total; i += chunkSize {
		var subscriptions []Subscription
		for j := 0; j < chunkSize && i+j < total; j++ {
			subscriptions = append(subscriptions, Subscription{
				UID:      uuid.New().String(),
				AppID:    rand.Intn(4) + 1, // Random app_id between 1 and 100
				Receipt:  fmt.Sprintf("receipt_%d", rand.Intn(1000000)),
				Status:   statuses[rand.Intn(len(statuses))],             // Random status
				ExpireAt: time.Now().UTC().AddDate(0, 0, -rand.Intn(30)), // Random expire_at in the past 30 days
			})
		}

		// Bulk insert
		if len(subscriptions) > 0 {
			if err := r.db.Create(&subscriptions).Error; err != nil {
				return fmt.Errorf("failed to insert mock subscriptions: %w", err)
			}
		}

		log.Printf("Inserted %d/%d subscriptions\n", i+chunkSize, total)
	}

	log.Println("Mock subscriptions generation completed!")
	return nil
}
