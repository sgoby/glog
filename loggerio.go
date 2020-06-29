package glog

import (
	"log"
	"os"
	"time"
	"fmt"
	"io"
	"bufio"
	"strings"
)

//flush 的时间间隔
var flushInterval = time.Millisecond * 100
//
func (l *Logger)getLogger() (*log.Logger,error){
	var err error
	//syslog
	if strings.ToLower(strings.TrimSpace(l.cnf.LogType)) == "syslog"{
		if l.logger != nil {
			return l.logger,nil
		}
		l.logger,err = l.createSyslogLogger()
		return l.logger,err
	}
	//
	currentSize := l.getCurrentFileSize()
	path,fileName := getLogFilePath(l.cnf,currentSize)
	//
	if l.logger == nil {
		l.logger, err = l.createFileLogger(path, fileName)
		if err != nil {
			return nil,err
		}
		return l.logger,nil
	}
	//
	if len(l.logFilePath) > 0 && l.logFilePath == fileName{
		return l.logger,nil
	}
	//
	l.logger, err = l.createFileLogger(path, fileName)
	if err != nil {
		return nil,err
	}
	return l.logger,nil
}

//
func (l *Logger) createFileLogger(path,fileName  string) (lg *log.Logger,err error){
	l.flushChann <- 1
	l.muCreate.Lock()
	defer func() {
		l.muCreate.Unlock()
		<-l.flushChann
	}()

	if len(l.logFilePath) > 0 && l.logFilePath == fileName{
		return l.logger,nil
	}
	if l.logger != nil && l.cnf.limitSize > 0 && l.getCurrentFileSize() < l.cnf.limitSize{
		return l.logger,nil
	}
	//
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	//
	//
	if l.bufWriter == nil {
		l.bufWriter = bufio.NewWriter(os.Stdout)
	}else {
		l.bufWriter.Flush()
	}
	//
	if l.fileHandler != nil{
		l.fileHandler.Close()
		//split with limit file size.
		if l.cnf.limitSize > 0 && len(l.logFilePath) > 0{
			err = os.Rename(l.logFilePath,fileName)
			if err != nil {
				return nil, err
			}
			fileName = l.logFilePath
		}
	}
	//

	logfile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	//
	if l.cnf.AlsoStdout {
		l.bufWriter.Reset(io.MultiWriter(logfile, os.Stdout))
	}else{
		l.bufWriter.Reset(logfile)
	}
	//
	l.fileHandler = logfile
	l.logFilePath = fileName
	//
	return log.New(l.bufWriter, "", log.Ldate|log.Ltime|log.Lshortfile),nil
}
//
func (l *Logger) getCurrentFileSize() int64{
	currentSize := int64(0)
	if l.fileHandler != nil {
		fi, err := l.fileHandler.Stat()
		if err != nil {
			return 0
		}
		currentSize = fi.Size()
	}
	return currentSize
}
//
func (l *Logger) flush()  {
	select {
	case l.flushChann <- 1:
		l.flushBuf()
		<-l.flushChann
	default:
		if time.Now().Sub(l.lastFlushTime) > flushInterval{
			if l.nextFlushTimer != nil{
				l.nextFlushTimer.Stop()
			}
			l.nextFlushTimer = time.AfterFunc(flushInterval,l.flush)
		}
	}
	return
}

//
func (l *Logger) flushBuf(){
	l.muCreate.Lock()
	defer l.muCreate.Unlock()
	//
	if l.bufWriter.Size() < 1{
		return
	}
	//
	if l.bufWriter.Size() > 1024 || time.Now().Sub(l.lastFlushTime) > flushInterval{
		err := l.bufWriter.Flush()
		if err != nil{
			fmt.Println(err)
		}
		//
		if l.nextFlushTimer != nil{
			l.nextFlushTimer.Stop()
		}
		l.lastFlushTime = time.Now()
		l.nextFlushTimer = time.AfterFunc(flushInterval,l.flush)
	}
	return
}