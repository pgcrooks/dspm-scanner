package scanner_int

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3BucketContents(client s3.Client, bucket_name string, contents *[]BucketObject) error {
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		return fmt.Errorf("can not list bucket. %s", err)
	}

	for _, object := range output.Contents {
		obj := BucketObject{
			Key:  aws.ToString(object.Key),
			Size: *object.Size,
		}
		*contents = append(*contents, obj)
	}
	return nil
}
