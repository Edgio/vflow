//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    hammer.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    03/01/2017
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

package hammer

import (
	"net"
	"strings"
	"testing"
)

func TestIPFIXGenPackets(t *testing.T) {
	ip := net.ParseIP("10.0.0.1")
	src := net.ParseIP("1.1.1.1")

	ipfix, err := NewIPFIX(ip)
	if err != nil {
		if !strings.Contains(err.Error(), "not permitted") {
			t.Error("unexpected error", err)
		} else {
			t.Skip(err)
		}
	}

	ipfix.srcs = append(ipfix.srcs, src)

	packets := ipfix.genPackets(dataType)
	if len(packets) < 1 {
		t.Error("expect to have packets, got", len(packets))
	}
	packets = ipfix.genPackets(templateType)
	if len(packets) < 1 {
		t.Error("expect to have tp; packets, got", len(packets))
	}
	packets = ipfix.genPackets(templateOptType)
	if len(packets) < 1 {
		t.Error("expect to have tpl opt packets, got", len(packets))
	}
}

func TestSFlowGenPackets(t *testing.T) {
	ip := net.ParseIP("10.0.0.1")
	src := net.ParseIP("1.1.1.1")

	sflow, err := NewSFlow(ip)
	if err != nil {
		if !strings.Contains(err.Error(), "not permitted") {
			t.Error("unexpected error", err)
		} else {
			t.Skip(err)
		}
	}

	sflow.srcs = append(sflow.srcs, src)

	packets := sflow.genPackets()
	if len(packets) < 1 {
		t.Error("expect to have packets, got", len(packets))
	}

}
