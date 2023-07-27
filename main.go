package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/islamyakin/tester-s3-filesystem/db"
	"github.com/islamyakin/tester-s3-filesystem/server"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func intiliazeDatabase() (*gorm.DB, error) {
	database, err := db.InitDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return database, nil
}

func startServer() {

	r := mux.NewRouter()
	r.HandleFunc("/s3", server.HandleS3Upload).Methods("POST")
	r.HandleFunc("/s3/{filename}", server.HandleS3Delete).Methods("DELETE")
	http.Handle("/", r)

	serve := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Println("Server Running on port 8080")
	if err := serve.ListenAndServe(); err != nil {
		panic(err)
	}

}

func main() {
	database, err := intiliazeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if sqlDB, err := database.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	startServer()
}
