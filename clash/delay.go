package clash

import (
	"context"
	"time"

	"github.com/Dreamacro/clash/tunnel"
)

type URLTestRequest struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}

type URLTestResponse struct {
	Name  string `json:"name"`
	Delay int64  `json:"delay"`
}

func URLTest(request URLTestRequest) URLTestResponse {

	proxies := tunnel.Proxies()
	proxy, exist := proxies[request.Name]

	if !exist {
		return URLTestResponse{
			Name:  request.Name,
			Delay: -1,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(request.Timeout))
	defer cancel()

	delay, err := proxy.URLTest(ctx, request.URL)
	if ctx.Err() != nil {
		return URLTestResponse{
			Name:  request.Name,
			Delay: -2,
		}
	}

	if err != nil || delay == 0 {
		return URLTestResponse{
			Name:  request.Name,
			Delay: -3,
		}
	}

	return URLTestResponse{
		Name:  request.Name,
		Delay: int64(delay),
	}
}
