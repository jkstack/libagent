package agent

import (
	"sync/atomic"
	"time"
)

func (app *app) report() {
	for {
		time.Sleep(app.a.Configure().Monitor.Interval.Duration())
	}
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
