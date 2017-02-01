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
	StatsHTTPPort int

	// sFlow options
	SFlowEnabled bool
	SFlowPort    int
	SFlowUDPSize int
	SFlowWorkers int

	// IPFIX options
	IPFIXEnabled bool
	IPFIXPort    int
	IPFIXUDPSize int
	IPFIXWorkers int
}

func NewOptions() *Options {
	return &Options{
		Verbose: true,
		Logger:  log.New(os.Stderr, "[vflow] ", log.Ldate|log.Ltime),

		StatsEnabled:  true,
		StatsHTTPPort: 8080,

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 10,

		IPFIXEnabled: true,
		IPFIXPort:    4739,
		IPFIXUDPSize: 1500,
		IPFIXWorkers: 10,
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
	flag.IntVar(&opts.StatsHTTPPort, "stats-http-port", opts.StatsHTTPPort, "stats port listener")

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

	flag.Parse()
}

func vFlowLoadCfg(file string, opts *Options) {
	// TODO
}
