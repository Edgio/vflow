// Package hammer generates ipfix packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    hammer.go
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
package hammer

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/VerizonDigital/vflow/mirror"
)

const (
	dataType = iota
	templateType
	templateOptType
	sFlowDataType
)

type Packet struct {
	payload []byte
	length  int
}

type IPFIX struct {
	conn  mirror.Conn
	ch    chan Packet
	srcs  []net.IP
	vflow net.IP

	MaxRouter int
	Tick      time.Duration
	Port      int
	RateLimit int
}

type SFlow struct {
	conn  mirror.Conn
	ch    chan Packet
	srcs  []net.IP
	vflow net.IP

	MaxRouter int
	Port      int
}

func NewIPFIX(raddr net.IP) (*IPFIX, error) {

	conn, err := mirror.NewRawConn(raddr)
	if err != nil {
		return nil, err
	}

	return &IPFIX{
		conn:      conn,
		ch:        make(chan Packet, 10000),
		vflow:     raddr,
		MaxRouter: 10,
	}, nil
}

func (i *IPFIX) Run() {
	var wg sync.WaitGroup

	for j := 1; j < i.MaxRouter; j++ {
		i.srcs = append(i.srcs, net.ParseIP(fmt.Sprintf("192.168.1.%d", j)))
	}

	wg.Add(1)
	go func() {
		var p Packet
		defer wg.Done()
		for {
			p = <-i.ch
			i.conn.Send(p.payload[:p.length])
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i.sendTemplate()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i.sendTemplateOpt()
	}()

	time.Sleep(1 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		i.sendData()
	}()

	wg.Wait()
}

func (i *IPFIX) sendData() {
	packets := i.genPackets(dataType)

	for {
		for j := range packets {
			i.ch <- packets[j]
		}
	}
}

func (i *IPFIX) sendTemplate() {
	c := time.Tick(i.Tick)
	packets := i.genPackets(templateType)

	for j := range packets {
		i.ch <- packets[j]
	}

	for range c {
		for j := range packets {
			i.ch <- packets[j]
		}
	}
}

func (i *IPFIX) sendTemplateOpt() {
	c := time.Tick(i.Tick)
	packets := i.genPackets(templateOptType)

	for j := range packets {
		i.ch <- packets[j]
	}

	for range c {
		for j := range packets {
			i.ch <- packets[j]
		}
	}
}

func (i *IPFIX) genPackets(typ int) []Packet {
	var (
		packets []Packet
		samples [][]byte
	)

	ipHLen := mirror.IPv4HLen
	udp := mirror.UDP{55117, i.Port, 0, 0}
	udpHdr := udp.Marshal()

	ip := mirror.NewIPv4HeaderTpl(mirror.UDPProto)
	ipHdr := ip.Marshal()

	switch typ {
	case dataType:
		samples = ipfixDataSamples
	case templateType:
		samples = ipfixTemplates
	case templateOptType:
		samples = ipfixTemplatesOpt
	case sFlowDataType:
		samples = sFlowDataSamples
	}

	for j := 0; j < len(samples); j++ {
		for _, src := range i.srcs {
			data := samples[j]
			payload := make([]byte, 1500)

			udp.SetLen(udpHdr, len(data))

			ip.SetAddrs(ipHdr, src, i.vflow)

			copy(payload[0:ipHLen], ipHdr)
			copy(payload[ipHLen:ipHLen+8], udpHdr)
			copy(payload[ipHLen+8:], data)

			packets = append(packets, Packet{
				payload: payload,
				length:  ipHLen + 8 + len(data),
			})

		}
	}

	return packets
}

func NewSFlow(raddr net.IP) (*SFlow, error) {

	conn, err := mirror.NewRawConn(raddr)
	if err != nil {
		return nil, err
	}

	return &SFlow{
		conn:      conn,
		ch:        make(chan Packet, 10000),
		vflow:     raddr,
		MaxRouter: 10,
	}, nil
}

func (s *SFlow) Run() {
	var wg sync.WaitGroup

	for j := 1; j < s.MaxRouter; j++ {
		s.srcs = append(s.srcs, net.ParseIP(fmt.Sprintf("192.168.1.%d", j)))
	}

	wg.Add(1)
	go func() {
		var p Packet
		defer wg.Done()
		for {
			p = <-s.ch
			s.conn.Send(p.payload[:p.length])
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.sendData()
	}()

	wg.Wait()
}

func (s *SFlow) genPackets() []Packet {
	var packets []Packet
	ipHLen := mirror.IPv4HLen
	udp := mirror.UDP{55117, s.Port, 0, 0}
	udpHdr := udp.Marshal()

	ip := mirror.NewIPv4HeaderTpl(mirror.UDPProto)
	ipHdr := ip.Marshal()

	for j := 0; j < len(sFlowDataSamples); j++ {
		for _, src := range s.srcs {
			data := sFlowDataSamples[j]
			payload := make([]byte, 1500)

			udp.SetLen(udpHdr, len(data))

			ip.SetAddrs(ipHdr, src, s.vflow)

			copy(payload[0:ipHLen], ipHdr)
			copy(payload[ipHLen:ipHLen+8], udpHdr)
			copy(payload[ipHLen+8:], data)

			packets = append(packets, Packet{
				payload: payload,
				length:  ipHLen + 8 + len(data),
			})

		}
	}

	return packets
}

func (s *SFlow) sendData() {
	packets := s.genPackets()

	for {
		for j := range packets {
			s.ch <- packets[j]
		}
	}
}
