// glogs project glogs.go
package glogs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type Log struct {
	logerr   *Logger      //error log
	loginfo  *Logger      //info log
	size     int64        //save file size
	fileName string       //file name
	file     *os.File     //file ptr
	maxSize  int64        //file max size
	mutex    sync.RWMutex //mutex
	ltype    string       //log type ,json or normal
	basepath string       //your project home dir
	flag     int          //flag =1 Once everyday; flag=2 Twice a day ;flag=0 means not use time 定时任务，后期增加
	logPath  string       //log save dir
}

var err error

const (
	// Max          = 1024 * 100 //文件最大10M
	Max             = 1024 * 1024 * 10 //文件最大10M
	CreateFullFile  = 0
	CreateOnceADay  = 1
	CreateTwiceADay = 2

	FormatString = "20060102150405"
	Err          = "ERROR"
	Info         = "INFO"
	NORMAL       = "normal"
	JSON         = "json"
)

//没有就new 有就打开
func (this *Log) initFile() {
	if this.file != nil {
		return
	}
	//search lastest file 搜索最新文件
	dirs, _ := ioutil.ReadDir(this.logPath) //return fileinfo-list
	tmps := 0                               //max
	index := -1
	for i, j := range dirs {

		//		j.Name()==string  == log_20190806193040.log
		if j.Name()[:4] == "log_" && j.Name()[len(j.Name())-4:len(j.Name())] == ".log" && len(j.Name()) == 22 {
			str := j.Name()[4 : len(j.Name())-4]
			// fmt.Println(str)
			// "2019"==>2019
			val, _ := strconv.Atoi(str)
			if tmps == 0 {
				tmps = val
				index = i
			} else {
				if tmps < val {
					tmps = val
					index = i
				}
			}
		}

	}
	if index == -1 {
		this.newFile()
		return
	}
	//get filename
	filename := dirs[index].Name()
	//open file
	file, _ := os.OpenFile(this.logPath+filename, os.O_WRONLY|os.O_APPEND, 0666)
	this.file = file
	this.fileName = filename
	stat, _ := os.Stat(this.logPath + this.fileName)
	this.size = stat.Size()
	this.logerr = New(io.MultiWriter(file, os.Stderr), Err+": ", Ldate|Ltime|Lshortfile, "error")
	this.loginfo = New(io.MultiWriter(file, os.Stderr), Info+": ", Ldate|Ltime|Lshortfile, "info")
	this.logerr.SetbasePathAndLtype(this.basepath, this.ltype)
	this.loginfo.SetbasePathAndLtype(this.basepath, this.ltype)
}

func (this *Log) newFile() {
	if this.file != nil {
		return
	}
	tim := time.Now().Format(FormatString) //20190806_181719
	filename := "log_" + tim + ".log"
	file, _ := os.Create(this.logPath + filename)
	this.fileName = filename
	stat, err := os.Stat(this.logPath + this.fileName)
	if err != nil {
		panic("find dir err :" + this.logPath)
	}
	this.size = stat.Size()
	this.file = file //new
	//该函数自动回关闭file
	this.logerr = New(io.MultiWriter(file, os.Stderr), Err+": ", Ldate|Ltime|Llongfile, "error")
	this.loginfo = New(io.MultiWriter(file, os.Stderr), Info+": ", Ldate|Ltime|Llongfile, "info")
	this.logerr.SetbasePathAndLtype(this.basepath, this.ltype)
	this.loginfo.SetbasePathAndLtype(this.basepath, this.ltype)
	return
}

//判断是否超出大小，是就new
func (this *Log) isBreakSize() {
	if this == nil {
		panic("log can not be nil")
		return
	}
	if this.flag == CreateFullFile {
		stat, _ := os.Stat(this.logPath + this.fileName)
		this.size = stat.Size()
		if this.flag != 0 {
			return
		}
		if this.size >= this.maxSize {
			this.file.Close()
			this.file = nil
			this.newFile()
		}
	}

}

