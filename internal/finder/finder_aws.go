package finder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type finderAWSS3 struct {
	Finder
	Client     *s3.Client
	BucketName string
}

func newFinderAWSS3(ctx context.Context, bucketName string, bucketChan chan<- BucketObjectBatch) (IFinder, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &finderAWSS3{
		Finder: Finder{
			Name:       "aws_s3",
			BucketChan: bucketChan,
		},
		Client:     client,
		BucketName: bucketName,
	}, nil
}

func (f finderAWSS3) Run(ctx context.Context) {
	slog.Info("running aws s3 finder", "obj", f)

	run := true
	for run {
		select {
		case <-ctx.Done():
			slog.Info("stopping finder local", "obj", f)
			run = false

		default:
			contents, err := listS3Bucket(ctx, f.Client, f.BucketName)
			if err != nil {
				slog.Error(err.Error())
			} else {
				slog.Info("first page results")
				for _, object := range contents {
					slog.Info("key=%s size=%d", object.Key, object.Size)
				}
			}
		}
	}
}

type S3API interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// 	client, err := newS3Client()
// 	if err != nil {
// 		slog.Error("unable to create AWS client", "err", err.Error())
// 	} else {
// 		contents, err := finder.ListS3Bucket(
// 			context.TODO(), client, config.Scraper.Aws.BucketName,
// 		)
// 		if err != nil {
// 			slog.Error(err.Error())
// 		} else {
// 			slog.Info("first page results")
// 			for _, object := range contents {
// 				slog.Info("key=%s size=%d", object.Key, object.Size)
// 			}
// 		}
// 	}

func listS3Bucket(ctx context.Context, api S3API, bucket_name string) (BucketObjectBatch, error) {
	output, err := api.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		return nil, fmt.Errorf("can not list bucket. %v", err)
	}

	var contents BucketObjectBatch
	for _, object := range output.Contents {
		obj := BucketObject{
			Key:  aws.ToString(object.Key),
			Size: *object.Size,
		}
		contents = append(contents, obj)
	}
	return contents, nil
}
