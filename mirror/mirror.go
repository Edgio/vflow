// Package mirror clones the IPFIX packets w/ spoofing feature
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    mirror.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------
package mirror

import (
	"net"
	"syscall"
)

const (
	// IPv4HLen is IP version 4 header length
	IPv4HLen = 20

	// IPv6HLen is IP version 6 header length
	IPv6HLen = 40

	// UDPHLen is UDP header length
	UDPHLen = 8

	// UDPProto is UDP protocol IANA number
	UDPProto = 17
)

// Conn represents socket connection properties
type Conn struct {
	family int
	sotype int
	proto  int
	fd     int
	raddr  syscall.Sockaddr
}

// IP is network layer corresponding to IPv4/IPv6
type IP interface {
	Marshal() []byte
	SetLen([]byte, int)
	SetAddrs([]byte, net.IP, net.IP)
}

// NewRawConn constructs new raw socket
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

// Send tries to put the bytes to wire
func (c *Conn) Send(b []byte) error {
	return syscall.Sendto(c.fd, b, 0, c.raddr)
}

// Close releases file descriptor
func (c *Conn) Close(b []byte) error {
	return syscall.Close(c.fd)
}
