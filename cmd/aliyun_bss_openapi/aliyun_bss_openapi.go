package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/jokimina/go-script/pkg/aliyun"
)

func main() {
	config := aliyun.NewConfigFromEnv()
	client, err := bssopenapi.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := bssopenapi.CreateQueryAccountBalanceRequest()
	resp, err := client.QueryAccountBalance(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
