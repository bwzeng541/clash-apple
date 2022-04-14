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
	"github.com/Dreamacro/clash/tunnel"
	T "github.com/Dreamacro/clash/tunnel"
	"github.com/Dreamacro/clash/tunnel/statistic"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

var (
	basic *config.Config
	stack core.LWIPStack
)

type PacketFlow interface {
	WritePacket(data []byte)
}

func ReadPacket(data []byte) {
	stack.Write(data)
}

func Setup(flow PacketFlow, homeDir string, config string) error {
	go fetchLogs()
	constant.SetHomeDir(homeDir)
	constant.SetConfig("")
	cfg, err := executor.ParseWithBytes(([]byte)(config))
	if err != nil {
		return err
	}
	basic = cfg
	executor.ApplyConfig(basic, true)
	stack = core.NewLWIPStack()
	core.RegisterTCPConnHandler(socks.NewTCPHandler("127.0.0.1", uint16(cfg.General.MixedPort)))
	core.RegisterUDPConnHandler(socks.NewUDPHandler("127.0.0.1", uint16(cfg.General.MixedPort), 30*time.Second))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		flow.WritePacket(data)
		return len(data), nil
	})
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
