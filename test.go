package elog

import (
	"fmt"
	"testing"
)

func Test(t *testing.T){
	cfg:=&Config{
		FileName: "Documents/all.log",//pathName
		Level:    InfoLevel,
		MaxAge:   MaxAgeWeek,
		CutTime:  CutTimeDay,
	}
	err:= Init(cfg)
	defer Close()
	if err!=nil{
		fmt.Println(err.Error())
	}
	Info("dsfff")
}