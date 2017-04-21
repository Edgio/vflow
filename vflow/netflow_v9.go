//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    netflow_v9.go
//: details: netflow decoders handler
//: author:  Mehrdad Arshad Rad
//: date:    04/21/2017
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
	"net"
)

// NetflowV9 represents netflow v9 collector
type NetflowV9 struct {
	port    int
	addr    string
	workers int
	stop    bool
	stats   NetflowV9Stats
	pool    chan chan struct{}
}

// NetflowV9UDPMsg represents netflow v9 UDP data
type NetflowV9UDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

// NetflowV9Stats represents netflow v9 stats
type NetflowV9Stats struct {
	UDPQueue     int
	MessageQueue int
	UDPCount     uint64
	DecodedCount uint64
	MQErrorCount uint64
	Workers      int32
}

// TODO
