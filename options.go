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

	// sFlow options
	SFlowEnabled bool
	SFlowPort    int
	SFlowUDPSize int
	SFlowWorkers int
}

func NewOptions() *Options {
	return &Options{
		Verbose: true,
		Logger:  log.New(os.Stderr, "[vflow] ", log.Ldate|log.Ltime),

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 10,
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

	// sflow options
	flag.BoolVar(&opts.SFlowEnabled, "sflow-enabled", opts.SFlowEnabled, "enable sflow listener")
	flag.IntVar(&opts.SFlowPort, "sflow-port", opts.SFlowPort, "sflow port number")
	flag.IntVar(&opts.SFlowUDPSize, "sflow-max-udp-size", opts.SFlowUDPSize, "sflow maximum UDP size")
	flag.IntVar(&opts.SFlowWorkers, "sflow-workers", opts.SFlowWorkers, "sflow workers / concurrency number")

	flag.Parse()
}

func vFlowLoadCfg(file string, opts *Options) {
	// TODO
}
