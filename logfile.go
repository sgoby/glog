// User: szh
// Date: 2020/9/25
// Time: 20:09

package glog

import (
	"os"
	"time"
	fmt "github.com/sgoby/glog/gfmt"
	"sync"
)

const(
	SplitDaily = "Daily"
	SplitHourly = "Hourly"
)

type LogFile struct {
	cnf             *Config
	wfile           *os.File
	writedSize      int64
	currentFileName string
	mu              sync.Mutex
}

//
func createLogFile(cnf *Config) (f *LogFile,err error){
	_, err = os.Stat(cnf.FileLogPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(cnf.FileLogPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	f = &LogFile{cnf:cnf}
	return f,nil
}

func (f *LogFile) Write(p []byte) (nn int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	//
	switch f.cnf.SplitType {
	case SplitDaily,SplitHourly:
		nfileName := f.getDataTimeSplitFileName()
		if len(f.currentFileName) < 1 || nfileName != f.currentFileName{
			if f.wfile != nil{
				err = f.wfile.Close()
				if err != nil{
					return 0,err
				}
			}
			//
			f.wfile,err = f.createFile(nfileName)
			if err != nil{
				return 0,err
			}
			f.currentFileName = nfileName
		}
		f.writedSize = 0
	default:
		if f.wfile ==  nil  || f.writedSize > f.cnf.limitSize{
			currentFileName := f.getDefaultFileName()
			if f.wfile != nil{
				err = f.wfile.Close()
				if err != nil{
					return 0,err
				}
				nfileName := f.getSizeSplitFileName()
				err = os.Rename(currentFileName,nfileName)
				if err != nil{
					return 0,err
				}
				os.Truncate(currentFileName,0)
			}
			f.wfile, err = f.createFile(currentFileName)
			if err != nil {
				return 0, err
			}
			f.writedSize = 0
		}
	}
	nn,err = f.wfile.Write(p)
	f.writedSize += int64(nn)
	return
}

//
func (f *LogFile) Close() (err error){
	if f.wfile != nil{
		err = f.wfile.Close()
		f.wfile = nil
	}
	return err
}

//
func (f *LogFile) getDataTimeSplitFileName() string{
	switch f.cnf.SplitType {
	case SplitDaily:
		return fmt.Sprintf("%s/%s_%s.log", f.cnf.FileLogPath,f.cnf.Tag, time.Now().Format("2006-01-02"))
	case SplitHourly:
		return fmt.Sprintf("%s/%s_%s.log", f.cnf.FileLogPath,f.cnf.Tag, time.Now().Format("2006-01-02_15"))
	}
	return f.getDefaultFileName()
}
//
func (f *LogFile) getDefaultFileName() string{
	return fmt.Sprintf("%s/%s.log", f.cnf.FileLogPath,f.cnf.Tag)
}
//
func (f *LogFile) getSizeSplitFileName() string{
	return fmt.Sprintf("%s/%s_%d.log", f.cnf.FileLogPath,f.cnf.Tag,time.Now().Unix())
}
//
func (f *LogFile) createFile(filePath string) (nf *os.File,err error){
	return  os.OpenFile(filePath,  os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
}