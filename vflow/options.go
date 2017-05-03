//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    options.go
//: details: vFlow options :: configuration and command line
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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"gopkg.in/yaml.v2"
)

var (
	version    string
	maxWorkers = runtime.NumCPU() * 1e4
)

// Options represents options
type Options struct {
	// global options
	Verbose    bool   `yaml:"verbose"`
	LogFile    string `yaml:"log-file"`
	PIDFile    string `yaml:"pid-file"`
	DynWorkers bool   `yaml:"dynamic-workers"`
	Logger     *log.Logger
	version    bool

	// stats options
	StatsEnabled  bool   `yaml:"stats-enabled"`
	StatsHTTPAddr string `yaml:"stats-http-addr"`
	StatsHTTPPort string `yaml:"stats-http-port"`

	// sFlow options
	SFlowEnabled bool   `yaml:"sflow-enabled"`
	SFlowPort    int    `yaml:"sflow-port"`
	SFlowUDPSize int    `yaml:"sflow-udp-size"`
	SFlowWorkers int    `yaml:"sflow-workers"`
	SFlowTopic   string `yaml:"sflow-topic"`

	// IPFIX options
	IPFIXEnabled       bool   `yaml:"ipfix-enabled"`
	IPFIXRPCEnabled    bool   `yaml:"ipfix-rpc-enabled"`
	IPFIXPort          int    `yaml:"ipfix-port"`
	IPFIXUDPSize       int    `yaml:"ipfix-udp-size"`
	IPFIXWorkers       int    `yaml:"ipfix-workers"`
	IPFIXTopic         string `yaml:"ipfix-topic"`
	IPFIXMirrorAddr    string `yaml:"ipfix-mirror-addr"`
	IPFIXMirrorPort    int    `yaml:"ipfix-mirror-port"`
	IPFIXMirrorWorkers int    `yaml:"ipfix-mirror-workers"`
	IPFIXTplCacheFile  string `yaml:"ipfix-tpl-cache-file"`

	// Netflow
	NetflowV9Enabled      bool   `yaml:"netflow9-enabled"`
	NetflowV9Port         int    `yaml:"netflow9-port"`
	NetflowV9UDPSize      int    `yaml:"netflow9-udp-size"`
	NetflowV9Workers      int    `yaml:"netflow9-workers"`
	NetflowV9Topic        string `yaml:"netflow9-topic"`
	NetflowV9TplCacheFile string `yaml:"netflow9-tpl-cache-file"`

	// producer
	MQName       string `yaml:"mq-name"`
	MQConfigFile string `yaml:"mq-config-file"`
}

func init() {
	if version == "" {
		version = "unknown"
	}
}

// NewOptions constructs new options
func NewOptions() *Options {
	return &Options{
		Verbose:    false,
		version:    false,
		DynWorkers: true,
		PIDFile:    "/var/run/vflow.pid",
		Logger:     log.New(os.Stderr, "[vflow] ", log.Ldate|log.Ltime),

		StatsEnabled:  true,
		StatsHTTPPort: "8081",
		StatsHTTPAddr: "",

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 200,
		SFlowTopic:   "vflow.sflow",

		IPFIXEnabled:       true,
		IPFIXRPCEnabled:    true,
		IPFIXPort:          4739,
		IPFIXUDPSize:       1500,
		IPFIXWorkers:       200,
		IPFIXTopic:         "vflow.ipfix",
		IPFIXMirrorAddr:    "",
		IPFIXMirrorPort:    4172,
		IPFIXMirrorWorkers: 5,
		IPFIXTplCacheFile:  "/tmp/vflow.templates",

		NetflowV9Enabled:      false,
		NetflowV9Port:         4729,
		NetflowV9UDPSize:      1500,
		NetflowV9Workers:      50,
		NetflowV9Topic:        "vflow.netflow9",
		NetflowV9TplCacheFile: "/tmp/netflowv9.templates",

		MQName:       "kafka",
		MQConfigFile: "/usr/local/vflow/etc/kafka.conf",
	}
}

// GetOptions gets options through cmd and file
func GetOptions() *Options {
	opts := NewOptions()

	opts.vFlowFlagSet()
	opts.vFlowVersion()

	if opts.LogFile != "" {
		f, err := os.OpenFile(opts.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			opts.Logger.Println(err)
		} else {
			opts.Logger.SetOutput(f)
		}
	}

	if ok := opts.vFlowIsRunning(); ok {
		opts.Logger.Fatal("the vFlow already is running!")
	}

	opts.vFlowPIDWrite()

	return opts
}

func (opts Options) vFlowPIDWrite() {
	f, err := os.OpenFile(opts.PIDFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		opts.Logger.Println(err)
		return
	}

	_, err = fmt.Fprintf(f, "%d", os.Getpid())
	if err != nil {
		opts.Logger.Println(err)
	}
}

func (opts Options) vFlowIsRunning() bool {
	b, err := ioutil.ReadFile(opts.PIDFile)
	if err != nil {
		return false
	}

	cmd := exec.Command("kill", "-0", string(b))
	_, err = cmd.Output()
	if err != nil {
		return false
	}

	return true
}

func (opts Options) vFlowVersion() {
	if opts.version {
		fmt.Printf("vFlow version: %s\n", version)
		os.Exit(0)
	}
}

