package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

const (
	s3Region = "ru-central1"
)

// Client клиент файлового хранилища Amazon S3.
type Client struct {
	endpoint  string
	accessID  string
	secretKey string
	client    *minio.Client
}

// NewClient creates new client for gibdd service.
func NewClient(endpoint, accessID, secretKey string) *Client {
	return &Client{
		endpoint:  endpoint,
		accessID:  accessID,
		secretKey: secretKey,
	}
}

// Auth авторизация в хранилище.
func (svc *Client) Auth() error {
	client, err := minio.New(svc.endpoint, &minio.Options{
		Secure: true,
		Region: s3Region,
		Creds:  credentials.NewStaticV2(svc.accessID, svc.secretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("s3. Auth failed: %w", err)
	}

	svc.client = client
	return nil
}

func (svc *Client) fileInfo(ctx context.Context, bucketName, fileName string) (*minio.ObjectInfo, error) {
	info, err := svc.client.StatObject(ctx, bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == 0 {
			return nil, err
		}

		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil
		}

		return nil, fmt.Errorf("s3. fileInfo: ошибка при проверке существования файла: name=%s: %w",
			fileName, minio.ToErrorResponse(err))
	}

	return &info, nil
}

// FileExists - проверка существования файла
func (svc *Client) FileExists(ctx context.Context, bucketName, fileName string) (bool, error) {
	info, err := svc.fileInfo(ctx, bucketName, fileName)
	if err != nil || info == nil {
		return false, err
	}

	return true, nil
}

// ReadFile - чтение документа из хранилища Amazon S3.
func (svc *Client) ReadFile(ctx context.Context, bucketName, fileKey string) (content []byte, err error) {
	reader, err := svc.client.GetObject(ctx, bucketName, fileKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("s3. ReadFile GetObject: %w, fileKey = %s", err, fileKey)
	}
	defer func() {
		if err2 := reader.Close(); err2 != nil {
			err = err2
		}
	}()

	content, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("s3. ReadFile ReadAll: %w, fileKey = %s", err, fileKey)
	}

	return content, nil
}

// UploadFile - загрузка файла в хранилище.
func (svc *Client) UploadFile(ctx context.Context, bucketName, fileName string, fileContent []byte) (response minio.UploadInfo, err error) {
	if svc.client == nil {
		return minio.UploadInfo{}, fmt.Errorf("s3. UploadFile. svc.client не инициализирован")
	}

	response, err = svc.client.PutObject(ctx, bucketName, fileName,
		bytes.NewBuffer(fileContent), int64(len(fileContent)),
		minio.PutObjectOptions{})
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("s3. UploadFile. Ошибка client.PutObject: %w", err)
	}

	return response, nil
}

// DeleteFile - удаление файла из хранилища.
func (svc *Client) DeleteFile(ctx context.Context, bucketName, fileName string) (err error) {
	if svc.client == nil {
		return fmt.Errorf("s3. DeleteFile. client не инициализирован")
	}

	err = svc.client.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3. DeleteFile. Ошибка %w", err)
	}

	return nil
}
