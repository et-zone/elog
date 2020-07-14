package elog

import "github.com/et-zone/elog/glogs"

var eLogs *glogs.Log

const (
	NORMAL          = glogs.NORMAL
	JSON            = glogs.JSON
	CreateFullFile  = glogs.CreateFullFile
	CreateOnceADay  = glogs.CreateOnceADay
	CreateTwiceADay = glogs.CreateTwiceADay
)

/*
* Create Log
 */
func CreateLog() {
	//json/normal
	eLogs = glogs.CreateLog()
}

/*
*  default normal
*  View type  json or normal
 */
func SetViewType(jsonOrNormal string) {
	eLogs.SetViewType(jsonOrNormal)
}

/*
*  设置日志创建类型:
*   CreateFullFile  = 超出文件大小就创建
*   CreateOnceADay  = 一天1次
*   CreateTwiceADay = 一天2次，
*
 */
func SetLogCreateFlag(flag int) {
	eLogs.SetLogCreateType(flag)
}
func InitLog() {
	eLogs.InitLog()
}

/*
*  Max = 1024 * 1024 * 100 //文件最大10M init之后执行，默认10M
 */
func SetMaxSize(size int) {
	eLogs.SetMaxSize(size)
}

//err 和info 有 path 必须原始调用
func Error(v ...interface{}) {
	eLogs.Error(v...)
}

func INFO(v ...interface{}) {
	eLogs.INFO(v...)
}

//不存储日志
func Println(v ...interface{}) {
	eLogs.Println()
}

//不存储日志
func Print(v ...interface{}) {
	eLogs.Print(v...)

}

//panic -->err==>panic只有err类型
func Panic(v ...interface{}) {
	eLogs.Panic(v...)
}

//print后退出程序
func Fatal(v ...interface{}) {
	eLogs.Fatal(v...)
}
