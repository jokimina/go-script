package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

func HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func downloadBucketFiles(client *oss.Client, bucketName string) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	if _, err := os.Stat(bucketName); os.IsNotExist(err) {
		os.Mkdir(bucketName, os.ModePerm)
	}

	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			HandleError(err)
		}
		for _, object := range lsRes.Objects {
			fmt.Println("Start download Bucket ",bucketName, object.Key)
			bucket.GetObjectToFile(object.Key, bucketName + "/" + object.Key)
			fmt.Println("Finished download Bucket ",bucketName, object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}

func main(){
	client, err := oss.New("oss-cn-hangzhou.aliyuncs.com", os.Getenv("ACCESS_KEY"), os.Getenv("ACCESS_SECRET"))
	if err != nil {
		HandleError(err)
	}
	marker := ""
	for {
		resp, err := client.ListBuckets(oss.Marker(marker))
		if err != nil {
			HandleError(err)
		}
		for _, bucket := range resp.Buckets {
			downloadBucketFiles(client, bucket.Name)
		}
		if resp.IsTruncated {
			marker = resp.NextMarker
		} else {
			break
		}
	}
}
