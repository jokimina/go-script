package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url string
	key string
)

func init() {
	flag.StringVar(&url, "url", "", "")
	flag.StringVar(&key, "key", "", "")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	resp, err := http.Post(url, "text/plain", bytes.NewBufferString(key))
	checkErr(err)
	b, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	log.Println(string(b))
}
