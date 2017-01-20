package sflow

import (
	"fmt"
)

type IPv4Header struct {
	Version  int    // protocol version
	TOS      int    // type-of-service
	TotalLen int    // packet total length
	ID       int    // identification
	Flags    int    // flags
	FragOff  int    // fragment offset
	TTL      int    // time-to-live
	Protocol int    // next protocol
	Checksum int    // checksum
	Src      string // source address
	Dst      string // destination address
	Options  []byte // options, extension headers
}

const (
	EtherTypeIPv4      = 0x0800
	EtherTypeIPv6      = 0x86DD
	EtherTypeIEEE8021Q = 0x8100
)

func decodeISO88023(b []byte) {

	if len(b) < 14 {
		return
	}

	etherType := uint16(b[13]) | uint16(b[12])<<8
	switch etherType {
	case EtherTypeIPv4:
		h := decodeIPv4Header(b[14:])
		fmt.Printf("IPv4: %#v\n", h)
	case EtherTypeIEEE8021Q:
		fmt.Println("802.1Q")
	}
}

func decodeIPv4Header(b []byte) IPv4Header {
	return IPv4Header{
		Version:  int(b[0] & 0xf0 >> 4),
		TOS:      int(b[1]),
		TotalLen: int(b[2])<<8 | int(b[3]),
		ID:       int(b[4])<<8 | int(b[5]),
		Flags:    int(b[6] & 0x07),
		TTL:      int(b[8]),
		Protocol: int(b[9]),
		Checksum: int(b[10])<<8 | int(b[11]),
		Src:      fmt.Sprintf("%d.%d.%d.%d", b[12], b[13], b[14], b[15]),
		Dst:      fmt.Sprintf("%d.%d.%d.%d", b[16], b[17], b[18], b[19]),
	}
}
