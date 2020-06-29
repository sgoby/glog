package glog

import (
	"log"
	"fmt"
	"sync"
	"bufio"
	"os"
	"time"
)

type logWriter interface {
	Close() error
	Write(b []byte) (int, error)
	Stat() (os.FileInfo, error)
	Sync() (err error)
}


type Logger struct {
	cnf            Config
	logger         *log.Logger
	muCreate       *sync.Mutex
	logFilePath    string
	fileHandler    logWriter
	calldepth      int
	logLevel       uint
	bufWriter      *bufio.Writer
	flushChann     chan uint
	lastFlushTime  time.Time
	nextFlushTimer *time.Timer
}

//
func newLogger(cnf Config) (*Logger,error){
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
	cnf.limitSize = getByteSizeByM(cnf.SplitType)
	//
	lg := &Logger{
		cnf:cnf,
		muCreate:new(sync.Mutex),
		calldepth:3,
		logLevel:uint(lv),
		flushChann:make(chan uint,1),
	}
	return lg,nil
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
//
func (l *Logger)output(lv uint, args ...interface{}) {
	lvTag, ok := logLevelMap[lv]
	if !ok {
		return
	}
	//
	gLogger,err := l.getLogger()
	if err != nil{
		log.Println(err)
		return
	}
	if gLogger != nil{
		select {
		case l.flushChann <- 1:
			err = gLogger.Output(l.calldepth,fmt.Sprintf("[%s] %s\n", lvTag, fmt.Sprint(args...)))
			if err != nil{
				log.Println(err)
			}
			<-l.flushChann
		}
		l.flush()
	}
}
//
func (l *Logger) Write(p []byte) (n int, err error){
	l.output(INFO,string(p))
	return len(p),nil
}