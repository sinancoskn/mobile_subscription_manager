package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"event-processor/internal/config"
	"event-processor/internal/services"
)

func StartWorkerManagerApp() error {
	container := BuildContainer()

	container.Invoke(func(workerManagerService *services.WorkerManagerService) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go workerManagerService.Heartbeat(ctx, 60*time.Second)
	})

	return container.Invoke(listenHttp)
}

func listenHttp(config *config.Config, service *services.WorkerManagerService) {
	http.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Invalid request method",
				"message": "Only POST method is allowed",
			})
			return
		}

		err := service.HandleTrigger()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Failed to handle trigger",
				"message": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Worker Manager triggered successfully.",
		})
	})

	port := fmt.Sprintf(":%d", config.Port)
	log.Println("Worker Manager is running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
