package logger

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config map[string]string) (log LogInterface, err error) {
	logLevel, ok := config["log_level"]
	if !ok {
		err = fmt.Errorf("can not find log_level")
		return
	}
	level := getLogLevel(logLevel)

	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelInfo
	}
	log = &ConsoleLogger{
		level: level,
	}
	return
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}
	logData := writeLog(LogLevelDebug, format, args...)

	//只是使用同一writeLog获取 logData方法，获取日志信息，然后直接在终端打印
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	logData := writeLog(LogLevelWarn, format, args...)
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	logData := writeLog(LogLevelError, format, args...)
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	logData := writeLog(LogLevelFatal, format, args...)
	fmt.Fprintf(os.Stdout, "%-24s %-6s (%-14s:%-14s:%-5d) %s\n",
		logData.TimeStr, logData.LevelStr, logData.Filename,
		logData.FuncName, logData.LineNum, logData.Message)
}

func (c *ConsoleLogger) Close() {

}

func (c *ConsoleLogger) Init() {

}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelInfo
	}
	c.level = level
}
