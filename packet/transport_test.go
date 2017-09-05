//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    transport_test.go
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

func TestDecoderUDP(t *testing.T) {
	b := []byte{
		0xa3, 0x6c, 0x0, 0x35, 0x0,
		0x3d, 0xc8, 0xdc, 0x81, 0x9f,
	}

	udp, err := decodeUDP(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if udp.SrcPort != 41836 {
		t.Error("expected src port:41836, got", udp.SrcPort)
	}

	if udp.DstPort != 53 {
		t.Error("expected dst port:53, got", udp.DstPort)
	}
}

func TestDecodeTCP(t *testing.T) {
	b := []byte{
		0xa5, 0x8e, 0x20, 0xfb, 0x54,
		0x1, 0x4f, 0x1c, 0x52, 0x7f,
		0x0, 0xf9, 0x50, 0x10, 0x1,
		0x2a, 0xbb, 0xde, 0x0, 0x0,
	}

	tcp, err := decodeTCP(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if tcp.SrcPort != 42382 {
		t.Error("expected src port:4382, got", tcp.SrcPort)
	}

	if tcp.DstPort != 8443 {
		t.Error("expected dst port:8443, got", tcp.DstPort)
	}

	if tcp.Flags != 16 {
		t.Error("expected flags:16, got", tcp.Flags)
	}
}

func TestDecodeTCP2(t *testing.T) {
	b := []byte{
		0xa5, 0x8e, 0x20, 0xfb, 0x54,
		0x1, 0x4f, 0x1c, 0x52, 0x7f,
		0x0, 0xf9, 0x51, 0x10, 0x1,
		0x2a, 0xbb, 0xde, 0x0, 0x0,
	}

	tcp, err := decodeTCP(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	// NS flag
	if tcp.Flags != 272 {
		t.Error("expected flags:272, got", tcp.Flags)
	}

	// check dataoffset
	if tcp.DataOffset != 5 {
		t.Error("expected dataoffset:5, got", tcp.DataOffset)
	}
}
