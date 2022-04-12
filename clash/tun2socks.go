package clash

import (
	"net"

	"github.com/xjasonlyu/tun2socks/v2/component/dialer"
	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
	"github.com/xjasonlyu/tun2socks/v2/engine/mirror"
	"github.com/xjasonlyu/tun2socks/v2/proxy"

	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"github.com/Dreamacro/clash/log"
)

const (
	name = "eno0"
	addr = "127.0.0.1:8080"
)

var (
	_stack *stack.Stack
)

func StartTun2Socks(fd int32, tcpModerateReceiveBuffer bool, tcpSendBufferSize int, tcpReceiveBufferSize int) error {

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return err
	}
	dialer.DefaultInterfaceName.Store(iface.Name)
	dialer.DefaultInterfaceIndex.Store(int32(iface.Index))

	_proxy, err := proxy.NewSocks5(addr, "", "")
	if err != nil {
		return err
	}
	proxy.SetDialer(_proxy)

	_device, err := createDeviceWithTunnelFileDescriptor(fd)
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
