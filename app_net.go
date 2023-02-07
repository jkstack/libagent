package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/libagent/internal/hostinfo"
	"github.com/jkstack/libagent/internal/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var dialer = websocket.Dialer{
	EnableCompression: true,
}

func (app *app) connect() (*websocket.Conn, error) {
	app.a.Configure().RewriteServer()
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s/ws/agent", app.a.Configure().Server), nil)
	if err != nil {
		logging.Error("dial: %v", err)
		return nil, err
	}

	app.sendCome(conn)
	logging.Info("send come message on server %s", app.a.Configure().Server)
	err = app.waitHandshake(conn, time.Minute)
	if err != nil {
		conn.Close()
		logging.Error("wait handshake: %v", err)
		return nil, err
	}
	logging.Info("wait handshake ok")
	return conn, nil
}

func (app *app) sendCome(conn *websocket.Conn) {
	var msg anet.Msg
	msg.Type = anet.TypeCome
	ip := utils.GetIP(conn)
	hostName, _ := os.Hostname()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cpus, err := cpu.InfoWithContext(ctx)
	if err != nil {
		cpus = []cpu.InfoStat{
			{ModelName: "unknown", Cores: 0},
		}
	}
	memory, _ := mem.VirtualMemory()
	host, err := hostinfo.Info()
	if err != nil {
		logging.Error("get host info failed, err=%v", err)
		return
	}
	app.a.Configure().RewriteID()
	msg.Come = &anet.ComePayload{
		ID:              app.a.Configure().ID,
		Name:            app.a.AgentName(),
		Version:         app.a.Version(),
		IP:              ip,
		MAC:             utils.GetMac(ip),
		HostName:        hostName,
		OS:              host.OS,
		Platform:        host.Platform,
		PlatformVersion: host.PlatformVersion,
		KernelVersion:   host.KernelVersion,
		Arch:            host.KernelArch,
		CPU:             cpus[0].ModelName,
		CPUCore:         uint64(cpus[0].Cores),
	}
	if memory != nil {
		msg.Come.Memory = memory.Total
	}
	err = conn.WriteJSON(msg)
	if err != nil {
		logging.Error("send come message: %v", err)
		return
	}
}

func (app *app) waitHandshake(conn *websocket.Conn, timeout time.Duration) error {
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{})
	var msg anet.Msg
	err := conn.ReadJSON(&msg)
	if err != nil {
		return err
	}
	if msg.Type != anet.TypeHandshake {
		return fmt.Errorf("unexpected message type(handshake): %d", msg.Type)
	}
	if !msg.Handshake.OK {
		return errors.New(msg.Handshake.Msg)
	}
	if len(msg.Handshake.ID) > 0 && msg.Handshake.ID != app.a.Configure().ID {
		logging.Info("agent_id reset rewrite configure file...")
		app.a.Configure().SetAgentID(msg.Handshake.ID)
		app.a.OnRewriteConfigure()
	}
	return nil
}
