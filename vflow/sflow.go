//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    sflow.go
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
}

// SFlowStats represents sflow stats
type SFlowStats struct {
	UDPQueue     int
	MessageQueue int
	UDPCount     uint64
	DecodedCount uint64
	MQErrorCount uint64
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
	}
}

func (s *SFlow) run() {
	var wg sync.WaitGroup

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

	wg.Add(s.workers)
	for i := 0; i < s.workers; i++ {
		go func() {
			defer wg.Done()
			s.sFlowWorker()

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

	wg.Wait()
}

func (s *SFlow) shutdown() {
	s.stop = true
	logger.Println("stopping sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("vFlow has been shutdown")
	close(sFlowUDPCh)
}

func (s *SFlow) sFlowWorker() {
	var (
		filter = []uint32{sflow.DataCounterSample}
		reader *bytes.Reader
		msg    SFUDPMsg
		ok     bool
		b      []byte
	)

	for {
		if msg, ok = <-sFlowUDPCh; !ok {
			break
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
	}
}
