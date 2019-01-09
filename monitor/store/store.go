//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    store.go
//: details: interface to other store back-end
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

package store

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Monitor is an interface to store system
// and netflow statistics
type Monitor interface {
	System(string) error
	Netflow(string) error
}

// IPFIX represents IPFIX metrics
type IPFIX struct {
	UDPQueue       int64
	UDPMirrorQueue int64
	MessageQueue   int64
	UDPCount       int64
	DecodedCount   int64
	MQErrorCount   int64
	Workers        int64
}

// SFlow represents SFlow metrics
type SFlow struct {
	UDPQueue     int64
	MessageQueue int64
	UDPCount     int64
	DecodedCount int64
	MQErrorCount int64
	Workers      int64
}

// NetflowV5 represents Netflow v5 metrics
type NetflowV5 struct {
	UDPQueue     int64
	MessageQueue int64
	UDPCount     int64
	DecodedCount int64
	MQErrorCount int64
	Workers      int64
}

// NetflowV9 represents Netflow v9 metrics
type NetflowV9 struct {
	UDPQueue     int64
	MessageQueue int64
	UDPCount     int64
	DecodedCount int64
	MQErrorCount int64
	Workers      int64
}

// Flow represents flow (IPFIX+sFlow) metrics
type Flow struct {
	StartTime int64
	Timestamp int64
	IPFIX     IPFIX
	SFlow     SFlow
	NetflowV5 NetflowV5
	NetflowV9 NetflowV9
}

// Sys represents system/go-runtime statistics
type Sys struct {
	MemHeapAlloc    int64
	MemAlloc        int64
	MCacheInuse     int64
	GCNext          int64
	MemTotalAlloc   int64
	GCSys           int64
	MemHeapSys      int64
	NumGoroutine    int64
	NumLogicalCPU   int64
	MemHeapReleased int64
}

// Client represents HTTP client
type Client struct {
	client *http.Client
}

// NewHTTP constructs HTTP client
func NewHTTP() *Client {
	return &Client{
		client: new(http.Client),
	}
}

// Get tries to get metrics through HTTP w/ get method
func (c *Client) Get(url string, s interface{}) error {

	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &s)
	if err != nil {
		return err
	}

	return nil
}

// Post tries to digest metrics through HTTP w/ post method
func (c *Client) Post(url string, cType, query string) ([]byte, error) {
	resp, err := c.client.Post(url, cType, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getFlow(vhost, host string) (*Flow, *Flow, error) {
	lastFlowFile := "/tmp/vflow.mon.lastflow." + host

	flow := new(Flow)
	lastFlow := new(Flow)

	client := NewHTTP()
	err := client.Get(vhost+"/flow", flow)
	if err != nil {
		return nil, nil, err
	}

	flow.Timestamp = time.Now().Unix()

	b, err := ioutil.ReadFile(lastFlowFile)
	if err != nil {
		b, _ = json.Marshal(flow)
		ioutil.WriteFile(lastFlowFile, b, 0644)
		return nil, nil, err
	}

	err = json.Unmarshal(b, &lastFlow)
	if err != nil {
		return nil, nil, err
	}

	b, err = json.Marshal(flow)
	if err != nil {
		return nil, nil, err
	}

	ioutil.WriteFile(lastFlowFile, b, 0644)

	// once the vFlow restarted
	if flow.StartTime != lastFlow.StartTime {
		os.Exit(1)
	}

	return flow, lastFlow, err
}
