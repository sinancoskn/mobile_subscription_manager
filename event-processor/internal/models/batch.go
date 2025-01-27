package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Batch struct {
	ID         int64      `gorm:"primaryKey" json:"id"`             // Primary key
	ActionID   int64      `gorm:"not null" json:"action_id"`        // Related action ID
	StartIndex int64      `gorm:"not null" json:"start_index"`      // Start index for the batch
	EndIndex   int64      `gorm:"not null" json:"end_index"`        // End index for the batch
	Status     string     `gorm:"default:pending" json:"status"`    // Batch status
	TryCount   int        `gorm:"not null" json:"try_count"`        // Number of processing attempts
	LockedBy   *string    `gorm:"default:null" json:"locked_by"`    // Worker that locked this batch
	LockedAt   *time.Time `gorm:"default:null" json:"locked_at"`    // When the batch was locked
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"` // Creation timestamp
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"` // Update timestamp
}

func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

type BatchRepository struct {
	db *gorm.DB
}

func (r *BatchRepository) CreateBatches(actionID int64, batchSize int64, totalRecords int64) error {
	var batches []Batch

	// Prepare batches for bulk insert
	for start := int64(1); start <= totalRecords; start += batchSize {
		end := start + batchSize - 1
		if end > totalRecords {
			end = totalRecords
		}

		batches = append(batches, Batch{
			ActionID:   actionID,
			StartIndex: start,
			EndIndex:   end,
			Status:     "pending",
		})
	}

	// Perform bulk insert
	if len(batches) > 0 {
		if err := r.db.Create(&batches).Error; err != nil {
			return err
		}
	}

	log.Printf("Inserted %d batches for action %d.\n", len(batches), actionID)
	return nil
}

func (r *BatchRepository) GetBatchesByActionIDs(actionIDs []int64) ([]Batch, error) {
	var batches []Batch
	query := r.db.Model(&Batch{})

	if len(actionIDs) > 0 {
		query = query.Where("action_id IN ?", actionIDs)
	}

	err := query.Find(&batches).Error
	if err != nil {
		return nil, err
	}

	return batches, nil
}

// LockNextBatch fetches and locks the next pending batch
func (r *BatchRepository) LockNextBatch(workerID string) (*Batch, error) {
	var batch Batch
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("status = ?", "pending").
			Order("id ASC").
			First(&batch).Error
		if err != nil {
			return err
		}

		// Lock the batch
		return tx.Model(&batch).
			Updates(map[string]interface{}{
				"status":    "processing",
				"locked_by": workerID,
				"locked_at": time.Now(),
			}).Error
	})
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// MarkBatchCompleted marks a batch as completed
func (r *BatchRepository) MarkBatchCompleted(batchID int64) error {
	return r.db.Model(&Batch{}).
		Where("id = ?", batchID).
		Updates(map[string]interface{}{
			"status":    "completed",
			"locked_by": nil,
			"locked_at": nil,
		}).Error
}

func (r *BatchRepository) MarkBatchAsStale(batchID int64) error {
	return r.db.Model(&Batch{}).
		Where("id = ?", batchID).
		Updates(map[string]interface{}{
			"status":    "stale",
			"locked_by": nil,
			"locked_at": nil,
		}).Error
}

func (r *BatchRepository) IncrementBatchTryCount(batchID int64) (*Batch, error) {
	var batch Batch
	err := r.db.Model(&Batch{}).
		Where("id = ?", batchID).
		UpdateColumn("try_count", gorm.Expr("try_count + ?", 1)).
		First(&batch).Error // Retrieve the updated batch record
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

// AreAllBatchesCompleted checks if all batches for an action are completed
func (r *BatchRepository) AreAllBatchesFinish(actionID int64) (bool, error) {
	var count int64
	err := r.db.Model(&Batch{}).
		Where("action_id = ? AND status IN ('pending', 'processing')", actionID).
		Count(&count).Error
	return count == 0, err
}

func (r *BatchRepository) GetCompletedBatchCount(actionID int64) (int64, error) {
	var count int64
	err := r.db.Model(&Batch{}).
		Where("action_id = ? AND status IN ('completed')", actionID).
		Count(&count).Error
	return count, err
}
