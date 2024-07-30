package storage

import (
	"context"
	"hash"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Cloud struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	ContentType     string
	Url             string
	Expire          time.Duration
}

type repository struct {
	client *minio.Client
	conf   Cloud
}

func NewCloud(conf Cloud) (RepositoryI, error) {
	creds := credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, "")
	opts := minio.Options{
		Creds:  creds,
		Secure: true,
		Region: "ap-southeast-3",
	}
	client, err := minio.New(conf.Endpoint, &opts)
	if err != nil {
		return nil, err
	}
	log.Println("[info] [NewCloud] succesed connected to NewCloud")

	return repository{
		conf:   conf,
		client: client,
	}, nil
}

type RepositoryI interface {
	Connect(ctx context.Context) error
	FileList(ctx context.Context, directory string) ([]string, error)
	Upload(ctx context.Context, path, objectName string, h hash.Hash) (minio.UploadInfo, error)
	Download(ctx context.Context, path, objectName string) error
	URLDownload(ctx context.Context, objectName string) (string, error)
}
