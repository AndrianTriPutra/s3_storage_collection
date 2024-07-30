package storage

import (
	"context"
	"encoding/base64"
	"errors"
	"hash"

	"github.com/minio/minio-go/v7"
)

func (r repository) Upload(ctx context.Context, path, objectName string, h hash.Hash) (minio.UploadInfo, error) {
	var info minio.UploadInfo

	err := r.isOffline(ctx)
	if err != nil {
		errN := errors.New("FPutObject:" + err.Error())
		return info, errN
	}

	// Calculate checksum
	b4se64 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	meta := map[string]string{"x-amz-checksum-sha256": b4se64}
	opts := minio.PutObjectOptions{
		UserMetadata: meta,
		ContentType:  r.conf.ContentType,
	}

	info, err = r.client.FPutObject(ctx, r.conf.BucketName, objectName, path, opts)
	if err != nil {
		errN := errors.New("FPutObject:" + err.Error())
		return info, errN
	}

	return info, nil
}
