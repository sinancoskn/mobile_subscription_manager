package app

import (
	"log"

	"event-processor/internal/services"
)

func StartWorkerApp() error {
	container := BuildContainer()

	// Invoke the application logic
	err := container.Invoke(func(workerService *services.WorkerService) {
		// Start the worker service
		log.Println("Starting worker service...")
		workerService.Start()
	})

	return err
}
