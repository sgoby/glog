package glog

import (
	"testing"
	"sync"
	"time"
	fmt "github.com/sgoby/glog/gfmt"
)

func Test_log(t *testing.T){
	//Info("hahaha1111111")
	//Tag("app").Info("asdfasdfasdf")

	//OnInit(Config{LogType:"syslog"})
	m := make(map[string]interface{})
	m["hello"] = 2022
	OnInit(Config{})

	DebugF("test json %j",m)

	wg := &sync.WaitGroup{}
	beginTime := time.Now()

	for i := 0; i < 1000;i++ {
		wg.Add(1)
		go writeLog(wg)
	}
	wg.Wait()
	ut := time.Now().Sub(beginTime)
	time.Sleep(time.Second * 1)
	//Info("f")
	fmt.Println("use time: ",ut)
}


func writeLog(wg *sync.WaitGroup){
	defer wg.Done()
	for i:=0;i < 10;i++ {
		Debug(i,"为什么博主会特意讲一")
		Info(i,"为什么博主会特意讲一下centos mini版的安装步骤呢，因为博主在VMware workstation上安装的非mini版本的centos,本想安装mini版1")
		Warn(i,"为什么博主会特意讲一")
		Error(i,"为什么博主会特意讲一22")
	}
}