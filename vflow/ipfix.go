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
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.edgecastcdn.net/vflow/ipfix"
	"git.edgecastcdn.net/vflow/mirror"
	"git.edgecastcdn.net/vflow/producer"
)

type IPFIX struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   IPFIXStats
}

type IPFIXUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

type IPFIXStats struct {
	UDPQueue       int
	UDPMirrorQueue int
	MessageQueue   int
	UDPCount       uint64
	DecodedCount   uint64
}

var (
	ipfixUdpCh         = make(chan IPFIXUDPMsg, 1000)
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

func NewIPFIX() *IPFIX {
	return &IPFIX{
		port:    opts.IPFIXPort,
		workers: opts.IPFIXWorkers,
	}
}

func (i *IPFIX) run() {
	var wg sync.WaitGroup

	// exit if the ipfix is disabled
	if !opts.IPFIXEnabled {
		logger.Println("ipfix disabled")
		return
	}

	hostPort := net.JoinHostPort(i.addr, strconv.Itoa(i.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {

	}

	for n := 0; n < i.workers; n++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			i.ipfixWorker()

		}()
	}

	logger.Printf("ipfix is running (workers#: %d)", i.workers)

	mCache = ipfix.GetCache(opts.IPFIXTemplateCacheFile)

	go func() {
		mirrorIPFIXDispatcher(ipfixMCh)
	}()

	go func() {
		p := producer.NewProducer(opts.MQName)

		p.MQConfigFile = opts.MQConfigFile
		p.Logger = logger
		p.Chan = ipfixMQCh
		p.Topic = "ipfix"

		if err := p.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	for !i.stop {
		b := ipfixBuffer.Get().([]byte)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		atomic.AddUint64(&i.stats.UDPCount, 1)
		ipfixUdpCh <- IPFIXUDPMsg{raddr, b[:n]}
	}

	wg.Wait()
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
	if err := mCache.Dump(opts.IPFIXTemplateCacheFile); err != nil {
		logger.Println("couldn't not dump template", err)
	}

	// logging and close UDP channel
	logger.Println("ipfix has been shutdown")
	close(ipfixUdpCh)
}

func (i *IPFIX) ipfixWorker() {
	var (
		decodedMsg *ipfix.Message
		msg        IPFIXUDPMsg
		err        error
		ok         bool
		b          []byte
	)

	for {
		if msg, ok = <-ipfixUdpCh; !ok {
			break
		}

		if opts.Verbose {
			logger.Printf("rcvd ipfix data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		if ipfixMirrorEnabled {
			ipfixMCh <- IPFIXUDPMsg{msg.raddr, append([]byte{}, msg.body...)}
		}

		d := ipfix.NewDecoder(msg.raddr.IP, msg.body)
		if decodedMsg, err = d.Decode(mCache); err != nil {
			logger.Println(err)
			continue
		}

		b, err = json.Marshal(decodedMsg)
		if err != nil {
			logger.Println(err)
			continue
		}

		atomic.AddUint64(&i.stats.DecodedCount, 1)

		if opts.Verbose {
			logger.Println(string(b))
		}

		select {
		case ipfixMQCh <- b:
		default:
		}

		ipfixBuffer.Put(msg.body[:opts.IPFIXUDPSize])
	}
}

func (i *IPFIX) status() *IPFIXStats {
	return &IPFIXStats{
		UDPQueue:       len(ipfixUdpCh),
		UDPMirrorQueue: len(ipfixMCh),
		MessageQueue:   len(ipfixMQCh),
		UDPCount:       atomic.LoadUint64(&i.stats.UDPCount),
		DecodedCount:   atomic.LoadUint64(&i.stats.DecodedCount),
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

	if opts.IPFIXMirror == "" {
		return
	}

	for _, mirrorHostPort := range strings.Split(opts.IPFIXMirror, ";") {
		host, port, err := net.SplitHostPort(mirrorHostPort)
		if err != nil {
			logger.Fatalf("wrong ipfix mirror address %s", opts.IPFIXMirror)
		}

		portNo, _ := strconv.Atoi(port)
		dst := net.ParseIP(host)

		if dst.To4() != nil {
			go mirrorIPFIX(dst, portNo, ch4)
		} else {
			go mirrorIPFIX(dst, portNo, ch6)
		}
	}

	ipfixMirrorEnabled = true
	logger.Println("ipfix mirror service is running ...")

	for {
		msg = <-ch
		if msg.raddr.IP.To4() != nil {
			ch4 <- msg
		} else {
			ch6 <- msg
		}
	}
}
