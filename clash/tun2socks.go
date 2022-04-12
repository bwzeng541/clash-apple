package clash

import (
	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/device"
	"github.com/xjasonlyu/tun2socks/v2/core/device/tun"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
	"github.com/xjasonlyu/tun2socks/v2/engine/mirror"
	"github.com/xjasonlyu/tun2socks/v2/proxy"

	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/Dreamacro/clash/log"
)

var (
	_device device.Device
	_stack  *stack.Stack
)

func SetupTun2Socks(fd int32, tcpModerateReceiveBuffer bool, tcpSendBufferSize int, tcpReceiveBufferSize int) error {

	_proxy, err := proxy.NewSocks5("127.0.0.1:8080", "", "")
	if err != nil {
		return err
	}
	proxy.SetDialer(_proxy)

	_device, err := tun.Open(fd)
	if err != nil {
		return err
	}

	var opts []option.Option
	if tcpModerateReceiveBuffer {
		opts = append(opts, option.WithTCPModerateReceiveBuffer(true))
	}
	if tcpSendBufferSize > 0 {
		opts = append(opts, option.WithTCPSendBufferSize(int(tcpSendBufferSize)))
	}
	if tcpReceiveBufferSize > 0 {
		opts = append(opts, option.WithTCPReceiveBufferSize(int(tcpReceiveBufferSize)))
	}
	_stack, err = core.CreateStack(&core.Config{
		LinkEndpoint:     _device,
		TransportHandler: &mirror.Tunnel{},
		PrintFunc: func(format string, v ...any) {
			log.Warnln(format, v...)
		},
		Options: opts,
	})

	if err != nil {
		return err
	}

	return nil
}
