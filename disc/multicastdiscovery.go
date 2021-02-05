//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    memcache_rpc.go
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

package disc

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Discovery represents vflow discovery
type MulticastDiscovery struct {
	discovery Discovery
	conn      interface{}
	group     net.IP
	port      int
	rcvdMsg   chan net.IP
	mu        sync.RWMutex
}

func (d *MulticastDiscovery) GetvFlowServers() map[string]vFlowServer {
	return d.discovery.GetvFlowServers()
}

func (d *MulticastDiscovery) GetRPCServers() []string {
	return BuildRpcServersList(d.GetvFlowServers())
}

func (disc *MulticastDiscovery) Setup(config *DiscoveryConfig) error {

	logger = config.Logger
	disc.discovery = NewDiscovery(config)
	disc.group = net.ParseIP("224.0.0.55")
	disc.port = 1024

	if err := disc.mConn(); err != nil {
		return err
	}

	if disc.group.To4() != nil {
		go disc.startV4()
	} else {
		go disc.startV6()
	}

	logger.Println("Multicast discovery started")

	return nil
}

func (d *MulticastDiscovery) mConn() error {
	addr := net.JoinHostPort("", strconv.Itoa(d.port))
	c, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	ifs, err := getMulticastIfs()
	if err != nil {
		return err
	}

	if d.group.To4() != nil {
		d.conn = ipv4.NewPacketConn(c)
		for _, i := range ifs {
			err = d.conn.(*ipv4.PacketConn).JoinGroup(
				&i,
				&net.UDPAddr{IP: d.group},
			)
			if err != nil {
				return err
			}
		}
	} else {
		d.conn = ipv6.NewPacketConn(c)
		for _, i := range ifs {
			err = d.conn.(*ipv6.PacketConn).JoinGroup(
				&i,
				&net.UDPAddr{IP: d.group},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *MulticastDiscovery) receiverV4() {
	var b = make([]byte, 1500)

	conn := d.conn.(*ipv4.PacketConn)
	laddrs, err := getLocalIPs()
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, _, addr, err := conn.ReadFrom(b)
		if err != nil {
			continue
		}

		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			continue
		}

		if _, ok := laddrs[host]; ok {
			continue
		}

		d.mu.Lock()
		d.GetvFlowServers()[host] = vFlowServer{time.Now().Unix()}
		d.mu.Unlock()
	}
}

func (d *MulticastDiscovery) startV4() {
	tick := time.NewTicker(1 * time.Second)

	b := []byte("Hello vFlow")
	conn := d.conn.(*ipv4.PacketConn)
	conn.SetTTL(2)
	go d.receiverV4()

	for {
		<-tick.C
		conn.WriteTo(b, nil, &net.UDPAddr{IP: d.group, Port: d.port})
	}
}

func (d *MulticastDiscovery) receiverV6() {
	var b = make([]byte, 1500)

	conn := d.conn.(*ipv6.PacketConn)
	laddrs, err := getLocalIPs()
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, _, addr, err := conn.ReadFrom(b)
		if err != nil {
			continue
		}

		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			continue
		}

		if _, ok := laddrs[host]; ok {
			continue
		}

		d.mu.Lock()
		d.GetvFlowServers()[host] = vFlowServer{time.Now().Unix()}
		d.mu.Unlock()
	}
}

func (d *MulticastDiscovery) startV6() {
	tick := time.NewTicker(1 * time.Second)

	b := []byte("Hello vFlow")
	conn := d.conn.(*ipv6.PacketConn)
	conn.SetHopLimit(2)
	go d.receiverV6()

	for {
		<-tick.C
		conn.WriteTo(b, nil, &net.UDPAddr{IP: d.group, Port: d.port})
	}
}

func getMulticastIfs() ([]net.Interface, error) {
	var out []net.Interface

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifs {
		if i.Flags == 19 {
			out = append(out, i)
		}
	}

	if len(out) < 1 {
		return nil, errMCInterfaceNotAvail
	}

	return out, nil
}
