package clash

import (
	"time"

	"github.com/Dreamacro/clash/tunnel/statistic"
)

type TrafficReceiver interface {
	ReceiveTraffic(up int64, down int64)
}

var (
	receiver TrafficReceiver
)

func SetTrafficReceiver(receive TrafficReceiver) {
	receiver = receive
}

func fetchTraffic() {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	t := statistic.DefaultManager
	for range tick.C {
		if receiver == nil {
			continue
		}
		up, down := t.Now()
		receiver.ReceiveTraffic(up, down)
	}
}
