package logx

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int16

const (
	UNKNOW LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var Log *LogX

type LogMsg struct {
	levels string
	date   string
	msg    string
}

type LogX struct {
	Level       LogLevel
	infoFile    *os.File
	errFile     *os.File
	chLogMsg    chan *LogMsg
	chErrLogMsg chan *LogMsg
}

func NewLogx(logLevel string) *LogX {
	lx := &LogX{
		Level: getLevel(logLevel),
	}
	if lx.Level > DEBUG {
		lx.infoFile = getFileObj(false)
		lx.errFile = getFileObj(true)
		// lx.chLogMsg = make(chan *LogMsg, 1000)
		// lx.chErrLogMsg = make(chan *LogMsg, 1000)
		// go lx.writeLogToFile()
	}
	return lx
}

func (lx *LogX) writeLogToFile() {
	for {
		select {
		// 判断日志写入的文件
		case ms := <-lx.chLogMsg:
			f := reOpenFile(lx.infoFile, ms.date)
			if f != nil {
				lx.infoFile = f
			}
			fmt.Fprintln(lx.infoFile, ms.msg)
		case ms := <-lx.chErrLogMsg:
			f := reOpenFile(lx.errFile, ms.date)
			if f != nil {
				lx.errFile = f
			}
			fmt.Fprintln(lx.errFile, ms.msg)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (lx *LogX) wlog(curr_f *os.File, msg string) {
	nows := time.Now().Format("20060102")
	f := reOpenFile(curr_f, nows)
	if f != nil {
		curr_f = f
	}
	fmt.Fprintln(curr_f, msg)
}

func (lx *LogX) log(logl LogLevel, args ...interface{}) string {
	if logl >= lx.Level {
		msg := getLogMsg(logl, args...)
		// 是否需要退出
		if logl == FATAL {
			if lx.Level > DEBUG {
				// lx.pushLog(lx.chErrLogMsg, &LogMsg{levels: levels, date: nows, msg: msg})
				// fmt.Fprintln(lx.errFile, msg)
				lx.wlog(lx.errFile, msg)
			}
			panic(msg)
		}
		if lx.Level > DEBUG {
			if logl >= ERROR {
				// lx.pushLog(lx.chErrLogMsg, &LogMsg{levels: levels, date: nows, msg: msg})
				// fmt.Fprintln(lx.errFile, msg)
				lx.wlog(lx.errFile, msg)
			} else {
				// lx.pushLog(lx.chLogMsg, &LogMsg{levels: levels, date: nows, msg: msg})
				// fmt.Fprintln(lx.infoFile, msg)
				lx.wlog(lx.infoFile, msg)
			}
		} else {
			fmt.Println(msg)
		}
		return msg
	}
	return fmt.Sprintf("%#v", args)
}

func (lx *LogX) pushLog(ch chan *LogMsg, logMsg *LogMsg) {
	select {
	case ch <- logMsg:
	default:
		// 要是chan满了,就把此条日志丢弃,保证不阻塞业务代码
	}

}

func (lx *LogX) Close() {
	lx.infoFile.Close()
	lx.errFile.Close()
}

// ---------------------------------------------------------------log fn

func (lx *LogX) Debug(args ...interface{}) string {
	return lx.log(DEBUG, args...)
}

func (lx *LogX) Info(args ...interface{}) string {
	return lx.log(INFO, args...)
}

func (lx *LogX) Warn(args ...interface{}) string {
	return lx.log(WARN, args...)
}

func (lx *LogX) Error(args ...interface{}) string {
	return lx.log(ERROR, args...)
}

func (lx *LogX) Fatal(args ...interface{}) string {
	return lx.log(FATAL, args...)
}
