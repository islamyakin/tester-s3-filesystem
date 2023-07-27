package server

import (
	"encoding/json"
	"net/http"

	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/models"
)

func HandleS3Cek(w http.ResponseWriter, _ *http.Request) {
	database := db.GetDB()

	var files []models.File
	if err := database.Find(&files).Error; err != nil {
		http.Error(w, "Failed to get files from database", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convert data to JSON
	jsonData, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Failed to convert data to JSON", http.StatusInternalServerError)
		return
	}
	// Write JSON data to response
	w.Write(jsonData)

}
