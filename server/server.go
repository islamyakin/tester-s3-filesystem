package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func StartServer() {

	r := mux.NewRouter()
	r.HandleFunc("/s3", HandleS3Upload).Methods("POST")
	r.HandleFunc("/s3/{filename}", HandleS3Delete).Methods("DELETE")
	r.HandleFunc("/s3/cek/local", HandleS3Cek).Methods("GET")
	r.HandleFunc("/s3/cek/s3", HandleListFilesS3).Methods("GET")
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
