package clash

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dreamacro/clash/common/batch"
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
	data, _ := json.Marshal(proxies)
	return data
}

func URLTest(names []string, url string, timeout int) {
	if basic == nil {
		return
	}
	if len(names) == 0 {
		return
	}
	ps := tunnel.Proxies()
	proxies := make(map[string]constant.Proxy)
	for _, name := range names {
		proxy, exist := ps[name]
		if exist {
			continue
		}
		proxies[name] = proxy
	}
	if len(proxies) == 0 {
		return
	}
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(10))
	for _, proxy := range proxies {
		p := proxy
		b.Go(p.Name(), func() (any, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
			defer cancel()
			p.URLTest(ctx, url)
			return nil, nil
		})
	}
	b.Wait()
}
