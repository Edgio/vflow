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

	"github.com/VerizonDigital/vflow/monitor/store"
)

type options struct {
	DBType       string
	VFlowHost    string
	InfluxDBAPI  string
	InfluxDBName string
	TSDBAPI      string
}

var opts = options{
	DBType:       "influxdb",
	VFlowHost:    "http://localhost:8080",
	InfluxDBAPI:  "http://localhost:8086",
	TSDBAPI:      "http://localhost:4242",
	InfluxDBName: "vflow",
}

func init() {

	flag.StringVar(&opts.DBType, "db-type", opts.DBType, "database type name to ingest")
	flag.StringVar(&opts.VFlowHost, "vflow-host", opts.VFlowHost, "vflow host address and port")
	flag.StringVar(&opts.InfluxDBAPI, "influxdb-api-addr", opts.InfluxDBAPI, "influxdb api address")
	flag.StringVar(&opts.InfluxDBName, "influxdb-db-name", opts.InfluxDBName, "influxdb database name")
	flag.StringVar(&opts.TSDBAPI, "tsdb-api-addr", opts.TSDBAPI, "tsdb api address")

	flag.Parse()
}

func main() {
	var m = make(map[string]store.Monitor)

	m["influxdb"] = store.InfluxDB{
		API:   opts.InfluxDBAPI,
		DB:    opts.InfluxDBName,
		VHost: opts.VFlowHost,
	}

	m["tsdb"] = store.TSDB{
		API:   opts.TSDBAPI,
		VHost: opts.VFlowHost,
	}

	if _, ok := m[opts.DBType]; !ok {
		log.Fatalf("the storage: %s is not available", opts.DBType)
	}

	if err := m[opts.DBType].Netflow(); err != nil {
		log.Println(err)
	}
	if err := m[opts.DBType].System(); err != nil {
		log.Println(err)
	}
}
