package conf

import (
	"github.com/jkstack/agent/internal/utils"
	"github.com/jkstack/agent/limit"
)

// Configure common configure
type Configure struct {
	// agent id, if empty generate by server
	//   $HOSTNAME: use agent id by hostname
	//   $IP: use agent id by ip of connect server
	//   ${env}: use envionment variables
	ID string `json:"id" yaml:"id" kv:"id"`

	// server address
	Server string `json:"server" yaml:"server" kv:"server"`

	// logging config
	Log struct {
		// log targets, supported: stdout, file
		Target logTarget `json:"target" yaml:"target" kv:"target"`
		// log file path
		Dir string `json:"dir" yaml:"dir" kv:"dir"`
		// rotate size for file
		Size utils.Bytes `json:"size" yaml:"size" kv:"size"`
		// rotate file count for save
		Rotate int `json:"rotate" yaml:"rotate" kv:"rotate"`
	} `json:"log" yaml:"log" kv:"log"`

	// monitor config
	Monitor struct {
		// enable report data
		Enabled bool `json:"enabled" yaml:"enabled" kv:"enabled"`
		// report interval
		Interval utils.Duration `json:"interval" yaml:"interval" kv:"interval"`
	} `json:"monitor" yaml:"monitor" kv:"monitor"`

	// limit config
	Limit limit.Configure `json:"limit" yaml:"limit" kv:"limit"`
}
