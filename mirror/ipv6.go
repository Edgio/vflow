//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipv6.go
//: details: mirror ipv6 handler
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

// IPv6 represents IP version 6 header
type IPv6 struct {
	Version       uint8
	TrafficClass  uint8
	FlowLabel     uint32
	PayloadLength uint16
	NextHeader    uint8
	HopLimit      uint8
}

// NewIPv6HeaderTpl returns a new IPv6 as template
func NewIPv6HeaderTpl(proto int) IPv6 {
	return IPv6{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		NextHeader:   uint8(proto),
		HopLimit:     64,
	}
}

// Marshal returns encoded IPv6
func (ip IPv6) Marshal() []byte {
	b := make([]byte, IPv6HLen)
	b[0] = byte((ip.Version << 4) | (ip.TrafficClass >> 4))
	b[1] = byte((ip.TrafficClass << 4) | uint8(ip.FlowLabel>>16))
	binary.BigEndian.PutUint16(b[2:], uint16(ip.FlowLabel))
	b[6] = byte(ip.NextHeader)
	b[7] = byte(ip.HopLimit)

	return b
}

// SetLen sets IPv6 length
func (ip IPv6) SetLen(b []byte, n int) {
	binary.BigEndian.PutUint16(b[4:], IPv6HLen+uint16(n))
}

// SetAddrs sets IPv6 src and dst addresses
func (ip IPv6) SetAddrs(b []byte, src, dst net.IP) {
	copy(b[8:], src)
	copy(b[24:], dst)
}
