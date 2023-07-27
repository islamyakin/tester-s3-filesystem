package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/islamyakin/tester-s3-filesystem/server"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/s3", server.HandleS3Upload).Methods("POST")
	r.HandleFunc("/s3/{filename}", server.HandleS3Delete).Methods("DELETE")
	http.Handle("/", r)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Println("Server Running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
