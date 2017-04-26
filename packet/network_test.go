//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    network_test.go
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

import "testing"

func TestDecodeIPv4Header(t *testing.T) {
	p := NewPacket()

	p.data = []byte{
		0x45, 0x0, 0x0, 0x4b, 0x8,
		0xf8, 0x0, 0x0, 0x3e, 0x11,
		0x82, 0x91, 0xc0, 0xe5, 0xd8,
		0x8f, 0xc0, 0xe5, 0x96, 0xbe,
		0x64, 0x9b, 0x0, 0x35, 0x0,
	}
	err := p.decodeIPv4Header()
	if err != nil {
		t.Error("unexpected error", err)
	}

	ipv4 := p.L3.(IPv4Header)

	if ipv4.Version != 4 {
		t.Error("unexpected version, got", ipv4.Version)
	}

	if ipv4.TOS != 0 {
		t.Error("unexpected TOS, got", ipv4.TOS)
	}
	if ipv4.TotalLen != 75 {
		t.Error("unexpected TotalLen, got", ipv4.TotalLen)
	}
	if ipv4.Flags != 0 {
		t.Error("unexpected Flags, got", ipv4.Flags)
	}

	if ipv4.FragOff != 0 {
		t.Error("unexpected FragOff", ipv4.FragOff)
	}

	if ipv4.TTL != 62 {
		t.Error("unexpected TTL, got", ipv4.TTL)
	}

	if ipv4.Protocol != 17 {
		t.Error("unexpected protocol, got", ipv4.Protocol)
	}

	if ipv4.Checksum != 33425 {
		t.Error("unexpected checksum, got", ipv4.Checksum)
	}
	if ipv4.Src != "192.229.216.143" {
		t.Error("unexpected src addr, got", ipv4.Src)
	}

	if ipv4.Dst != "192.229.150.190" {
		t.Error("unexpected dst addr, got", ipv4.Dst)
	}
}
