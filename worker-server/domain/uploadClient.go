package domain

import "github.com/aws/aws-sdk-go-v2/service/s3"

type UploadClient struct {
	S3Client		*s3.Client
	BucketName  	string
	PublicBaseURL 	string
}


type UploadUsecase interface {
	UploadImage(filePath string) (string,error)
}
