// Package store ingest monitoring time series data points
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    store.go
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
package store

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Monitor interface {
	System() error
	Netflow() error
}

type IPFIX struct {
	UDPQueue       int64
	UDPMirrorQueue int64
	MessageQueue   int64
	UDPCount       int64
	DecodedCount   int64
	MQErrorCount   int64
}

type SFlow struct {
	UDPQueue     int64
	MessageQueue int64
	UDPCount     int64
	DecodedCount int64
	MQErrorCount int64
}

type Flow struct {
	StartTime int64
	Timestamp int64
	IPFIX     IPFIX
	SFlow     SFlow
}

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

type Client struct {
	client *http.Client
}

func NewHTTP() *Client {
	return &Client{
		client: new(http.Client),
	}
}

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

func (c *Client) Post(url string, cType, query string) (error, []byte) {
	resp, err := c.client.Post(url, cType, bytes.NewBufferString(query))
	if err != nil {
		return err, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	return nil, body
}

func getFlow(host string) (error, *Flow, *Flow) {
	lastFlowFile := "/tmp/vflow.mon.lastflow"

	flow := new(Flow)
	lastFlow := new(Flow)

	client := NewHTTP()
	err := client.Get(host+"/flow", flow)
	if err != nil {
		return err, nil, nil
	}

	flow.Timestamp = time.Now().Unix()

	b, err := ioutil.ReadFile(lastFlowFile)
	if err != nil {
		b, _ = json.Marshal(flow)
		ioutil.WriteFile(lastFlowFile, b, 0644)
		return err, nil, nil
	}

	err = json.Unmarshal(b, &lastFlow)
	if err != nil {
		return err, nil, nil
	}

	b, err = json.Marshal(flow)
	if err != nil {
		return err, nil, nil
	}

	ioutil.WriteFile(lastFlowFile, b, 0644)

	// once the vFlow restarted
	if flow.StartTime != lastFlow.StartTime {
		os.Exit(1)
	}

	return nil, flow, lastFlow
}
