package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"testing"
)

func Test_Post_file(t *testing.T) {
	// 创建 AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		log.Fatal(err)
	}

	// 创建一个 Uploader 实例
	uploader := s3manager.NewUploader(sess)

	// 打开要上传的文件
	file, err := os.Open(`C:\Users\liyy2023\Desktop\demo.txt`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 上传文件
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("chatcro-test-file"),
		Key:    aws.String("chatcro-test-file"),
		Body:   file,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully uploaded to", "yourBucketName")

}
