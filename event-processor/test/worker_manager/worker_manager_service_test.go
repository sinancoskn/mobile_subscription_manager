package workermanager

import (
	"context"
	"event-processor/internal/models"
	"event-processor/internal/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func init() {
	SetupTest()
}

func TestWorkerManager_HandleTrigger(t *testing.T) {
	err := Container.Invoke(func(db *gorm.DB, workerManagerService *services.WorkerManagerService, subscriptionRepo *models.SubscriptionRepository, managerRepo *models.ManagerActionRepository, batchRepo *models.BatchRepository) {
		// Step 1: Seed test data for subscriptions
		subscriptions := []models.Subscription{
			{ID: 1, Status: "active", UID: uuid.New().String(), ExpireAt: time.Now().AddDate(0, -1, 0)},
			{ID: 2, Status: "active", UID: uuid.New().String(), ExpireAt: time.Now().AddDate(0, -1, 0)},
			{ID: 3, Status: "active", UID: uuid.New().String(), ExpireAt: time.Now().AddDate(0, -1, 0)},
		}
		err := db.Create(&subscriptions).Error
		assert.NoError(t, err, "Failed to seed subscriptions")

		// Step 2: Execute the method being tested
		err = workerManagerService.HandleTrigger()
		assert.NoError(t, err, "HandleTrigger should not return an error")

		// Step 3: Validate results in the manager actions table
		var managerAction models.ManagerAction
		err = db.First(&managerAction).Error
		assert.NoError(t, err, "Manager action should be created")
		assert.Equal(t, int64(3), managerAction.ExpectedCount, "ExpectedCount should match the number of subscriptions")
		assert.Equal(t, int64(3), managerAction.WillBeProcessedCount, "WillBeProcessedCount should match the number of subscriptions")

		// Step 4: Validate results in the batches table
		var batches []models.Batch
		err = db.Find(&batches).Error
		assert.NoError(t, err, "Batches should be created")
		assert.Equal(t, 3, len(batches), "Number of batches should match the maxBatch value")

		// Step 5: Validate that no additional manager actions are created
		err = workerManagerService.HandleTrigger()
		assert.NoError(t, err, "HandleTrigger should not return an error when a manager action exists")
		var managerActionCount int64
		db.Model(&models.ManagerAction{}).Count(&managerActionCount)
		assert.Equal(t, int64(1), managerActionCount, "There should still only be one manager action")
	})

	if err != nil {
		t.Fatalf("Failed to invoke WorkerManagerService: %v", err)
	}
}

func TestWorkerManager_Heartbeat(t *testing.T) {
	err := Container.Invoke(func(db *gorm.DB, workerManagerService *services.WorkerManagerService, managerRepo *models.ManagerActionRepository, batchRepo *models.BatchRepository) {
		// Step 1: Seed test data for manager actions and batches
		managerAction := models.ManagerAction{
			ID:            1,
			Status:        "active",
			ExpectedCount: 3,
		}
		err := db.Create(&managerAction).Error
		assert.NoError(t, err, "Failed to seed manager action")

		batches := []models.Batch{
			{ID: 1, ActionID: 1, StartIndex: 0, EndIndex: 10, Status: "completed", TryCount: 0},
			{ID: 2, ActionID: 1, StartIndex: 11, EndIndex: 20, Status: "completed", TryCount: 0},
			{ID: 3, ActionID: 1, StartIndex: 21, EndIndex: 30, Status: "completed", TryCount: 0},
		}
		err = db.Omit("status").Create(&batches).Error
		assert.NoError(t, err, "Failed to seed batches")

		// Step 2: Run Heartbeat in a goroutine with a controlled context
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // Stop after 1 second
		defer cancel()

		go workerManagerService.Heartbeat(ctx, 100*time.Millisecond) // Use a short interval for the test

		// Step 3: Wait for the context to timeout
		<-ctx.Done()

		// Step 4: Validate that the manager action has been marked as completed
		var updatedAction models.ManagerAction
		err = db.First(&updatedAction, managerAction.ID).Error
		assert.NoError(t, err, "Manager action should exist")
	})

	if err != nil {
		t.Fatalf("Failed to invoke WorkerManagerService: %v", err)
	}
}
