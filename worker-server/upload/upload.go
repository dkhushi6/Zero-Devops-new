package upload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"Zero_Devops/worker_server/domain"
)

type clientUsecase struct {
	uploadClient *domain.UploadClient
}

func NewUploadUsecase(client *s3.Client, bucketName string, publicBaseURL string) domain.UploadUsecase {

	return &clientUsecase{
		uploadClient: &domain.UploadClient {
			S3Client: client, 
			BucketName: bucketName, 
			PublicBaseURL: publicBaseURL,
		},
	}
}

func (c *clientUsecase) UploadImage(filePath string) (string, error) {

	s3Client := c.uploadClient;

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	key := "images/" + filename

	_, err = s3Client.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s3Client.BucketName,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	if s3Client.PublicBaseURL == "" {
		return fmt.Sprintf("s3://%s/%s", s3Client.BucketName, key), nil
	}

	return fmt.Sprintf("%s/%s", s3Client.PublicBaseURL, key), nil
}
