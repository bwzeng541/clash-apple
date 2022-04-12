package clash

import (
	"fmt"
	"os"

	"github.com/xjasonlyu/tun2socks/v2/core/device"
	"github.com/xjasonlyu/tun2socks/v2/core/device/iobased"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/tun"
)

const (
	offset     = 4
	defaultMTU = 1500
)

type appleTun struct {
	*iobased.Endpoint
	nt     *tun.NativeTun
	offset int
}

func createDeviceWithTunnelFileDescriptor(fd int32) (_ device.Device, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("open tun: %v", r)
		}
	}()

	t := &appleTun{offset: offset}

	dupTunFd, err := unix.Dup(int(fd))
	if err != nil {
		return nil, err
	}

	err = unix.SetNonblock(dupTunFd, true)
	if err != nil {
		unix.Close(dupTunFd)
		return nil, err
	}

	nt, err := tun.CreateTUNFromFile(os.NewFile(uintptr(dupTunFd), "/dev/tun"), defaultMTU)
	if err != nil {
		unix.Close(dupTunFd)
		return nil, err
	}

	t.nt = nt.(*tun.NativeTun)

	mtu, err := nt.MTU()
	if err != nil {
		unix.Close(dupTunFd)
		return nil, fmt.Errorf("get mtu: %w", err)
	}

	ep, err := iobased.New(t, uint32(mtu), offset)
	if err != nil {
		unix.Close(dupTunFd)
		return nil, fmt.Errorf("create endpoint: %w", err)
	}
	t.Endpoint = ep

	return t, nil
}

func (t *appleTun) Read(packet []byte) (int, error) {
	return t.nt.Read(packet, t.offset)
}

func (t *appleTun) Write(packet []byte) (int, error) {
	return t.nt.Write(packet, t.offset)
}

func (t *appleTun) Name() string {
	name, _ := t.nt.Name()
	return name
}

func (t *appleTun) Close() error {
	defer t.Endpoint.Close()
	return t.nt.Close()
}

func (t *appleTun) Type() string {
	return "tun"
}
