package tun

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/Kr328/tun2socket"
	"golang.org/x/sys/unix"

	"github.com/Dreamacro/clash/adapter/inbound"
	"github.com/Dreamacro/clash/common/pool"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/context"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/transport/socks5"
	"github.com/Dreamacro/clash/tunnel"
)

var _, ipv4LoopBack, _ = net.ParseCIDR("127.0.0.0/8")

func StartTun2Socket(fd int, gateway, portal string) (io.Closer, error) {

	log.Debugln("TUN: fd = %d, gateway = %s, portal = %s", fd, gateway, portal)

	dupTunFd, err := unix.Dup(int(fd))
	if err != nil {
		return nil, err
	}

	err = unix.SetNonblock(dupTunFd, true)
	if err != nil {
		unix.Close(dupTunFd)
		return nil, err
	}

	device := os.NewFile(uintptr(dupTunFd), "/dev/tun")

	ip, network, err := net.ParseCIDR(gateway)
	if err != nil {
		panic(err.Error())
	} else {
		network.IP = ip
	}

	stack, err := tun2socket.StartTun2Socket(device, network, net.ParseIP(portal))
	if err != nil {
		_ = device.Close()
		return nil, err
	}

	tcp := func() {

		defer stack.TCP().Close()
		defer log.Debugln("TCP: closed")

		for stack.TCP().SetDeadline(time.Time{}) == nil {
			conn, err := stack.TCP().Accept()
			if err != nil {
				log.Debugln("Accept connection: %v", err)
				continue
			}

			lAddr := conn.LocalAddr().(*net.TCPAddr)
			rAddr := conn.RemoteAddr().(*net.TCPAddr)

			if ipv4LoopBack.Contains(rAddr.IP) {
				conn.Close()
				continue
			}

			tunnel.TCPIn() <- context.NewConnContext(conn, createMetadata(lAddr, rAddr))
		}
	}

	udp := func() {

		defer stack.UDP().Close()
		defer log.Debugln("UDP: closed")

		for {
			buf := pool.Get(pool.UDPBufferSize)

			n, lRAddr, rRAddr, err := stack.UDP().ReadFrom(buf)
			if err != nil {
				return
			}

			raw := buf[:n]
			lAddr := lRAddr.(*net.UDPAddr)
			rAddr := rRAddr.(*net.UDPAddr)

			if ipv4LoopBack.Contains(rAddr.IP) {
				pool.Put(buf)
				continue
			}

			pkt := &packet{
				local: lAddr,
				data:  raw,
				writeBack: func(b []byte, addr net.Addr) (int, error) {
					return stack.UDP().WriteTo(b, addr, lAddr)
				},
				drop: func() {
					pool.Put(buf)
				},
			}

			tunnel.UDPIn() <- inbound.NewPacket(socks5.ParseAddrToSocksAddr(rAddr), pkt, constant.SOCKS5)
		}
	}

	go tcp()
	go udp()

	return stack, nil
}
