package agent

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGetblockID(t *testing.T) {
	blocks, _ := filepath.Glob("/sys/block/*")
	for _, block := range blocks {
		data, _ := ioutil.ReadFile(filepath.Join(block, "dev"))
		fmt.Println(filepath.Base(block), string(data))
	}
}
