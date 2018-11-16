package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/jokimina/go-script/pkg/aliyun"
)

func main() {
	config := aliyun.NewConfigFromEnv()
	client, err := cs.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := cs.CreateGetClusterProjectsRequest()
	//request.SetDomain("cs.aliyuncs.com")
	//request.SetScheme(requests.HTTPS)
	request.ClusterId = "c8418affe75db4861a836a9dd26316afa"
	resp, err := client.GetClusterProjects(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
