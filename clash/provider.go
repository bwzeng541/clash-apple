package clash

import (
	"encoding/json"

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
	data, _ := json.Marshal(proxies)
	return data
}
