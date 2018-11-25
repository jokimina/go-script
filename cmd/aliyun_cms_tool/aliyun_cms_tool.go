package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/jokimina/go-script/pkg/aliyun"
)

func main() {
	config := aliyun.NewConfigFromEnv()
	client, err := cms.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := cms.CreateQueryMetricListRequest()
	request.Project = "acs_ecs_dashboard"
	request.Metric = "cpu_total"
	request.Length = "100"
	response, err := client.QueryMetricList(request)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(response.Datapoints, "", "")
	fmt.Println(string(b))
	var data interface{}
	json.Unmarshal(bytes.NewBufferString(response.Datapoints).Bytes(), &data)
	fmt.Println(data.([]interface{})[0].(map[string]interface{}))
}
