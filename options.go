package main

type Options struct {
	// global options
	Verbose bool

	// sFlow options
	SFlowEnabled bool
	SFlowPort    int
	SFlowUDPSize int
	SFlowWorkers int
}

func NewOptions() *Options {
	return &Options{
		Verbose: true,

		SFlowEnabled: true,
		SFlowPort:    6343,
		SFlowUDPSize: 1500,
		SFlowWorkers: 10,
	}
}
