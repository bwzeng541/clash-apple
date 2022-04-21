package clash

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dreamacro/clash/tunnel"
)

type _URLTestRequest struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Timeout int64  `json:"timeout"`
}

type _URLTestResponse struct {
	Name  string `json:"name"`
	Delay int64  `json:"delay"`
}

func ABCURLTest() ([]byte, error) {
	return nil, nil
}

func URLTest(request []byte) []byte {
	var (
		req  *_URLTestRequest
		resp _URLTestResponse
	)
	for {
		req = new(_URLTestRequest)
		json.Unmarshal(request, req)

		proxies := tunnel.Proxies()
		proxy, exist := proxies[req.Name]

		if !exist {
			resp = _URLTestResponse{
				Name:  req.Name,
				Delay: -1,
			}
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(req.Timeout))
		defer cancel()

		delay, err := proxy.URLTest(ctx, req.URL)
		if ctx.Err() != nil {
			resp = _URLTestResponse{
				Name:  req.Name,
				Delay: -2,
			}
			break
		}

		if err != nil || delay == 0 {
			resp = _URLTestResponse{
				Name:  req.Name,
				Delay: -3,
			}
			break
		}

		resp = _URLTestResponse{
			Name:  req.Name,
			Delay: int64(delay),
		}
		break
	}
	data, _ := json.Marshal(resp)
	return data
}
