package glog

import (
	"log"
	"os"
)

//
func (l *Logger) createSyslogLogger() (lg *log.Logger,err error){
	return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),nil
}
