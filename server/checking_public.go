package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

type FileData struct {
	Name     string `json:"name"`
	FileType string `json:"file_type"`
	S3URL    string `json:"s3_url"`
}

func HandleListFilesS3(w http.ResponseWriter, _ *http.Request) {
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("AWS_BUCKET")
	region := os.Getenv("AWS_REGION")

	// Configure session AWS
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Initialize service S3
	svc := s3.New(sess)

	// Get list of files from S3
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		http.Error(w, "Failed to get list of files from S3", http.StatusInternalServerError)
		return
	}

	// Convert the S3 file list to a list of FileData
	var files []FileData
	for _, item := range resp.Contents {
		file := FileData{
			Name:     *item.Key,
			FileType: "", // You can set the file type based on your logic
			S3URL:    endpoint + "/" + bucket + "/" + *item.Key,
		}
		files = append(files, file)
	}

	// Convert files slice to JSON
	jsonData, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "Failed to convert files to JSON", http.StatusInternalServerError)
		return
	}

	// Set content type and write JSON data to response body
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
