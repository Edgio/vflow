//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix.go
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
package main

import (
	"net"
	"strconv"
	"sync"
	"time"

	"git.edgecastcdn.net/vflow/ipfix"
	"git.edgecastcdn.net/vflow/mirror"
)

type IPFIX struct {
	port    int
	addr    string
	udpSize int
	workers int
	stop    bool
}

type IPFIXUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

var (
	ipfixUdpCh         = make(chan IPFIXUDPMsg, 1000)
	ipfixMCh           = make(chan IPFIXUDPMsg, 1000)
	ipfixMirrorEnabled bool
)

func NewIPFIX() *IPFIX {
	return &IPFIX{
		port:    opts.IPFIXPort,
		udpSize: opts.IPFIXUDPSize,
		workers: opts.IPFIXWorkers,
	}
}

func (i *IPFIX) run() {
	var wg sync.WaitGroup

	hostPort := net.JoinHostPort(i.addr, strconv.Itoa(i.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {

	}

	for n := 0; n < i.workers; n++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			ipfixWorker()

		}()
	}

	logger.Printf("ipfix is running (workers#: %d)", i.workers)

	go func() {
		mirrorIPFIX(ipfixMCh)
	}()

	for !i.stop {
		b := make([]byte, i.udpSize)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		ipfixUdpCh <- IPFIXUDPMsg{raddr, b[:n]}
	}

	wg.Wait()
}

func (i *IPFIX) shutdown() {
	i.stop = true
	logger.Println("stopped ipfix service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("ipfix has been shutdown")
	close(ipfixUdpCh)
}

func ipfixWorker() {
	var (
		msg IPFIXUDPMsg
		ok  bool
	)

	for {
		if msg, ok = <-ipfixUdpCh; !ok {
			break
		}

		if verbose {
			logger.Printf("rcvd ipfix data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		if ipfixMirrorEnabled {
			ipfixMCh <- IPFIXUDPMsg{msg.raddr, append([]byte{}, msg.body...)}
		}

		d := ipfix.NewDecoder(msg.body)
		d.Decode()
	}
}

func mirrorIPFIXv4(dst net.IP, port int, ch chan IPFIXUDPMsg) error {
	var (
		packet = make([]byte, 1500)
		msg    IPFIXUDPMsg
		pLen   int
		err    error
	)

	conn, err := mirror.NewRawConn(dst)
	if err != nil {
		return err
	}

	udp := mirror.UDP{4041, port, 0, 0}
	udpHdr := udp.Marshal()

	ip := mirror.NewIPv4HeaderTpl(mirror.UDPProto)
	ipHdr := ip.Marshal()

	for {
		msg = <-ch
		pLen = len(msg.body)

		ip.SetAddrs(ipHdr, msg.raddr.IP, dst)
		ip.SetLen(ipHdr, pLen+mirror.UDPHLen)
		udp.SetLen(udpHdr, pLen)

		copy(packet[0:20], ipHdr)
		copy(packet[20:28], udpHdr)
		copy(packet[28:], msg.body)

		if err = conn.Send(packet[0 : 28+pLen]); err != nil {
			return err
		}
	}
	return nil
}

func mirrorIPFIX(ch chan IPFIXUDPMsg) {
	var (
		ch4 = make(chan IPFIXUDPMsg, 1000)
		msg IPFIXUDPMsg
	)

	if opts.IPFIXMirror == "" {
		return
	}

	host, port, err := net.SplitHostPort(opts.IPFIXMirror)
	if err != nil {
		logger.Fatalf("wrong ipfix mirror address%s", opts.IPFIXMirror)
	}

	portNo, _ := strconv.Atoi(port)
	dst := net.ParseIP(host)

	if dst.To4() != nil {
		go mirrorIPFIXv4(dst, portNo, ch4)
	}

	ipfixMirrorEnabled = true
	logger.Println("ipfix mirror service is running ...")

	for {
		msg = <-ch
		if msg.raddr.IP.To4() != nil {
			ch4 <- msg
		}
	}
}
