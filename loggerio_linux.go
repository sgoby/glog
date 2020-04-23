package glog

import (
	"log"
	"os"
	"log/syslog"
	"bufio"
	"io"
)

/*
syslog config example:

# $template logFormat,"%TIMESTAMP:::date-rfc3339% %syslogtag% - %programname% - %msg%\n"
# $template localPath,"/var/log/%syslogtag%/local-%$year%-%$month%-%$day%.log"
$template localPath,"/var/log/%programname%/local-%$year%-%$month%-%$day%.log"
# local5.*   -?localPath;logFormat
local5.*   -?localPath;

*/
type sysLogWriter struct {
	mSyslogWriter *syslog.Writer
}
//
func (s *sysLogWriter) Sync() (err error){return nil }
func (s *sysLogWriter) Stat() (os.FileInfo, error){
	return nil,nil
}

func (s *sysLogWriter) Write(b []byte) (int, error){
	if s.mSyslogWriter == nil{
		return 0,nil
	}
	return s.mSyslogWriter.Write(b)
}
func (s *sysLogWriter) Close() (err error){
	if s.mSyslogWriter == nil{
		return nil
	}
	return s.mSyslogWriter.Close()
}
//
func newSysLogWriter(network, raddr string, tag string) (*sysLogWriter,error) {
	w,err := syslog.Dial(network,raddr,syslog.LOG_LOCAL5|syslog.LOG_DEBUG,tag)
	if err != nil{
		return nil,err
	}
	//
	sw := &sysLogWriter{
		mSyslogWriter:w,
	}
	return sw,nil
}

//
func (l *Logger) createSyslogLogger() (lg *log.Logger,err error){
	if l.bufWriter == nil {
		l.bufWriter = bufio.NewWriter(os.Stdout)
	}else {
		l.bufWriter.Flush()
	}
	//
	if len(l.cnf.SysLogAddr) < 1{
		l.cnf.SysLogAddr = "127.0.0.1:514"
	}
	//
	mSysLogWriter,err := newSysLogWriter("udp",l.cnf.SysLogAddr,l.cnf.Tag)
	if err != nil{
		return nil,err
	}
	//
	if l.cnf.AlsoStdout {
		l.bufWriter.Reset(io.MultiWriter(mSysLogWriter, os.Stdout))
	}else{
		l.bufWriter.Reset(mSysLogWriter)
	}
	//
	l.fileHandler = mSysLogWriter
	return log.New(l.bufWriter, "", log.Ldate|log.Ltime|log.Lshortfile),nil
}