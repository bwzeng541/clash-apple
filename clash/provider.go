package clash

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dreamacro/clash/common/batch"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/tunnel"
)

func MergedProxyData() []byte {
	if basic == nil {
		return nil
	}
	proxies := tunnel.Proxies()
	providers := tunnel.Providers()
	mapping := make(map[string]any)
	mapping["proxies"] = proxies
	mapping["providers"] = providers
	data, _ := json.Marshal(mapping)
	return data
}

func PatchData() []byte {
	if basic == nil {
		return nil
	}
	proxies := tunnel.Proxies()
	data, _ := json.Marshal(proxies)
	return data
}

func HealthCheck(name string, url string, timeout int) {
	if basic == nil {
		return
	}
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
