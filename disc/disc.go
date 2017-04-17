// Package discovery handles finding vFlow nodes through multicasting
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    disc.go
//: details: discovery vFlow nodes by multicasting
//: author:  Mehrdad Arshad Rad
//: date:    04/17/2017
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
package discovery

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type vFlowServer struct {
	timestamp int64
}

// Discovery represents vflow discovery
type Discovery struct {
	vFlowServers map[string]vFlowServer
	mu           sync.RWMutex
}

var errMCInterfaceNotAvail = errors.New("multicast interface not available")

// Run starts sending multicast hello packet
func Run(ip, port string) error {
	tick := time.NewTicker(1 * time.Second)

	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: p,
	})

	b := []byte("hello vflow")

	if err != nil {
		return err
	}

	for {
		<-tick.C
		c.Write(b)
	}
}

// Listen receives discovery hello packet
func Listen(ip, port string) (*Discovery, error) {
	var (
		conn interface{}
		buff = make([]byte, 1500)
		disc = &Discovery{
			vFlowServers: make(map[string]vFlowServer, 10),
		}
	)

	c, err := net.ListenPacket("udp", net.JoinHostPort(
		ip,
		port,
	))

	if err != nil {
		return nil, err
	}

	ifs, err := getMulticastIfs()
	if err != nil {
		return nil, err
	}

	if net.ParseIP(ip).To4() != nil {
		conn = ipv4.NewPacketConn(c)
		for _, i := range ifs {
			err = conn.(*ipv4.PacketConn).JoinGroup(
				&i,
				&net.UDPAddr{IP: net.ParseIP(ip)},
			)
			if err != nil {
				return nil, err
			}
		}
	} else {
		conn = ipv6.NewPacketConn(c)
		for _, i := range ifs {
			err = conn.(*ipv4.PacketConn).JoinGroup(
				&i,
				&net.UDPAddr{IP: net.ParseIP(ip)},
			)
			if err != nil {
				return nil, err
			}
		}
	}

	laddrs, err := getLocalIPs()
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		var (
			addr net.Addr
			err  error
		)

		for {

			if net.ParseIP(ip).To4() != nil {
				_, _, addr, err = conn.(*ipv4.PacketConn).ReadFrom(buff)
			} else {
				_, _, addr, err = conn.(*ipv6.PacketConn).ReadFrom(buff)
			}

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

			disc.mu.Lock()
			disc.vFlowServers[host] = vFlowServer{time.Now().Unix()}
			disc.mu.Unlock()

		}

	}()

	return disc, nil
}

// Nodes returns a slice of available vFlow nodes
func (d *Discovery) Nodes() []string {
	var servers []string

	now := time.Now().Unix()
	for ip, server := range d.vFlowServers {
		if now-server.timestamp < 300 {
			servers = append(servers, ip)
		} else {
			delete(d.vFlowServers, ip)
		}
	}
	return servers
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

func getLocalIPs() (map[string]struct{}, error) {
	ips := make(map[string]struct{})

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifs {
		addrs, err := i.Addrs()
		if err != nil || i.Flags != 19 {
			continue
		}
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			ips[ip.String()] = struct{}{}
		}
	}

	return ips, nil
}
