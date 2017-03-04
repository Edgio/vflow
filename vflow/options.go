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
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Options represents options
type Options struct {
	// global options
	Verbose bool `yaml:"verbose"`
	Logger  *log.Logger

	// stats options
	StatsEnabled  bool   `yaml:"stats-enabled"`
	StatsHTTPAddr string `yaml:"stats-http-addr"`
	StatsHTTPPort string `yaml:"stats-http-port"`

	// sFlow options
	SFlowEnabled bool `yaml:"sflow-enabled"`
	SFlowPort    int  `yaml:"sflow-port"`
	SFlowUDPSize int  `yaml:"sflow-udp-size"`
	SFlowWorkers int  `yaml:"sflow-workers"`

	// IPFIX options
	IPFIXEnabled       bool   `yaml:"ipfix-enabled"`
	IPFIXRPCEnabled    bool   `yaml:"ipfix-rpc-enabled"`
	IPFIXPort          int    `yaml:"ipfix-port"`
	IPFIXUDPSize       int    `yaml:"ipfix-udp-size"`
	IPFIXWorkers       int    `yaml:"ipfix-workers"`
	IPFIXMirrorAddr    string `yaml:"ipfix-mirror-addr"`
	IPFIXMirrorPort    int    `yaml:"ipfix-mirror-port"`
	IPFIXMirrorWorkers int    `yaml:"ipfix-mirror-workers"`
	IPFIXTplCacheFile  string `yaml:"ipfix-tpl-cache-file"`

	// producer
	MQName       string `yaml:"mq-name"`
	MQConfigFile string `yaml:"mq-config-file"`
}

// NewOptions constructs new options
func NewOptions() *Options {
	return &Options{
		Verbose: false,
		Logger:  log.New(os.Stderr, "[vflow] ", log.Ldate|log.Ltime),

		StatsEnabled:  true,
		StatsHTTPPort: "8080",
		StatsHTTPAddr: "",

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 10,

		IPFIXEnabled:       true,
		IPFIXRPCEnabled:    true,
		IPFIXPort:          4739,
		IPFIXUDPSize:       1500,
		IPFIXWorkers:       10,
		IPFIXMirrorAddr:    "",
		IPFIXMirrorPort:    4172,
		IPFIXMirrorWorkers: 5,
		IPFIXTplCacheFile:  "/tmp/vflow.templates",

		MQName:       "kafka",
		MQConfigFile: "/usr/local/vflow/etc/kafka.conf",
	}
}

// GetOptions gets options through cmd and file
func GetOptions() *Options {
	opts := NewOptions()
	vFlowFlagSet(opts)

	return opts
}

func vFlowFlagSet(opts *Options) {

	var config string

	flag.StringVar(&config, "config", "/usr/local/vflow/etc/vflow.conf", "path to config file")

	vFlowLoadCfg(config, opts)

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
	flag.IntVar(&opts.SFlowWorkers, "sflow-workers", opts.SFlowWorkers, "sflow workers number")

	// ipfix options
	flag.BoolVar(&opts.IPFIXEnabled, "ipfix-enabled", opts.IPFIXEnabled, "enable IPFIX listener")
	flag.BoolVar(&opts.IPFIXRPCEnabled, "ipfix-rpc-enabled", opts.IPFIXRPCEnabled, "enable RPC IPFIX")
	flag.IntVar(&opts.IPFIXPort, "ipfix-port", opts.IPFIXPort, "IPFIX port number")
	flag.IntVar(&opts.IPFIXUDPSize, "ipfix-max-udp-size", opts.IPFIXUDPSize, "IPFIX maximum UDP size")
	flag.IntVar(&opts.IPFIXWorkers, "ipfix-workers", opts.IPFIXWorkers, "IPFIX workers number")
	flag.StringVar(&opts.IPFIXTplCacheFile, "ipfix-tpl-cache-file", opts.IPFIXTplCacheFile, "IPFIX template cache file")
	flag.StringVar(&opts.IPFIXMirrorAddr, "ipfix-mirror-addr", opts.IPFIXMirrorAddr, "IPFIX mirror destination address")
	flag.IntVar(&opts.IPFIXMirrorPort, "ipfix-mirror-port", opts.IPFIXMirrorPort, "IPFIX mirror destination port number")
	flag.IntVar(&opts.IPFIXMirrorWorkers, "ipfix-mirror-workers", opts.IPFIXMirrorWorkers, "IPFIX mirror workers number")

	// producer options
	flag.StringVar(&opts.MQName, "mqueue", opts.MQName, "producer message queue name")
	flag.StringVar(&opts.MQConfigFile, "mqueue-conf", opts.MQConfigFile, "producer message queue configuration file")

	flag.Parse()
}

func vFlowLoadCfg(f string, opts *Options) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		opts.Logger.Println(err)
		return
	}
	err = yaml.Unmarshal(b, opts)
	if err != nil {
		opts.Logger.Println(err)
	}
}
