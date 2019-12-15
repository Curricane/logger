package logger

import "fmt"

//想要被外部使用，需要一个同名的go文件

var log LogInterface

/*
file 初始化一个文件日志
console 初始化控制台日志
name: file or console
config: log_level log_name log_path
*/
func InitLogger(name string, config map[string]string) (err error) {
	switch name {
	case "file":
		log, err = NewFileLoger(config)
	case "console":
		log, err = NewConsoleLogger(config)
	default:
		err = fmt.Errorf("unsupport logger name: %s", name)
	}
	return
}

func Debug(format string, args ...interface{}) {
	log.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	log.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatal(format, args...)
}
