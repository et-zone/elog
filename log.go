package elog

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	MaxAgeThreeDay = time.Hour * time.Duration(3*24)
	MaxAgeWeek     = time.Hour * time.Duration(7*24)
	MaxAgeMonth    = time.Hour * time.Duration(30*24)
	CutTimeHour    = time.Hour * time.Duration(1)
	CutTimeDay     = time.Hour * time.Duration(24)

)

var (
	errCfgNil = errors.New("config is nil")
	errLogNil = errors.New("log not init succ,it is nil")
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

type Level zapcore.Level

type Log struct {
	logger *zap.Logger
	*zap.SugaredLogger
}

type Config struct {
	FileName string //pathName
	Level    Level
	MaxAge   time.Duration //max save time
	CutTime  time.Duration //cut time
}

func NewLog(cfg *Config) *Log {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "Time",
		LevelKey:         "Level",
		NameKey:          "Logger",
		CallerKey:        "Code",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "Msg",
		StacktraceKey:    "Stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     encodeCaller,
		EncodeLevel:      encodeLevel,
		EncodeTime:       encodeTime,
		ConsoleSeparator: " ",
	})

	writerSyncer := zapcore.AddSync(getCuter(cfg.FileName, cfg.MaxAge, cfg.CutTime))
	consoleSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writerSyncer, consoleSyncer), zapcore.Level(cfg.Level)),
		//zapcore.NewCore(encoder, consoleSyncer, zapcore.DebugLevel),
	)
	log := zap.New(core,
		zap.AddCaller(),                       // code line
		zap.Development(),                     // develop
		zap.AddStacktrace(zapcore.ErrorLevel), // open  stack model where level = error
		zap.AddCallerSkip(1),                  // reporting the wrapper code as the caller.
	)
	return &Log{log,log.Sugar()}
}


func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}
