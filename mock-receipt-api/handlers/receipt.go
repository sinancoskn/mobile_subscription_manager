package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ReceiptResponse struct {
	Status     bool   `json:"status"`
	ExpireDate string `json:"expire_date"`
}

func ValidateReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody["receipt"] == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	receipt := requestBody["receipt"]

	lastChar := receipt[len(receipt)-1:]
	if _, err := strconv.Atoi(lastChar); err != nil || len(lastChar) != 1 {
		response := ReceiptResponse{
			Status:     false,
			ExpireDate: "",
		}
		writeJSONResponse(w, response)
		return
	}

	// location, _ := time.LoadLocation("America/Mexico_City")
	expireDate := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")

	response := ReceiptResponse{
		Status:     true,
		ExpireDate: expireDate,
	}
	writeJSONResponse(w, response)
}

func writeJSONResponse(w http.ResponseWriter, response ReceiptResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
