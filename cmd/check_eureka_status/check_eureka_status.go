package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	. "github.com/jokimina/go-script/pkg/alarm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Applications struct {
	Application []struct {
		Instance []struct {
			ActionType     string `json:"actionType"`
			App            string `json:"app"`
			CountryID      int    `json:"countryId"`
			DataCenterInfo struct {
				Class string `json:"@class"`
				Name  string `json:"name"`
			} `json:"dataCenterInfo"`
			HealthCheckURL                string `json:"healthCheckUrl"`
			HomePageURL                   string `json:"homePageUrl"`
			HostName                      string `json:"hostName"`
			InstanceID                    string `json:"instanceId"`
			IPAddr                        string `json:"ipAddr"`
			IsCoordinatingDiscoveryServer string `json:"isCoordinatingDiscoveryServer"`
			LastDirtyTimestamp            string `json:"lastDirtyTimestamp"`
			LastUpdatedTimestamp          string `json:"lastUpdatedTimestamp"`
			LeaseInfo                     struct {
				DurationInSecs        int   `json:"durationInSecs"`
				EvictionTimestamp     int   `json:"evictionTimestamp"`
				LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
				RegistrationTimestamp int64 `json:"registrationTimestamp"`
				RenewalIntervalInSecs int   `json:"renewalIntervalInSecs"`
				ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
			} `json:"leaseInfo"`
			Metadata struct {
				ManagementPort string `json:"management.port"`
			} `json:"metadata"`
			OverriddenStatus string `json:"overriddenStatus"`
			Port             struct {
				Port    int    `json:"$"`
				Enabled string `json:"@enabled"`
			} `json:"port"`
			SecurePort struct {
				Port    int    `json:"$"`
				Enabled string `json:"@enabled"`
			} `json:"securePort"`
			SecureVipAddress string `json:"secureVipAddress"`
			Status           string `json:"status"`
			StatusPageURL    string `json:"statusPageUrl"`
			VipAddress       string `json:"vipAddress"`
		} `json:"instance"`
		Name string `json:"name"`
	} `json:"application"`
	AppsHashcode  interface{} `json:"apps__hashcode"`
	VersionsDelta interface{} `json:"versions__delta"`
}

type arrayFlags []string

func (a *arrayFlags) String() string {
	return fmt.Sprintf("%s", *a)
}

func (a *arrayFlags) Set(value string) error {
	if len(*a) > 0 {
		return errors.New("interval flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		*a = append(*a, dt)
	}
	return nil
}

var (
	EurekaUrl   string
	DingDingUrl string
	Receivers   arrayFlags
	Comment     string
	atLeastNum  int
)

func init() {
	flag.StringVar(&EurekaUrl, "eUrl", "", "eureka url")
	flag.StringVar(&DingDingUrl, "dUrl", "", "dingding url")
	flag.Var(&Receivers, "r", "receivers")
	flag.StringVar(&Comment, "m", "", "comment")
	flag.IntVar(&atLeastNum, "n", 0, "At latest service number")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stdout, "Usage %s -eUrl xxx.eureka.xxx.com -dUrl https://oapi.dingtalk.com/robot/send?access_token=xxxx -r 13611111111 -d 华北集群\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if len(Receivers) == 0 || len(EurekaUrl) == 0 || len(DingDingUrl) == 0 {
		flag.Usage()
	}
	fmt.Println(EurekaUrl, DingDingUrl, Receivers, atLeastNum)
	url := EurekaUrl + "/eureka/apps"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Accept", "application/json")
	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var apps map[string]Applications
	json.Unmarshal([]byte(body), &apps)
	//r, _ := json.MarshalIndent(apps,""," ")
	//fmt.Println(string(r))
	var exceptList []string

	applications := apps["applications"].Application
	if len(applications) < atLeastNum {
		exceptList = append(exceptList, fmt.Sprintf("- 服务注册数量异常 应注册: %d, 目前: %d", atLeastNum, len(applications)))
	}
	for _, appValMap := range applications {
		for _, instance := range appValMap.Instance{
			if "UP" != instance.Status {
				exceptList = append(exceptList, fmt.Sprintf("- %s %s %s\n", instance.App, instance.HostName, instance.Status))
			}
		}
	}
	if len(exceptList) > 0 {
		d := DingDing{
			Url: DingDingUrl,
		}
		d.Report(Receivers,
			fmt.Sprintf("[%s eureka注册异常](%s)\n%s", Comment, EurekaUrl, strings.Join(exceptList, "")))
	}
}
