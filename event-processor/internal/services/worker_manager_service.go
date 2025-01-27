package services

import (
	"context"
	"event-processor/internal/models"
	"fmt"
	"log"
	"time"
)

type WorkerManagerService struct {
	managerRepo         *models.ManagerActionRepository
	batchRepo           *models.BatchRepository
	subscriptionRepo    *models.SubscriptionRepository
	workerRepo          *models.WorkerRepository
	maxProcessableCount int64
	maxBatch            int
}

func NewWorkerManagerService(
	managerRepo *models.ManagerActionRepository,
	batchRepo *models.BatchRepository,
	subscriptionRepo *models.SubscriptionRepository,
	workerRepo *models.WorkerRepository,
) *WorkerManagerService {
	return &WorkerManagerService{
		managerRepo:         managerRepo,
		batchRepo:           batchRepo,
		subscriptionRepo:    subscriptionRepo,
		workerRepo:          workerRepo,
		maxProcessableCount: 1000000,
		maxBatch:            100,
	}
}

func (s *WorkerManagerService) GetActiveWorkers() ([]models.Worker, error) {
	return s.workerRepo.GetActiveWorkers()
}

func (s *WorkerManagerService) GetActiveActions() ([]models.ManagerAction, error) {
	// Fetch active actions
	actions, err := s.managerRepo.GetActiveActions()
	if err != nil {
		return nil, err
	}

	actionIDs := make([]int64, len(actions))
	for i, action := range actions {
		actionIDs[i] = action.ID
	}

	batches, err := s.batchRepo.GetBatchesByActionIDs(actionIDs)
	if err != nil {
		return nil, err
	}

	batchesByActionID := make(map[int64][]models.Batch)
	for _, batch := range batches {
		batchesByActionID[batch.ActionID] = append(batchesByActionID[batch.ActionID], batch)
	}

	for i := range actions {
		actions[i].Batches = batchesByActionID[actions[i].ID]
	}

	return actions, nil
}

// Heartbeat checks the status of active manager actions periodically
func (s *WorkerManagerService) Heartbeat(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Worker Manager Heartbeat: Checking actions and batches...")

			// Fetch all active manager actions
			actions, err := s.managerRepo.GetActiveActions()
			if err != nil {
				log.Printf("Failed to fetch active manager actions: %v\n", err)
				continue
			}

			for _, action := range actions {
				// Check if all batches for this action are completed
				isCompleted, err := s.batchRepo.AreAllBatchesFinish(action.ID)
				if err != nil {
					log.Printf("Failed to check batch statuses for action %d: %v\n", action.ID, err)
					continue
				}

				if isCompleted {
					completedBatchCount, err := s.batchRepo.GetCompletedBatchCount(action.ID)
					if err != nil {
						log.Printf("Failed to check completed batch count for action %d: %v\n", action.ID, err)
						continue
					}

					log.Printf("All batches for action %d are completed. Marking action as completed.\n", action.ID)
					err = s.managerRepo.MarkActionAsCompleted(action.ID, int(completedBatchCount))
					if err != nil {
						log.Printf("Failed to mark action %d as completed: %v\n", action.ID, err)
					}
				}
			}
		case <-ctx.Done():
			log.Println("Heartbeat stopped.")
			return
		}
	}
}

func (s *WorkerManagerService) HandleTrigger() error {
	// Check if there's a pending manager action
	actions, err := s.managerRepo.GetActiveActions()
	if err != nil {
		return fmt.Errorf("failed to check for pending actions: %w", err)
	}
	if len(actions) > 0 {
		log.Println("A pending manager action already exists. Skipping new action creation.")
		return nil
	}

	// Step 1: Calculate expected_count from subscriptions table
	expectedCount, err := s.subscriptionRepo.GetCountForProcessing()
	if err != nil {
		return fmt.Errorf("failed to calculate expected count: %w", err)
	}
	log.Printf("Expected count of subscriptions to process: %d\n", expectedCount)

	// Step 2: Limit will_be_processed_count to a static maximum
	willBeProcessedCount := expectedCount
	if expectedCount > s.maxProcessableCount {
		willBeProcessedCount = s.maxProcessableCount
	}
	log.Printf("Will process a maximum of %d subscriptions.\n", willBeProcessedCount)

	// Step 4: Calculate batch size and create batches
	batchSize := willBeProcessedCount / int64(s.maxBatch)
	if willBeProcessedCount%int64(s.maxBatch) != 0 {
		batchSize++ // Adjust for any remainder
	}

	// Step 3: Create a new manager action
	action, err := s.managerRepo.CreateNewAction(expectedCount, willBeProcessedCount, s.maxBatch, batchSize)
	if err != nil {
		return fmt.Errorf("failed to create new manager action: %w", err)
	}

	err = s.batchRepo.CreateBatches(action.ID, batchSize, willBeProcessedCount)
	if err != nil {
		return fmt.Errorf("failed to create batches: %w", err)
	}

	log.Printf("Created new manager action %d with %d batches.\n", action.ID, s.maxBatch)
	return nil
}
