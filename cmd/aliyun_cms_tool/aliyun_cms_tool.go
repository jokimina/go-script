package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jokimina/go-script/pkg/aliyun"
	"os"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type Datapoints []Datapoint
type Datapoint struct {
	Timestamp  int64   `json:"timestamp"`
	UserID     string  `json:"userId"`
	InstanceID string  `json:"instanceId"`
	Maximum    int     `json:"Maximum"`
	Minimum    int     `json:"Minimum"`
	Average    float64 `json:"Average"`
}

func getNames(l []string) (lName []string) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	defer db.Close()
	if err != nil {
		panic(err)
	}

	var (
		t   string
		sql string = fmt.Sprintf(`select hostname from detail_aliyun_ecs where instance_id in ("%s") and hostname not like "%%proxy%%"`, strings.Join(l, `","`))
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
	return
}

func getLeisuresConditional(metric string, f func(d Datapoint) bool) mapset.Set {
	var leisures []interface{}
	config := aliyun.NewConfigFromEnv()
	client, err := cms.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := cms.CreateQueryMetricListRequest()
	request.Project = "acs_ecs_dashboard"
	//request.Metric = "cpu_total"
	//request.Metric = "memory_usedutilization"
	request.Metric = metric
	request.Period = strconv.Itoa(60 * 60 * 24 * 1)
	request.Length = "100"
	request.StartTime = time.Now().AddDate(0, 0, -7).Format(TimeFormat)
	request.Dimensions = `{"regionid": "cn-beijing"}`

	var datas Datapoints
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
		fmt.Printf("%s: cursor [%s]\n", metric, response.Cursor)
		json.Unmarshal([]byte(response.Datapoints), &datas)

		for _, d := range datas {
			if f(d) && func() bool {
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
	return mapset.NewSetFromSlice(leisures)
}

func main() {
	cpuLeisures := getLeisuresConditional("cpu_total", func(d Datapoint) bool {
		return d.Average < 10
	})
	memLeisures := getLeisuresConditional("memory_usedutilization", func(d Datapoint) bool {
		return d.Average < 30
	})
	leisures := cpuLeisures.Intersect(memLeisures)
	fmt.Println(getNames(func() (l []string) {
		for v := range leisures.Iterator().C {
			l = append(l, v.(string))
		}
		return
	}()))
	//fmt.Println(len(leisures), dup_count(leisures), leisures)
	//b, err := json.MarshalIndent(getNames(leisures), "", " ")
	//fmt.Println(string(b))
}
