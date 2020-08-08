package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

//实现日志输出的结构体
type Clog struct {
	prefix     string
	w          io.Writer
	level      LogLevel
	format     string
	caller     int
	dataFormat DataFormat
	support    map[string]FormatFunc
}

var exp = regexp.MustCompile(`\$([a-zA-Z0-9_]+)`)
var dataexp = regexp.MustCompile(`\$\{data\}`)
var support = map[string]FormatFunc{
	"fn": func(level LogLevel, skip int) string {
		fn, _, _, _ := runtime.Caller(skip)
		return path.Base(runtime.FuncForPC(fn).Name())
	},
	"FN": func(level LogLevel, skip int) string {
		fn, _, _, _ := runtime.Caller(skip)
		return runtime.FuncForPC(fn).Name()
	},
	"ln": func(level LogLevel, skip int) string {
		_, _, ln, _ := runtime.Caller(skip)
		return fmt.Sprintf("%d", ln)
	},
	"t": func(level LogLevel, skip int) string {
		return time.Now().Format("2006-01-02 03:04:05")
	},
	"T": func(level LogLevel, skip int) string {
		return time.Now().Format("2006-01-02 15:04:05")
	},
	"f": func(level LogLevel, skip int) string {
		_, f, _, _ := runtime.Caller(skip)
		return path.Base(f)
	},
	"F": func(level LogLevel, skip int) string {
		_, f, _, _ := runtime.Caller(skip)
		return f
	},
	"l": func(level LogLevel, skip int) string {
		return fmt.Sprintf("%s", level)
	},
}

func (c *Clog) parseFormat(lvl LogLevel) {

	//取默认函数，但是优先级是结构体内部的函数优先于默认的函数
	var supports = make(map[string]FormatFunc, 0)
	for k, v := range support {
		supports[k] = v
	}
	for k, v := range c.support {
		supports[k] = v
	}

	tmp := strings.Trim(c.format, " \r\n\t")

	kv := exp.FindAllStringSubmatch(tmp, -1)
	//匹配结果 [[匹配的值,子模式],...]
	var skv = innerSort(kv)
	sort.Sort(skv)

	for i := 0; i < len(skv); i++ {
		fn, ok := supports[skv[i][1]]
		var replace = skv[i][0]
		if ok {
			replace = fn(lvl, c.caller+3)
		}
		tmp = strings.ReplaceAll(tmp, skv[i][0], replace)
	}
	c.prefix += " " + tmp + " "
}

func (c *Clog) SetDataFormat(dfmt DataFormat) {
	c.dataFormat = dfmt
}

func (c *Clog) SetShowLevel(level LogLevel) {
	c.level = level
}

func (c Clog) Println(level LogLevel, v ...interface{}) {
	if level&c.level == level {
		c.parseFormat(level)

		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v", t)
			}
		}

		data = dataexp.ReplaceAllString(c.prefix, data)
		data = data + "\n"
		c.w.Write([]byte(data))
		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprintln(c.w, v...)
	}
}

func (c Clog) Print(level LogLevel, v ...interface{}) {
	if level&c.level == level {
		c.parseFormat(level)
		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v", t)
			}
		}
		data = dataexp.ReplaceAllString(c.prefix, data)
		c.w.Write([]byte(data))

		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprint(c.w, v...)
	}
}

func (c Clog) Printf(level LogLevel, format string, v ...interface{}) {
	if level&c.level == level {
		c.parseFormat(level)
		data := c.prefix + fmt.Sprintf(format, v...)
		c.w.Write([]byte(data))
	}
}

func (c Clog) Info(v ...interface{}) {
	if Lvl_Info&c.level == Lvl_Info {
		c.parseFormat(Lvl_Info)

		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v ", t)
			}
		}

		data = dataexp.ReplaceAllString(c.prefix, data)
		data = data + "\n"
		c.w.Write([]byte(data))
		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprintln(c.w, v...)
	}
}

