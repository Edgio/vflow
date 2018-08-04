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

package ipfix

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// IRPC represents IPFIX RPC
type IRPC struct {
	reqCount uint64
	mCache   MemCache
}

// RPCClient represents RPC client
type RPCClient struct {
	conn *rpc.Client
}

// RPCConfig represents RPC config
type RPCConfig struct {
	Enabled bool
	Port    int
	Addr    net.IP
	Logger  *log.Logger
}

// RPCRequest represents RPC request
type RPCRequest struct {
	ID uint16
	IP net.IP
}

type vFlowServer struct {
	timestamp int64
}

// Discovery represents vflow discovery
type Discovery struct {
	conn         interface{}
	group        net.IP
	port         int
	rcvdMsg      chan net.IP
	vFlowServers map[string]vFlowServer
	mu           sync.RWMutex
}

var (
	errNotAvail            = errors.New("the template is not available")
	errMCInterfaceNotAvail = errors.New("multicast interface not available")
)

// NewRPC constructs RPC
func NewRPC(mCache MemCache) *IRPC {
	return &IRPC{
		reqCount: 0,
		mCache:   mCache,
	}
}

// Get retrieves a request from mCache
func (r *IRPC) Get(req RPCRequest, resp *TemplateRecord) error {
	var ok bool

	*resp, ok = r.mCache.retrieve(req.ID, req.IP)
	if !ok {
		return errNotAvail
	}

	return nil
}

// RPCServer runs the RPC server
func RPCServer(mCache MemCache, config *RPCConfig) error {
	rpc.Register(NewRPC(mCache))
	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		return err
	}

	rpc.Accept(l)

	return nil
}

// NewRPCClient initializes a new client connection
func NewRPCClient(r string) (*RPCClient, error) {
	raddr := net.JoinHostPort(r, "8085")

	conn, err := net.DialTimeout("tcp", raddr, 1*time.Second)
	if err != nil {
		return nil, err
	}

	return &RPCClient{conn: rpc.NewClient(conn)}, nil
}

// Get tries to get a request from remote server
func (c *RPCClient) Get(req RPCRequest) (*TemplateRecord, error) {
	var tr *TemplateRecord

	err := c.conn.Call("IRPC.Get", req, &tr)

	return tr, err
}

// RPC handles RPC with discovery
func RPC(m MemCache, config *RPCConfig) {
	if !config.Enabled {
		return
	}

	disc, err := vFlowDiscovery()
	if err != nil {
		config.Logger.Println(err)
		config.Logger.Println("RPC has been disabled")
		return
	}

	go RPCServer(m, config)

	config.Logger.Println("ipfix RPC enabled")
	throttle := time.Tick(time.Duration(1e6/10) * time.Microsecond)

	for {
		req := <-rpcChan

		for _, rpcServer := range disc.rpcServers() {
			r, err := NewRPCClient(rpcServer)
			if err != nil {
				config.Logger.Println(err)
				continue
			}

			tr, err := r.Get(req)
			r.conn.Close()

			if err != nil {
				continue
			}

			m.insert(req.ID, req.IP, *tr)
			break
		}

		<-throttle
	}
}
func vFlowDiscovery() (*Discovery, error) {
	// TODO
	disc := &Discovery{
		group: net.ParseIP("224.0.0.55"),
		port:  1024,
	}

	if err := disc.mConn(); err != nil {
		return nil, err
	}

	if disc.group.To4() != nil {
		go disc.startV4()
	} else {
		go disc.startV6()
	}

	return disc, nil
}

func (d *Discovery) mConn() error {
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

func (d *Discovery) receiverV4() {
	var b = make([]byte, 1500)

	d.vFlowServers = make(map[string]vFlowServer, 10)
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
		d.vFlowServers[host] = vFlowServer{time.Now().Unix()}
		d.mu.Unlock()
	}
}

func (d *Discovery) startV4() {
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

func (d *Discovery) receiverV6() {
	var b = make([]byte, 1500)

	d.vFlowServers = make(map[string]vFlowServer, 10)
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
		d.vFlowServers[host] = vFlowServer{time.Now().Unix()}
		d.mu.Unlock()
	}
}

func (d *Discovery) startV6() {
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

func (d *Discovery) rpcServers() []string {
	var servers []string

	now := time.Now().Unix()

	d.mu.Lock()

	for ip, server := range d.vFlowServers {
		if now-server.timestamp < 300 {
			servers = append(servers, ip)
		} else {
			delete(d.vFlowServers, ip)
		}
	}

	d.mu.Unlock()

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
