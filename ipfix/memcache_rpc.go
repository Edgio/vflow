// Package ipfix decodes IPFIX packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder.go
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

type RPC struct {
	reqCount uint64
	mCache   MemCache
}

type RPCClient struct {
	conn *rpc.Client
}

type Request struct {
	ID uint16
	IP net.IP
}

var errNotAvail = errors.New("the template is not available")

func NewRPC(mCache MemCache) *RPC {
	return &RPC{
		reqCount: 0,
		mCache:   mCache,
	}
}

func (r *RPC) Get(req Request, resp *TemplateRecords) error {
	var ok bool

	*resp, ok = r.mCache.retrieve(req.ID, req.IP)
	if !ok {
		return errNotAvail
	}

	return nil
}

func RPCServer(mCache MemCache) error {
	rpc.Register(NewRPC(mCache))
	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		return err
	}

	rpc.Accept(l)

	return nil
}

func NewRPCClient(rAddr string) *RPCClient {
	conn, _ := net.DialTimeout("tcp", rAddr, 1*time.Second)
	return &RPCClient{conn: rpc.NewClient(conn)}
}

func (c *RPCClient) Get(req Request) (*TemplateRecords, error) {
	var tr *TemplateRecords
	err := c.conn.Call("RPC.Get", req, &tr)

	return tr, err
}
