package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	scanner_int "github.com/pgcrooks/dspm-scanner/internal"
)

func main() {
	slog.Info("starting scraper")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if viper.GetBool("aws.enabled") {
		slog.Debug("aws enabled")

		bucketName := viper.GetString("aws.bucket_name")

		client, err := newS3Client()
		if err != nil {
			slog.Error("unable to create AWS client", "err", err.Error())
		} else {
			contents, err := scanner_int.ListS3Bucket(context.TODO(), client, bucketName)
			if err != nil {
				slog.Error(err.Error())
			} else {
				slog.Info("first page results")
				for _, object := range contents {
					slog.Info("key=%s size=%d", object.Key, object.Size)
				}
			}
		}

	} else {
		slog.Debug("aws disabled")
	}

	if viper.GetBool("local.enabled") {
		slog.Info("local enabled")

		directory := viper.GetString("local.path")

		contents, err := scanner_int.ListLocalBucket(context.TODO(), directory)
		if err != nil {
			slog.Error(err.Error())
		} else {
			for _, object := range contents {
				slog.Info("found", "key", object.Key, "size", object.Size)
			}
		}
	} else {
		slog.Debug("local disabled")
	}
}

func newS3Client() (*s3.Client, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}
