//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    options.go
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
	"os"
)

type Options struct {
	// global options
	Verbose bool
	Logger  *log.Logger

	// stats options
	StatsEnabled  bool
	StatsHTTPPort string
	StatsHTTPAddr string

	// sFlow options
	SFlowEnabled bool
	SFlowPort    int
	SFlowUDPSize int
	SFlowWorkers int

	// IPFIX options
	IPFIXEnabled           bool
	IPFIXPort              int
	IPFIXUDPSize           int
	IPFIXWorkers           int
	IPFIXMirror            string
	IPFIXTemplateCacheFile string

	// producer
	MQName       string
	MQConfigFile string
}

func NewOptions() *Options {
	return &Options{
		Verbose: true,
		Logger:  log.New(os.Stderr, "[vflow] ", log.Ldate|log.Ltime),

		StatsEnabled:  true,
		StatsHTTPPort: "8080",
		StatsHTTPAddr: "",

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 10,

		IPFIXEnabled:           true,
		IPFIXPort:              4739,
		IPFIXUDPSize:           1500,
		IPFIXWorkers:           10,
		IPFIXMirror:            "139.49.193.73:4172",
		IPFIXTemplateCacheFile: "/tmp/vflow.templates",

		MQName:       "kafka",
		MQConfigFile: "/usr/local/vflow/etc/kafka.conf",
	}
}

func GetOptions() *Options {
	opts := NewOptions()
	vFlowFlagSet(opts)

	return opts
}

func vFlowFlagSet(opts *Options) {

	var config string

	flag.StringVar(&config, "config", "", "path to config file")

	if config != "" {
		vFlowLoadCfg(config, opts)
	}

	// global options
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "enable verbose logging")

	// stats options
	flag.BoolVar(&opts.StatsEnabled, "stats-enabled", opts.StatsEnabled, "enable stats listener")
	flag.StringVar(&opts.StatsHTTPPort, "stats-http-port", opts.StatsHTTPPort, "stats port listener")
	flag.StringVar(&opts.StatsHTTPAddr, "stats-http-addr", opts.StatsHTTPAddr, "stats bind address listener")

	// sflow options
	flag.BoolVar(&opts.SFlowEnabled, "sflow-enabled", opts.SFlowEnabled, "enable sflow listener")
	flag.IntVar(&opts.SFlowPort, "sflow-port", opts.SFlowPort, "sflow port number")
	flag.IntVar(&opts.SFlowUDPSize, "sflow-max-udp-size", opts.SFlowUDPSize, "sflow maximum UDP size")
	flag.IntVar(&opts.SFlowWorkers, "sflow-workers", opts.SFlowWorkers, "sflow workers / concurrency number")

	// ipfix options
	flag.BoolVar(&opts.IPFIXEnabled, "ipfix-enabled", opts.IPFIXEnabled, "enable IPFIX listener")
	flag.IntVar(&opts.IPFIXPort, "ipfix-port", opts.IPFIXPort, "IPFIX port number")
	flag.IntVar(&opts.IPFIXUDPSize, "ipfix-max-udp-size", opts.IPFIXUDPSize, "IPFIX maximum UDP size")
	flag.IntVar(&opts.IPFIXWorkers, "ipfix-workers", opts.IPFIXWorkers, "IPFIX workers / concurrency number")
	flag.StringVar(&opts.IPFIXTemplateCacheFile, "ipfix-tpl-cache-file", opts.IPFIXTemplateCacheFile, "IPFIX template cache file")

	// producer
	flag.StringVar(&opts.MQName, "mqueue", opts.MQName, "producer message queue name")
	flag.StringVar(&opts.MQConfigFile, "mqueue-conf", opts.MQConfigFile, "producer message queue configuration file")

	flag.Parse()
}

func vFlowLoadCfg(file string, opts *Options) {
	// TODO
}
