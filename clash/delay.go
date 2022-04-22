package clash

import (
	"context"
	"time"

	"github.com/Dreamacro/clash/common/batch"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"
)

var (
	_ticker   *time.Ticker
	_callback URLTestCallback
)

type URLTestCallback interface {
	URLTest(name string, delay int64)
}

func URLTests(callback URLTestCallback, duration int64) {
	_callback = callback
	if _ticker != nil {
		_ticker.Stop()
	}
	_URLTests()
	_ticker = time.NewTicker(time.Duration(duration) * time.Second)
	for range _ticker.C {
		_URLTests()
		log.Warnln("-----------> providers : %d", len(tunnel.Providers()))
		for _, provider := range tunnel.Providers() {
			provider.HealthCheck()
		}
	}
}

func _URLTests() {
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(10))
	proxies := tunnel.Proxies()
	for _, proxy := range proxies {
		p := proxy
		b.Go(p.Name(), func() (any, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			delay, err := p.URLTest(ctx, "http://www.gstatic.com/generate_204")
			if ctx.Err() != nil {
				_callback.URLTest(p.Name(), 0)
			} else if err != nil || delay == 0 {
				_callback.URLTest(p.Name(), -1)
			} else {
				_callback.URLTest(p.Name(), int64(delay))
			}
			return nil, nil
		})
	}
	b.Wait()
}
