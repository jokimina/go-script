package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	TotalNum   int = 600000000
	PartNum    int = 10000
	TimeFormat     = "2006-01-02 15:04:05"
)

var (
	tmpSlice []int
	t        []int
	lrw      sync.RWMutex
	minTime  = time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC).Unix()
	maxTime  = time.Now().Unix()
	delta    = maxTime - minTime
	chNumCtrl = make(chan bool, runtime.NumCPU())
	blockLock = make(chan bool, 1)
)

func genSql(tmpSlice []int, chNumCtrl chan bool, file *os.File) {
	slice := make([]int, PartNum)
	copy(slice, tmpSlice)
	blockLock <- true
	startPos := slice[0]
	endPos := slice[len(slice)-1]
	fmt.Printf("start %d - %d...\n", startPos, endPos)

	date := time.Unix(rand.Int63n(delta)+minTime, 0).Format(TimeFormat)

	for i, v := range slice {
		if i > 10 {
			break
		}
		sql := fmt.Sprintf("'%d', '%d', '%f', '%s'\n", v, rand.Int31n(5000), float32(rand.Int31n(100)+1), date)
		fmt.Print(sql)
		lrw.Lock()
		file.Write([]byte(sql))
		file.Sync()
		lrw.Unlock()
	}
	fmt.Printf("end %d - %d...\n", startPos, endPos)
	runtime.GC()
	<- chNumCtrl
}

func main() {
	file, err := os.Create("data.sql")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 1; i <= TotalNum; i++ {
		tmpSlice = append(tmpSlice, i)
		if i%PartNum == 0 {
			chNumCtrl <- true
			go genSql(tmpSlice, chNumCtrl, file)
			<- blockLock
			fmt.Println(runtime.NumGoroutine())
			tmpSlice = tmpSlice[:0]
		}
	}
}
