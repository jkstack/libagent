package conf

import (
	"github.com/jkstack/jkframe/utils"
	"github.com/jkstack/libagent/limit"
)

// Configure 基础配置
type Configure struct {
	// agent id, 为空则为服务器端分配
	//   $HOSTNAME: 使用当前主机名作为agent id
	//   $IP: 使用连接到服务器端的网卡IP作为agent id
	//   ${env}: 使用环境变量作为agent id
	ID string `json:"id" yaml:"id" kv:"id"`

	// server 服务器端地址，支持环境变量
	Server string `json:"server" yaml:"server" kv:"server"`

	// 日志相关配置
	Log struct {
		// 日志输出目标，支持stdout和file
		Target logTarget `json:"target" yaml:"target" kv:"target"`
		// 日志文件保存路径
		Dir string `json:"dir" yaml:"dir" kv:"dir"`
		// 日志文件滚动生成时的文件大小
		Size utils.Bytes `json:"size" yaml:"size" kv:"size"`
		// 日志文件滚动生成时的保留数量
		Rotate int `json:"rotate" yaml:"rotate" kv:"rotate"`
	} `json:"log" yaml:"log" kv:"log"`

	// 监控配置
	Monitor struct {
		// 是否启用监控数据上报
		Enabled bool `json:"enabled" yaml:"enabled" kv:"enabled"`
		// 监控数据上报间隔
		Interval utils.Duration `json:"interval" yaml:"interval" kv:"interval"`
	} `json:"monitor" yaml:"monitor" kv:"monitor"`

	// 资源限制配置
	Limit limit.Configure `json:"limit" yaml:"limit" kv:"limit"`
}
