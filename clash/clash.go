package clash

import (
	"encoding/json"
	"path/filepath"

	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
	T "github.com/Dreamacro/clash/tunnel"
	"github.com/Dreamacro/clash/tunnel/statistic"
)

var (
	basic *config.Config
)

func Setup(homeDir string, config string) error {
	go fetchLogs()
	constant.SetHomeDir(homeDir)
	constant.SetConfig("")
	cfg, err := executor.ParseWithBytes(([]byte)(config))
	if err != nil {
		return err
	}
	basic = cfg
	executor.ApplyConfig(basic, true)
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
		log.Infoln("patch proxy name %s", name)
		selected, exist := mapping[name]
		if !exist {
			log.Warnln("patch proxy name not found: %s", name)
			continue
		}
		outbound, ok := proxy.(*adapter.Proxy)
		if !ok {
			log.Warnln("patch proxy name convert failed: %s", name)
			continue
		}
		selector, ok := outbound.ProxyAdapter.(*outboundgroup.Selector)
		if !ok {
			log.Warnln("patch proxy name no selector: %s", name)
			continue
		}
		err := selector.Set(selected)
		if err != nil {
			log.Warnln("patch proxy failed: %s", err.Error())
		}
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
