package main

import (
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
