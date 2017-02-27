// Package ipfix decodes IPFIX packets
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
	"net"
	"net/rpc"
	"time"
)

// RPC represents RPC
type RPC struct {
	reqCount uint64
	mCache   MemCache
}

// RPCClient represents RPC client
type RPCClient struct {
	conn *rpc.Client
}

// RPCConfig represents RPC config
type RPCConfig struct {
	enabled bool
	port    int
	addr    net.IP
}

// RPCRequest represents RPC request
type RPCRequest struct {
	ID uint16
	IP net.IP
}

var (
	vFlowServers []string
	errNotAvail  = errors.New("the template is not available")
)

// NewRPC constructs RPC
func NewRPC(mCache MemCache) *RPC {
	return &RPC{
		reqCount: 0,
		mCache:   mCache,
	}
}

// Get retrieves a request from mCache
func (r *RPC) Get(req RPCRequest, resp *TemplateRecords) error {
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
func NewRPCClient(rAddr string) *RPCClient {
	conn, _ := net.DialTimeout("tcp", rAddr, 1*time.Second)
	return &RPCClient{conn: rpc.NewClient(conn)}
}

// Get tries to get a request from remote server
func (c *RPCClient) Get(req RPCRequest) (*TemplateRecords, error) {
	var tr *TemplateRecords
	err := c.conn.Call("RPC.Get", req, &tr)

	return tr, err
}
