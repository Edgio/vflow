package mirror

import (
	"net"
	"syscall"
)

const (
	IPv4HLen = 20
	IPv6HLen = 40
	UDPHLen  = 8
)

type Conn struct {
	family int
	sotype int
	proto  int
	fd     int
	raddr  syscall.Sockaddr
}

func NewRawConn(raddr net.IP) (Conn, error) {
	var err error
	conn := Conn{
		sotype: syscall.SOCK_RAW,
		proto:  syscall.IPPROTO_RAW,
	}

	if ipv4 := raddr.To4(); ipv4 != nil {
		ip := [4]byte{}
		copy(ip[:], ipv4)

		conn.family = syscall.AF_INET
		conn.raddr = &syscall.SockaddrInet4{
			Port: 0,
			Addr: ip,
		}
	} else if ipv6 := raddr.To16(); ipv6 != nil {
		ip := [16]byte{}
		copy(ip[:], ipv6)

		conn.family = syscall.AF_INET6
		conn.raddr = &syscall.SockaddrInet6{
			Addr: ip,
		}

	}

	conn.fd, err = syscall.Socket(conn.family, conn.sotype, conn.proto)

	return conn, err
}

func (c *Conn) Send(b []byte) error {
	return syscall.Sendto(c.fd, b, 0, c.raddr)
}

func (c *Conn) Close(b []byte) error {
	return syscall.Close(c.fd)
}
