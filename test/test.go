package main

import (
	"github.com/et-zone/elog"
	"sync"
	"time"
)

var gc sync.WaitGroup

func test3() {
	elog.CreateLog()
	elog.SetLogCreateFlag(elog.CreateTwiceADay) //文件超出大小就新建
	elog.SetMaxSize(1024 * 1024 * 1000)         //存储大小100M
	elog.SetViewType(elog.NORMAL)               //json格式
	elog.InitLog()
	// elog.Error("aaa", "111")
	// elog.INFO("aaa", "111")
}
func test4() {
	for {
		elog.INFO("aaa", "111")
		time.Sleep(time.Second * 2)
	}
}
func main() {
	// test2()
	test3()
	test4()
}

//INFO: 2020/07/13 16:18:08 /home/gzy/Desktop/go/src/logs/glogs/glogs.go:145: 222222222 adf
//INFO: 2020/07/13 16:18:35 /home/gzy/Desktop/go/src/logs/main.go:34: 222222222 adf
//INFO: 2020/07/13 16:19:02 /usr/local/go/src/runtime/asm_amd64.s:1337: 222222222 adf
//INFO: 2020/07/13 16:19:24 ???:0: 222222222 adf
