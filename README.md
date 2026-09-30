# glog

轻量的 Go 日志库，集成**文件自动分割、多 Tag 日志管理、JSON 格式化和 Linux Syslog 输出**。无需额外引入日志轮转组件，即可按模块记录日志、按时间或大小拆分文件。

## 特点

- **内置文件分割**：支持按天、按小时或按文件大小分割，自动创建日志目录，文件名包含 Tag，便于查找和归档。
- **多 Tag 独立配置**：为不同模块设置各自的日志目录、输出级别、分割策略和输出方式。
- **直接输出 JSON 内容**：在格式化日志中使用 `%j`，将 map、struct 等可 JSON 序列化的数据写入消息，无需手动调用 `json.Marshal`。
- **文件与终端同步输出**：通过 `AlsoStdout` 同时写入日志文件和标准输出，方便本地调试。
- **Linux Syslog 集成**：通过 UDP 发送至 Syslog 服务，便于接入集中式日志收集。
- **常用接口齐全**：提供分级日志、格式化日志、文件名与行号、调用栈记录；`Logger` 实现 `io.Writer`，可对接标准库 `log`。
- **无第三方运行时依赖**：模块仅依赖 Go 标准库和库内的 `gfmt` 包；文件输出使用内置缓冲写入。

## 安装

```sh
go get github.com/sgoby/glog
```

模块在 `go.mod` 中声明的 Go 版本为 `1.14`。当前提供 Linux 和 Windows 平台实现。

## 快速开始

在应用启动时完成初始化，再开始记录日志：

```go
package main

import (
    "log"
    "net/http"

    "github.com/sgoby/glog"
)

func main() {
    if err := glog.OnInit(glog.Config{
        Tag:         "app",
        FileLogPath: "logs",
        AlsoStdout:  true,
        Level:       "Info",
        SplitType:   glog.SplitDaily,
    }); err != nil {
        log.Fatal(err)
    }

    glog.Info("服务启动")
    glog.InfoF("监听地址：%s", ":8080")
    glog.InfoF("启动配置：%j", map[string]interface{}{
        "port": 8080,
        "env":  "dev",
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

文件默认包含日期、时间、日志级别和调用位置，输出示意：

```text
2026/09/30 10:00:00 [INFO]   main.go:24  服务启动
2026/09/30 10:00:00 [INFO]   main.go:25  监听地址：:8080
2026/09/30 10:00:00 [INFO]   main.go:26  启动配置：{"env":"dev","port":8080}
```

按天分割时，日志写入 `logs/app_2026-09-30.log`。`%j` 只将消息中的对应参数转为 JSON，整条日志仍为文本格式。

## 日志级别与接口

级别从低到高为：

```text
Debug < Info < Warn < Error < Fatal < Off
```

`Level` 设置最低输出级别，例如 `Info` 会过滤 `Debug`；`Off` 关闭常规分级日志。级别名称不区分大小写。

| 级别 | 普通日志 | 格式化日志 |
| --- | --- | --- |
| Debug | `glog.Debug(args...)` | `glog.DebugF(format, args...)` |
| Info | `glog.Info(args...)` | `glog.InfoF(format, args...)` |
| Warn | `glog.Warn(args...)` | `glog.WarnF(format, args...)` |
| Error | `glog.Error(args...)` | `glog.ErrorF(format, args...)` |
| Fatal | `glog.Fatal(args...)` | `glog.FatalF(format, args...)` |

普通接口采用 `Sprint` 风格拼接参数；格式化接口支持 `%s`、`%d`、`%v` 等常用格式符，并扩展了 `%j`。

**`Fatal` / `FatalF` 仅记录 FATAL 级别日志，不会退出进程或触发 panic。**

## 文件自动分割

通过 `SplitType` 选择分割方式：

| 配置值 | 分割方式 | 文件名示例（Tag 为 `app`） |
| --- | --- | --- |
| `glog.SplitDaily` 或 `"Daily"` | 按天，默认值 | `app_2026-09-30.log` |
| `glog.SplitHourly` 或 `"Hourly"` | 按小时 | `app_2026-09-30_10.log` |
| `"4M"` | 按大小，约 4 MiB | 当前文件 `app.log`，轮转后为 `app_<Unix时间戳>.log` |

大小配置支持 `K`、`M`、`G` 后缀（大小写均可，按 1024 进制），也支持直接填写字节数，例如 `"4096"`。请使用 `"4M"`，不要写成 `"4mb"`；时间策略需使用准确的 `Daily` / `Hourly` 大小写。

分割在实际写入时触发。大小模式按本次打开文件后的写入量判断，超过阈值后在后续写入时轮转，并非文件大小的严格上限。库不会自动清理历史日志，保留周期需由应用或外部工具管理。

## 多 Tag 日志

对每个 Tag 调用一次 `OnInit`，即可为业务模块建立独立配置。以下代码在初始化默认日志后执行：

```go
if err := glog.OnInit(glog.Config{
    Tag:         "payment",
    FileLogPath: "logs/payment",
    Level:       "Info",
    SplitType:   "100M",
}); err != nil {
    log.Fatal(err)
}

