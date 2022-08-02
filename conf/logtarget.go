package conf

import "strings"

type logTarget []string

// MarshalKV marshal log target
func (value logTarget) MarshalKV() (string, error) {
	return strings.Join([]string(value), ","), nil
}

// UnmarshalKV unmarshal log target
func (value *logTarget) UnmarshalKV(data string) error {
	for _, v := range strings.Split(data, ",") {
		if v == "stdout" || v == "file" {
			*value = append(*value, v)
		}
	}
	return nil
}

// SupportStdout check is supported stdout
func (value logTarget) SupportStdout() bool {
	for _, v := range value {
		if v == "stdout" {
			return true
		}
	}
	return false
}

// SupportFile check is supported file
func (value logTarget) SupportFile() bool {
	for _, v := range value {
		if v == "file" {
			return true
		}
	}
	return false
}
