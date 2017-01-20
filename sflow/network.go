package sflow

import (
	"errors"
	"fmt"
)

// IPv4Header represents an IPv4 header
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
}

// IPv6Header represents an IPv6 header
type IPv6Header struct {
	Version      int    // protocol version
	TrafficClass int    // traffic class
	FlowLabel    int    // flow label
	PayloadLen   int    // payload length
	NextHeader   int    // next header
	HopLimit     int    // hop limit
	Src          string // source address
	Dst          string // destination address
}

type Datalink struct {
	SrcMAC    string
	DstMAC    string
	Vlan      int
	EtherType uint16
}

type Packet struct {
	L2 Datalink
	L3 interface{}
	L4 interface{}
}

const (
	EtherTypeIPv4      = 0x0800
	EtherTypeIPv6      = 0x86DD
	EtherTypeIEEE8021Q = 0x8100
)

var (
	errShortEthernetHeaderLength = errors.New("the ethernet header is too small")
	errShortIPv4HeaderLength     = errors.New("short ipv4 header length")
	errShortIPv6HeaderLength     = errors.New("short ipv6 header length")
	errShortEthernetLength       = errors.New("short ethernet header length")
	errUnknownEtherType          = errors.New("unknown ether type")
)

func decodeISO88023(b []byte) (Packet, error) {
	var (
		err    error
		packet Packet
	)

	if len(b) < 14 {
		return packet, errShortEthernetHeaderLength
	}

	packet.L2, err = decodeIEEE802(b)
	if err != nil {
		return packet, err
	}

	switch packet.L2.EtherType {
	case EtherTypeIPv4:
		b = b[14:]

		packet.L3, err = decodeIPv4Header(b)
		if err != nil {
			return packet, err
		}

		packet.L4, err = decodeTransportLayer(int(b[9]), b[20:])
		if err != nil {
			return packet, err
		}

		return packet, nil

	case EtherTypeIPv6:
	case EtherTypeIEEE8021Q:
		vlan := int(b[14])<<8 | int(b[15])
		b[12], b[13] = b[16], b[17]
		b = append(b[:14], b[18:]...)
		packet, err = decodeISO88023(b)
		if err != nil {
			return packet, err
		}
		packet.L2.Vlan = vlan
		return packet, nil

	default:
		return packet, errUnknownEtherType
	}

	return packet, nil
}

func decodeIPv6Header(b []byte) (IPv6Header, error) {
	if len(b) < 40 {
		return IPv6Header{}, errShortIPv6HeaderLength
	}

	return IPv6Header{}, nil
}

func decodeIPv4Header(b []byte) (IPv4Header, error) {
	if len(b) < 20 {
		return IPv4Header{}, errShortIPv4HeaderLength
	}

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
	}, nil
}

func decodeIEEE802(b []byte) (Datalink, error) {
	var d Datalink

	if len(b) < 14 {
		return d, errShortEthernetLength
	}

	d.EtherType = uint16(b[13]) | uint16(b[12])<<8

	if d.EtherType != EtherTypeIEEE8021Q {
		d.SrcMAC = fmt.Sprintf("%0.2x:%0.2x:%0.2x:%0.2x:%0.2x:%0.2x", b[0], b[1], b[2], b[3], b[4], b[5])
		d.DstMAC = fmt.Sprintf("%0.2x:%0.2x:%0.2x:%0.2x:%0.2x:%0.2x", b[6], b[7], b[8], b[9], b[10], b[11])
	}

	return d, nil
}
