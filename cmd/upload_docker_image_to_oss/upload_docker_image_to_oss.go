package main

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/docker/docker/client"
	"github.com/jokimina/go-script/pkg/aliyun"
	"os"
	"time"
)

var (
	ossEndpoint = "http://oss-ap-southeast-5.aliyuncs.com"
	ossBucketName = "xxxx-indonesia"
	ossDir = "image/"
	testImageName = "319971b89281"
	ossObjectKey = ossDir + testImageName
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	r, err := cli.ImageSave(ctx, []string{testImageName})
	defer r.Close()

	config := aliyun.NewConfigFromEnv()
	client, err := oss.New(ossEndpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	bucket, err := client.Bucket(ossBucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	options := []oss.Option{
		oss.Expires(time.Now().Add(time.Hour * 24 * 7)),
		oss.Meta("type", "image"),
	}
	fmt.Printf("%s --> %s (%s)\n", testImageName, ossEndpoint, ossObjectKey)
	err = bucket.PutObject(ossObjectKey, r, options...)
	//err = bucket.UploadFile("", "testfile", 1*1024*1024*1024, oss.Routines(10), oss.Checkpoint(true, ""))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}
