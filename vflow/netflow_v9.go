//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    netflow_v9.go
//: details: netflow decoders handler
//: author:  Mehrdad Arshad Rad
//: date:    04/21/2017
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
	"sync/atomic"
	"time"

	"github.com/VerizonDigital/vflow/producer"
)

// NetflowV9 represents netflow v9 collector
type NetflowV9 struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   NetflowV9Stats
	pool    chan chan struct{}
}

// NetflowV9UDPMsg represents netflow v9 UDP data
type NetflowV9UDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

// NetflowV9Stats represents netflow v9 stats
type NetflowV9Stats struct {
	UDPQueue     int
	MessageQueue int
	UDPCount     uint64
	DecodedCount uint64
	MQErrorCount uint64
	Workers      int32
}

var (
	netflowV9UDPCh = make(chan NetflowV9UDPMsg, 1000)
	netflowV9MQCh  = make(chan []byte, 1000)

	// ipfix udp payload pool
	netflowV9Buffer = &sync.Pool{
		New: func() interface{} {
			return make([]byte, opts.NetflowV9UDPSize)
		},
	}
)

// NewNetflowV9 constructs NetflowV9
func NewNetflowV9() *NetflowV9 {
	return &NetflowV9{
		port:    opts.IPFIXPort,
		workers: opts.IPFIXWorkers,
		pool:    make(chan chan struct{}, maxWorkers),
	}
}

func (i *NetflowV9) run() {
	//TODO
	// exit if the ipfix is disabled
	if !opts.NetflowV9Enabled {
		logger.Println("netflowv9 has been disabled")
		return
	}

	hostPort := net.JoinHostPort(i.addr, strconv.Itoa(i.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal(err)
	}

	atomic.AddInt32(&i.stats.Workers, int32(i.workers))
	for n := 0; n < i.workers; n++ {
		go func() {
			wQuit := make(chan struct{})
			i.pool <- wQuit
			i.netflowV9Worker(wQuit)
		}()
	}

	logger.Printf("netflow v9 is running (workers#: %d)", i.workers)

	go func() {
		p := producer.NewProducer(opts.MQName)

		p.MQConfigFile = opts.MQConfigFile
		p.MQErrorCount = &i.stats.MQErrorCount
		p.Logger = logger
		p.Chan = netflowV9MQCh
		p.Topic = opts.NetflowV9Topic

		if err := p.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		if !opts.DynWorkers {
			logger.Println("netflow v9 dynamic worker disabled")
			return
		}

		i.dynWorkers()
	}()

	for !i.stop {
		b := netflowV9Buffer.Get().([]byte)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		atomic.AddUint64(&i.stats.UDPCount, 1)
		netflowV9UDPCh <- NetflowV9UDPMsg{raddr, b[:n]}
	}

}

func (i *NetflowV9) shutdown() {
	//TODO
}

func (i *NetflowV9) netflowV9Worker(wQuit chan struct{}) {
	// TODO
}

func (i *NetflowV9) status() *NetflowV9Stats {
	//TODO
	return &NetflowV9Stats{}
}

func (i *NetflowV9) dynWorkers() {
	//TODO
}
