package logx

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func getLevel(levels string) LogLevel {
	lowLevals := strings.ToLower(levels)
	m := map[string]LogLevel{
		"debug": 1,
		"info":  2,
		"warn":  3,
		"error": 4,
		"fatal": 5,
	}
	ll, ok := m[lowLevals]
	if ok {
		return ll
	}
	msg := fmt.Sprintf("NewLogX()中参数的日志级别: %v不存在, 只能在debug, info, warn, error, fatal中选择一个!", levels)
	panic(msg)
}

func getLevelByIdx(idx LogLevel) string {
	m := map[LogLevel]string{
		1: "debug",
		2: "info",
		3: "warn",
		4: "error",
		5: "fatal",
	}
	return m[idx]
}

func createPath(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
	}
}

func getFileObj(isErr bool) (f *os.File) {
	_, mfile, _, ok := runtime.Caller(2)
	if ok {
		farr := strings.Split(mfile, "/")
		farr = farr[:len(farr)-1]
		logFileName := time.Now().Format("20060102")
		if isErr {
			farr = append(farr, "logs", "error")
			createPath(strings.Join(farr, "/"))
			farr = append(farr, logFileName+".log")
		} else {
			farr = append(farr, "logs", "info")
			createPath(strings.Join(farr, "/"))
			farr = append(farr, logFileName+".log")
		}
		fullPath := strings.Join(farr, "/")
		fObj, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		f = fObj
	}
	return
}

func reOpenFile(f *os.File, name string) (fo *os.File) {
	farr := strings.Split(f.Name(), "/")
	name = name + ".log"
	if farr[len(farr)-1] != name {
		f.Close()
		farr[len(farr)-1] = name
		fullPath := strings.Join(farr, "/")
		fObj, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			// panic(err)
			fmt.Println("重新打开新的日志文件错误!")
		}
		fo = fObj
	}
	return
}

func getWhere(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		// fary := strings.Split(file, "/")
		// file = fary[len(fary)-1]
		fn := runtime.FuncForPC(pc).Name()
		arr := strings.Split(fn, ".")
		if strings.Contains(fn, "(") {
			fn = fmt.Sprintf("%v.%v", arr[1], arr[2])
		} else {
			fn = arr[1]
		}
		return fmt.Sprintf("%v:%v():line(%v)", file, fn, line)
	}
	return ""
}

func getLogMsg(logl LogLevel, args ...interface{}) string {
	msg := ""
	if len(args) < 1 {
		return msg
	}
	format, ok := args[0].(string)
	if ok {
		if strings.Contains(format, "%") {
			args = args[1:]
			msg = fmt.Sprintf(format, args...)
		} else {
			msg = fmt.Sprintf("%v", args)
		}
	} else {
		msg = fmt.Sprintf("%v", args)
	}
	levels := getLevelByIdx(logl)
	msg = fmt.Sprintf("[%v][%v][%v]  %v", time.Now().Format("20060102 15:04:05"), levels, getWhere(5), msg)
	return msg
}

// ------------------------------------------------------is called

func Debug(args ...interface{}) string {
	return Log.Debug(args...)
}

func Info(args ...interface{}) string {
	return Log.Info(args...)
}

func Warn(args ...interface{}) string {
	return Log.Warn(args...)
}

func Error(args ...interface{}) string {
	return Log.Error(args...)
}

func Fatal(args ...interface{}) string {
	return Log.Fatal(args...)
}

func Fn(args ...interface{}) {
	msg := Info(args...)
	fmt.Println("msg = ", msg)
}
