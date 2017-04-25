//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    tsdb.go
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

// Package store ingest monitoring time series data points
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// TSDB represents TSDB ingestion
type TSDB struct {
	API   string
	VHost string
}

// TSDBDataPoint represents single TSDB data point
type TSDBDataPoint struct {
	Metric    string `json:"metric"`
	Timestamp int64  `json:"timestamp"`
	Value     int64  `json:"value"`
	Tags      struct {
		Host string `json:"host"`
		Type string `json:"type"`
	}
}

// TSDBResp represents TSDP response
type TSDBResp struct {
	Failed  int `json:"failed"`
	Success int `json:"success"`
}

// Netflow ingests flow's stats to TSDB
func (t TSDB) Netflow() error {
	var (
		dps    []TSDBDataPoint
		values []int64
	)

	flow, lastFlow, err := getFlow(t.VHost)
	if err != nil {
		return err
	}

	delta := flow.Timestamp - lastFlow.Timestamp
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	metrics := [][]string{
		{"ipfix", "udp.rate"},
		{"sflow", "udp.rate"},
		{"ipfix", "decode.rate"},
		{"sflow", "decode.rate"},
		{"ipfix", "mq.error.rate"},
		{"sflow", "mq.error.rate"},
		{"ipfix", "workers"},
		{"sflow", "workers"},
	}

	values = append(values, abs((flow.IPFIX.UDPCount-lastFlow.IPFIX.UDPCount)/delta))
	values = append(values, abs((flow.SFlow.UDPCount-lastFlow.SFlow.UDPCount)/delta))
	values = append(values, abs((flow.IPFIX.DecodedCount-lastFlow.IPFIX.DecodedCount)/delta))
	values = append(values, abs((flow.SFlow.DecodedCount-lastFlow.SFlow.DecodedCount)/delta))
	values = append(values, abs((flow.IPFIX.MQErrorCount-lastFlow.IPFIX.MQErrorCount)/delta))
	values = append(values, abs((flow.SFlow.MQErrorCount-lastFlow.SFlow.MQErrorCount)/delta))
	values = append(values, flow.IPFIX.Workers)
	values = append(values, flow.SFlow.Workers)

	for i, m := range metrics {
		dps = append(dps, TSDBDataPoint{
			Metric:    m[1],
			Timestamp: time.Now().Unix(),
			Value:     values[i],
			Tags: struct {
				Host string `json:"host"`
				Type string `json:"type"`
			}{
				Host: hostname,
				Type: m[0],
			},
		})

	}

	if err = t.put(dps); err != nil {
		return err
	}

	return nil
}

// System ingests system's stats to TSDB
func (t TSDB) System() error {
	var dps []TSDBDataPoint

	sys := new(Sys)
	client := NewHTTP()
	err := client.Get(t.VHost+"/sys", sys)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	metrics := []string{
		"mem.heap.alloc",
		"mem.alloc",
		"mcache.inuse",
		"mem.total.alloc",
		"mem.heap.sys",
		"num.goroutine",
	}

	values := []int64{
		sys.MemHeapAlloc,
		sys.MemAlloc,
		sys.MCacheInuse,
		sys.MemTotalAlloc,
		sys.MemHeapSys,
		sys.NumGoroutine,
	}

	for i, m := range metrics {
		dps = append(dps, TSDBDataPoint{
			Metric:    m,
			Timestamp: time.Now().Unix(),
			Value:     values[i],
			Tags: struct {
				Host string `json:"host"`
				Type string `json:"type"`
			}{
				Host: hostname,
			},
		})

	}

	if err = t.put(dps); err != nil {
		return err
	}

	return nil
}

func (t TSDB) put(dps []TSDBDataPoint) error {
	b, err := json.Marshal(dps)
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/api/put", t.API)
	client := NewHTTP()
	b, err = client.Post(api, "application/json", string(b))
	if err != nil {
		return err
	}

	resp := TSDBResp{}
	json.Unmarshal(b, resp)

	if resp.Failed > 0 {
		return errors.New("TSDB error")
	}

	return nil
}
