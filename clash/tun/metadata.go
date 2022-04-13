package tun

import (
	"net"
	"strconv"

	"github.com/Dreamacro/clash/constant"
)

func createMetadata(lAddr, rAddr *net.TCPAddr) *constant.Metadata {
	return &constant.Metadata{
		NetWork:  constant.TCP,
		Type:     constant.SOCKS5,
		SrcIP:    lAddr.IP,
		DstIP:    rAddr.IP,
		SrcPort:  strconv.Itoa(lAddr.Port),
		DstPort:  strconv.Itoa(rAddr.Port),
		AddrType: constant.AtypIPv4,
		Host:     "",
	}
}