func (c Clog) Warning(v ...interface{}) {
	if Lvl_Warning&c.level == Lvl_Warning {
		c.parseFormat(Lvl_Warning)
		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v ", t)
			}
		}

		data = dataexp.ReplaceAllString(c.prefix, data)
		data = data + "\n"
		c.w.Write([]byte(data))
		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprintln(c.w, v...)
	}
}

func (c Clog) Error(v ...interface{}) {
	if Lvl_Error&c.level == Lvl_Error {
		c.parseFormat(Lvl_Error)
		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v ", t)
			}
		}

		data = dataexp.ReplaceAllString(c.prefix, data)
		data = data + "\n"
		c.w.Write([]byte(data))
		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprintln(c.w, v...)
	}
}

func (c Clog) Debug(v ...interface{}) {
	if Lvl_Debug&c.level == Lvl_Debug {
		c.parseFormat(Lvl_Debug)
		var data = ""
		if c.dataFormat == FMT_Json {
			tmp, _ := json.Marshal(v)
			data = string(tmp)
		} else {
			for _, t := range v {
				data += fmt.Sprintf("%+v ", t)
			}
		}

		data = dataexp.ReplaceAllString(c.prefix, data)
		data = data + "\n"
		c.w.Write([]byte(data))
		//fmt.Fprint(c.w, c.prefix)
		//fmt.Fprintln(c.w, v...)
	}
}

func (c *Clog) SetWriter(wd io.Writer) {
	c.w = wd
}

func (c *Clog) SetPrefix(prefix string) {
	c.prefix = prefix
}

func (c *Clog) AddCallerLevel() {
	c.caller++
}

func (c *Clog) SetFormat(format string) {
	c.format = format
}

func (c Clog) GetFormat() string {
	return c.format
}

//添加自定义format函数
func (c *Clog) AddCustomFormatFunc(name string, fn FormatFunc) {
	if c.support == nil {
		c.support = make(map[string]FormatFunc, 0)
	}
	c.support[name] = fn
}

func NewClog() *Clog {
	return &Clog{
		w:      os.Stdout,
		level:  Lvl_Warning | Lvl_Info,
		format: "[$l] $T file:$f line:$ln func:$fn ${data}",
	}
}

type innerSort [][]string

func (s innerSort) Len() int {
	return len(s)
}

func (s innerSort) Less(i, j int) bool {
	return len(s[i][0]) > len(s[j][0])
}

func (s innerSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//提供包导出函数
var dfLog = Clog{
	w:      os.Stdout,
	level:  Lvl_Warning | Lvl_Info,
	format: "[$l] $T file:$f line:$ln func:$fn ${data}",
	caller: 1,
}

func SetShowLevel(level LogLevel) {
	dfLog.SetShowLevel(level)
}

func Println(level LogLevel, v ...interface{}) {
	dfLog.Println(level, v...)
}

func Print(level LogLevel, v ...interface{}) {
	dfLog.Print(level, v...)
}

func Printf(level LogLevel, format string, v ...interface{}) {
	dfLog.Printf(level, format, v...)
}

func Info(v ...interface{}) {
	dfLog.Info(v...)
}

func Warning(v ...interface{}) {
	dfLog.Warning(v...)
}

func Error(v ...interface{}) {
	dfLog.Error(v...)
}

func Debug(v ...interface{}) {
	dfLog.Debug(v...)
}

func SetWriter(wd io.Writer) {
	dfLog.SetWriter(wd)
}

func SetPrefix(prefix string) {
	dfLog.SetPrefix(prefix)
}

func AddCallerLevel() {
	dfLog.AddCallerLevel()
}

func SetFormat(format string) {
	dfLog.SetFormat(format)
}

func GetFormat() string {
	return dfLog.GetFormat()
}

//添加自定义format函数
func AddCustomFormatFunc(name string, fn FormatFunc) {
	dfLog.AddCustomFormatFunc(name, fn)
}
