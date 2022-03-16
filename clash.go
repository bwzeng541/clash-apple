package clash

import (
	"path/filepath"
	"time"

	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	L "github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
	"github.com/Dreamacro/clash/tunnel/statistic"
)

var (
	trafficReceiver TrafficReceiver
	logger          RealTimeLogger
	primaryConfig   *config.Config
)

type TrafficReceiver interface {
	ReceiveTraffic(up int64, down int64)
}

type RealTimeLogger interface {
	Log(level string, payload string)
}

func Setup(homeDir string, config string) error {
	go fetchLogs()
	constant.SetHomeDir(homeDir)
	constant.SetConfig("")
	cfg, err := executor.ParseWithBytes(([]byte)(config))
	if err != nil {
		return err
	}
	primaryConfig = cfg
	executor.ApplyConfig(primaryConfig, true)
	go fetchTraffic()
	return nil
}

func SetConfig(uuid string) error {
	if primaryConfig == nil {
		return nil
	}
	path := filepath.Join(constant.Path.HomeDir(), uuid, "config.yaml")
	cfg, err := executor.ParseWithPath(path)
	if err != nil {
		return err
	}
	constant.SetConfig(path)
	CloseAllConnections()
	cfg.General = primaryConfig.General
	cfg.DNS = primaryConfig.DNS
	executor.ApplyConfig(cfg, false)
	return nil
}

func SetTunnelMode(mode string) {
	if primaryConfig == nil {
		return
	}
	CloseAllConnections()
	T.SetMode(T.ModeMapping[mode])
}

func CloseAllConnections() {
	snapshot := statistic.DefaultManager.Snapshot()
	for _, c := range snapshot.Connections {
		c.Close()
	}
}

func SetTrafficReceiver(receive TrafficReceiver) {
	trafficReceiver = receive
}

func fetchTraffic() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	t := statistic.DefaultManager
	for range tick.C {
		if trafficReceiver == nil {
			continue
		}
		up, down := t.Now()
		trafficReceiver.ReceiveTraffic(up, down)
	}
}

func SetLogLevel(level string) {
	if primaryConfig == nil {
		return
	}
	L.SetLevel(L.LogLevelMapping[level])
}

func SetRealTimeLogger(l RealTimeLogger) {
	logger = l
}

func fetchLogs() {
	sub := L.Subscribe()
	defer L.UnSubscribe(sub)
	for elm := range sub {
		if logger == nil {
			continue
		}
		log := elm.(*L.Event)
		if log.LogLevel < L.Level() {
			continue
		}
		logger.Log(log.Type(), log.Payload)
	}
}
