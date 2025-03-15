package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// KODE INI AKAN MENGUPLOAD FILE PADA DIRECTORY YANG SAMA DENGAN FILE INI
// KODE INI JUGA AKAN MENAMPILKAN LIST OBJECTS DAN BUCKETS YANG TERDAPAT PADA BUCKET YANG DITENTUKAN

func main() {
	var bucketName = "<bucketName>"           // Change to your bucket name
	var accessKeyId = "<accessKeyId>"         // Change to your access key ID
	var accessKeySecret = "<accessKeySecret>" // Change to your access key secret

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("<region>"), // Change to your region
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("<endpoint>") // Change to your endpoint
		o.UsePathStyle = false
	})

	// List objects in the bucket
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listObjectsOutput.Contents {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

	// List available buckets
	listBucketsOutput, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listBucketsOutput.Buckets {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}

	// Upload an image from the current directory
	uploadDir := "."                                               // Current directory
	files, err := filepath.Glob(filepath.Join(uploadDir, "*.jpg")) // Change to other formats if needed
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println("Uploading:", file)
		err := uploadFile(client, bucketName, file)
		if err != nil {
			log.Println("Failed to upload", file, err)
		} else {
			fmt.Println("Successfully uploaded", file)
		}
	}
}

func uploadFile(client *s3.Client, bucketName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   file,
	})
	return err
}
