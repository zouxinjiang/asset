package log

import (
	"io"
)

type Logger interface {
	SetShowLevel(level LogLevel)
	SetPrefix(prefix string)
	SetWriter(writer io.Writer)

	Println(level LogLevel, v ...interface{})
	Print(level LogLevel, v ...interface{})
	Printf(level LogLevel, fmt string, v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	//设置日志格式
	SetFormat(format string)
	GetFormat() string
	//添加自定义日志格式函数
	AddCustomFormatFunc(name string, fn FormatFunc)
}

type (
	LogLevel   uint64
	DataFormat string
)

const (
	Lvl_Info LogLevel = 1 << iota
	Lvl_Warning
	Lvl_Error
	Lvl_Debug

	Lvl_All = Lvl_Info | Lvl_Warning | Lvl_Error | Lvl_Debug

	FMT_Json  DataFormat = "json"
	FMT_Plain DataFormat = "plain"
)

var str2lvl = map[LogLevel]string{
	Lvl_Debug:   "DEBUG",
	Lvl_Error:   "ERROR",
	Lvl_Info:    "INFO",
	Lvl_Warning: "WARNING",
}

func (ll LogLevel) String() string {
	return str2lvl[ll]
}

type FormatFunc func(level LogLevel, skip int) string
