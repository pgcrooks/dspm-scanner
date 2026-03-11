package scanner_int

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func ListS3Bucket(ctx context.Context, api S3API, bucket_name string) ([]BucketObject, error) {
	output, err := api.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		return nil, fmt.Errorf("can not list bucket. %v", err)
	}

	var contents []BucketObject
	for _, object := range output.Contents {
		obj := BucketObject{
			Key:  aws.ToString(object.Key),
			Size: *object.Size,
		}
		contents = append(contents, obj)
	}
	return contents, nil
}
