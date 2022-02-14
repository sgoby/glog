package glog

import (
	"runtime"
	"strings"
	"strconv"
	fmt "github.com/sgoby/glog/gfmt"
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
func getByteSizeByM(fSize string) (v int64){
	fSize = strings.Trim(fSize, " \r\n")
	num := 1
	mult := 1024
	if len(fSize) > 1 {
		switch fSize[len(fSize)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			fSize = fSize[0 : len(fSize)-1]
		}
	}
	parsed, _ := strconv.Atoi(fSize)
	return int64(parsed * num)
}