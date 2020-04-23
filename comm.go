package glog

import (
	"runtime"
	"strings"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
)

// private
func readRuntimeCaller() string {
	rootPath := runtime.GOROOT()
	rootPath = strings.TrimSpace(rootPath)
	rootPath = strings.ToLower(rootPath)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)
	message := "    "
	skip := 1
	for {
		skip += 1
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		file = strings.TrimSpace(file)
		file = strings.ToLower(file)
		file = strings.Replace(file, "\\", "/", -1)
		if strings.Index(file, rootPath) == 0 {
			continue
		}
		message += fmt.Sprintf("\n    %s: %d", file, line)
	}
	return message
}

//
func getLogFilePath(cnf Config,currentSize int64) (string,string){
	path := cnf.FileLogPath
	tag := cnf.Tag
	if len(path) < 1 {
		localPath, _ := os.Getwd()
		path = localPath + "/logs"
	}
	//
	if len(tag) < 1{
		tag = "app"
	}
	//split with filezie limit
	if cnf.limitSize > 0 && currentSize >= cnf.limitSize {
		return path,fmt.Sprintf("%s/%s_%d.log", path,tag,time.Now().UnixNano())
	}
	//
	splitType := "Daily"
	if len(cnf.SplitType) > 0{
		splitType = cnf.SplitType
	}
	switch splitType {
	case "Daily":
		return path,fmt.Sprintf("%s/%s_%s.log", path,tag, time.Now().Format("2006-01-02"))
	case "Hourly":
		return path,fmt.Sprintf("%s/%s_%s.log", path,tag, time.Now().Format("2006-01-02_15"))
	}
	return path,fmt.Sprintf("%s/%s.log", path,tag)
}

//
func getByteSizeByM(fSize string) (v int64){
	reg,err := regexp.Compile(`\d+[m|M|mb|MB|mB|Mb]`)
	if err != nil{
		fmt.Println(err)
		return 0
	}
	if !reg.MatchString(fSize){
		return 0
	}
	//
	reg,err = regexp.Compile(`[m|M|b|B]`)
	fSize = reg.ReplaceAllString(fSize,"")
	v,err = strconv.ParseInt(fSize,10,64)
	if err != nil{
		return 0
	}
	v = v * 1024 * 1024
	return v
}