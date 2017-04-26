//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    packet.go
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

// Packet represents layer 2,3,4 available info
type Packet struct {
	L2   Datalink
	L3   interface{}
	L4   interface{}
	data []byte
}

var (
	errUnknownEtherType = errors.New("unknown ether type")
)

// NewPacket constructs a packet object
func NewPacket() Packet {
	return Packet{}
}

// Decoder decodes packet's layers
func (p *Packet) Decoder(data []byte) (*Packet, error) {
	var (
		err error
	)

	p.data = data
	err = p.decodeEthernet()
	if err != nil {
		return p, err
	}

	switch p.L2.EtherType {
	case EtherTypeIPv4:

		err = p.decodeIPv4Header()
		if err != nil {
			return p, err
		}

		err = p.decodeNextLayer()
		if err != nil {
			return p, err
		}

	case EtherTypeIPv6:

		err = p.decodeIPv6Header()
		if err != nil {
			return p, err
		}

		err = p.decodeNextLayer()
		if err != nil {
			return p, err
		}

	default:
		return p, errUnknownEtherType
	}

	return p, nil
}
