package scanner_int

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

type mockListObjectAPI func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)

func (m mockListObjectAPI) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m(ctx, params, optFns...)
}

func TestListS3ucket(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) S3API
		bucket string
		expect BucketObjectBatch
	}{
		{
			client: func(t *testing.T) S3API {
				return mockListObjectAPI(func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatal("expect bucket to not be nil")
					}
					if e, a := "test-bucket", *params.Bucket; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					var testContents []types.Object
					oneKey := "one"
					var oneSize int64 = 42
					testContents = append(testContents, types.Object{
						Key:  &oneKey,
						Size: &oneSize,
					})

					return &s3.ListObjectsV2Output{
						Contents: testContents,
					}, nil
				})
			},
			bucket: "test-bucket",
			expect: BucketObjectBatch{{
				Key:  "one",
				Size: 42,
			}},
		},
		{
			client: func(t *testing.T) S3API {
				return mockListObjectAPI(func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
					t.Helper()
					if params.Bucket == nil {
						t.Fatal("expect bucket to not be nil")
					}
					if e, a := "test-bucket-2", *params.Bucket; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}

					var testContents []types.Object
					oneKey := "one"
					var oneSize int64 = 42
					twoKey := "two"
					var twoSize int64 = 100
					testContents = append(testContents, types.Object{
						Key:  &oneKey,
						Size: &oneSize,
					})
					testContents = append(testContents, types.Object{
						Key:  &twoKey,
						Size: &twoSize,
					})

					return &s3.ListObjectsV2Output{
						Contents: testContents,
					}, nil
				})
			},
			bucket: "test-bucket-2",
			expect: BucketObjectBatch{
				{
					Key:  "one",
					Size: 42,
				},
				{
					Key:  "two",
					Size: 100,
				}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := ListS3Bucket(ctx, tt.client(t), tt.bucket)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			assert.ElementsMatch(t, content, tt.expect)
		})
	}
}
