package main

import (
	"atp/storage/collection"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Println("Usage: ")
		log.Printf("      go run . methode file")
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	methode := flag.Args()[0]
	log.Printf("methode:%s", methode)

	file := ""
	if len(flag.Args()) > 1 {
		file = flag.Args()[1]
		log.Printf("file   :%s", file)
	}

	conf := collection.Cloud{
		Endpoint:        "", //api url from minio
		AccessKeyID:     "",
		SecretAccessKey: "",
		BucketName:      "atp",
		ContentType:     "application/octet-stream",
		Url:             "https://min.io/", //if you use minio
		Expire:          15 * time.Minute,  //expire url download
	}

	repo, err := collection.NewCloud(conf)
	if err != nil {
		log.Fatalf("[main] NewCloud:%s", err.Error())
	}

	ctx := context.Background()
	upload := "./example/upload/"
	download := "./example/download/"
	images := "images/"

	switch methode {
	case "list": // go run . list
		list, err := repo.FileList(ctx, images)
		if err != nil {
			log.Fatalf("[main] list:%s", err.Error())
		}
		if len(list) == 0 {
			log.Println("[main] file is empty")
		}
		for i, name := range list {
			log.Printf("[%d] name:%s", i+1, name)
		}

	case "upload":
		path := upload + file
		object := images + file
		log.Printf("path     :%s", path)
		log.Printf("object   :%s", object)

		hash, err := hasher(path)
		if err != nil {
			log.Fatalf("[main] Upload hash:%s", err.Error())
		}
		hash256 := hex.EncodeToString(hash.Sum(nil))
		log.Printf("Hash [%s]", hash256)

		info, err := repo.Upload(ctx, path, object, hash)
		if err != nil {
			log.Fatalf("[main] Upload:%s", err.Error())
		}

		buffer, _ := base64.StdEncoding.DecodeString(info.ChecksumSHA256)
		sha256 := fmt.Sprintf("%x", buffer)
		if hash256 != sha256 {
			log.Fatalf("[main] Hash [%s] is NOT match", hash256)
		}

		log.Printf("Hash [%s] is match", hash256)
		js, _ := json.MarshalIndent(info, " ", " ")
		log.Printf("info   :%s", string(js))

	case "download":
		path := download + file
		object := images + file
		log.Printf("path     :%s", path)
		log.Printf("object   :%s", object)

		err := repo.Download(ctx, path, object)
		if err != nil {
			log.Fatalf("[main] Download:%s", err.Error())
		}
		log.Printf("check file at here:%s", path)

	case "url":
		object := images + file
		url, err := repo.URLDownload(ctx, object)
		if err != nil {
			log.Fatalf("[main] URLDownload:%s", err.Error())
		}
		log.Printf("url:>>>>>>>>>>%s<<<<<<<<<", url)
	}
}

func hasher(path string) (hash.Hash, error) {
	h := sha256.New()

	f, err := os.Open(path)
	if err != nil {
		errN := errors.New("open->" + err.Error())
		return h, errN
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		errN := errors.New("copy->" + err.Error())
		return h, errN
	}

	return h, nil
}
