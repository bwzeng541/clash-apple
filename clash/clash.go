package clash

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/Dreamacro/clash/common/batch"
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
	"github.com/Dreamacro/clash/tunnel/statistic"
)

type Client interface {
	Traffic(up, down int64)
	Log(level, message string)
}

var (
	base   *config.Config
	client Client
)

func Setup(homeDir, config string, c Client) {
	client = c
	go fetchLogs()
	go fetchTraffic()
	constant.SetHomeDir(homeDir)
	constant.SetConfig("")
	cfg, err := executor.ParseWithBytes(([]byte)(config))
	if err != nil {
		panic(err)
	}
	base = cfg
	executor.ApplyConfig(base, true)
}

func SetConfig(uuid string) error {
	path := filepath.Join(constant.Path.HomeDir(), uuid, "config.yaml")
	cfg, err := executor.ParseWithPath(path)
	if err != nil {
		constant.SetConfig("")
		CloseAllConnections()
		executor.ApplyConfig(base, true)
		return err
	}
	constant.SetConfig(path)
	CloseAllConnections()
	cfg.General = base.General
	executor.ApplyConfig(cfg, false)
	return nil
}

func SetLogLevel(level string) {
	log.SetLevel(log.LogLevelMapping[level])
}

func SetTunnelMode(mode string) {
	CloseAllConnections()
	tunnel.SetMode(tunnel.ModeMapping[mode])
}

func CloseAllConnections() {
	snapshot := statistic.DefaultManager.Snapshot()
	for _, c := range snapshot.Connections {
		c.Close()
	}
}

func fetchLogs() {
	ch := make(chan log.Event, 1024)
	sub := log.Subscribe()
	defer log.UnSubscribe(sub)
	go func() {
		for elm := range sub {
			l := elm.(log.Event)
			select {
			case ch <- l:
			default:
			}
		}
		close(ch)
	}()
	for l := range ch {
		if l.LogLevel < log.Level() {
			continue
		}
		client.Log(l.Type(), l.Payload)
	}
}

func fetchTraffic() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	t := statistic.DefaultManager
	for range tick.C {
		up, down := t.Now()
		client.Traffic(up, down)
	}
}

func ProvidersData() []byte {
	providers := tunnel.Providers()
	data, _ := json.Marshal(providers)
	return data
}

func ProxiesData() []byte {
	proxies := tunnel.Proxies()
	data, _ := json.Marshal(proxies)
	return data
}

func HealthCheck(name string, url string, timeout int) {
	providers := tunnel.Providers()
	provider, exist := providers[name]
	if !exist {
		return
	}
	ps := provider.Proxies()
	proxies := make(map[string]constant.Proxy)
	for _, proxy := range ps {
		if isURLTestAdapterType(proxy.Type()) {
			proxies[proxy.Name()] = proxy
		}
	}
	if len(proxies) == 0 {
		return
	}
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(10))
	for _, proxy := range proxies {
		p := proxy
		b.Go(p.Name(), func() (any, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()
			p.URLTest(ctx, url)
			return nil, nil
		})
	}
	b.Wait()
}

func isURLTestAdapterType(at constant.AdapterType) bool {
	switch at {
	case constant.Direct:
		return false
	case constant.Reject:
		return false

	case constant.Shadowsocks:
		return true
	case constant.ShadowsocksR:
		return true
	case constant.Socks5:
		return true
	case constant.Http:
		return true
	case constant.Vmess:
		return true
	case constant.Trojan:
		return true

	case constant.Relay:
		return false
	case constant.Selector:
		return false
	case constant.Fallback:
		return false
	case constant.URLTest:
		return false
	case constant.LoadBalance:
		return false

	default:
		return false
	}
}
