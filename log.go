package glog

import (
	"errors"
	"fmt"
	"sync"
)

const calldepth  =4

//
var loggerMap map[string]*Logger
var logMu *sync.RWMutex

func init(){
	loggerMap = make(map[string]*Logger)
	logMu = new(sync.RWMutex)
}

var  defautlTag = "app"
//
func OnInit(v interface{}) error{
	cnf,ok := v.(Config)
	if !ok{
		return errors.New("is no log config")
	}
	//
	lg,err := newLogger(cnf)
	if err != nil{
		return err
	}
	if defautlTag == "app" {
		defautlTag = lg.cnf.Tag
	}
	//
	lg.Info("==== log init finish ====")
	lg.calldepth = calldepth
	loggerMap[lg.cnf.Tag] = lg
	return nil
}

//
func Tag(tags ...string) *Logger{
	tag := defautlTag
	if len(tags) > 0{
		tag = tags[0]
	}
	logMu.RLock()
	lg,ok := loggerMap[tag]
	logMu.RUnlock()
	if !ok{
		if tag == defautlTag {
			logMu.Lock()
			lg,err := newLogger(Config{Tag:defautlTag,AlsoStdout:false})
			if err != nil{
				logMu.Unlock()
				panic(err)
			}
			loggerMap[tag] = lg
			logMu.Unlock()
			return lg
		}
		panic(fmt.Sprintf("log tag '%s' not config.",tag))
	}
	lg.calldepth = calldepth -1
	return lg
}

//
func Debug(args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.Debug(args...)

}

//
func DebugF(format string, args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.DebugF(format,args...)
}

//
func Info(args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.Info(args...)
}
func InfoF(format string, args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.InfoF(format,args...)
}

//
func Warn(args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.Warn(args...)
}
func WarnF(format string, args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.WarnF(format,args...)
}

//
func Error(args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.Error(args...)
}
func ErrorF(format string, args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.ErrorF(format,args...)
}

//
func Fatal(args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.Fatal(args...)
}
func FatalF(format string, args ...interface{}) {
	lg := Tag()
	lg.calldepth = calldepth
	lg.FatalF(format,args...)
}
//
func PanicRuntimeCaller(args ...interface{}){
	lg := Tag()
	lg.calldepth = calldepth
	lg.PanicRuntimeCaller(args...)
}