func (opts *Options) vFlowFlagSet() {

	var config string
	flag.StringVar(&config, "config", "/usr/local/vflow/etc/vflow.conf", "path to config file")

	vFlowLoadCfg(opts)

	// global options
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "enable/disable verbose logging")
	flag.BoolVar(&opts.DynWorkers, "dynamic-workers", opts.DynWorkers, "enable/disable dynamic workers")
	flag.BoolVar(&opts.version, "version", opts.version, "show version")
	flag.StringVar(&opts.LogFile, "log-file", opts.LogFile, "log file name")
	flag.StringVar(&opts.PIDFile, "pid-file", opts.PIDFile, "pid file name")

	// stats options
	flag.BoolVar(&opts.StatsEnabled, "stats-enabled", opts.StatsEnabled, "enable/disable stats listener")
	flag.StringVar(&opts.StatsHTTPPort, "stats-http-port", opts.StatsHTTPPort, "stats port listener")
	flag.StringVar(&opts.StatsHTTPAddr, "stats-http-addr", opts.StatsHTTPAddr, "stats bind address listener")

	// sflow options
	flag.BoolVar(&opts.SFlowEnabled, "sflow-enabled", opts.SFlowEnabled, "enable/disable sflow listener")
	flag.IntVar(&opts.SFlowPort, "sflow-port", opts.SFlowPort, "sflow port number")
	flag.IntVar(&opts.SFlowUDPSize, "sflow-max-udp-size", opts.SFlowUDPSize, "sflow maximum UDP size")
	flag.IntVar(&opts.SFlowWorkers, "sflow-workers", opts.SFlowWorkers, "sflow workers number")
	flag.StringVar(&opts.SFlowTopic, "sflow-topic", opts.SFlowTopic, "sflow topic name")

	// ipfix options
	flag.BoolVar(&opts.IPFIXEnabled, "ipfix-enabled", opts.IPFIXEnabled, "enable/disable IPFIX listener")
	flag.BoolVar(&opts.IPFIXRPCEnabled, "ipfix-rpc-enabled", opts.IPFIXRPCEnabled, "enable/disable RPC IPFIX")
	flag.IntVar(&opts.IPFIXPort, "ipfix-port", opts.IPFIXPort, "IPFIX port number")
	flag.IntVar(&opts.IPFIXUDPSize, "ipfix-max-udp-size", opts.IPFIXUDPSize, "IPFIX maximum UDP size")
	flag.IntVar(&opts.IPFIXWorkers, "ipfix-workers", opts.IPFIXWorkers, "IPFIX workers number")
	flag.StringVar(&opts.IPFIXTopic, "ipfix-topic", opts.IPFIXTopic, "ipfix topic name")
	flag.StringVar(&opts.IPFIXTplCacheFile, "ipfix-tpl-cache-file", opts.IPFIXTplCacheFile, "IPFIX template cache file")
	flag.StringVar(&opts.IPFIXMirrorAddr, "ipfix-mirror-addr", opts.IPFIXMirrorAddr, "IPFIX mirror destination address")
	flag.IntVar(&opts.IPFIXMirrorPort, "ipfix-mirror-port", opts.IPFIXMirrorPort, "IPFIX mirror destination port number")
	flag.IntVar(&opts.IPFIXMirrorWorkers, "ipfix-mirror-workers", opts.IPFIXMirrorWorkers, "IPFIX mirror workers number")

	// netflow version 9
	flag.BoolVar(&opts.NetflowV9Enabled, "netflow9-enabled", opts.NetflowV9Enabled, "enable/disable netflow version 9 listener")
	flag.IntVar(&opts.NetflowV9Port, "netflow9-port", opts.NetflowV9Port, "Netflow Version 9 port number")
	flag.IntVar(&opts.NetflowV9UDPSize, "netflow9-max-udp-size", opts.NetflowV9UDPSize, "Netflow version 9 maximum UDP size")
	flag.IntVar(&opts.NetflowV9Workers, "netflow9-workers", opts.NetflowV9Workers, "Netflow version 9 workers number")
	flag.StringVar(&opts.NetflowV9Topic, "netflow9-topic", opts.NetflowV9Topic, "Netflow version 9 topic name")
	flag.StringVar(&opts.NetflowV9TplCacheFile, "netflow9-tpl-cache-file", opts.NetflowV9TplCacheFile, "Netflow version 9 template cache file")

	// producer options
	flag.StringVar(&opts.MQName, "mqueue", opts.MQName, "producer message queue name")
	flag.StringVar(&opts.MQConfigFile, "mqueue-conf", opts.MQConfigFile, "producer message queue configuration file")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
    Example:
	# set workers
	vflow -sflow-workers 15 -ipfix-workers 20

	# set 3rd party ipfix collector
	vflow -ipfix-mirror-addr 192.168.1.10 -ipfix-mirror-port 4319

	# enaable verbose logging
	vflow -verbose=true

	# for more information
	https://github.com/VerizonDigital/vflow/blob/master/docs/config.md

    `)

	}

	flag.Parse()
}

func vFlowLoadCfg(opts *Options) {
	var file = "/usr/local/vflow/etc/vflow.conf"

	for i, flag := range os.Args {
		if flag == "-config" {
			file = os.Args[i+1]
			break
		}
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		opts.Logger.Println(err)
		return
	}
	err = yaml.Unmarshal(b, opts)
	if err != nil {
		opts.Logger.Println(err)
	}
}
