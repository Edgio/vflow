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
	"time"

	"github.com/VerizonDigital/vflow/disc"
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
	Enabled                     bool
	Port                        int
	Addr                        net.IP
	DiscoveryStrategy           string
	DiscoveryStrategyConfigFile string
	Logger                      *log.Logger
}

// RPCRequest represents RPC request
type RPCRequest struct {
	ID uint16
	IP net.IP
}

var (
	errNotAvail = errors.New("the template is not available")
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

	var disc_config disc.DiscoveryConfig
	disc_config.DiscoveryStrategy = config.DiscoveryStrategy
	disc_config.LoadConfig(config.DiscoveryStrategyConfigFile)
	disc_config.Logger = config.Logger

	disc, err := disc.BuildDiscovery(&disc_config)

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

		for _, rpcServer := range disc.GetRPCServers() {
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
