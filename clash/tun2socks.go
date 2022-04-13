package clash

import (
	"io"

	"ClashKit/clash/tun"
)

var (
	closer io.Closer
)

func StartTun2Socket(fd int, gateway, portal string) error {
	stack, err := tun.StartTun2Socket(fd, gateway, portal)
	if err != nil {
		return err
	}
	closer = stack
	return nil
}

func StopTun2Socket() error {
	return closer.Close()
}
