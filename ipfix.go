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
	ipfixUdpCh = make(chan IPFIXUDPMsg, 1000)
)

func NewIPFIX(opts *Options) *IPFIX {
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

		d := ipfix.NewDecoder(msg.body)
		d.Decode()
	}
}
