package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DeleteFileFromS3(accessKeyID, secretAccessKey, endpoint, bucket, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-west-2"),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(false),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	_, err = svc.DeleteObject(deleteInput)
	if err != nil {
		return err
	}

	return nil
}
