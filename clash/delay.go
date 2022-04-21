package clash

import (
	"context"
	"time"

	"github.com/Dreamacro/clash/tunnel"
)

func URLTest(name string, url string, timeout int64) int64 {

	proxies := tunnel.Proxies()
	proxy, exist := proxies[name]

	if !exist {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
	defer cancel()

	delay, err := proxy.URLTest(ctx, url)
	if ctx.Err() != nil {
		return -1
	}
	if err != nil || delay == 0 {
		return -2
	}

	return int64(delay)
}
