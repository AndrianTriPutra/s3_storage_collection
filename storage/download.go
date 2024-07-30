package storage

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
)

func (r repository) Download(ctx context.Context, path, objectName string) error {
	err := r.isOffline(ctx)
	if err != nil {
		errN := errors.New("Download:" + err.Error())
		return errN
	}

	opts := minio.GetObjectOptions{
		Checksum: true,
	}

	err = r.client.FGetObject(ctx, r.conf.BucketName, objectName, path, opts)
	if err != nil {
		errN := errors.New("FGetObject:" + err.Error())
		return errN
	}

	return nil
}

func (r repository) URLDownload(ctx context.Context, objectName string) (string, error) {
	err := r.isOffline(ctx)
	if err != nil {
		errN := errors.New("URLDownload:" + err.Error())
		return "", errN
	}

	fileName := objectName[strings.LastIndex(objectName, "/")+1 : len(objectName)]
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+fileName+"\"")

	// Generates a presigned url which expires
	presignedURL, err := r.client.PresignedGetObject(ctx, r.conf.BucketName, objectName, r.conf.Expire, reqParams)
	if err != nil {
		return "", err
	}
	url := presignedURL.String()

	return url, nil
}
