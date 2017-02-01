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
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"git.edgecastcdn.net/vflow/packet"
	"git.edgecastcdn.net/vflow/sflow"
)

type SFUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

var (
	sFlowUdpCh = make(chan SFUDPMsg, 1000)
	logger     *log.Logger
	verbose    bool
)

type SFlow struct {
	port        int
	addr        string
	laddr       *net.UDPAddr
	readTimeout time.Duration
	udpSize     int
	workers     int
	stop        bool
}

func NewSFlow(opts *Options) *SFlow {
	logger = opts.Logger
	verbose = opts.Verbose

	return &SFlow{
		port:    opts.SFlowPort,
		udpSize: opts.SFlowUDPSize,
		workers: opts.SFlowWorkers,
	}
}

func (s *SFlow) run() {
	var wg sync.WaitGroup

	hostPort := net.JoinHostPort(s.addr, strconv.Itoa(s.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal(err)
	}

	for i := 0; i < s.workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			sFlowWorker()

		}()
	}

	logger.Printf("sFlow is running (workers#: %d)", s.workers)

	for !s.stop {
		b := make([]byte, s.udpSize)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		sFlowUdpCh <- SFUDPMsg{raddr, b[:n]}
	}

	wg.Wait()
}

func (s *SFlow) shutdown() {
	s.stop = true
	logger.Println("stopped sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("vFlow has been shutdown")
	close(sFlowUdpCh)
}

func sFlowWorker() {
	var (
		msg    SFUDPMsg
		ok     bool
		reader *bytes.Reader
		filter = []uint32{sflow.DataCounterSample}
	)

	for {
		if msg, ok = <-sFlowUdpCh; !ok {
			break
		}

		if verbose {
			logger.Printf("rcvd sflow data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		reader = bytes.NewReader(msg.body)
		d := sflow.NewSFDecoder(reader, filter)
		records, err := d.SFDecode()
		if err != nil {
			logger.Println(err)
		}
		for _, data := range records {
			switch data.(type) {
			case *packet.Packet:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			case *sflow.ExtSwitchData:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			case *sflow.FlowSample:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			}
		}
	}
}
