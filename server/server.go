package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/islamyakin/tester-s3-filesystem/middleware"
	"net/http"
	"time"
)

func StartServer() {

	r := mux.NewRouter()
	r.HandleFunc("/s3", HandleS3Upload).Methods("POST")
	r.HandleFunc("/s3/{filename}", HandleS3Delete).Methods("DELETE")
	r.HandleFunc("/s3/cek/local", middleware.BasicAuthMiddleware(HandleS3Cek)).Methods("GET", "OPTIONS")
	r.HandleFunc("/s3/cek/s3", HandleListFilesS3).Methods("GET")
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		http.ServeFile(w, request, "index.html")
	})

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
