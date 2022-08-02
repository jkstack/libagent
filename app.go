package agent

import (
	"context"
	"fmt"
	"os"
	rt "runtime"
	"time"

	"github.com/jkstack/agent/utils"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

type app struct {
	a         App
	shortExit int
	chRead    chan *anet.Msg
	chWrite   chan *anet.Msg

	// runtime
	ctx    context.Context
	cancel context.CancelFunc

	// monitor
	inPackets, inBytes   uint64
	outPackets, outBytes uint64
}

func newApp(a App) *app {
	ctx, cancel := context.WithCancel(context.Background())
	return &app{
		a:       a,
		chRead:  make(chan *anet.Msg, 1024*1024),
		chWrite: make(chan *anet.Msg, 1024*1024),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (app *app) start() {
	app.initLogging()
	defer logging.Flush()

	defer utils.Recover("app_start")

	if app.a.Configure().Monitor.Enabled {
		go app.report()
	}
	go app.readCallback()
	go app.a.LoopWrite(app.ctx, app.chWrite)
	go app.debug(app.ctx)

	var nextSleep time.Duration
	for {
		if nextSleep > 0 {
			time.Sleep(nextSleep)
			nextSleep <<= 1
			if nextSleep > 30*time.Second {
				nextSleep = 30 * time.Second
			}
		}
		if app.shortExit > 20 {
			logging.Error("short exit more than 20 times")
			logging.Flush()
			os.Exit(255)
		}
		conn, err := app.connect()
		if err != nil {
			app.shortExit++
			continue
		}

		ctx, cancel := context.WithCancel(app.ctx)
		go app.read(ctx, cancel, conn)
		go app.write(ctx, cancel, conn)
		go app.keepalive(ctx, conn)

		<-ctx.Done()
		conn.Close()
	}
}

func (app *app) stop() {
	app.cancel()
}

func (app *app) initLogging() {
	var cfg logging.SizeRotateConfig
	cfg.WriteStdout = app.a.Configure().Log.Target.SupportStdout()
	cfg.WriteFile = app.a.Configure().Log.Target.SupportFile()
	if rt.GOOS == "windows" {
		cfg.WriteStdout = false
	}
	if !cfg.WriteStdout && !cfg.WriteFile {
		fmt.Printf("[WARN]no log target set, default to stdout")
		cfg.WriteStdout = true
	}
	cfg.Dir = app.a.Configure().Log.Dir
	cfg.Name = app.a.AgentName()
	cfg.Size = int64(app.a.Configure().Log.Size.Bytes())
	cfg.Rotate = app.a.Configure().Log.Rotate
	logging.SetSizeRotate(cfg)
}

func (app *app) debug(ctx context.Context) {
	tk := time.NewTicker(10 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tk.C:
			if len(app.chRead) > 0 {
				logging.Info("read channel size: %d", len(app.chRead))
			}
			if len(app.chWrite) > 0 {
				logging.Info("write channel size: %d", len(app.chWrite))
			}
		}
	}
}
