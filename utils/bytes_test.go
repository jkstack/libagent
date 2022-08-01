package utils

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestJSON(t *testing.T) {
	bytes := Bytes(1024 * 1024)
	data, err := json.Marshal(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"1.0 MiB"` {
		t.Fatal("unexpected value")
	}
	err = json.Unmarshal(data, &bytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Bytes() != 1024*1024 {
		t.Fatal("unexpected value")
	}
}

func TestYAML(t *testing.T) {
	bytes := Bytes(1024 * 1024)
	data, err := yaml.Marshal(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "1.0 MiB\n" {
		t.Fatal("unexpected value")
	}
	bytes = Bytes(0)
	err = yaml.Unmarshal(data, &bytes)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Bytes() != 1024*1024 {
		t.Fatal("unexpected value")
	}
}
