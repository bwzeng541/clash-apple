package clash

import (
	"net"

	"github.com/xjasonlyu/tun2socks/v2/component/dialer"
	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/device"
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

func SetupTun2Socks(fd int32) error {

	iface, err := net.InterfaceByName("en0")
	if err != nil {
		return err
	}
	dialer.DefaultInterfaceName.Store(iface.Name)
	dialer.DefaultInterfaceIndex.Store(int32(iface.Index))

	_proxy, err := proxy.NewSocks5("127.0.0.1:8080", "", "")
	if err != nil {
		return err
	}
	proxy.SetDialer(_proxy)

	_device, err = Open(fd)
	if err != nil {
		return err
	}

	var opts []option.Option

	_stack, err = core.CreateStack(&core.Config{
		LinkEndpoint:     _device,
		TransportHandler: &mirror.Tunnel{},
		PrintFunc: func(format string, v ...any) {
			log.Warnln(format, v...)
		},
		Options: opts,
	})

	log.Warnln(
		"[STACK] %s://%s <-> %s://%s",
		_device.Type(), _device.Name(),
		_proxy.Proto(), _proxy.Addr(),
	)

	return nil
}
