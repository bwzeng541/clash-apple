package clash

import (
	"encoding/json"

	"github.com/Dreamacro/clash/adapter"
	"github.com/Dreamacro/clash/adapter/outboundgroup"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/tunnel"
)

func HealthChecks() {
	if basic == nil {
		return
	}
	providers := tunnel.Providers()
	for _, provider := range providers {
		provider.HealthCheck()
	}
}

func HealthCheck(name string) {
	if basic == nil {
		return
	}
	providers := tunnel.Providers()
	provider, exist := providers[name]
	if !exist {
		return
	}
	provider.HealthCheck()
}

func MergedProxyData() []byte {
	if basic == nil {
		return nil
	}
	proxies := tunnel.Proxies()
	providers := tunnel.Providers()
	mapping := make(map[string]interface{})
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
	mapping := make(map[string]interface{})
	for _, proxy := range proxies {
		temp := make(map[string]interface{})
		temp["current"] = func(proxy constant.Proxy) string {
			outbound, ok := proxy.(*adapter.Proxy)
			if !ok {
				return ""
			}
			selector, ok := outbound.ProxyAdapter.(*outboundgroup.Selector)
			if !ok {
				return ""
			}
			return selector.Now()
		}(proxy)
		temp["histories"] = proxy.DelayHistory()
		mapping[proxy.Name()] = temp
	}
	data, _ := json.Marshal(mapping)
	return data
}
