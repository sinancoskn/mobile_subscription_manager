package models

import (
	"time"

	"gorm.io/gorm"
)

type ManagerAction struct {
	ID                   int64     `gorm:"primaryKey"`
	ExpectedCount        int64     `gorm:"not null"`
	WillBeProcessedCount int64     `gorm:"not null"`
	MaxBatch             int       `gorm:"not null"`
	BatchCount           int64     `gorm:"not null;default:0"` // Total number of batches
	CompletedBatchCount  int       `gorm:"not null;default:0"` // Number of completed batches
	TriggeredAt          time.Time `gorm:"not null"`
	Status               string    `gorm:"size:20;default:pending"` // "pending", "running", "completed"
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}

func NewManagerActionRepository(db *gorm.DB) *ManagerActionRepository {
	return &ManagerActionRepository{db: db}
}

type ManagerActionRepository struct {
	db *gorm.DB
}

// HasPendingAction checks if there are any pending manager actions
func (r *ManagerActionRepository) HasPendingAction() (bool, error) {
	var count int64
	err := r.db.Model(&ManagerAction{}).
		Where("status = ?", "pending").
		Count(&count).Error
	return count > 0, err
}

// CreateNewAction creates a new manager action with initial values
func (r *ManagerActionRepository) CreateNewAction(expectedCount, willBeProcessedCount int64, maxBatch int, batchCount int64) (*ManagerAction, error) {
	action := &ManagerAction{
		ExpectedCount:        expectedCount,
		WillBeProcessedCount: willBeProcessedCount,
		MaxBatch:             maxBatch,
		BatchCount:           batchCount, // Default value
		CompletedBatchCount:  0,          // Default value
		TriggeredAt:          time.Now(),
		Status:               "pending",
	}
	err := r.db.Create(action).Error
	return action, err
}

// GetActiveActions fetches all active (pending or running) manager actions
func (r *ManagerActionRepository) GetActiveActions() ([]ManagerAction, error) {
	var actions []ManagerAction
	err := r.db.Where("status IN ('pending', 'running')").Find(&actions).Error
	return actions, err
}

// MarkActionAsCompleted updates the status of a manager action to "completed"
func (r *ManagerActionRepository) MarkActionAsCompleted(actionID int64, completedBatchCount int) error {
	return r.db.Model(&ManagerAction{}).
		Where("id = ?", actionID).
		Update("completed_batch_count", completedBatchCount).
		Update("status", "completed").Error
}
