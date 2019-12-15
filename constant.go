package logger

const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

//日志切分模式
const (
	LogSplitTypeHour = iota
	LogSplitTypeSize
)

var levelText []string = []string{
	"Debug",
	"Trance",
	"Info",
	"Warn",
	"Error",
	"Fatal",
}
