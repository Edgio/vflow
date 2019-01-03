//: ----------------------------------------------------------------------------
//: Copyright (C) 2018 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    netflow_v5.go
//: details: netflow v5 decoders handler
//: author:  Christopher Noel
//: date:    12/10/2018
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
	"bytes"
	"net"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VerizonDigital/vflow/netflow/v5"
	"github.com/VerizonDigital/vflow/producer"
)

// NetflowV5 represents netflow v5 collector
type NetflowV5 struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   NetflowV5Stats
	pool    chan chan struct{}
}

// NetflowV5UDPMsg represents netflow v5 UDP data
type NetflowV5UDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

// NetflowV5Stats represents netflow v5 stats
type NetflowV5Stats struct {
	UDPQueue     int
	MessageQueue int
	UDPCount     uint64
	DecodedCount uint64
	MQErrorCount uint64
	Workers      int32
}

var (
	netflowV5UDPCh = make(chan NetflowV5UDPMsg, 1000)
	netflowV5MQCh  = make(chan []byte, 1000)

	// ipfix udp payload pool
	netflowV5Buffer = &sync.Pool{
		New: func() interface{} {
			return make([]byte, opts.NetflowV5UDPSize)
		},
	}
)

// NewNetflowV5 constructs NetflowV5
func NewNetflowV5() *NetflowV5 {
	return &NetflowV5{
		port:    opts.NetflowV5Port,
		workers: opts.NetflowV5Workers,
		pool:    make(chan chan struct{}, maxWorkers),
	}
}

func (i *NetflowV5) run() {
	// exit if the netflow v5 is disabled
	if !opts.NetflowV5Enabled {
		logger.Println("netflowv5 has been disabled")
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
			i.netflowV5Worker(wQuit)
		}()
	}

	logger.Printf("netflow v5 is running (UDP: listening on [::]:%d workers#: %d)", i.port, i.workers)

	go func() {
		p := producer.NewProducer(opts.MQName)

		p.MQConfigFile = path.Join(opts.VFlowConfigPath, opts.MQConfigFile)
		p.MQErrorCount = &i.stats.MQErrorCount
		p.Logger = logger
		p.Chan = netflowV5MQCh
		p.Topic = opts.NetflowV5Topic

		if err := p.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		if !opts.DynWorkers {
			logger.Println("netflow v5 dynamic worker disabled")
			return
		}

		i.dynWorkers()
	}()

	for !i.stop {
		b := netflowV5Buffer.Get().([]byte)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		atomic.AddUint64(&i.stats.UDPCount, 1)
		netflowV5UDPCh <- NetflowV5UDPMsg{raddr, b[:n]}
	}

}

func (i *NetflowV5) shutdown() {
	// exit if the netflow v5 is disabled
	if !opts.NetflowV5Enabled {
		logger.Println("netflow v5 disabled")
		return
	}

	// stop reading from UDP listener
	i.stop = true
	logger.Println("stopping netflow v5 service gracefully ...")
	time.Sleep(1 * time.Second)

	// logging and close UDP channel
	logger.Println("netflow v5 has been shutdown")
	close(netflowV9UDPCh)
}

func (i *NetflowV5) netflowV5Worker(wQuit chan struct{}) {
	var (
		decodedMsg *netflow5.Message
		msg        = NetflowV5UDPMsg{body: netflowV5Buffer.Get().([]byte)}
		buf        = new(bytes.Buffer)
		err        error
		ok         bool
		b          []byte
	)

LOOP:
	for {

		netflowV5Buffer.Put(msg.body[:opts.NetflowV5UDPSize])
		buf.Reset()

		select {
		case <-wQuit:
			break LOOP
		case msg, ok = <-netflowV5UDPCh:
			if !ok {
				break LOOP
			}
		}

		if opts.Verbose {
			logger.Printf("rcvd netflow v5 data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		d := netflow5.NewDecoder(msg.raddr.IP, msg.body)
		if decodedMsg, err = d.Decode(); err != nil {
			logger.Println(err)
			if decodedMsg == nil {
				continue
			}
		}

		atomic.AddUint64(&i.stats.DecodedCount, 1)

		if decodedMsg.Flows != nil {
			b, err = decodedMsg.JSONMarshal(buf)
			if err != nil {
				logger.Println(err)
				continue
			}

			select {
			case netflowV5MQCh <- append([]byte{}, b...):
			default:
			}
		}

		if opts.Verbose {
			logger.Println(string(b))
		}

	}

}

func (i *NetflowV5) status() *NetflowV5Stats {
	return &NetflowV5Stats{
		UDPQueue:     len(netflowV5UDPCh),
		MessageQueue: len(netflowV5MQCh),
		UDPCount:     atomic.LoadUint64(&i.stats.UDPCount),
		DecodedCount: atomic.LoadUint64(&i.stats.DecodedCount),
		MQErrorCount: atomic.LoadUint64(&i.stats.MQErrorCount),
		Workers:      atomic.LoadInt32(&i.stats.Workers),
	}

}

func (i *NetflowV5) dynWorkers() {
	var load, nSeq, newWorkers, workers, n int

	tick := time.Tick(120 * time.Second)

	for {
		<-tick
		load = 0

		for n = 0; n < 30; n++ {
			time.Sleep(1 * time.Second)
			load += len(netflowV5UDPCh)
		}

		if load > 15 {

			switch {
			case load > 300:
				newWorkers = 100
			case load > 200:
				newWorkers = 60
			case load > 100:
				newWorkers = 40
			default:
				newWorkers = 30
			}

			workers = int(atomic.LoadInt32(&i.stats.Workers))
			if workers+newWorkers > maxWorkers {
				logger.Println("netflow v5 :: max out workers")
				continue
			}

			for n = 0; n < newWorkers; n++ {
				go func() {
					atomic.AddInt32(&i.stats.Workers, 1)
					wQuit := make(chan struct{})
					i.pool <- wQuit
					i.netflowV5Worker(wQuit)
				}()
			}

		}

		if load == 0 {
			nSeq++
		} else {
			nSeq = 0
			continue
		}

		if nSeq > 15 {
			for n = 0; n < 10; n++ {
				if len(i.pool) > i.workers {
					atomic.AddInt32(&i.stats.Workers, -1)
					wQuit := <-i.pool
					close(wQuit)
				}
			}

			nSeq = 0
		}
	}

}
