package helpers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var bucketName = aws.String(os.Getenv("AWS_BUCKET_NAME"))

var s3Client *s3.Client

func InitS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func UploadToS3(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	key := fmt.Sprintf("%d%s", time.Now().Local().Unix(), filepath.Ext(header.Filename))

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      bucketName,
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(header.Header.Get("Content-Type")),
		ACL:         "public-read",
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
	return url, nil
}
