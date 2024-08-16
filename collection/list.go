package collection

import (
	"context"
	"errors"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (r repository) FileList(ctx context.Context, directory string) ([]string, error) {
	err := r.isOffline(ctx)
	if err != nil {
		errN := errors.New("FileList:" + err.Error())
		return nil, errN
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    directory,
		Recursive: true,
	}

	var nameFile []string
	objectCh := r.client.ListObjects(ctx, r.conf.BucketName, opts)
	for object := range objectCh {
		if object.Err != nil {
			errN := errors.New("ListObjects:" + object.Err.Error())
			return nil, errN
		}
		name := strings.TrimPrefix(object.Key, directory)
		nameFile = append(nameFile, name)
	}

	return nameFile, nil
}
