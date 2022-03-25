package main

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {

	for {
		s3Client, err := minio.New("10.100.3.30:9000", &minio.Options{
			Creds:  credentials.NewStaticV4("minio", "minio123", ""),
			Secure: false,
		})
		if err != nil {
			log.Fatalln(err)
		}

		buckets, err := s3Client.ListBuckets(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("buckets lens", len(buckets))
		time.Sleep(1 * time.Second)
	}
}
