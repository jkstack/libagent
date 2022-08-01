package agent

// Config
type Config struct {
	// agent id, if empty generate by server
	//   $HOSTNAME: use agent id by hostname
	//   $IP: use agent id by ip of connect server
	//   ${env}: use envionment variables
	ID string `json:"id" yaml:"id" kv:"id"`

	// server address
	Server string `json:"server" yaml:"server" kv:"server"`

	// logging config
	Log struct {
		Target []string `json:"target" yaml:"target" kv:"target"`
		Dir    string   `json:"dir" yaml:"dir" kv:"dir"`
	} `json:"log" yaml:"log" kv:"log"`
}
