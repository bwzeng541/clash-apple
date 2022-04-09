package clash

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	L "github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
	T "github.com/Dreamacro/clash/tunnel"
	"github.com/Dreamacro/clash/tunnel/statistic"
	"github.com/xjasonlyu/tun2socks/v2/engine"
)

var (
	key      = new(engine.Key)
	receiver TrafficReceiver
	logger   RealTimeLogger
	basic    *config.Config
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

type TrafficReceiver interface {
	ReceiveTraffic(up int64, down int64)
}

type RealTimeLogger interface {
	Log(level string, payload string)
}

func ReadPacket(data []byte) {

}

func Setup(name string, homeDir string, config string) error {
	go fetchLogs()
	constant.SetHomeDir(homeDir)
	constant.SetConfig("")
	cfg, err := executor.ParseWithBytes(([]byte)(config))
	if err != nil {
		return err
	}
	basic = cfg
	executor.ApplyConfig(basic, true)

	key.Device = name
	key.Proxy = "socks5://127.0.0.1:8080"
	key.Interface = "en0"
	key.MTU = 1500
	engine.Insert(key)
	engine.Start()

	go fetchTraffic()
	return nil
}

func SetConfig(uuid string) error {
	if basic == nil {
		return nil
	}
	path := filepath.Join(constant.Path.HomeDir(), uuid, "config.yaml")
	cfg, err := executor.ParseWithPath(path)
	if err != nil {
		return err
	}
	constant.SetConfig(path)
	CloseAllConnections()
	cfg.General = basic.General
	cfg.DNS.Enable = false
	cfg.Profile.StoreSelected = false
	executor.ApplyConfig(cfg, false)
	return nil
}

func PatchSelectGroup(data []byte) {
	if basic == nil {
		return
	}
	mapping := make(map[string]string)
	err := json.Unmarshal(data, &mapping)
	if err != nil {
		return
	}
	proxies := tunnel.Proxies()
	for name, proxy := range proxies {
		selected, exist := mapping[name]
		if !exist {
			continue
		}
		outbound, ok := proxy.(*adapter.Proxy)
		if !ok {
			continue
		}
		selector, ok := outbound.ProxyAdapter.(*outboundgroup.Selector)
		if !ok {
			continue
		}
		selector.Set(selected)
	}
}

func SetTunnelMode(mode string) {
	if basic == nil {
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
	receiver = receive
}

func fetchTraffic() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	t := statistic.DefaultManager
	for range tick.C {
		if receiver == nil {
			continue
		}
		up, down := t.Now()
		receiver.ReceiveTraffic(up, down)
	}
}

func SetLogLevel(level string) {
	if basic == nil {
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
