package main

import (
	"event-processor/internal/app"
	"log"
)

func main() {
	err := app.StartWorkerManagerApp()
	if err != nil {
		log.Fatalf("Failed to start callback app: %v", err)
	}
}
