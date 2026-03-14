package scanner_int

import (
	"context"
	"fmt"
	"log"
	"os"
)

type LocalAPI interface {
	ReadDir(path string) ([]os.DirEntry, error)
}

func ListLocalBucket(ctx context.Context, path string) ([]BucketObject, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can not list bucket. %v", err)
	}

	var contents []BucketObject
	for _, file := range files {
		// Ignore directories
		if !file.IsDir() {
			var fileSize int64 = 0
			fileInfo, err := os.Stat(path + file.Name())
			if err != nil {
				log.Printf("failed to get file stat for %s", file.Name())
			} else {
				fileSize = fileInfo.Size()
			}
			obj := BucketObject{
				Key:  file.Name(),
				Size: fileSize,
			}
			contents = append(contents, obj)
		}
	}
	return contents, nil
}
