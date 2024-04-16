package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
)

const (
	AccessKeyEnv    = "AWS_ACCESS_KEY"
	SecretAccessKey = "AWS_SECRET_KEY"
	Region          = "AWS_REGION"
	BucketName      = "AWS_BUCKET"
)

func FileProcessing(file multipart.File, header *multipart.FileHeader, productID int) (filename, contentType string) {
	fileNameParts := strings.Split(header.Filename, ".")
	ext := ""
	if len(fileNameParts) == 2 {
		ext = fileNameParts[1]
	}
	filename = fmt.Sprint(productID, ".", ext)
	contentType = "image/jpeg"
	if ext == "pdf" {
		contentType = "application/pdf"
	}
	return
}
