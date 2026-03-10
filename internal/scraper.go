package scanner_int

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func list_s3(client s3.Client) ([]bucket_object, error) {
	return nil, errors.New("can not list bucket")
}
