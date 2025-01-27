package app

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"event-processor/internal/config"
	"event-processor/internal/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func StartWorkerManagerApp() error {
	container := BuildContainer()

	container.Invoke(func(workerManagerService *services.WorkerManagerService) {
		ctx := context.Background()
		go workerManagerService.Heartbeat(ctx, 5*time.Second)
	})

	return container.Invoke(listenHttp)
}

func listenHttp(config *config.Config, service *services.WorkerManagerService) {
	// Dynamically render the HTML file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./static/index.html")
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}

		// Inject the WebSocket URL dynamically
		data := struct {
			WebSocketURL string
		}{
			WebSocketURL: fmt.Sprintf("ws://%s:%d/ws", config.ManagerHost, config.ManagerPort),
		}

		tmpl.Execute(w, data)
	})

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket Upgrade Error:", err)
			return
		}
		defer conn.Close()

		for {
			data := map[string]interface{}{}
			actions, err := service.GetActiveActions()
			if err == nil {
				data["actions"] = actions
			}

			workers, err := service.GetActiveWorkers()
			if err == nil {
				data["workers"] = workers
			}

			err = conn.WriteJSON(data)
			if err != nil {
				log.Println("WebSocket Write Error:", err)
				break
			}

			// Publish data every 5 seconds
			time.Sleep(5 * time.Second)
		}
	})

	// Trigger endpoint
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

	// Start the server
	port := fmt.Sprintf(":%d", config.ManagerPort)
	log.Println("Worker Manager is running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
