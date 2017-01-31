package mirror

import (
	"net"
	"syscall"
)

type Conn struct {
	family int
	sotype int
	proto  int
	fd     int
	raddr  *net.UDPAddr
}

func DialUDP(raddr *net.UDPAddr) (Conn, error) {
	var err error

	conn := Conn{
		sotype: syscall.SOCK_DGRAM,
		proto:  syscall.IPPROTO_UDP,
		raddr:  raddr,
	}

	if raddr.IP.To4 != nil {
		conn.family = syscall.AF_INET
	} else {
		conn.family = syscall.AF_INET6
	}

	conn.fd, err = syscall.Socket(conn.family, conn.sotype, conn.proto)

	return conn, err
}
