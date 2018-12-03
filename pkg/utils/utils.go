package utils

import (
	"strconv"
	"time"
)

func TimeUnix2Human(t int64) string {
	t, err := strconv.ParseInt(strconv.FormatInt(t, 10)[:10], 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
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
