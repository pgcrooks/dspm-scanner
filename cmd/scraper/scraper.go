package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
)

func main() {
	log.Println("starting scraper")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if viper.GetBool("aws.enabled") {
		log.Println("aws enabled")

		bucketName := viper.GetString("aws.bucket_name")

		client := newS3Client()

		contents, err := scanner_int.ListS3Bucket(context.TODO(), client, bucketName)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("first page results")
		for _, object := range contents {
			log.Printf("key=%s size=%d", object.Key, object.Size)
		}
	} else {
		log.Println("aws disabled")
	}

	if viper.GetBool("local.enabled") {
		log.Println("local enabled")

		directory := viper.GetString("local.path")

		log.Printf("reading %s", directory)

		contents, err := scanner_int.ListLocalBucket(context.TODO(), directory)
		if err != nil {
			log.Fatal(err)
		}

		for _, object := range contents {
			log.Printf("key=%s size=%d", object.Key, object.Size)
		}
	} else {
		log.Println("local disabled")
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