// 包级接口使用第一个成功初始化的 Tag。
glog.Info("应用日志")

// 指定 Tag，使用该模块的独立配置。
glog.Tag("payment").InfoF("订单支付成功：order_id=%s", "ORD-1001")
glog.Tag("payment").ErrorF("支付失败：%s", "余额不足")
```

`glog.Tag()` 返回默认 Logger，`glog.Tag("payment")` 返回指定 Logger。自定义 Tag 应先初始化，否则访问时会 panic。建议在启动阶段完成全部 Tag 配置，再启动业务 goroutine。

## 配置参考

`OnInit` 接收 `glog.Config` **值**，不是 `*glog.Config` 指针。

| 字段 | 默认值 | 说明 |
| --- | --- | --- |
| `Tag` | `app` | 日志实例标识，用于文件命名及 Syslog Tag |
| `LogType` | 文件输出 | 设置为 `Syslog` 启用 Syslog，不区分大小写；其余值走文件输出 |
| `FileLogPath` | 当前工作目录下的 `logs` | 文件输出目录，不存在时自动创建 |
| `SysLogAddr` | `127.0.0.1:514` | Linux Syslog 的 UDP 目标地址 |
| `AlsoStdout` | `false` | 是否同时输出至标准输出 |
| `Level` | `Debug` | 最低输出级别；无法识别的值也按 Debug 处理 |
| `SplitType` | `Daily` | 文件分割方式：`Daily`、`Hourly` 或大小字符串 |

### INI 配置映射

`Config` 字段带有 `ini` 标签，可由应用选用的 INI 解析库映射。glog 本身不读取 INI 文件；解析后将配置值传给 `glog.OnInit(cfg)` 即可。

```ini
[Glog]
Tag = app
LogType = File
FileLogPath = logs
AlsoStdout = true
Level = Info
SplitType = Daily
SysLogAddr = 127.0.0.1:514
```

## Linux Syslog

```go
if err := glog.OnInit(glog.Config{
    Tag:        "app",
    LogType:    "Syslog",
    SysLogAddr: "127.0.0.1:514",
    AlsoStdout: true,
    Level:      "Info",
}); err != nil {
    log.Fatal(err)
}

glog.Info("发送到 Syslog")
```

Linux 下使用 UDP，Syslog facility 为 `local5`，协议 severity 固定为 `debug`；glog 的实际日志级别以 `[INFO]`、`[ERROR]` 等文本写入消息。`FileLogPath` 和 `SplitType` 不参与 Syslog 输出。

Windows 下将 `LogType` 设为 `Syslog` 时，会输出到标准输出，不会向 Syslog 服务发送数据。

## 对接标准库与记录调用栈

`Logger` 实现了 `io.Writer`，写入内容按 Info 级别处理：

```go
// 先通过 OnInit 初始化 app。
stdLogger := log.New(glog.Tag("app"), "[http] ", 0)
stdLogger.Print("request received")
```

需要在恢复 panic 时记录调用栈，可使用：

```go
defer func() {
    if err := recover(); err != nil {
        glog.PanicRuntimeCaller(err)
    }
}()
```

`PanicRuntimeCaller` 记录 FATAL 日志和调用栈，本身不执行 `recover`、不触发 panic，也不受常规 `Level` 过滤控制。

## 缓冲与退出行为

文件模式使用缓冲并异步刷新；Linux Syslog 模式使用 `bufio.Writer` 缓冲。当前高层 `Logger` API 未提供公开的 `Flush` / `Close` 方法，不保证进程退出前所有缓冲日志均已写出；Syslog 的少量日志也可能停留在缓冲区中。需要严格退出落盘保障时，应先评估这一行为。
