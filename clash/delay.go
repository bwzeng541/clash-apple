package clash

import (
	"encoding/json"

	"github.com/Dreamacro/clash/tunnel"
)

func HealthCheck() {
	if basic == nil {
		return
	}
	providers := tunnel.Providers()
	for _, provider := range providers {
		provider.HealthCheck()
	}
}

func Proxies() []byte {
	if basic == nil {
		return nil
	}
	proxies := tunnel.Proxies()
	data, _ := json.Marshal(proxies)
	return data
}

func Providers() []byte {
	if basic == nil {
		return nil
	}
	providers := tunnel.Providers()
	data, _ := json.Marshal(providers)
	return data
}

func Provider(name string) []byte {
	if basic == nil {
		return nil
	}
	providers := tunnel.Providers()
	provider, exist := providers[name]
	if !exist {
		return nil
	}
	data, _ := json.Marshal(provider)
	return data
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
