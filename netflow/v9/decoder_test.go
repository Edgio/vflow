//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder_test.go
//: details: netflow v9 decoder tests and benchmarks
//: author:  Mehrdad Arshad Rad
//: date:    05/05/2017
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

package netflow9

import (
	"net"
	"testing"
)

func TestDecodeNoData(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")
	body := []byte{}
	d := NewDecoder(ip, body, false)
	if _, err := d.Decode(mCache); err == nil {
		t.Error("expected err but nothing")
	}
}

func TestUnknownElement(t *testing.T) {
	var template = []byte{
		0x00, 0x09, 0x00, 0x01, // v9, set count = 1
		0x00, 0x00, 0x00, 0x01, // uptime
		0x63, 0x99, 0xe9, 0x21, // timestamp
		0x00, 0x00, 0xff, 0x01, // sequence
		0x00, 0x00, 0x00, 0x01, // source id
		0x00, 0x00, // template set
		0x00, 0x10, // length
		0xee, 0xee, // template id
		0x00, 0x02, // field count
		0x00, 0x01, // element id 1
		0x00, 0x08, // length 8
		0xde, 0xad, // element id 57005
		0x00, 0x04, // length 4
	}

	var payload = []byte{
		0x00, 0x09, 0x00, 0x01, // v9, set count = 1
		0x00, 0x00, 0x00, 0x02, // uptime
		0x63, 0x99, 0xe9, 0x22, // timestamp
		0x00, 0x00, 0xff, 0x02, // sequence
		0x00, 0x00, 0x00, 0x01, // source id
		0xee, 0xee, // template id
		0x00, 0x10, // length
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x11, // element #1 data
		0x00, 0x22, 0x22, 0x22, // element #2 data
	}

	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")
	d := NewDecoder(ip, template, false)
	_, err := d.Decode(mCache)
	if err != nil {
		t.Error(err)
	}

	// Parse data with unknown element
	d = NewDecoder(ip, payload, false)
	m, err := d.Decode(mCache)
	if err == nil {
		t.Error("Expected error due to unknown element, but got nil")
	}
	if len(m.DataSets) != 0 {
		t.Error("Did not expect any result datasets, but got", m.DataSets)
	}

	// Now parse again, skip unknown element
	d = NewDecoder(ip, payload, true)
	m, err = d.Decode(mCache)
	if err != nil {
		t.Error(err)
	}
	if len(m.DataSets) != 1 {
		t.Error("Expected 1 dataset, but got", m.DataSets)
	}
}
