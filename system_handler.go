package agent

import (
	"archive/zip"
	"crypto/md5"
	"io"
	"os"
	"path/filepath"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/compress"
	"github.com/jkstack/jkframe/logging"
)

const blockSize = 32_000

func (app *app) handleSystemPacket(msg *anet.Msg) bool {
	switch msg.Type {
	case anet.TypeLogLsReq:
		app.handleLogLs(msg.TaskID)
		return false
	case anet.TypeLogDownloadReq:
		go app.downloadLogs(msg.TaskID, msg.LogDownload.Files)
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

func (app *app) downloadLogs(id string, files []string) {
	tmpDir, err := zipFiles(app.a.Configure().Log.Dir, files)
	if len(tmpDir) != 0 {
		defer os.Remove(tmpDir)
	}
	if err != nil {
		app.downloadLogFailed(id, err.Error())
		return
	}
	app.downloadLogOK(id, tmpDir)
}

func zipFiles(dir string, files []string) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "log")
	if err != nil {
		return "", err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()
	for _, file := range files {
		src, err := os.Open(filepath.Join(dir, file))
		if err != nil {
			return "", err
		}
		defer src.Close()
		fi, err := src.Stat()
		if err != nil {
			return "", err
		}
		hdr, err := zip.FileInfoHeader(fi)
		if err != nil {
			return "", err
		}
		dst, err := zw.CreateHeader(hdr)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return "", err
		}
	}
	return f.Name(), nil
}

func (app *app) downloadLogFailed(id, msg string) {
	var m anet.Msg
	m.Type = anet.TypeLogDownloadInfo
	m.TaskID = id
	m.LogDownloadInfo = &anet.LogDownloadInfo{
		OK:     false,
		ErrMsg: msg,
	}
	app.chWrite <- &m
}

func (app *app) downloadLogOK(id, dir string) {
	err := app.responseLogInfo(id, dir)
	if err != nil {
		logging.Error("response log info: %v", err)
		return
	}
	app.responseLogData(id, dir)
}

func (app *app) responseLogInfo(id, dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	enc := md5.New()
	_, err = io.Copy(enc, f)
	if err != nil {
		return err
	}
	var sum [md5.Size]byte
	copy(sum[:], enc.Sum(nil))
	var msg anet.Msg
	msg.Type = anet.TypeLogDownloadInfo
	msg.TaskID = id
	msg.LogDownloadInfo = &anet.LogDownloadInfo{
		OK:        true,
		Size:      uint64(fi.Size()),
		BlockSize: blockSize,
		MD5:       sum,
	}
	app.chWrite <- &msg
	return nil
}

func (app *app) responseLogData(id, dir string) {
	f, err := os.Open(dir)
	if err != nil {
		logging.Error("open file: %v", err)
		return
	}
	buf := make([]byte, blockSize)
	var msg anet.Msg
	msg.Type = anet.TypeLogDownloadData
	msg.TaskID = id
	var offset uint64
	for {
		n, err := f.Read(buf)
		if err != nil {
			logging.Error("read file: %v", err)
			return
		}
		msg.LogDownloadData = &anet.LogDownloadData{
			Offset: offset,
			Data:   compress.Compress(buf[:n]),
		}
		app.chWrite <- &msg
	}
}
