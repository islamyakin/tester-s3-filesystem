package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/s3", handleS3Upload).Methods("POST")
	http.Handle("/", r)

	fmt.Println("Server Running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleS3Upload(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env file", err)
		return
	}
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_REGION")

	authHeader := r.Header.Get("Authorization")
	if authHeader != "-" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
	}

	svc := s3.New(sess)

	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(header.Filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	}
	_, err = svc.PutObject(uploadInput)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, header.Filename)
	fmt.Fprintf(w, "uploaded : %s\n", url)
}
