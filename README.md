# libagent

[![libagent](https://github.com/jkstack/libagent/actions/workflows/build.yml/badge.svg)](https://github.com/jkstack/libagent/actions/workflows/build.yml)
[![license](https://img.shields.io/github/license/jkstack/libagent)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jkstack/libagent)](https://goreportcard.com/report/github.com/jkstack/libagent)
[![go-mod](https://img.shields.io/github/go-mod/go-version/jkstack/libagent)](https://github.com/jkstack/libagent)

agent封装类库，用于快速开发agent

## 开发方式

请查看[example-agent](https://github.com/jkstack/example-agent#agent%E5%BC%80%E5%8F%91)

## 已支持功能

1. 与agent-server创建websocket连接并进行握手分配agent id
2. 标准化日志输出
   - 支持stdout和日志文件双目标输出
   - 支持DEBUG、WARNING、INFO、ERROR级别的标准格式日志输出
   - 支持日志文件的滚动存储
3. 支持json、yaml或kv格式的配置文件，以下是一个kv格式配置文件示例

        id = example-01
        server = 127.0.0.1:13081
        # log
        log.target = stdout
        #log.target = stdout,file
        #log.dir = ./logs
        #log.size = 10M
        #log.rotate = 7
        # monitor
        monitor.enabled = true
        monitor.interval = 10s
        # limit
        limit.cpu_quota = 100
        limit.memory_limit = 1G
        limit.disk_limit = [{"dev":"8:0","read_bytes":"1M","write_bytes":"1M","read_iops":4000,"write_iops":4000}]
4. 支持通过[cgroups](https://github.com/containerd/cgroups)或[runtime](https://pkg.go.dev/runtime)库进行agent的资源限制，可通过配置文件中的limit相关选项进行配置
5. 支持agent的自监控数据上报，可通过配置文件中的monitor相关选项进行配置
6. 支持断线重连功能
7. 支持系统服务注册