package mirror

import (
	"encoding/binary"
	"net"
)

type IPv4 struct {
	Version  uint8
	IHL      uint8
	TOS      uint8
	Length   uint16
	Id       uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	SrcAddr  net.IP
	DstAddr  net.IP
}

func NewIPv4Header(src, dst net.IP, proto int) IPv4 {
	return IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: uint8(proto),
		SrcAddr:  src,
		DstAddr:  dst,
	}
}

func (ip *IPv4) SetLen(b []byte, n int) {
	binary.BigEndian.PutUint16(b[2:], uint16(n))
}

func (ip *IPv4) Marshal() ([]byte, error) {
	b := make([]byte, IPv4HLen)
	b[0] = byte((ip.Version << 4) | ip.IHL)
	b[1] = byte(ip.TOS)
	binary.BigEndian.PutUint16(b[2:], ip.Length)
	b[6] = byte(0)
	b[7] = byte(0)
	b[8] = byte(ip.TTL)
	b[9] = byte(ip.Protocol)

	copy(b[12:16], ip.SrcAddr[12:16])
	copy(b[16:20], ip.DstAddr[12:16])

	return b, nil
}
