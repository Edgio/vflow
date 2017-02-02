//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipv4.go
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
}

func NewIPv4HeaderTpl(proto int) IPv4 {
	return IPv4{
		Version:  4,
		IHL:      5,
		TOS:      0,
		TTL:      64,
		Protocol: uint8(proto),
	}
}

func (ip *IPv4) Marshal() []byte {
	b := make([]byte, IPv4HLen)
	b[0] = byte((ip.Version << 4) | ip.IHL)
	b[1] = byte(ip.TOS)
	binary.BigEndian.PutUint16(b[2:], ip.Length)
	b[6] = byte(0)
	b[7] = byte(0)
	b[8] = byte(ip.TTL)
	b[9] = byte(ip.Protocol)

	return b
}

func (ip *IPv4) SetLen(b []byte, n int) {
	binary.BigEndian.PutUint16(b[2:], IPv4HLen+uint16(n))
}

func (ip *IPv4) SetAddrs(b []byte, src, dst net.IP) {
	copy(b[12:16], src[12:16])
	copy(b[16:20], dst[12:16])
}
