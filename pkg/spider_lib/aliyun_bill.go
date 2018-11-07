package spider_lib

import (
	"encoding/json"
	"github.com/henrylee2cn/pholcus/app/downloader/request"
	. "github.com/henrylee2cn/pholcus/app/spider"
	"io/ioutil"
	"net/url"
	"os"
)

type CostOverviewResponse struct {
	Code string `json:"code"`
	Data []struct {
		BeforeDiscountTotalMoney   string `json:"beforeDiscountTotalMoney"`
		CardMoney                  string `json:"cardMoney"`
		CashMoney                  string `json:"cashMoney"`
		CashMoneyNoCreditPayAmount string `json:"cashMoneyNoCreditPayAmount"`
		CommodityCode              string `json:"commodityCode,omitempty"`
		CouponMoney                string `json:"couponMoney"`
		CreditPayAmount            string `json:"creditPayAmount"`
		ItemInfo                   string `json:"itemInfo,omitempty"`
		MonthProductConsume        bool   `json:"monthProductConsume"`
		OweMoney                   string `json:"oweMoney"`
		ParentCommodityCode        string `json:"parentCommodityCode,omitempty"`
		ParentGroup                bool   `json:"parentGroup"`
		RefundMoney                string `json:"refundMoney"`
		TotalDiscountMoney         string `json:"totalDiscountMoney"`
		TotalMoney                 string `json:"totalMoney"`
		YouhuiquanMoney            string `json:"youhuiquanMoney"`
		ChargeType                 string `json:"chargeType,omitempty"`
	} `json:"data"`
	Success    bool `json:"success"`
	TotalCount int  `json:"totalCount"`
}

func init() {
	AliyunBill.Register()
}

var AliyunBill = &Spider{
	Name:            "阿里云账单",
	Description:     "阿里云费用中心账单",
	EnableCookie:    true,
	NotDefaultField: false,
	Namespace:       func(self *Spider) string { return "aliyun_bill" },
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "https://signin.aliyun.com/login.htm", Rule: "登录页"})
		},
		Trunk: map[string]*Rule{
			"登录页": {
				ParseFunc: func(ctx *Context) {
					//NewForm(
					//	ctx,
					//	"费用总览",
					//	"https://signin.aliyun.com/xxx/login.htm?callback=https://usercenter.aliyun.com/?defaultMonth=2018-11#/consumptionByMonth?_k=l6if7b",
					//	ctx.GetDom().Find(".login-form"),
					//).Inputs(map[string]string{
					//	"user_principal_name": os.Getenv("USERNAME"),
					//	"password_ims": os.Getenv("PASSWORD"),
					//}).Submit()

					ctx.AddQueue(&request.Request{
						Url:    "https://signin.aliyun.com/login.htm",
						Rule:   "登陆后",
						Method: "POST",
						PostData: url.Values{
							//"sec_token": {""},
							"callback":            {"https://home.console.aliyun.com/"},
							"user_principal_name": {os.Getenv("USERNAME")},
							"password_ims": {os.Getenv("PASSWORD")},
						}.Encode(),
					})
				},
			},
			"登陆后": {
				ParseFunc: func(ctx *Context) {
					ctx.AddQueue(&request.Request{
						//Url: "https://usercenter.aliyun.com/?defaultMonth=2018-11#/consumptionByMonth",
						Url:  "https://usercenter.aliyun.com/ajax/postchargeConsumptionSummaryAjax/getUserProductConsumeInfo.json?billCycle=2018-11",
						Rule: "费用总览",
					})
				},
			},
			"费用总览": {
				ItemFields: []string{
					"产品名称",
					"应付金额",
					"现金支付",
					"信用额度结算",
					"代金券抵扣",
					"储值卡地阔",
				},
				ParseFunc: func(ctx *Context) {
					//query := ctx.GetDom()
					//logs.Log.Informational(query.Html())
					//query.Find(".overview-product-table .next-table-row").Each(func(i int, selection *goquery.Selection) {
					//	outMap := map[int]interface{}{}
					//	tds := selection.Find("td")
					//	for i := 0; i < tds.Length(); i++ {
					//		outMap[i] = strings.Trim(tds.Eq(i).Text(), "¥")
					//	}
					//	ctx.Output(outMap)
					//})
					//fmt.Println("xxxxxxx")
					//logs.Log.Informational(ctx.GetText())
					content, err := ioutil.ReadAll(ctx.GetResponse().Body)
					if err != nil {
						ctx.Log().Error(err.Error())
					}
					ctx.GetResponse().Body.Close()
					r := new(CostOverviewResponse)
					err = json.Unmarshal(content, r)
					if err != nil {
						ctx.Log().Error(err.Error())
					}
					for _, e := range r.Data {
						ctx.Output(map[int]interface{}{
							0: e.ItemInfo,
							1: e.OweMoney,
							2: e.CashMoney,
							3: e.CreditPayAmount,
							4: e.CouponMoney,
							5: e.CardMoney,
						})
					}
				},
			},
		},
	},
}
