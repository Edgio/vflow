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

	config := flag.Bool("verbose", opts.Verbose, "enable verbose logging")

	if config != "" {
		vFlowLoadCfg(config, opts)
	}

	// global
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "enable verbose logging")

	flag.Parse()

	log.Printf("%#v\n", opts)
}

func vFlowLoadCfg(file string, opts *Options) {
	// TODO
}
