package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Worker represents a worker in the system
type Worker struct {
	ID             int64     `gorm:"primaryKey"`
	WorkerID       string    `gorm:"type:uuid;not null;unique"` // Unique worker instance identifier
	Status         string    `gorm:"size:20;not null"`          // Worker status: "idle", "processing", "stale"
	LastHeartbeat  time.Time `gorm:"type:timestamptz"`          // Timestamp of the last heartbeat
	ActionID       *int64    `gorm:"default:null"`              // Reference to the current manager action
	CurrentBatchID *int64    `gorm:"default:null"`              // Reference to the batch being processed
	CreatedAt      time.Time `gorm:"autoCreateTime"`            // Creation timestamp
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`            // Update timestamp
}

// WorkerRepository handles operations related to the Worker model
type WorkerRepository struct {
	db *gorm.DB
}

// NewWorkerRepository creates a new instance of WorkerRepository
func NewWorkerRepository(db *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

// RegisterWorker registers a new worker in the database
func (r *WorkerRepository) RegisterWorker(workerID string) error {
	worker := Worker{
		WorkerID:      workerID,
		Status:        "idle",
		LastHeartbeat: time.Now(),
	}

	if err := r.db.Create(&worker).Error; err != nil {
		return err
	}

	log.Printf("Worker registered with ID: %s", workerID)
	return nil
}

// UpdateHeartbeat updates the last heartbeat of a worker
func (r *WorkerRepository) UpdateHeartbeat(workerID string) error {
	return r.db.Model(&Worker{}).
		Where("worker_id = ?", workerID).
		Update("last_heartbeat", time.Now()).Error
}

// UpdateWorkerStatus updates the status of a worker
func (r *WorkerRepository) UpdateWorkerStatus(workerID string, status string) error {
	return r.db.Model(&Worker{}).
		Where("worker_id = ?", workerID).
		Update("status", status).Error
}

// UpdateCurrentBatch updates the current batch being processed by a worker
func (r *WorkerRepository) UpdateCurrentBatch(workerID string, batchID *int64) error {
	return r.db.Model(&Worker{}).
		Where("worker_id = ?", workerID).
		Update("current_batch_id", batchID).Error
}

// SetStaleWorkers marks workers as stale if they have not sent a heartbeat recently
func (r *WorkerRepository) SetStaleWorkers(timeout time.Duration) error {
	staleThreshold := time.Now().Add(-timeout)

	err := r.db.Model(&Worker{}).
		Where("last_heartbeat < ?", staleThreshold).
		Update("status", "stale").Error

	if err == nil {
		log.Printf("Marked workers as stale if heartbeat was older than: %s", staleThreshold)
	}

	return err
}

// GetAvailableWorker retrieves an available worker
func (r *WorkerRepository) GetAvailableWorker() (*Worker, error) {
	var worker Worker

	err := r.db.Model(&Worker{}).
		Where("status = ?", "idle").
		Order("last_heartbeat ASC"). // Get the oldest idle worker
		First(&worker).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No available workers
		}
		return nil, err
	}

	return &worker, nil
}
