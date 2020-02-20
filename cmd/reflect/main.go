package main

import (
	"fmt"
	"reflect"
)

type UnknownType struct {
	S1     int `tag:"a"`
	S2, S3 string
	F      func(a string)
}

type A int
type B = int

func (u UnknownType) func1(a string) {
	fmt.Println(a)
}

func main() {
	secret := UnknownType{1, "A", "B", func(a string) {
		//fmt.Println(a)
	}}

	getType := reflect.TypeOf(secret)
	getValue := reflect.ValueOf(secret)

	for i := 0; i < getValue.NumField(); i++ {
		field := getValue.Field(i)
		tField := getType.Field(i)
		fType := field.Type()
		fValue := field.Interface()
		fmt.Printf("Num %d, Name: %s, Type: %s, Value: %v, Tag: %v\n", i, tField.Name, fType, fValue, tField.Tag)
	}

	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("Method %s: %v\n", m.Name, m.Type)
		args := []reflect.Value{reflect.ValueOf("reflect")}
		getValue.Method(i).Call(args)
	}

	var a A = 1
	var b B = 2
	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	fmt.Println(ra.Type(), rb.Type())
	fmt.Println(ra.Kind(), rb.Kind())
}
