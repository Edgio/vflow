package packet

import (
	"errors"
	_ "fmt"
)

const (
	IANAProtoICMP = 1
	IANAProtoTCP  = 6
	IANAProtoUDP  = 17
)

type TCPHeader struct {
	SrcPort    int
	DstPort    int
	DataOffset int
	Reserved   int
	Flags      int
}

type UDPHeader struct {
	SrcPort int
	DstPort int
}

var (
	errShortTCPHeaderLength = errors.New("short TCP header length")
	errShortUDPHeaderLength = errors.New("short UDP header length")
)

func decoderTCP(b []byte) (TCPHeader, error) {
	if len(b) < 20 {
		return TCPHeader{}, errShortTCPHeaderLength
	}

	tmp := int(b[12])

	return TCPHeader{
		SrcPort:    int(b[0])<<8 | int(b[1]),
		DstPort:    int(b[2])<<8 | int(b[3]),
		DataOffset: (tmp & 0xf000) >> 12,
		Reserved:   (tmp & 0x0e00) >> 8,
		Flags:      (tmp & 0x01ff),
	}, nil
}

func decoderUDP(b []byte) (UDPHeader, error) {
	if len(b) < 8 {
		return UDPHeader{}, errShortUDPHeaderLength
	}

	return UDPHeader{
		SrcPort: int(b[0])<<8 | int(b[1]),
		DstPort: int(b[2])<<8 | int(b[3]),
	}, nil
}
