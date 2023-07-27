package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/islamyakin/tester-s3-filesystem/server"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/s3", HandleS3Upload).Methods("POST")
	r.HandleFunc("/s3/{filename}", HandleS3Delete).Methods("DELETE")
	http.Handle("/", r)

	fmt.Println("Server Running on port 8080")
	http.ListenAndServe(":8080", nil)

}
