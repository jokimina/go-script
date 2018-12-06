package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	TotalNum   int = 600000000
	PartNum    int = 100000
	TimeFormat     = "2006-01-02 15:04:05"
	FileName       = "data.sql"
)

var (
	wg            sync.WaitGroup
	tmpSlice      []int
	lrw           sync.RWMutex
	minTime       = time.Date(2018, 11, 1, 0, 0, 0, 0, time.UTC).Unix()
	maxTime       = time.Now().Unix()
	delta         = maxTime - minTime
	blockLock     = make(chan bool, 1)
	maxRoutineNum = runtime.NumCPU()
	chNumCtrl     chan bool
)

func init() {
	flag.IntVar(&maxRoutineNum, "n", runtime.NumCPU(), "goroutine num")
	flag.Parse()
	chNumCtrl = make(chan bool, maxRoutineNum)
	runtime.GOMAXPROCS(maxRoutineNum)
}

func genSql(tmpSlice []int, chNumCtrl chan bool, file *os.File) {
	wg.Add(1)
	slice := make([]int, PartNum)
	copy(slice, tmpSlice)
	blockLock <- true
	startPos := slice[0]
	endPos := slice[len(slice)-1]
	fmt.Printf("start %d - %d...\n", startPos, endPos)

	date := time.Unix(rand.Int63n(delta)+minTime, 0).Format(TimeFormat)
	bf := bytes.NewBuffer([]byte{})
	for _, v := range slice {
		sql := fmt.Sprintf("'%d', '%d', '%f', '%s'\n", v, rand.Int31n(5000), float32(rand.Int31n(100)+1), date)
		//fmt.Print(sql)
		bf.WriteString(sql)
	}
	lrw.Lock()
	file.Write(bf.Bytes())
	file.Sync()
	lrw.Unlock()
	<-chNumCtrl
	fmt.Printf("end %d - %d...\n", startPos, endPos)
	wg.Done()
}

func main() {
	file, err := os.Create(FileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 1; i <= TotalNum; i++ {
		tmpSlice = append(tmpSlice, i)
		if i%PartNum == 0 {
			chNumCtrl <- true
			go genSql(tmpSlice, chNumCtrl, file)
			<-blockLock
			fmt.Println(runtime.NumGoroutine())
			tmpSlice = tmpSlice[:0]
		}
	}
	wg.Wait()
}
