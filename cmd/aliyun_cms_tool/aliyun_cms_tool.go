package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jokimina/go-script/pkg/aliyun"
	"os"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type Datapoints []struct {
	Timestamp  int64   `json:"timestamp"`
	UserID     string  `json:"userId"`
	InstanceID string  `json:"instanceId"`
	Maximum    int     `json:"Maximum"`
	Minimum    int     `json:"Minimum"`
	Average    float64 `json:"Average"`
}

func timeUnix2Human(t int64) string {
	t, err := strconv.ParseInt(strconv.FormatInt(t, 10)[:10], 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(t, 0).Format(TimeFormat)
}

func getNames(l []string) (lName []string) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	defer db.Close()
	if err != nil {
		panic(err)
	}

	var (
		t   string
		sql string = fmt.Sprintf(`select hostname from detail_aliyun_ecs where instance_id in ("%s")`, strings.Join(l, `","`))
	)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&t)
		if err != nil {
			panic(err)
		}
		lName = append(lName, t)
	}
	fmt.Println(t)
	return
}

func dup_count(list []string) map[string]int {
	duplicate_frequency := make(map[string]int)
	for _, item := range list {
		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1
		} else {
			duplicate_frequency[item] = 1
		}
	}
	return func() (m map[string]int) {
		for k, v := range duplicate_frequency {
			if v > 1 {
				m[k] = v
			}
		}
		return
	}()
}

func main() {
	//fmt.Println(timeUnix2Human(154091520000))
	//os.Exit(0)
	config := aliyun.NewConfigFromEnv()
	client, err := cms.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := cms.CreateQueryMetricListRequest()
	request.Project = "acs_ecs_dashboard"
	request.Metric = "cpu_total"
	request.Period = strconv.Itoa(60 * 60 * 24 * 1)
	request.Length = "1000"
	request.StartTime = time.Now().AddDate(0, 0, -3).Format(TimeFormat)
	request.Dimensions = `{"regionid": "cn-beijing"}`

	var (
		datas    Datapoints
		leisures []string
	)
	for {
		response, err := client.QueryMetricList(request)
		if response.Cursor != "" {
			request.Cursor = response.Cursor
		} else {
			request.Cursor = ""
		}

		if err != nil {
			panic(err)
		}
		json.Unmarshal([]byte(response.Datapoints), &datas)

		for _, d := range datas {
			if d.Average < 10 && func() bool {
				for _, e := range leisures {
					if e == d.InstanceID {
						return false
					}
				}
				return true
			}() {
				leisures = append(leisures, d.InstanceID)
			}
		}
		if request.Cursor == "" {
			break
		}
	}
	fmt.Println(len(leisures), dup_count(leisures), leisures)
	fmt.Println(getNames(leisures))

	//b, err := json.MarshalIndent(datas, "", " ")
	//fmt.Println(string(b))
}
