# Testing Upload file to S3
### Usage
- copy file .env.example to .env
- ``go run main.go``
- endpoint ``/s3``

### Client side
- ```curl -XPOST localhost:8080/s3 -H "Authorization: " -F "file=@path"```
- Result = ``endpoint, bucket, filename``