//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    sflow.go
//: details: sflow decoding raw data samples
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
	"encoding/json"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VerizonDigital/vflow/packet"
	"github.com/VerizonDigital/vflow/producer"
	"github.com/VerizonDigital/vflow/sflow"
)

// SFUDPMsg represents sFlow UDP message
type SFUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

// SFlow represents sFlow collector
type SFlow struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   SFlowStats
	pool    chan chan struct{}
}

// SFlowStats represents sflow stats
type SFlowStats struct {
	UDPQueue     int
	MessageQueue int
	UDPCount     uint64
	DecodedCount uint64
	MQErrorCount uint64
	Workers      int32
}

var (
	sFlowUDPCh = make(chan SFUDPMsg, 1000)
	sFlowMQCh  = make(chan []byte, 1000)

	// sflow udp payload pool
	sFlowBuffer = &sync.Pool{
		New: func() interface{} {
			return make([]byte, opts.SFlowUDPSize)
		},
	}
)

// NewSFlow constructs sFlow collector
func NewSFlow() *SFlow {
	logger = opts.Logger

	return &SFlow{
		port:    opts.SFlowPort,
		workers: opts.SFlowWorkers,
		pool:    make(chan chan struct{}, maxWorkers),
	}
}

func (s *SFlow) run() {
	// exit if the sflow is disabled
	if !opts.SFlowEnabled {
		logger.Println("sflow has been disabled")
		return
	}

	hostPort := net.JoinHostPort(s.addr, strconv.Itoa(s.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal(err)
	}

	atomic.AddInt32(&s.stats.Workers, int32(s.workers))
	for i := 0; i < s.workers; i++ {
		go func() {
			wQuit := make(chan struct{})
			s.pool <- wQuit
			s.sFlowWorker(wQuit)
		}()
	}

	logger.Printf("sFlow is running (workers#: %d)", s.workers)

	go func() {
		p := producer.NewProducer(opts.MQName)

		p.MQConfigFile = opts.MQConfigFile
		p.MQErrorCount = &s.stats.MQErrorCount
		p.Logger = logger
		p.Chan = sFlowMQCh
		p.Topic = "sflow"

		if err := p.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		if !opts.DynWorkers {
			logger.Println("sFlow dynamic worker disabled")
			return
		}

		s.dynWorkers()
	}()

	for !s.stop {
		b := sFlowBuffer.Get().([]byte)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}

		atomic.AddUint64(&s.stats.UDPCount, 1)
		sFlowUDPCh <- SFUDPMsg{raddr, b[:n]}
	}

}

func (s *SFlow) shutdown() {
	s.stop = true
	logger.Println("stopping sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("vFlow has been shutdown")
	close(sFlowUDPCh)
}

func (s *SFlow) sFlowWorker(wQuit chan struct{}) {
	var (
		filter = []uint32{sflow.DataCounterSample}
		reader *bytes.Reader
		msg    SFUDPMsg
		ok     bool
		b      []byte
	)

LOOP:
	for {

		select {
		case <-wQuit:
			break LOOP
		case msg, ok = <-sFlowUDPCh:
			if !ok {
				break LOOP
			}
		}

		if opts.Verbose {
			logger.Printf("rcvd sflow data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		reader = bytes.NewReader(msg.body)
		d := sflow.NewSFDecoder(reader, filter)
		records, err := d.SFDecode()
		if err != nil || len(records) < 1 {
			sFlowBuffer.Put(msg.body[:opts.SFlowUDPSize])
			continue
		}

		decodedMsg := sflow.Message{}

		for _, data := range records {
			switch data.(type) {
			case *packet.Packet:
				decodedMsg.Packet = data.(*packet.Packet)
			case *sflow.ExtSwitchData:
				decodedMsg.ExtSWData = data.(*sflow.ExtSwitchData)
			case *sflow.FlowSample:
				decodedMsg.Sample = data.(*sflow.FlowSample)
			case *sflow.SFDatagram:
				decodedMsg.Header = data.(*sflow.SFDatagram)
			}
		}

		b, err = json.Marshal(decodedMsg)
		if err != nil {
			sFlowBuffer.Put(msg.body[:opts.SFlowUDPSize])
			logger.Println(err)
			continue
		}

		atomic.AddUint64(&s.stats.DecodedCount, 1)

		if opts.Verbose {
			logger.Println(string(b))
		}

		select {
		case sFlowMQCh <- b:
		default:
		}

		sFlowBuffer.Put(msg.body[:opts.SFlowUDPSize])
	}
}

func (s *SFlow) status() *SFlowStats {
	return &SFlowStats{
		UDPQueue:     len(sFlowUDPCh),
		MessageQueue: len(sFlowMQCh),
		UDPCount:     atomic.LoadUint64(&s.stats.UDPCount),
		DecodedCount: atomic.LoadUint64(&s.stats.DecodedCount),
		MQErrorCount: atomic.LoadUint64(&s.stats.MQErrorCount),
		Workers:      atomic.LoadInt32(&s.stats.Workers),
	}
}

func (s *SFlow) dynWorkers() {
	var load, normalSeq, newWorkers, n int

	tick := time.Tick(120 * time.Second)

	for {
		<-tick
		load = 0

		for n = 0; n < 30; n++ {
			time.Sleep(1 * time.Second)
			load += len(sFlowUDPCh)
		}

		if load > 15 && s.stats.Workers < int32(maxWorkers) {

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

			for n = 0; n < newWorkers; n++ {
				go func() {
					atomic.AddInt32(&s.stats.Workers, 1)
					wQuit := make(chan struct{})
					s.pool <- wQuit
					s.sFlowWorker(wQuit)
				}()
			}

		}

		if load == 0 {
			normalSeq++
		} else {
			normalSeq = 0
			continue
		}

		if normalSeq > 15 {
			for n = 0; n < 10; n++ {
				if len(s.pool) > s.workers {
					atomic.AddInt32(&s.stats.Workers, -1)
					wQuit := <-s.pool
					close(wQuit)
				}
			}
			normalSeq = 0
		}
	}
}
