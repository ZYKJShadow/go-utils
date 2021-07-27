package main

import (
	"fmt"
	"go-utils/regex/util"
	"log"
)

func main() {
	matchNumberTest()
}

func matchAllParamsTest() {
	s := "OnSettle origin.num:0 oriUID:85 sourceUID:108 isForce:false"
	reg := "OnSettle origin.num:([0-9]+) oriUID:([0-9]+) sourceUID:([0-9]+) isForce:(.+)"
	var num int
	var oriUID int64
	var sourceUID int32
	var isForce bool
	params, err := util.MatchAllParams(s, reg, num, oriUID, sourceUID, isForce)
	if err != nil {
		log.Fatal(err)
		return
	}
	num = params[0].(int)
	oriUID = params[1].(int64)
	sourceUID = params[2].(int32)
	isForce = params[3].(bool)
	fmt.Println(num, oriUID, sourceUID, isForce)
}

func matchNumberTest() {
	s := "OnSettle origin.num:0 oriUID:85 sourceUID:108 isForce:false"
	res, err := util.MatchNumber(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