//err 和info 有 path 必须原始调用
func (this *Log) Error(v ...interface{}) {
	this.isBreakSize()
	this.logerr.Output(3, fmt.Sprintln(v...))
}

func (this *Log) INFO(v ...interface{}) {
	this.isBreakSize()
	this.loginfo.Output(3, fmt.Sprintln(v...))
}

func (this *Log) Println(v ...interface{}) {
	this.isBreakSize()
	fmt.Println(v...)
}

//panic -->err==>panic只有err类型
func (this *Log) Panic(v ...interface{}) {
	this.isBreakSize()
	s := fmt.Sprint(v...)
	this.logerr.Output(3, s)
	panic(s)
}

//print后退出程序，都可以用
func (this *Log) Fatal(v ...interface{}) {
	this.isBreakSize()
	this.logerr.Output(3, fmt.Sprint(v[1:]...))
	os.Exit(1)

}

//都可以用
func (this *Log) Print(v ...interface{}) {
	this.isBreakSize()
	fmt.Print(v...)

}

//大于该max自动创建新文件 ==>
func CreateLog() *Log {
	basepath, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	eLogs := &Log{
		flag:     CreateFullFile,
		maxSize:  Max,
		basepath: basepath + "/",
		ltype:    NORMAL,
		logPath:  basepath + "/logger/",
	}
	return eLogs
}

func (this *Log) InitLog() {
	switch this.flag {
	case CreateOnceADay:
		cron := newCron()
		//凌晨2点清理7天前的日志
		cron.addFunc("0 0 2 * * *", func() {
			cleanfile(this)
		})
		//凌早上8点创建新日志
		cron.addFunc("0 0 8 * * *", func() {
			docron(this)
		})
		cron.start()

		this.initFile()
	case CreateTwiceADay:
		cron := newCron()
		cron.addFunc("0 0 2 * * *", func() {
			cleanfile(this)
		})
		//每天8点，20点更新一次日志
		cron.addFunc("0 0 8,20 * * *", func() {
			docron(this)
		})
		cron.start()

		this.initFile()
	case CreateFullFile:
		cron := newCron()
		cron.addFunc("0 0 2 * * *", func() {
			cleanfile(this)
		})
		cron.start()
		this.initFile()
	}

}

/*
*  Max = 1024 * 1024 * 100 //文件最大10M init之后执行，默认10M
 */
func (this *Log) SetMaxSize(size int) {
	this.maxSize = int64(size)
}

/*
*  View type  json or numal
 */
func (this *Log) SetViewType(jsonOrNormal string) {
	if jsonOrNormal == JSON || jsonOrNormal == NORMAL {
		this.ltype = jsonOrNormal
	} else {
		panic("jsonOrNormal data error ,please use 'normal' or 'json'")
	}

}

/*
*  设置日志创建类型:
*   CreateFullFile  = 超出文件大小就创建
*   CreateOnceADay  = 一天1次
*   CreateTwiceADay = 一天2次，
*
 */
func (this *Log) SetLogCreateType(flag int) {
	this.flag = flag
}

func docron(log *Log) {
	if log == nil {
		panic("docron err log can not be nil")
		return
	}
	log.file.Close()
	log.file = nil
	log.newFile()
}

//定期清理日志，清理1周前的日志
func cleanfile(log *Log) {
	if log == nil {
		panic("docron err log can not be nil")
		return
	}
	dirs, _ := ioutil.ReadDir(log.logPath) //return fileinfo-list

	for _, j := range dirs {

		//		j.Name()==string  == log_20190806193040.log
		if j.Name()[:4] == "log_" && j.Name()[len(j.Name())-4:len(j.Name())] == ".log" && len(j.Name()) == 22 {
			str := j.Name()[4 : len(j.Name())-4]
			// fmt.Println(str)
			t, _ := time.ParseInLocation(FormatString, str, time.Local)
			t = t.AddDate(0, 0, 7)
			if t.Before(time.Now()) { //+7天后还比现在时间小
				// fmt.Println("删除:", j.Name())
				err := os.Remove(log.logPath + j.Name())
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}

	}

}
