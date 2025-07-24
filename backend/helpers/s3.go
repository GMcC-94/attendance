package helpers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appCfg "github.com/gmcc94/attendance-go/config"
)

var s3Client *s3.Client

func InitS3() {
	region := appCfg.AppConfig.AWSRegion

	accessKey := appCfg.AppConfig.AWSAccessKeyID
	secretKey := appCfg.AppConfig.AWSSecretAccessKey

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
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
		Bucket:      aws.String(appCfg.AppConfig.AWSBucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(header.Header.Get("Content-Type")),
		ACL:         "public-read",
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload: %v", err)
	}

	// Generate pre-signed URL for temporary access (1 hour expiration)
	req := &s3.GetObjectInput{
		Bucket: aws.String(appCfg.AppConfig.AWSBucketName),
		Key:    aws.String(key),
	}

	presignClient := s3.NewPresignClient(s3Client)

	presignedURL, err := presignClient.PresignGetObject(context.TODO(), req, s3.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.URL, nil
}
