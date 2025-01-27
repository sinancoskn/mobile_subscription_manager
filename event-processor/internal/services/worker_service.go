package services

import (
	"context"
	"event-processor/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type WorkerService struct {
	workerRepo       *models.WorkerRepository
	batchRepo        *models.BatchRepository
	subscriptionRepo *models.SubscriptionRepository
	storeApiService  *StoreApiService
	workerID         string
}

// NewWorkerService creates a new WorkerService instance
func NewWorkerService(workerRepo *models.WorkerRepository, batchRepo *models.BatchRepository, subscriptionRepo *models.SubscriptionRepository, storeApiService *StoreApiService) *WorkerService {
	return &WorkerService{
		workerRepo:       workerRepo,
		batchRepo:        batchRepo,
		subscriptionRepo: subscriptionRepo,
		storeApiService:  storeApiService,
		workerID:         uuid.New().String(), // Generate a unique worker ID
	}
}

// Start initializes the worker and begins processing
func (s *WorkerService) Start() {
	log.Printf("Worker started with ID: %s", s.workerID)

	// Register the worker in the database
	err := s.workerRepo.RegisterWorker(s.workerID)
	if err != nil {
		log.Fatalf("Failed to register worker: %v", err)
	}

	// Start heartbeat updater in a separate goroutine
	go s.startHeartbeat()

	// Main processing loop
	for {
		// time.Sleep(5 * time.Second)

		// Fetch and lock an available batch
		batch, err := s.fetchAndLockBatch()
		if err != nil {
			log.Printf("Error fetching batch: %v", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// If no batch is available, wait and try again
		if batch == nil {
			log.Println("No available batches. Waiting...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Process the batch
		log.Printf("Processing batch: %d (Action ID: %d)", batch.ID, batch.ActionID)
		success, err := s.processBatch(batch)
		if err != nil {
			log.Printf("Failed to process batch %d: %v", batch.ID, err)
		}

		// Handle batch completion or retry
		if success {
			// Mark the batch as completed
			err = s.batchRepo.MarkBatchCompleted(batch.ID)
			if err != nil {
				log.Printf("Failed to mark batch %d as completed: %v", batch.ID, err)
			}
		} else {
			// Increment the batch try count and get the updated record
			updatedBatch, err := s.batchRepo.IncrementBatchTryCount(batch.ID)
			if err != nil {
				log.Printf("Failed to increment try count for batch %d: %v", batch.ID, err)
				continue
			}

			// Check if the batch should be marked as stale
			if updatedBatch.TryCount > 5 {
				err = s.batchRepo.MarkBatchAsStale(batch.ID)
				if err != nil {
					log.Printf("Failed to mark batch %d as stale: %v", batch.ID, err)
				}
			}
		}
	}
}

// startHeartbeat periodically updates the worker's heartbeat
func (s *WorkerService) startHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := s.workerRepo.UpdateHeartbeat(s.workerID)
		if err != nil {
			log.Printf("Failed to update heartbeat for worker %s: %v", s.workerID, err)
		}
	}
}

// fetchAndLockBatch fetches and locks an available batch
func (s *WorkerService) fetchAndLockBatch() (*models.Batch, error) {
	return s.batchRepo.LockNextBatch(s.workerID)
}

// processBatch processes the records in the batch
func (s *WorkerService) processBatch(batch *models.Batch) (bool, error) {
	log.Printf("Processing batch: ID %d, ActionID %d\n", batch.ID, batch.ActionID)

	// Fetch records in the batch
	subscriptions, err := s.subscriptionRepo.FetchSubscriptionsForBatch(batch.StartIndex, batch.EndIndex)
	if err != nil {
		return false, fmt.Errorf("failed to fetch subscriptions for batch %d: %w", batch.ID, err)
	}

	ctx := context.TODO()

	successCount := 0
	failureCount := 0

	// Prepare slices for batch updates
	var activeSubscriptions []models.Subscription
	var expiredSubscriptions []models.Subscription

	// Process each subscription
	for _, sub := range subscriptions {
		// Skip if subscription is canceled
		if sub.Status == "canceled" {
			log.Printf("Skipping subscription ID %d: status is canceled", sub.ID)
			continue
		}

		// Skip if subscription was updated in the last 30 minutes
		if time.Since(sub.UpdatedAt) < 30*time.Minute {
			log.Printf("Skipping subscription ID %d: updated within 30 minutes", sub.ID)
			continue
		}

		// Check if the subscription is expired and not canceled
		if time.Now().After(sub.ExpireAt) && sub.Status != "canceled" {
			// Request the Store API to validate the receipt
			result, err := s.storeApiService.ValidateReceipt(ctx, sub.Receipt)
			if err != nil {
				log.Printf("Failed to validate receipt for subscription ID %d: %v", sub.ID, err)
				failureCount++
				continue
			}

			// Process the Store API result
			status, ok := result["status"].(bool)
			if !ok {
				log.Printf("Invalid status received for subscription ID %d", sub.ID)
				failureCount++
				continue
			}

			expireDate, ok := result["expire_date"].(string)
			if !ok {
				log.Printf("Invalid expireDate received for subscription ID %d", sub.ID)
				failureCount++
				continue
			}

			// Parse expireDate
			expireTime, err := time.Parse("2006-01-02 15:04:05", expireDate)
			if err != nil {
				log.Printf("Failed to parse expireDate for subscription ID %d: %v", sub.ID, err)
				failureCount++
				continue
			}

			// Add to the appropriate batch
			if status {
				sub.ExpireAt = expireTime
				sub.Status = "active"
				activeSubscriptions = append(activeSubscriptions, sub)
				successCount++
			} else {
				sub.Status = "expired"
				expiredSubscriptions = append(expiredSubscriptions, sub)
				successCount++
			}
		}
	}

	// Perform batch updates
	if len(activeSubscriptions) > 0 {
		if err := s.subscriptionRepo.BulkUpdateSubscriptions(activeSubscriptions); err != nil {
			log.Printf("Failed to update active subscriptions in batch: %v", err)
		}
	}

	if len(expiredSubscriptions) > 0 {
		if err := s.subscriptionRepo.BulkUpdateSubscriptions(expiredSubscriptions); err != nil {
			log.Printf("Failed to update expired subscriptions in batch: %v", err)
		}
	}

	log.Printf("Completed processing batch: ID %d, Success: %d, Failures: %d\n", batch.ID, successCount, failureCount)

	// If all subscriptions were processed successfully, return true
	return failureCount == 0, nil
}
