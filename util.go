package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

/*
	用于存储打印日志信息
	结合channel使用
	实现异步打印日志到文件提高并发量
*/
type LogData struct {
	Message  string
	TimeStr  string
	LevelStr string
	Filename string
	FuncName string
	LineNum  int
	//File     *os.File
}

/*
func sendLogToChan(logChan chan *LogData, file *os.File, level int, format string, args ...interface{}) {
	//参数检查
	if level < LogLevelDebug || level > LogLevelFatal {
		fmt.Println("invalid log level")
		return
	}
	if file == nil {
		return
	}
	//日志内容
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	levelStr := levelText[level]
	fileN, funcN, lineNum := GetLineIfo()
	fileName := path.Base(fileN)
	funcName := path.Base(funcN)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		Message:  msg,
		TimeStr:  nowStr,
		LevelStr:    level,
		Filename: fileName,
		FuncName: funcName,
		LineNum:  lineNum,
		File:     file,
	}

	select {
	case logChan <- logData:
	default:
	}
}
*/

func GetLineIfo() (fileName string, funcName string, lineNum int) {
	var skip int = 4 //栈帧
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNum = line
	}
	return
}

func getLogLevel(level string) int {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

func getLogSplitType(str string) int {
	switch str {
	case "size":
		return LogSplitTypeSize
	case "house":
		return LogSplitTypeHour
	default:
		return LogSplitTypeHour
	}
}

func writeLog(level int, format string, args ...interface{}) *LogData {
	//参数检查
	if level < LogLevelDebug || level > LogLevelFatal {
		fmt.Println("invalid log level")
		return nil
	}

	//日志内容
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.999")
	levelStr := levelText[level]
	fileN, funcN, lineNum := GetLineIfo()
	fileName := path.Base(fileN)
	funcName := path.Base(funcN)
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		Message:  msg,
		TimeStr:  nowStr,
		LevelStr: levelStr,
		Filename: fileName,
		FuncName: funcName,
		LineNum:  lineNum,
	}
	return logData
	//fmt.Fprintf(file, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n", nowStr, levelStr, fileName, funcName, lineNum, msg)
}
