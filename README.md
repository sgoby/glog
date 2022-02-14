# glog 介绍

高性能日志 logger 处理库，自带日志文件分割(按天，按小时，按文件大小)功能，在linux环境下支持日志转发到 Syslog。

支持日志级别（低 > 高）：Debug, Info , Warn , Error , Fatal , Off

可配置是否需要回显，支持多标签日志tag, 支持ini配置



## Installation

```go
go get github.com/sgoby/glog
```



#### 配置信息：

```go
//
type Config struct {
	Tag         string `ini:"Tag"`         // default app
	LogType     string `ini:"LogType"`     // support type:[File,Syslog] default File,
	FileLogPath string `ini:"FileLogPath"` // defult path id logs
	SysLogAddr  string `ini:"SysLogAddr"`  // support when type = Syslog ex: 127.0.0.1:514
	AlsoStdout  bool   `ini:"AlsoStdout"`  // default false
	Level       string `ini:"Level"`       // support level, low -> heigh: Debug, Info, Warn , Error ,Fatal, Off. default Debug
	SplitType   string `ini:"SplitType"`   // support type: Daily,Hourly,4mb . default Daily
}
```



#### INI配置：

```ini
[Glog]
# default app
Tag = "app"
# defult path id logs 当 LogType = File 生效
FileLogPath = logs
# 是否需回显 default false
AlsoStdout = true
# support level, low -> heigh: Debug, Info, Warn , Error ,Fatal, Off
# default Debug
Level = Debug
# support type: Daily(每天),Hourly(每小时),4mb(按文件大小分割,单位:M) 当 LogType = File 生效
# default Daily
SplitType = Daily
# 日志记录方式[File,Syslog] 默认是file, (注意：Syslog 在linux环境才有效)
# LogType = File
# syslog 服务器地址,当 LogType = Syslog 生效, 默认:127.0.0.1:154,（备注:local5,debug）
# SysLogAddr = 127.0.0.1:154
```



#### 示例：

```go
glog.OnInit(Config{LogType:"syslog"})

.....
m := make(map[string]interface{})
m["hello"] = 2022
//format json
DebugF("test json %j",m)
//output   2022/02/14 12:31:59 [DEBUG]  log_test.go:19  test json {"hello":2022}
glog.Info(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")
glog.Debug(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")
glog.Warn(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")
glog.Error(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")
glog.Fatal(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")

//任意参数
glog.Info(1,"index","为什么博主会特意讲一下centos mini版的安装步骤呢")

//格式化
glog.InfoF("Index:%d 为什么博主会特意讲一下centos mini版的安装步骤呢",1)

.....
//自定义标计
glog.Tag("game").Info(1,"为什么博主会特意讲一下centos mini版的安装步骤呢")


```


