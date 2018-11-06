package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var (
	loginUrl = "https://signin.aliyun.com/login.htm"
	client http.Client
)

func main(){
	jar, err := cookiejar.New(nil)
	client.Jar = jar
	// 跳转到登录页面
	_, err = client.Get(loginUrl)
	if err != nil {
		panic(err)
	}

	params := url.Values{
		//"sec_token": {""},
		"callback": {"https://home.console.aliyun.com/"},
		"user_principal_name": {os.Getenv("USERNAME")},
		"password_ims": {os.Getenv("PASSWORD")},
	}
	resp, err := client.PostForm(loginUrl, params)
	if err != nil {
		panic(err)
	}
	// 直接获取接口信息
	resp, err = client.Get("https://usercenter.aliyun.com/?defaultMonth=2018-11#/consumptionByMonth?_k=btb6si")
	resp, err = client.Get("https://usercenter.aliyun.com/ajax/postchargeConsumptionSummaryAjax/getUserInvoiceAmount.json?billCycle=201811")
	resp, err = client.Get("https://usercenter.aliyun.com/ajax/postchargeConsumptionSummaryAjax/getUserProductConsumeInfo.json?billCycle=2018-11")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}
