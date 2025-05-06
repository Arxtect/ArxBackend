package xminio

import (
	"github.com/arxtect/ArxBackend/golangp/config"
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	AccessKey string
	SecretKey string
	S3URL     string
}

type S3Manager struct {
	Client     *minio.Client
	BucketName string
}

func NewS3Client(accessKey, secretKey, s3URL string) (*minio.Client, error) {
	if accessKey == "" || secretKey == "" || s3URL == "" {
		accessKey = config.Env.MinioAccessKey
		secretKey = config.Env.MinioSecretKey
		s3URL = config.Env.MinioBucketUrl
	}

	client, err := minio.New(s3URL, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: config.Env.MinioSecure,
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewS3Manager(bucketName, accessKey, secretKey, s3URL string) *S3Manager {
	client, err := NewS3Client(accessKey, secretKey, s3URL)

	if err != nil {
		log.Fatalln(err)
	}

	found, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	if !found {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Created bucket %s\n", bucketName)
	} else {
		log.Printf("Bucket %s already exists\n", bucketName)
	}

	return &S3Manager{client, bucketName}
}

func (s *S3Manager) DownloadFile(minioFileKeyList []string, localFilePath string) {
	for _, item := range minioFileKeyList {
		localFile := localFilePath + "/" + item
		if _, err := os.Stat(localFile); os.IsNotExist(err) {
			err = s.Client.FGetObject(context.Background(), s.BucketName, item, localFile, minio.GetObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("File %s downloaded successfully to %s\n", item, localFile)
		} else {
			log.Printf("File %s already exists at %s, skipping download.\n", item, localFile)
		}
	}
}

func (s *S3Manager) UploadFile(localFileList []string, localFilePath string) {
	for _, item := range localFileList {
		_, err := s.Client.FPutObject(context.Background(), s.BucketName, item, localFilePath+"/"+item, minio.PutObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("File %s uploaded successfully to %s\n", item, item)
	}
}

func (s *S3Manager) GeneratePresignedURL(objectName string, expiryTime time.Duration) string {
	presignedURL, err := s.Client.PresignedGetObject(context.Background(), s.BucketName, objectName, expiryTime, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return presignedURL.String()
}

// Generate a put presigned URL for the object
func (s *S3Manager) GeneratePresignedPutURL(objectName string, expiryTime time.Duration) string {
	preSignedURL, err := s.Client.PresignedPutObject(context.Background(), s.BucketName, objectName, expiryTime)

	if err != nil {
		log.Fatalln(err)
	}
	return preSignedURL.String()
}

func (s *S3Manager) ServeObject(w http.ResponseWriter, objectName string) {
	// Generate a presigned URL for the object
	presignedURL, err := s.Client.PresignedGetObject(context.Background(), s.BucketName, objectName, time.Hour, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make a GET request to the presigned URL
	resp, err := http.Get(presignedURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set the Content-Disposition header to inline
	w.Header().Set("Content-Disposition", "inline")

	// Copy the response body to the HTTP response
	io.Copy(w, resp.Body)
}

// DownloadObject 直接下载
func (s *S3Manager) DownloadObject(w http.ResponseWriter, objectName string) {
	// Generate a presigned URL for the object
	presignedURL, err := s.Client.PresignedGetObject(context.Background(), s.BucketName, objectName, time.Hour, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make a GET request to the presigned URL
	resp, err := http.Get(presignedURL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set the Content-Disposition header to attachment
	w.Header().Set("Content-Disposition", "attachment; filename="+objectName)

	// Copy the response body to the HTTP response
	io.Copy(w, resp.Body)
}

func (s *S3Manager) UploadByteData(data []byte, objectName string) minio.UploadInfo {
	reader := bytes.NewReader(data)
	info, err := s.Client.PutObject(context.Background(), s.BucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	return info
}

// CheckFileExistence 判断一个文件是否存在
func (s *S3Manager) CheckFileExistence(fileName string) (bool, error) {
	objectInfo, err := s.Client.StatObject(context.Background(), s.BucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	if objectInfo.Size > 0 {
		return true, nil
	}

	return false, nil
}
