//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    mirror_test.go
//: details: provides support for automated testing of mirror methods
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
	"net"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func TestNewRawConn(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	c, err := NewRawConn(ip)
	if err != nil {
		if strings.Contains(err.Error(), "not permitted") {
			t.Log(err)
			return
		}
		t.Error("unexpected error", err)
	}
	if c.family != syscall.AF_INET {
		t.Error("expected family# 2, got", c.family)
	}
}

func TestNewRawConnIPv6(t *testing.T) {
	ip := net.ParseIP("2001:0db8:0000:0000:0000:ff00:0042:8329")
	c, err := NewRawConn(ip)
	if err != nil {
		if strings.Contains(err.Error(), "not permitted") {
			t.Log(err)
			return
		}
		t.Error("unexpected error", err)
	}
	if c.family != syscall.AF_INET6 {
		t.Error("expected family# 10, got", c.family)
	}
}

func TestIPv4Header(t *testing.T) {
	ipv4RawHeader := NewIPv4HeaderTpl(17)
	b := ipv4RawHeader.Marshal()
	h, err := ipv4.ParseHeader(b)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if h.Version != 4 {
		t.Error("expect version: 4, got", h.Version)
	}
	if h.Protocol != 17 {
		t.Error("expect protocol: 17, got", h.Protocol)
	}
	if h.TTL != 64 {
		t.Error("expect TTL: 64, got", h.TTL)
	}
	if h.Len != 20 {
		t.Error("expect Len: 20, got", h.Len)
	}
	if h.Checksum != 0 {
		t.Error("expect Checksum: 0, got", h.Checksum)
	}
}

func TestIPv6Header(t *testing.T) {
	ipv6RawHeader := NewIPv6HeaderTpl(17)
	b := ipv6RawHeader.Marshal()
	h, err := ipv6.ParseHeader(b)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if h.Version != 6 {
		t.Error("expect version: 4, got", h.Version)
	}
	if h.NextHeader != 17 {
		t.Error("expect protocol: 17, got", h.NextHeader)
	}
	if h.HopLimit != 64 {
		t.Error("expect TTL: 64, got", h.HopLimit)
	}
}

func TestSetAddrs(t *testing.T) {
	src := net.ParseIP("10.11.12.13")
	dst := net.ParseIP("192.17.11.1")
	ipv4RawHeader := NewIPv4HeaderTpl(17)
	b := ipv4RawHeader.Marshal()
	ipv4RawHeader.SetAddrs(b, src, dst)
	h, err := ipv4.ParseHeader(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if h.Src.String() != "10.11.12.13" {
		t.Error("expect src 10.11.12.13, got", h.Src.String())
	}
	if h.Dst.String() != "192.17.11.1" {
		t.Error("expect dst 192.17.11.1, got", h.Src.String())
	}
}

func TestSetLen(t *testing.T) {
	ipv4RawHeader := NewIPv4HeaderTpl(17)
	b := ipv4RawHeader.Marshal()
	ipv4RawHeader.SetLen(b, 15)
	h, err := ipv4.ParseHeader(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if h.TotalLen != 35 {
		t.Error("expect total len 35, got", h.TotalLen)
	}
}
