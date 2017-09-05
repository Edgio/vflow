//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    transport.go
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

package packet

import "errors"

const (
	// IANAProtoICMP is IANA Internet Control Message number
	IANAProtoICMP = 1

	// IANAProtoTCP is IANA Transmission Control number
	IANAProtoTCP = 6

	// IANAProtoUDP is IANA User Datagram number
	IANAProtoUDP = 17
)

// TCPHeader represents TCP header
type TCPHeader struct {
	SrcPort    int
	DstPort    int
	DataOffset int
	Reserved   int
	Flags      int
}

// UDPHeader represents UDP header
type UDPHeader struct {
	SrcPort int
	DstPort int
}

var (
	errShortTCPHeaderLength = errors.New("short TCP header length")
	errShortUDPHeaderLength = errors.New("short UDP header length")
)

func decodeTCP(b []byte) (TCPHeader, error) {
	if len(b) < 20 {
		return TCPHeader{}, errShortTCPHeaderLength
	}

	return TCPHeader{
		SrcPort:    int(b[0])<<8 | int(b[1]),
		DstPort:    int(b[2])<<8 | int(b[3]),
		DataOffset: int(b[12]) >> 4,
		Reserved:   0,
		Flags:      ((int(b[12])<<8 | int(b[13])) & 0x01ff),
	}, nil
}

func decodeUDP(b []byte) (UDPHeader, error) {
	if len(b) < 8 {
		return UDPHeader{}, errShortUDPHeaderLength
	}

	return UDPHeader{
		SrcPort: int(b[0])<<8 | int(b[1]),
		DstPort: int(b[2])<<8 | int(b[3]),
	}, nil
}
