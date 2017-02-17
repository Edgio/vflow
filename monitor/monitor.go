//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    monitor.go
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
package main

import (
	"flag"
	"log"
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
}

type SFlow struct {
	UDPQueue     int64
	MessageQueue int64
	UDPCount     int64
	DecodedCount int64
}

type Flow struct {
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

type Options struct {
	DBType       string
	VFlowHost    string
	InfluxDBAPI  string
	InfluxDBName string
}

var opts = Options{
	DBType:       "influxdb",
	VFlowHost:    "http://localhost:8080",
	InfluxDBAPI:  "http://localhost:8086",
	InfluxDBName: "vflow",
}

func init() {

	flag.StringVar(&opts.DBType, "db-type", opts.DBType, "database type name to ingest")
	flag.StringVar(&opts.VFlowHost, "vflow-host", opts.VFlowHost, "vflow host address and port")
	flag.StringVar(&opts.InfluxDBAPI, "influxdb-api-addr", opts.InfluxDBAPI, "influxdb api address")
	flag.StringVar(&opts.InfluxDBName, "influxdb-db-name", opts.InfluxDBName, "influxdb database name")

	flag.Parse()
}

func main() {
	var m = make(map[string]Monitor)

	m["influxdb"] = InfluxDB{
		API:   opts.InfluxDBAPI,
		DB:    opts.InfluxDBName,
		VHost: opts.VFlowHost,
	}

	switch opts.DBType {
	case "influxdb":
		if err := m["influxdb"].Netflow(); err != nil {
			log.Println(err)
		}
		if err := m["influxdb"].System(); err != nil {
			log.Println(err)
		}
	}
}
