//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix.go
//: details: Read IPFIX and Netflow v9 data fields based on the type
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

package ipfix

import (
	"encoding/binary"
	"math"
	"net"
)

// Interpret read data fields based on the type - big endian
func Interpret(b *[]byte, t FieldType) interface{} {
	if len(*b) < t.minLen() {
		return *b
	}

	switch t {
	case Boolean:
		return (*b)[0] == 1
	case Uint8:
		return (*b)[0]
	case Uint16:
		return binary.BigEndian.Uint16(*b)
	case Uint32:
		return binary.BigEndian.Uint32(*b)
	case Uint64:
		return binary.BigEndian.Uint64(*b)
	case Int8:
		return int8((*b)[0])
	case Int16:
		return int16(binary.BigEndian.Uint16(*b))
	case Int32:
		return int32(binary.BigEndian.Uint32(*b))
	case Int64:
		return int64(binary.BigEndian.Uint64(*b))
	case Float32:
		return math.Float32frombits(binary.BigEndian.Uint32(*b))
	case Float64:
		return math.Float64frombits(binary.BigEndian.Uint64(*b))
	case MacAddress:
		return net.HardwareAddr(*b)
	case String:
		return string(*b)
	case Ipv4Address, Ipv6Address:
		return net.IP(*b)
	case DateTimeSeconds:
		return binary.BigEndian.Uint32(*b)
	case DateTimeMilliseconds, DateTimeMicroseconds, DateTimeNanoseconds:
		return binary.BigEndian.Uint64(*b)
	case Unknown, OctetArray:
		return *b
	}
	return *b
}

func (t FieldType) minLen() int {
	switch t {
	case Boolean:
		return 1
	case Uint8, Int8:
		return 1
	case Uint16, Int16:
		return 2
	case Uint32, Int32, Float32:
		return 4
	case DateTimeSeconds:
		return 4
	case Uint64, Int64, Float64:
		return 8
	case DateTimeMilliseconds, DateTimeMicroseconds, DateTimeNanoseconds:
		return 8
	case MacAddress:
		return 6
	case Ipv4Address:
		return 4
	case Ipv6Address:
		return 16
	default:
		return 0
	}
}
