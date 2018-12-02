package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
)

var config *Config

func init() {
	config = NewConfigFromEnv()
}

func NewClientFromEnv() (client *sdk.Client) {
	client, err := sdk.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	return
}

func NewRamClientFromEnv() (ramClient *ram.Client) {
	ramClient, err := ram.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	return
}
