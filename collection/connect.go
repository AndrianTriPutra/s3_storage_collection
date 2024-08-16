package collection

import (
	"context"
	"errors"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// region
// "ap-southeast-1"=>Asia Pasifik (Singapura),
// "ap-southeast-2"=>Asia Pasifik (Sydney),
// "ap-southeast-3"=>Asia Pasifik (Jakarta)

func (r repository) Connect(ctx context.Context) error {
	creds := credentials.NewStaticV4(r.conf.AccessKeyID, r.conf.SecretAccessKey, "")
	opts := minio.Options{
		Creds:  creds,
		Secure: true,
		Region: "ap-southeast-3",
	}
	client, err := minio.New(r.conf.Endpoint, &opts)
	if err != nil {
		return err
	}
	r.client = client
	log.Println("[info] [Connect] succesed connected to MinioConnect")

	return nil
}

func (r repository) isOffline(ctx context.Context) error {
	offline := r.client.IsOffline()
	if offline {
		err := r.Connect(ctx)
		if err != nil {
			errN := errors.New("isOffline->" + err.Error())
			return errN
		}
	}
	log.Println("[info] [isOffline] minio:", "still connect")
	return nil
}
