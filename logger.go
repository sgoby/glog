package glog

import (
	"log"
	"fmt"
	"bufio"
	"os"
	"strings"
	"io"
)

type logWriter interface {
	Close() error
	Write(b []byte) (int, error)
	Stat() (os.FileInfo, error)
	Sync() (err error)
}

type Logger struct {
	cnf         Config
	mlogger     ILogger
	fileHandler logWriter
	calldepth   int
	logLevel    uint
	bufWriter   *bufio.Writer
}

//
func newLogger(cnf Config) (lg *Logger,err error){
	if len(cnf.Tag) < 1{
		cnf.Tag = defautlTag
	}
	if len(cnf.Level) < 1{
		cnf.Level = "Debug"
	}
	lv := 0
	for i,val := range logLevelMap{
		if val == cnf.Level{
			lv = int(i)
			break
		}
	}
	if len(cnf.SplitType) < 1{
		cnf.SplitType = SplitDaily
	}
	cnf.limitSize = getByteSizeByM(cnf.SplitType)
	lg = &Logger{
		cnf:cnf,
		calldepth:3,
		logLevel:uint(lv),
	}
	lg.mlogger,err = lg.getLogger()
	return lg,err
}


//
func (l *Logger)Debug(args ...interface{}) {
	if l.logLevel > DEBUG {
		return
	}
	l.output(DEBUG,args...)
}
func (l *Logger)DebugF(format string, args ...interface{}) {
	if l.logLevel > DEBUG {
		return
	}
	l.output(DEBUG,fmt.Sprintf(format, args...))
}

//
func (l *Logger)Info(args ...interface{}) {
	if l.logLevel > INFO {
		return
	}
	l.output(INFO,args...)
}
func (l *Logger)InfoF(format string, args ...interface{}) {
	if l.logLevel > INFO {
		return
	}
	l.output(INFO,fmt.Sprintf(format, args...))
}

//
func (l *Logger)Warn(args ...interface{}) {
	if l.logLevel > WARN {
		return
	}
	l.output(WARN,args...)
}
func (l *Logger)WarnF(format string, args ...interface{}) {
	if l.logLevel > WARN {
		return
	}
	l.output(WARN,fmt.Sprintf(format, args...))
}

//
func (l *Logger)Error(args ...interface{}) {
	if l.logLevel > ERROR {
		return
	}
	l.output(ERROR,args...)
}
func (l *Logger)ErrorF(format string, args ...interface{}) {
	if l.logLevel > ERROR {
		return
	}
	l.output(ERROR,fmt.Sprintf(format, args...))
}

//
func (l *Logger)Fatal(args ...interface{}) {
	if l.logLevel > FATAL {
		return
	}
	l.output(FATAL,args...)
}
func (l *Logger)FatalF(format string, args ...interface{}) {
	if l.logLevel > FATAL {
		return
	}
	l.output(FATAL,fmt.Sprintf(format, args...))
}
//
func (l *Logger)PanicRuntimeCaller(args ...interface{}){
	runMsg := readRuntimeCaller()
	argsStr := fmt.Sprintf("Panic: runtime error. %v", args...)
	l.output(FATAL,fmt.Sprintf("%s%s", argsStr,runMsg))
}
//======================================
func (l *Logger)output(lv uint, args ...interface{}) {
	lvTag, ok := logLevelMap[lv]
	if !ok {
		return
	}
	//
	var err error
	if l.mlogger != nil{
		if mlogio,ok := l.mlogger.(*Logio);ok {
			err = mlogio.OutputByLv(l.calldepth, lvTag, fmt.Sprintf("%s\n", fmt.Sprint(args...)))
		}else{
			err = l.mlogger.Output(l.calldepth, fmt.Sprintf("[%s] %s\n", lvTag, fmt.Sprint(args...)))
		}
	}
	if err != nil {
		log.Println(err)
	}
}

//
func (l *Logger) Write(p []byte) (n int, err error){
	l.output(INFO,string(p))
	return len(p),nil
}

//
func (l *Logger)getLogger() (g ILogger,err error){
	if strings.ToLower(strings.TrimSpace(l.cnf.LogType)) == "syslog"{
		return l.createSyslogLogger()
	}
	//
	lg := NewLogio(Ldate|Ltime|Lshortfile)
	logfile,err := createLogFile(&l.cnf)
	if err != nil{
		return nil,err
	}
	if l.cnf.AlsoStdout {
		lg.ResetWriter(io.MultiWriter(logfile, os.Stdout))
	}else{
		lg.ResetWriter(logfile)
	}
	return lg,nil
}