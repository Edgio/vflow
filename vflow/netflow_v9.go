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
	"bytes"
	"net"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	netflow9 "github.com/guardicore/vflow/netflow/v9"
	"github.com/guardicore/vflow/producer"
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

	mCacheNF9 netflow9.MemCache

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
		port:    opts.NetflowV9Port,
		addr:    opts.NetflowV9Addr,
		workers: opts.NetflowV9Workers,
	}
}

func (i *NetflowV9) run() {
	// exit if the netflow v9 is disabled
	if !opts.NetflowV9Enabled {
		logger.Println("netflow v9 has been disabled")
		return
	}

	i.pool = make(chan chan struct{}, maxWorkers)

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

	logger.Printf("netflow v9 is running (UDP: listening on [::]:%d workers#: %d)", i.port, i.workers)

	mCacheNF9 = netflow9.GetCache(opts.NetflowV9TplCacheFile)

	go func() {
		if !opts.ProducerEnabled {
			return
		}

		p := producer.NewProducer(opts.MQName)
		p.MQConfigFile = path.Join(opts.VFlowConfigPath, opts.MQConfigFile)
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
	// exit if the netflow v9 is disabled
	if !opts.NetflowV9Enabled {
		return
	}

	// stop reading from UDP listener
	i.stop = true
	logger.Println("stopping netflow v9 service gracefully ...")
	time.Sleep(1 * time.Second)

	// dump the templates to storage
	if err := mCacheNF9.Dump(opts.NetflowV9TplCacheFile); err != nil {
		logger.Println("couldn't not dump template", err)
	}

	// logging and close UDP channel
	logger.Println("netflow v9 has been shutdown")
	close(netflowV9UDPCh)
}

func (i *NetflowV9) netflowV9Worker(wQuit chan struct{}) {
	var (
		decodedMsg *netflow9.Message
		msg        = NetflowV9UDPMsg{body: netflowV9Buffer.Get().([]byte)}
		buf        = new(bytes.Buffer)
		err        error
		ok         bool
		b          []byte
	)

LOOP:
	for {

		netflowV9Buffer.Put(msg.body[:opts.NetflowV9UDPSize])
		buf.Reset()

		select {
		case <-wQuit:
			break LOOP
		case msg, ok = <-netflowV9UDPCh:
			if !ok {
				break LOOP
			}
		}

		if opts.Verbose {
			logger.Printf("rcvd netflow v9 data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		d := netflow9.NewDecoder(msg.raddr.IP, msg.body)
		if decodedMsg, err = d.Decode(mCacheNF9); err != nil {
			logger.Println(err)
			if decodedMsg == nil {
				continue
			}
		}

		atomic.AddUint64(&i.stats.DecodedCount, 1)

		if decodedMsg.DataSets != nil {
			b, err = decodedMsg.JSONMarshal(buf)
			if err != nil {
				logger.Println(err)
				continue
			}

			select {
			case netflowV9MQCh <- append([]byte{}, b...):
			default:
			}
		}

		if opts.Verbose {
			logger.Println(string(b))
		}

	}

}

func (i *NetflowV9) status() *NetflowV9Stats {
	return &NetflowV9Stats{
		UDPQueue:     len(netflowV9UDPCh),
		MessageQueue: len(netflowV9MQCh),
		UDPCount:     atomic.LoadUint64(&i.stats.UDPCount),
		DecodedCount: atomic.LoadUint64(&i.stats.DecodedCount),
		MQErrorCount: atomic.LoadUint64(&i.stats.MQErrorCount),
		Workers:      atomic.LoadInt32(&i.stats.Workers),
	}

}

func (i *NetflowV9) dynWorkers() {
	var load, nSeq, newWorkers, workers, n int

	tick := time.Tick(120 * time.Second)

	for {
		<-tick
		load = 0

		for n = 0; n < 30; n++ {
			time.Sleep(1 * time.Second)
			load += len(netflowV9UDPCh)
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
				logger.Println("netflow v9 :: max out workers")
				continue
			}

			for n = 0; n < newWorkers; n++ {
				go func() {
					atomic.AddInt32(&i.stats.Workers, 1)
					wQuit := make(chan struct{})
					i.pool <- wQuit
					i.netflowV9Worker(wQuit)
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
