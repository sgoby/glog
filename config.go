package glog

const (
	DEBUG uint = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF
)
//
var logLevelMap map[uint]string = map[uint]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	OFF:   "OFF",
}

//
type Config struct {
	Tag         string `ini:"Tag"`         // default app
	LogType     string `ini:"LogType"`     // support type:[File,Syslog] default File,
	FileLogPath string `ini:"FileLogPath"` // defult path id logs
	SysLogAddr  string `ini:"SysLogAddr"`  // support when type = Syslog ex: 127.0.0.1:514
	AlsoStdout  bool   `ini:"AlsoStdout"`  // default false
	Level       string `ini:"Level"`       // support level, low -> heigh: Debug, Info, Warn , Error ,Fatal, Off. default Debug
	SplitType   string `ini:"SplitType"`   // support type: Daily,Hourly,4mb . default Daily
	limitSize   int64
}

