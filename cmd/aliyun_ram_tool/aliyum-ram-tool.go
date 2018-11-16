package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/jokimina/go-script/pkg/aliyun"
	"os"
	"sync"
)

var (
	usernameList []string
	ackMap       map[string]string
)

func searchAccessId(accessKeyId string) {
	ramClient := aliyun.NewRamClientFromEnv()
	luRequest := ram.CreateListUsersRequest()
	luRequest.SetScheme(requests.HTTPS)
	for {
		response, err := ramClient.ListUsers(luRequest)
		if err != nil {
			panic(err)
		}
		for _, user := range response.Users.User {
			usernameList = append(usernameList, user.UserName)
		}
		if !response.IsTruncated {
			break
		} else {
			luRequest.Marker = response.Marker
		}
	}

	//usernameList = usernameList[:30]
	wg := &sync.WaitGroup{}
	c := make(chan bool, 20)
	for i, username := range usernameList {
		wg.Add(1)
		c <- true
		go func(i int, c chan bool, username string) {
			defer wg.Done()
			//fmt.Println("Get ack of username: ", username, "-- ", i)
			lkRequest := ram.CreateListAccessKeysRequest()
			lkRequest.SetScheme(requests.HTTPS)
			lkRequest.UserName = username
			resp, err := ramClient.ListAccessKeys(lkRequest)
			if err != nil {
				//fmt.Println("Error: ", username, "-- ", i)
			}
			if len(resp.AccessKeys.AccessKey) > 0 {
				//fmt.Println(resp.AccessKeys.AccessKey, username, "-- ", i)
			} else {
				//fmt.Println("no ack of ", username, "-- ", i)
			}
			for _, ack := range resp.AccessKeys.AccessKey {
				if ack.AccessKeyId == accessKeyId {
					fmt.Printf("ack [%s] -> user [%s]", accessKeyId, username)
					os.Exit(0)
				}
			}

			<-c
		}(i, c, username)
	}
	wg.Wait()

	//quit := make(chan bool)
	//c := make(chan string, 10)
	//go func() {
	//	for {
	//		select {
	//		  	case v := <- c:
	//		  		fmt.Println("Get ack of username: ", v)
	//				lkRequest := ram.CreateListAccessKeysRequest()
	//				lkRequest.SetScheme(requests.HTTPS)
	//		  		lkRequest.UserName = v
	//		  		resp, err := ramClient.ListAccessKeys(lkRequest)
	//		  		if err != nil {
	//		  			fmt.Println("Error: ", v)
	//				}
	//		  		fmt.Println(resp.AccessKeys, v)
	//			case <- time.After(5 * time.Second):
	//				fmt.Println("Timeout!")
	//		}
	//	}
	//}()

	//fmt.Println(accessKeyId)
	//<-quit
}

func main() {
	searchAccessId("xxx")
}
