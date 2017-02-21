// Package ipfix decodes IPFIX packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    marshal_test.go
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
package ipfix

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"
)

var mockDecodedMsg = Message{
	AgentID: "10.10.10.10",
	Header: MessageHeader{
		Version:    10,
		Length:     420,
		ExportTime: 1483484756,
		SequenceNo: 2563920489,
		DomainID:   34560,
	},
	DataSets: [][]DecodedField{
		[]DecodedField{
			DecodedField{ID: 0x8, Value: net.IP{0x5b, 0x7d, 0x82, 0x79}},
			DecodedField{ID: 0xc, Value: net.IP{0xc0, 0xe5, 0xdc, 0x85}},
			DecodedField{ID: 0x5, Value: 0x0},
			DecodedField{ID: 0x4, Value: 0x6},
			DecodedField{ID: 0x7, Value: 0xecba},
			DecodedField{ID: 0xb, Value: 0x1bb},
			DecodedField{ID: 0x20, Value: 0x0},
			DecodedField{ID: 0xa, Value: 0x503},
			DecodedField{ID: 0x3a, Value: 0x0},
			DecodedField{ID: 0x9, Value: 0x10},
			DecodedField{ID: 0xd, Value: 0x18},
			DecodedField{ID: 0x10, Value: 0x1ad7},
			DecodedField{ID: 0x11, Value: 0x3b1d},
			DecodedField{ID: 0xf, Value: net.IP{0xc0, 0x10, 0x1c, 0x58}},
			DecodedField{ID: 0x6, Value: []uint8{0x10}},
			DecodedField{ID: 0xe, Value: 0x4f6},
			DecodedField{ID: 0x1, Value: 0x28},
			DecodedField{ID: 0x2, Value: 0x1},
			DecodedField{ID: 0x34, Value: 0x3a},
			DecodedField{ID: 0x35, Value: 0x3a},
			DecodedField{ID: 0x98, Value: 1483484685331},
			DecodedField{ID: 0x99, Value: 1483484685331},
			DecodedField{ID: 0x88, Value: 0x1},
			DecodedField{ID: 0xf3, Value: 0x0},
			DecodedField{ID: 0xf5, Value: 0x0},
		},
	},
}

func TestJSONMarshal(t *testing.T) {
	buf := bytes.NewBufferString("")
	msg := Message{}

	b, err := mockDecodedMsg.JSONMarshal(buf)
	if err != nil {
		t.Error("unexpected error", err)
	}

	err = json.Unmarshal(b, &msg)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if msg.AgentID != "10.10.10.10" {
		t.Error("expect AgentID 10.10.10.10, got", msg.AgentID)
	}
	if msg.Header.Version != 10 {
		t.Error("expect Version 10, got", msg.Header.Version)
	}
	for _, ds := range msg.DataSets {
		for _, f := range ds {
			switch f.ID {
			case 1:
				if f.Value.(float64) != 40 {
					t.Error("expect ID 1 value 40, got", f.Value)
				}
			case 2:
				if f.Value.(float64) != 1 {
					t.Error("expect ID 2 value 1, got", f.Value)
				}
			case 4:
				if f.Value.(float64) != 6 {
					t.Error("expect ID 4 value 6, got", f.Value)
				}
			case 5:
				if f.Value.(float64) != 0 {
					t.Error("expect ID 5 value 0, got", f.Value)
				}
			case 6:
				if f.Value.(string) != "0x10" {
					t.Error("expect ID 6 value 0x10, got", f.Value)
				}
			case 8:
				if f.Value != "91.125.130.121" {
					t.Error("expect ID 8 value 91.125.130.121, got", f.Value)
				}
			case 12:
				if f.Value != "192.229.220.133" {
					t.Error("expect ID 12 value 192.229.220.133, got", f.Value)
				}
			case 13:
				if f.Value.(float64) != 24 {
					t.Error("expect ID 13 value 24, got", f.Value)
				}
			case 14:
				if f.Value.(float64) != 1270 {
					t.Error("expect ID 14 value 1270, got", f.Value)
				}
			case 152:
				if f.Value.(float64) != 1483484685331 {
					t.Error("expect ID 152 value 1483484685331, got", f.Value)
				}
			}
		}
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	buf := bytes.NewBufferString("")

	for i := 0; i < b.N; i++ {
		mockDecodedMsg.JSONMarshal(buf)
	}

}
