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

func getLeisuresConditional(project string, metric string, f func(d Datapoint) bool) mapset.Set {
	var leisures []interface{}
	config := aliyun.NewConfigFromEnv()
	client, err := cms.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	request := cms.CreateQueryMetricListRequest()
	request.Project = project
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

func findEcsLeisures() {
	var (
		cpuLeisures, memLeisures mapset.Set
		Project = "acs_ecs_dashboard"
	)
	ch := make(chan int, 2)
	go func() {
		cpuLeisures = getLeisuresConditional(Project, "cpu_total", func(d Datapoint) bool {
			return d.Average < 10
		})
		ch <- 0
	}()
	go func() {
		memLeisures = getLeisuresConditional(Project, "memory_usedutilization", func(d Datapoint) bool {
			return d.Average < 30
		})
		ch <- 0
	}()
	for i := 0; i < cap(ch); i++ {
		<-ch
	}
	leisures := cpuLeisures.Intersect(memLeisures)
	fmt.Println(getNames(func() (l []string) {
		for v := range leisures.Iterator().C {
			l = append(l, v.(string))
		}
		return
	}()))
}

func findRdsLeisures(){
	var (
		Project = "acs_rds_dashboard"
	)
	rdsLeisures := getLeisuresConditional(Project, "IOPSUsage", func(d Datapoint) bool {
		return  d.Average < 0.1
	})
	fmt.Println(rdsLeisures)
}
// Aliyun cms api reference: https://help.aliyun.com/document_detail/28619.html?spm=a2c4g.11174283.6.670.75298f4foxUvfw#h2-url-3
func main() {
	findRdsLeisures()
	//findEcsLeisures()
}
