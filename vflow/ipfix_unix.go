//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix.go
//: author:  Mehrdad Arshad Rad - copied by Jeremy Rossi, but not important
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
// +build !windows

package main

import (
	"github.com/VerizonDigital/vflow/mirror"

	"net"
)

func mirrorIPFIXDispatcher(ch chan IPFIXUDPMsg) {
	var (
		ch4 = make(chan IPFIXUDPMsg, 1000)
		ch6 = make(chan IPFIXUDPMsg, 1000)
		msg IPFIXUDPMsg
	)

	if opts.IPFIXMirrorAddr == "" {
		return
	}

	for w := 0; w < opts.IPFIXMirrorWorkers; w++ {
		dst := net.ParseIP(opts.IPFIXMirrorAddr)

		if dst.To4() != nil {
			go mirrorIPFIX(dst, opts.IPFIXMirrorPort, ch4)
		} else {
			go mirrorIPFIX(dst, opts.IPFIXMirrorPort, ch6)
		}
	}

	ipfixMirrorEnabled = true
	logger.Printf("ipfix mirror service is running (workers#: %d) ...", opts.IPFIXMirrorWorkers)

	for {
		msg = <-ch
		if msg.raddr.IP.To4() != nil {
			ch4 <- msg
		} else {
			ch6 <- msg
		}
	}
}

func mirrorIPFIX(dst net.IP, port int, ch chan IPFIXUDPMsg) error {
	var (
		packet = make([]byte, opts.IPFIXUDPSize)
		msg    IPFIXUDPMsg
		pLen   int
		err    error
		ipHdr  []byte
		ipHLen int
		ipv4   bool
		ip     mirror.IP
	)

	conn, err := mirror.NewRawConn(dst)
	if err != nil {
		return err
	}

	udp := mirror.UDP{55117, port, 0, 0}
	udpHdr := udp.Marshal()

	if dst.To4() != nil {
		ipv4 = true
	}

	if ipv4 {
		ip = mirror.NewIPv4HeaderTpl(mirror.UDPProto)
		ipHdr = ip.Marshal()
		ipHLen = mirror.IPv4HLen
	} else {
		ip = mirror.NewIPv6HeaderTpl(mirror.UDPProto)
		ipHdr = ip.Marshal()
		ipHLen = mirror.IPv6HLen
	}

	for {
		msg = <-ch
		pLen = len(msg.body)

		ip.SetAddrs(ipHdr, msg.raddr.IP, dst)
		ip.SetLen(ipHdr, pLen+mirror.UDPHLen)

		udp.SetLen(udpHdr, pLen)
		// IPv6 checksum mandatory
		if !ipv4 {
			udp.SetChecksum()
		}

		copy(packet[0:ipHLen], ipHdr)
		copy(packet[ipHLen:ipHLen+8], udpHdr)
		copy(packet[ipHLen+8:], msg.body)

		ipfixBuffer.Put(msg.body[:opts.IPFIXUDPSize])

		if err = conn.Send(packet[0 : ipHLen+8+pLen]); err != nil {
			return err
		}
	}
}
