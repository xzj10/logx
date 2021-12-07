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
		logFileName := time.Now().Format("20060102.log")
		if isErr {
			farr = append(farr, "logs", "error")
			createPath(strings.Join(farr, "/"))
			farr = append(farr, logFileName)
		} else {
			farr = append(farr, "logs", "info")
			createPath(strings.Join(farr, "/"))
			farr = append(farr, logFileName)
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
		fary := strings.Split(file, "/")
		file = fary[len(fary)-1]
		fn := runtime.FuncForPC(pc).Name()
		// fn = strings.Split(fn, ".")[1]
		return fmt.Sprintf("%v:%v():%v", file, fn, line)
	}
	return ""
}
