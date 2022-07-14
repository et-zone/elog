package elog

import (
	cut "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

const (
	logTmFmt = "2006-01-02 15:04:05"
	constHourFormat = ".%Y%m%d%H"
	constDayFormat = ".%Y%m%d00"
)


// EncodeLevel 自定义日志级别显示
func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// EncodeTime 自定义时间格式显示
func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTmFmt) + "]")
}

// EncodeCaller 自定义行号显示
func encodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

// GetConsoleEncoder 输出日志到控制台
func GetConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

func getCuter(fileName string,maxAge time.Duration,cutTime time.Duration) io.Writer {
	dir,name:= getFileNameAndPath(fileName)
	if dir!=""&&dir!="."&&dir!="./"{
		if _, err := os.Stat(dir);err!=nil{
			_ = os.Mkdir(dir, os.ModePerm)
		}
	}
	if cutTime>= CutTimeDay {
		if writer, err := cut.New(dir+name+constDayFormat,
			cut.WithLinkName(dir+name), cut.WithRotationTime(cutTime),cut.WithMaxAge(maxAge));err!=nil{
			panic(err)
		}else {
			return writer
		}
	}else {
		if writer, err := cut.New(dir+name+constHourFormat,
			cut.WithLinkName(dir+name), cut.WithRotationTime(cutTime),cut.WithMaxAge(maxAge));err!=nil{
			panic(err)
		}else {
			return writer
		}
	}

}
func getFileNameAndPath(f string)(string,string){

	list:=strings.SplitAfter(f,"/")
	if len(list)==1{
		return "",list[0]
	}else {
		return strings.Join(list[:len(list)-1],""),list[len(list)-1]
	}
}