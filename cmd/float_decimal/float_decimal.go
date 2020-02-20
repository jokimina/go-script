package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

func main() {
	var a,b float32
	a = 810000
	b = 810110.1
	c := a - b
	fmt.Println(a,b,c)
	fmt.Printf("%.2f\n", c)

	fmt.Println("---- decimal ----")
	da := decimal.NewFromFloat32(a)
	db := decimal.NewFromFloat32(b)
	dc := da.Sub(db)
	fmt.Println(dc)
	fmt.Println(strconv.ParseFloat(dc.String(), 32))
}
