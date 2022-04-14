package clash

import (
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

type PacketFlow interface {
	WritePacket(data []byte)
}

var (
	stack core.LWIPStack
)

func startTun2Socks(flow PacketFlow, port uint16) {
	stack = core.NewLWIPStack()
	core.RegisterTCPConnHandler(socks.NewTCPHandler("127.0.0.1", port))
	core.RegisterUDPConnHandler(socks.NewUDPHandler("127.0.0.1", port, 30*time.Second))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		flow.WritePacket(data)
		return len(data), nil
	})
}
