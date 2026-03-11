package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
)

func main() {
	fmt.Println("Running dspm-scanner")

	client := newS3Client()

	contents, err := scanner_int.ListS3Bucket(context.TODO(), client, "pgcrooks-dspm")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results")
	for _, object := range contents {
		log.Printf("key=%s size=%d", object.Key, object.Size)
	}
}

func newS3Client() *s3.Client {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)
	return client
}
