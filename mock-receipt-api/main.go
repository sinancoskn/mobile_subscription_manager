package main

import (
	"fmt"
	"log"
	"mock-receipt-api/handlers"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/validate-receipt", handlers.ValidateReceipt)

	port := fmt.Sprintf(":%s", getEnv("PORT", "1234"))

	log.Println("Mock Receipt API running on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
