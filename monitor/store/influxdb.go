// Package store ingest monitoring time series data points
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    influxdb.go
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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type InfluxDB struct {
	API   string
	DB    string
	VHost string
}

const (
	iFlowFile = "/tmp/vflow.mon.flow"
)

func (i InfluxDB) Netflow() error {
	var (
		flow     = new(Flow)
		lastFlow = new(Flow)
		client   = new(http.Client)
		err      error
		b        []byte
	)

	resp, err := client.Get(i.VHost + "/flow")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &flow)
	if err != nil {
		return err
	}

	flow.Timestamp = time.Now().Unix()

	b, err = ioutil.ReadFile(iFlowFile)
	if err != nil {
		b, _ = json.Marshal(flow)
		ioutil.WriteFile(iFlowFile, b, 0644)
		return err
	}

	err = json.Unmarshal(b, &lastFlow)
	if err != nil {
		return err
	}

	b, err = json.Marshal(flow)
	if err != nil {
		return err
	}

	ioutil.WriteFile(iFlowFile, b, 0644)

	delta := flow.Timestamp - lastFlow.Timestamp
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("udp.rate,type=ipfix,host=%s value=%d\n", hostname, (flow.IPFIX.UDPCount-lastFlow.IPFIX.UDPCount)/delta)
	query += fmt.Sprintf("udp.rate,type=sflow,host=%s value=%d\n", hostname, (flow.SFlow.UDPCount-lastFlow.SFlow.UDPCount)/delta)
	query += fmt.Sprintf("decode.rate,type=ipfix,host=%s value=%d\n", hostname, (flow.IPFIX.DecodedCount-lastFlow.IPFIX.DecodedCount)/delta)
	query += fmt.Sprintf("decode.rate,type=sflow,host=%s value=%d\n", hostname, (flow.SFlow.DecodedCount-lastFlow.SFlow.DecodedCount)/delta)

	api := fmt.Sprintf("%s/write?db=%s", i.API, i.DB)
	resp, err = client.Post(api, "text/plain", bytes.NewBufferString(query))
	if err != nil {
		return err
	}

	if err = chkInfluxDBResp(resp.Body); err != nil {
		return err
	}

	return nil
}
func (i InfluxDB) System() error {
	var (
		sys    = new(Sys)
		client = new(http.Client)
		err    error
	)

	resp, err := client.Get(i.VHost + "/sys")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &sys)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("mem.heap.alloc,host=%s value=%d\n", hostname, sys.MemHeapAlloc)
	query += fmt.Sprintf("mem.alloc,host=%s value=%d\n", hostname, sys.MemAlloc)
	query += fmt.Sprintf("mcache.inuse,host=%s value=%d\n", hostname, sys.MCacheInuse)
	query += fmt.Sprintf("mem.total.alloc,host=%s value=%d\n", hostname, sys.MemTotalAlloc)
	query += fmt.Sprintf("mem.heap.sys,host=%s value=%d\n", hostname, sys.MemHeapSys)
	query += fmt.Sprintf("num.goroutine,host=%s value=%d\n", hostname, sys.NumGoroutine)

	api := fmt.Sprintf("%s/write?db=%s", i.API, i.DB)
	resp, err = client.Post(api, "text/plain", bytes.NewBufferString(query))
	if err != nil {
		return err
	}

	if err = chkInfluxDBResp(resp.Body); err != nil {
		return err
	}

	return nil
}

func chkInfluxDBResp(r io.Reader) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		return errors.New("influxdb error: " + string(body))
	}

	return nil
}
