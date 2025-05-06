package handle

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"testing"
	"time"
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

func NewS3Manager(bucketName, accessKey, secretKey, s3URL string) *S3Manager {
	if accessKey == "" || secretKey == "" || s3URL == "" {
		accessKey = "jancsitech"
		secretKey = "jancsitech"
		s3URL = "dev.nas.corp.jancsitech.net:9000"
	}

	config := Config{accessKey, secretKey, s3URL}

	client, err := minio.New(config.S3URL, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: false,
	})

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

func (s *S3Manager) UploadByteData(data []byte, objectName string) {
	reader := bytes.NewReader(data)
	info, err := s.Client.PutObject(context.Background(), s.BucketName, objectName, reader, int64(reader.Len()), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Byte data uploaded successfully to %s\n", objectName)
	fmt.Println(info)

}

func Test_PostFile(t *testing.T) {
	manager := NewS3Manager("chatcro-test-file", "", "", "")
	manager.UploadFile([]string{"myfile.txt"}, "E://demo")

	//manager.DownloadFile([]string{"myfile.txt"}, "/path/to/download/directory")

	// 获取预下载的URL
	//url := manager.GeneratePresignedURL("myfile.txt", 24*time.Hour)
	//log.Println(url)
}

func Test_byte_data(t *testing.T) {
	manager := NewS3Manager("chatcro-test-file", "", "", "")
	data := []byte("Hello, world!")
	manager.UploadByteData(data, "my.txt")

	url := manager.GeneratePresignedURL("my.txt", 24*time.Hour)
	log.Println(url)

	log.Println("done")
}
