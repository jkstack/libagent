package agent

import (
	"fmt"
	rt "runtime"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/jkstack/anet"
)

func (app *app) report() {
	for {
		time.Sleep(app.a.Configure().Monitor.Interval.Duration())
		if !app.connected {
			continue
		}
		app.a.OnReportMonitor()
		app.chWrite <- app.buildReport()
	}
}

func (app *app) buildReport() *anet.Msg {
	var msg anet.Msg
	msg.Type = anet.TypeAgentInfo

	var info anet.AgentInfo
	info.Version = app.a.Version()
	info.GoVersion = rt.Version()

	cpu, _ := app.process.CPUPercent()
	info.CpuUsage = float32(cpu)
	mem, _ := app.process.MemoryPercent()
	info.MemoryUsage = mem

	n, _ := rt.ThreadCreateProfile(nil)
	info.Threads = n
	info.Routines = rt.NumGoroutine()
	info.Startup = app.startup

	var stats rt.MemStats
	rt.ReadMemStats(&stats)
	info.HeapInuse = stats.HeapInuse

	var gc debug.GCStats
	gc.PauseQuantiles = make([]time.Duration, 5)
	debug.ReadGCStats(&gc)

	quantiles := make(map[string]float64)
	for i := 0; i < 5; i++ {
		quantiles[fmt.Sprintf("%d", i*25)] = gc.PauseQuantiles[i].Seconds()
	}
	info.GC = quantiles

	info.InPackets = app.inPackets
	info.InBytes = app.inBytes
	info.OutPackets = app.outPackets
	info.OutBytes = app.outBytes
	info.ReconnectCount = app.reconnectCount

	info.ReadChanSize = len(app.chRead)
	info.WriteChanSize = len(app.chWrite)

	msg.AgentInfo = &info
	return &msg
}

func (app *app) incInPackets() {
	atomic.AddUint64(&app.inPackets, 1)
}

func (app *app) incInBytes(n uint64) {
	atomic.AddUint64(&app.inBytes, n)
}

func (app *app) incOutPackets() {
	atomic.AddUint64(&app.outPackets, 1)
}

func (app *app) incOutBytes(n uint64) {
	atomic.AddUint64(&app.outBytes, n)
}
