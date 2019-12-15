package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	LogChan       chan *LogData
	logSplitType  int
	logSplitHour  int
	logSplitSize  int64
	lastSplitHour int
}

func NewFileLoger(config map[string]string) (log LogInterface, err error) {
	logPath, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found log_path")
		return
	}

	logName, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found log_name")
		return
	}

	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("not found log_level")
		return
	}

	logSplitStr, ok := config["log_split_type"]
	if !ok {
		logSplitStr = "hour"
	}
	logSplitInt := getLogSplitType(logSplitStr)

	logSplitHourStr, ok := config["log_split_hour"]
	if !ok {
		logSplitHourStr = "1"
	}
	logSplitHourInt, err := strconv.Atoi(logSplitHourStr)
	if err != nil {
		logSplitHourInt = 1
	}

	logSplitSizeStr, ok := config["log_split_size"]
	if !ok {
		logSplitSizeStr = "104857600" //100M
	}
	logSplitSizeInt, err := strconv.ParseInt(logSplitSizeStr, 10, 64)
	if err != nil {
		logSplitSizeInt = 104857600
	}

	//level, err := strconv.Atoi(logLevel)
	level := getLogLevel(logLevel)
	if err != nil {
		return
	}

	logChanSize, ok := config["log_chan_size"]
	if !ok {
		logChanSize = "50000"
	}
	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}

	log = &FileLogger{
		level:        level,
		logPath:      logPath,
		logName:      logName,
		LogChan:      make(chan *LogData, chanSize),
		logSplitType: logSplitInt,
		logSplitHour: logSplitHourInt,
		logSplitSize: logSplitSizeInt,
	}
	log.Init()
	return log, nil
}

func (f *FileLogger) Init() {
	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.file = file

	//写error fatal 日志的文件
	filename = fmt.Sprintf("%s/%s_wef.log", f.logPath, f.logName)
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.warnFile = file

	go f.writeLogToFile()
}

func (f *FileLogger) checkSplitFile(weflog bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileByHour(weflog)
	} else {
		f.splitFileBySize(weflog)
	}
}

func (f *FileLogger) splitFileByHour(weflog bool) {
	now := time.Now()
	hour := now.Hour()
	if hour == f.lastSplitHour {
		return
	}

	f.lastSplitHour = hour
	var oldFileName string
	var backupFilename string
	if weflog {
		backupFilename = fmt.Sprintf("%s%s_wef_%04d%02d%02d%02d.log",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		oldFileName = fmt.Sprintf("%s%s_wef.log",
			f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s/%s_%04d%02d%02d%02d.log",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		oldFileName = fmt.Sprintf("%s%s.log",
			f.logPath, f.logName)
	}

	file := f.file
	if weflog {
		file = f.warnFile
	}
	file.Close()
	os.Rename(oldFileName, backupFilename)

	file, err := os.OpenFile(oldFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if weflog {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) splitFileBySize(weflog bool) {
	file := f.file
	if weflog {
		file = f.warnFile
	}

	statInfo, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	now := time.Now()
	var backupFilename string
	var oldFileName string
	if weflog {
		backupFilename = fmt.Sprintf("%s%s_wef_%04d%02d%02d_%02d%2d%2d.log",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		oldFileName = fmt.Sprintf("%s%s_wef.log",
			f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s%s_%04d%02d%02d_%02d%2d%2d.log",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		oldFileName = fmt.Sprintf("%s%s.log",
			f.logPath, f.logName)
	}

	file.Close()
	os.Rename(oldFileName, backupFilename)

	file, err = os.OpenFile(oldFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}

	if weflog {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) writeLogToFile() {
	if f.LogChan == nil {
		return
	}

	for logData := range f.LogChan {
		if logData.LevelStr == "fatal" ||
			logData.LevelStr == "warn" ||
			logData.LevelStr == "error" {
			f.checkSplitFile(true)
			fmt.Fprintf(f.warnFile, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
				logData.TimeStr, logData.LevelStr, logData.Filename,
				logData.FuncName, logData.LineNum, logData.Message)
			continue
		}
		f.checkSplitFile(false)
		fmt.Fprintf(f.file, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
			logData.TimeStr, logData.LevelStr, logData.Filename,
			logData.FuncName, logData.LineNum, logData.Message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	f.level = level
}
func (f *FileLogger) Debug(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelDebug, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}

}
func (f *FileLogger) Trace(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}
}
func (f *FileLogger) Info(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}
}
func (f *FileLogger) Warn(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}
}
func (f *FileLogger) Error(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}
}
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	//判断是否满足打印级别
	if f.level > LogLevelFatal {
		return
	}
	logData := writeLog(LogLevelFatal, format, args...)
	select {
	case f.LogChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	//检查一遍
	if f.file != nil {
		err := f.file.Close()
		if err != nil {
			fmt.Printf("close %s failed \n", f.file.Name())
		}
	}
	if f.warnFile != nil {
		err := f.warnFile.Close()
		if err != nil {
			fmt.Printf("close %s failed \n", f.warnFile.Name())
		}
	}
}
