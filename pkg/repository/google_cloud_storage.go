package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type GoogleStorageConnection struct {
	Client *storage.Client
}

func SetupGoogleStorageConnection() (*GoogleStorageConnection, error) {

	log.Printf("Connecting to Cloud Storage\n")
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	storage, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error while initializing Cloud Storage Client: %w", err)
	}

	return &GoogleStorageConnection{Client: storage}, nil
}

func (gsc *GoogleStorageConnection) Close() error {
	if err := gsc.Client.Close(); err != nil {
		return fmt.Errorf("Error while closing Google Storage Client: %w", err)
	}

	return nil
}

type GCImageRepository struct {
	Storage    *storage.Client
	BucketName string
}

func NewImageRepository(gcClient *storage.Client, bucketName string) *GCImageRepository {
	return &GCImageRepository{Storage: gcClient, BucketName: bucketName}
}

func (ir *GCImageRepository) UploadProfileImage(ctx context.Context, objName string, imageFile multipart.File) (string, error) {
	bucket := ir.Storage.Bucket(ir.BucketName)

	object := bucket.Object(objName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, imageFile); err != nil {
		log.Printf("Unable to write file to Google Cloud Storage: %v", err)
		return "", nil
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	bucketUrl := os.Getenv("PUBLIC_URL")
	imageUrl := fmt.Sprintf(bucketUrl, ir.BucketName, objName)

	return imageUrl, nil
}

func (ir *GCImageRepository) DeleteProfileImage(ctx context.Context, objName string) error {
	bucket := ir.Storage.Bucket(ir.BucketName)

	object := bucket.Object(objName)

	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}

///////////////////////