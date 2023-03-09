package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	rt "runtime"
	"sync"
	"time"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/libagent/internal/utils"
	"github.com/shirou/gopsutil/v3/process"
)

type app struct {
	a       App
	startup int64
	process *process.Process
	chRead  chan *anet.Msg
	chWrite chan *anet.Msg

	// runtime
	ctx        context.Context
	cancel     context.CancelFunc
	connected  bool
	mWriteLock sync.Mutex

	// monitor
	inPackets, inBytes   uint64
	outPackets, outBytes uint64
	reconnectCount       int
}

func newApp(a App) *app {
	ctx, cancel := context.WithCancel(context.Background())
	p, _ := process.NewProcess(int32(os.Getpid()))
	return &app{
		a:       a,
		process: p,
		startup: time.Now().Unix(),
		chRead:  make(chan *anet.Msg, 1024*1024),
		chWrite: make(chan *anet.Msg, 1024*1024),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (app *app) start() {
	app.rel2abs()

	app.initLogging()
	defer logging.Flush()

	app.a.Configure().Limit.Do(app.a.AgentName())

	defer utils.Recover("app_start")

	if app.a.Configure().Monitor.Enabled {
		go app.report()
	}
	go app.readCallback()
	go app.a.LoopWrite(app.ctx, app.chWrite)
	go app.debug(app.ctx)

	i := 0
	nextSleep := time.Second
	for {
		i++
		if i > 1 {
			time.Sleep(nextSleep)
			nextSleep <<= 1
			if nextSleep > 5*time.Second {
				nextSleep = 5 * time.Second
			}
		}

		conn, err := app.connect()
		if err != nil {
			continue
		}

		app.connected = true
		deferCallback("on_connect", app.a.OnConnect)

		ctx, cancel := context.WithCancel(app.ctx)
		go app.read(ctx, cancel, conn)
		go app.write(ctx, cancel, conn)
		go app.keepalive(ctx, conn)

		<-ctx.Done()

		app.connected = false
		conn.Close()

		deferCallback("dis_connect", app.a.OnDisconnect)
		app.reconnectCount++
		nextSleep = time.Second
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
		fmt.Printf("[WARN]no log target set, default to stdout\n")
		cfg.WriteStdout = true
	}
	cfg.Dir = app.a.Configure().Log.Dir
	if !filepath.IsAbs(cfg.Dir) {
		dir, err := os.Executable()
		if err != nil {
			fmt.Printf("can not change log.dir to absolute path: %v\n", err)
		} else {
			cfg.Dir = filepath.Join(filepath.Dir(dir), cfg.Dir)
		}
	}
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

func (app *app) rel2abs() {
	exe, err := os.Executable()
	if err != nil {
		logging.Error("get executable file dir: %v", err)
		return
	}
	trans := func(dir string) string {
		if !filepath.IsAbs(dir) {
			return filepath.Join(filepath.Dir(exe), dir)
		}
		return dir
	}
	cfg := app.a.Configure()
	if len(cfg.Log.Dir) > 0 {
		cfg.Log.Dir = trans(cfg.Log.Dir)
	}
}
