package agent

import (
	"os"
	"path/filepath"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

func (app *app) handleSystemPacket(msg *anet.Msg) bool {
	switch msg.Type {
	case anet.TypeLogLsReq:
		app.handleLogLs(msg.TaskID)
		return false
	}
	return true
}

func (app *app) handleLogLs(id string) {
	files := logging.Files()
	var fs []anet.LogFile
	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			continue
		}
		fs = append(fs, anet.LogFile{
			Name:    filepath.Base(file),
			Size:    uint64(fi.Size()),
			ModTime: fi.ModTime(),
		})
	}
	var msg anet.Msg
	msg.Type = anet.TypeLogLsRep
	msg.TaskID = id
	msg.LsLog = &anet.LsLogPayload{
		Files: fs,
	}
	app.chWrite <- &msg
}
