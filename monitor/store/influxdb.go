//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    influxdb.go
//: details: influx ingest handler
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
	"errors"
	"fmt"
	"os"
)

// InfluxDB represents InfluxDB backend
type InfluxDB struct {
	API   string
	DB    string
	VHost string
}

// Netflow ingests flow's stats to InfluxDB
func (i InfluxDB) Netflow() error {
	flow, lastFlow, err := getFlow(i.VHost)
	if err != nil {
		return err
	}

	delta := flow.Timestamp - lastFlow.Timestamp
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	value := abs((flow.IPFIX.UDPCount - lastFlow.IPFIX.UDPCount) / delta)
	query := fmt.Sprintf("udp.rate,type=ipfix,host=%s value=%d\n", hostname, value)
	value = abs((flow.SFlow.UDPCount - lastFlow.SFlow.UDPCount) / delta)
	query += fmt.Sprintf("udp.rate,type=sflow,host=%s value=%d\n", hostname, value)
	value = abs((flow.IPFIX.DecodedCount - lastFlow.IPFIX.DecodedCount) / delta)
	query += fmt.Sprintf("decode.rate,type=ipfix,host=%s value=%d\n", hostname, value)
	value = abs((flow.SFlow.DecodedCount - lastFlow.SFlow.DecodedCount) / delta)
	query += fmt.Sprintf("decode.rate,type=sflow,host=%s value=%d\n", hostname, value)
	value = abs((flow.IPFIX.MQErrorCount - lastFlow.IPFIX.MQErrorCount) / delta)
	query += fmt.Sprintf("mq.error.rate,type=ipfix,host=%s value=%d\n", hostname, value)
	value = abs((flow.SFlow.MQErrorCount - lastFlow.SFlow.MQErrorCount) / delta)
	query += fmt.Sprintf("mq.error.rate,type=sflow,host=%s value=%d\n", hostname, value)

	query += fmt.Sprintf("workers,type=ipfix,host=%s value=%d\n", hostname, flow.IPFIX.Workers)
	query += fmt.Sprintf("workers,type=sflow,host=%s value=%d\n", hostname, flow.SFlow.Workers)
	query += fmt.Sprintf("udp.queue,type=ipfix,host=%s value=%d\n", hostname, flow.IPFIX.UDPQueue)
	query += fmt.Sprintf("udp.queue,type=sflow,host=%s value=%d\n", hostname, flow.SFlow.UDPQueue)
	query += fmt.Sprintf("mq.queue,type=ipfix,host=%s value=%d\n", hostname, flow.IPFIX.MessageQueue)
	query += fmt.Sprintf("mq.queue,type=sflow,host=%s value=%d\n", hostname, flow.SFlow.MessageQueue)
	query += fmt.Sprintf("udp.mirror.queue,type=ipfix,host=%s value=%d\n", hostname, flow.IPFIX.UDPMirrorQueue)

	api := fmt.Sprintf("%s/write?db=%s", i.API, i.DB)
	client := NewHTTP()
	b, err := client.Post(api, "text/plain", query)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		return errors.New("influxdb error: " + string(b))
	}

	return nil
}

// System ingests system's stats to InfluxDB
func (i InfluxDB) System() error {
	sys := new(Sys)
	client := NewHTTP()
	err := client.Get(i.VHost+"/sys", sys)
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
	b, err := client.Post(api, "text/plain", query)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		return errors.New("influxdb error: " + string(b))
	}

	return nil
}

func abs(a int64) int64 {
	if a < 0 {
		os.Exit(1)
	}

	return a
}
