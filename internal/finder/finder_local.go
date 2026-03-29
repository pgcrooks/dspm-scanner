package finder

import (
	"context"
	"fmt"
	"log"
	"os"
)

func listLocalBucket(ctx context.Context, path string) (BucketObjectBatch, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can not list bucket. %v", err)
	}

	var contents BucketObjectBatch
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
