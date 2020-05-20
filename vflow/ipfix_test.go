//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix_test.go
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

package main

import (
	"net"
	"strings"
	"testing"
	"time"
)

func init() {
	opts = &Options{}
}

func TestMirrorIPFIX(t *testing.T) {
	var (
		msg   = make(chan IPFIXUDPMsg, 1)
		fb    = make(chan IPFIXUDPMsg)
		dst   = net.ParseIP("127.0.0.1")
		ready = make(chan struct{})
	)

	go func() {
		err := mirrorIPFIX(dst, 10024, msg)
		if err != nil {
			if strings.Contains(err.Error(), "not permitted") {
				t.Log(err)
				ready <- struct{}{}
			} else {
				t.Fatal("unexpected error", err)
			}
		}
	}()

	time.Sleep(1 * time.Second)

	go func() {
		b := make([]byte, 1500)
		laddr := &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 10024,
		}

		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			t.Fatal("unexpected error", err)
		}

		close(ready)

		n, raddr, err := conn.ReadFrom(b)
		if err != nil {
			t.Error("unexpected error", err)
		}

		host, _, err := net.SplitHostPort(raddr.String())
		if err != nil {
			t.Error("unexpected error", err)
		}

		fb <- IPFIXUDPMsg{
			body:  b[:n],
			raddr: &net.UDPAddr{IP: net.ParseIP(host)},
		}

	}()

	_, ok := <-ready
	if ok {
		return
	}

	body := ipfixBuffer.Get().([]byte)
	copy(body, []byte("hello"))
	body = body[:5]

	msg <- IPFIXUDPMsg{
		body: body,
		raddr: &net.UDPAddr{
			IP: net.ParseIP("192.1.1.1"),
		},
	}

	feedback := <-fb

	if string(feedback.body) != "hello" {
		t.Error("expect body is hello, got", string(feedback.body))
	}

	if feedback.raddr.IP.String() != "192.1.1.1" {
		t.Error("expect raddr is 192.1.1.1, got", feedback.raddr.IP.String())
	}
}
