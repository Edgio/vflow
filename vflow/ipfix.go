//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix.go
//: details: ipfix decoders handler
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
	"bytes"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VerizonDigital/vflow/ipfix"
	"github.com/VerizonDigital/vflow/mirror"
	"github.com/VerizonDigital/vflow/producer"
)

// IPFIX represents IPFIX collector
type IPFIX struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   IPFIXStats
	pool    chan chan struct{}
}

// IPFIXUDPMsg represents IPFIX UDP data
type IPFIXUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

// IPFIXStats represents IPFIX stats
type IPFIXStats struct {
	UDPQueue       int
	UDPMirrorQueue int
	MessageQueue   int
	UDPCount       uint64
	DecodedCount   uint64
	MQErrorCount   uint64
	Workers        int32
}

var (
	ipfixUDPCh         = make(chan IPFIXUDPMsg, 1000)
	ipfixMCh           = make(chan IPFIXUDPMsg, 1000)
	ipfixMQCh          = make(chan []byte, 1000)
	ipfixMirrorEnabled bool

	// templates memory cache
	mCache ipfix.MemCache

	// ipfix udp payload pool
	ipfixBuffer = &sync.Pool{
		New: func() interface{} {
			return make([]byte, opts.IPFIXUDPSize)
		},
	}
)

// NewIPFIX constructs IPFIX
func NewIPFIX() *IPFIX {
	return &IPFIX{
		port:    opts.IPFIXPort,
		workers: opts.IPFIXWorkers,
		pool:    make(chan chan struct{}, maxWorkers),
	}
}

func (i *IPFIX) run() {
	// exit if the ipfix is disabled
	if !opts.IPFIXEnabled {
		logger.Println("ipfix has been disabled")
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
			i.ipfixWorker(wQuit)
		}()
	}

	logger.Printf("ipfix is running (workers#: %d)", i.workers)

	mCache = ipfix.GetCache(opts.IPFIXTplCacheFile)
	go ipfix.RPC(mCache, &ipfix.RPCConfig{
		Enabled: opts.IPFIXRPCEnabled,
		Logger:  logger,
	})

	go mirrorIPFIXDispatcher(ipfixMCh)

	go func() {
		p := producer.NewProducer(opts.MQName)

		p.MQConfigFile = opts.MQConfigFile
		p.MQErrorCount = &i.stats.MQErrorCount
		p.Logger = logger
		p.Chan = ipfixMQCh
		p.Topic = "ipfix"

		if err := p.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		if !opts.DynWorkers {
			logger.Println("IPFIX dynamic worker disabled")
			return
		}

		i.dynWorkers()
	}()

	for !i.stop {
		b := ipfixBuffer.Get().([]byte)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		atomic.AddUint64(&i.stats.UDPCount, 1)
		ipfixUDPCh <- IPFIXUDPMsg{raddr, b[:n]}
	}

}

func (i *IPFIX) shutdown() {
	// exit if the ipfix is disabled
	if !opts.IPFIXEnabled {
		logger.Println("ipfix disabled")
		return
	}

	// stop reading from UDP listener
	i.stop = true
	logger.Println("stopping ipfix service gracefully ...")
	time.Sleep(1 * time.Second)

	// dump the templates to storage
	if err := mCache.Dump(opts.IPFIXTplCacheFile); err != nil {
		logger.Println("couldn't not dump template", err)
	}

	// logging and close UDP channel
	logger.Println("ipfix has been shutdown")
	close(ipfixUDPCh)
}

func (i *IPFIX) ipfixWorker(wQuit chan struct{}) {
	var (
		decodedMsg *ipfix.Message
		mirror     IPFIXUDPMsg
		msg        = IPFIXUDPMsg{body: ipfixBuffer.Get().([]byte)}
		buf        = new(bytes.Buffer)
		err        error
		ok         bool
		b          []byte
	)

LOOP:
	for {

		ipfixBuffer.Put(msg.body[:opts.IPFIXUDPSize])
		buf.Reset()

		select {
		case <-wQuit:
			break LOOP
		case msg, ok = <-ipfixUDPCh:
			if !ok {
				break LOOP
			}
		}

		if opts.Verbose {
			logger.Printf("rcvd ipfix data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		if ipfixMirrorEnabled {
			mirror.body = ipfixBuffer.Get().([]byte)
			mirror.raddr = msg.raddr
			mirror.body = append(mirror.body[:0], msg.body...)

			select {
			case ipfixMCh <- mirror:
			default:
			}
		}

		d := ipfix.NewDecoder(msg.raddr.IP, msg.body)
		if decodedMsg, err = d.Decode(mCache); err != nil {
			logger.Println(err)
			continue
		}

		atomic.AddUint64(&i.stats.DecodedCount, 1)

		if decodedMsg.DataSets != nil {
			b, err = decodedMsg.JSONMarshal(buf)
			if err != nil {
				logger.Println(err)
				continue
			}

			select {
			case ipfixMQCh <- b:
			default:
			}
		}

		if opts.Verbose {
			logger.Println(string(b))
		}

	}
}

func (i *IPFIX) status() *IPFIXStats {
	return &IPFIXStats{
		UDPQueue:       len(ipfixUDPCh),
		UDPMirrorQueue: len(ipfixMCh),
		MessageQueue:   len(ipfixMQCh),
		UDPCount:       atomic.LoadUint64(&i.stats.UDPCount),
		DecodedCount:   atomic.LoadUint64(&i.stats.DecodedCount),
		MQErrorCount:   atomic.LoadUint64(&i.stats.MQErrorCount),
		Workers:        atomic.LoadInt32(&i.stats.Workers),
	}
}

func mirrorIPFIX(dst net.IP, port int, ch chan IPFIXUDPMsg) error {
	var (
		packet = make([]byte, 1500)
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

func (i *IPFIX) dynWorkers() {
	var load, nSeq, newWorkers, workers, n int

	tick := time.Tick(120 * time.Second)

	for {
		<-tick
		load = 0

		for n = 0; n < 30; n++ {
			time.Sleep(1 * time.Second)
			load += len(sFlowUDPCh)
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
				logger.Println("sflow :: max out workers")
				continue
			}

			for n = 0; n < newWorkers; n++ {
				go func() {
					atomic.AddInt32(&i.stats.Workers, 1)
					wQuit := make(chan struct{})
					i.pool <- wQuit
					i.ipfixWorker(wQuit)
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
