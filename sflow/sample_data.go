package sflow

import (
	"fmt"
	"io"
)

type FlowSample struct {
	SequenceNo   uint32 // Incremented with each flow sample
	SourceId     byte
	SamplingRate uint32
	SamplePool   uint32
	Drops        uint32
	Input        uint32
	Output       uint32
}

func decodeFlowSample(r io.ReadSeeker) error {
	var (
		fs  = new(FlowSample)
		err error
	)

	if err = read(r, &fs.SequenceNo); err != nil {
		return err
	}

	if err = read(r, &fs.SourceId); err != nil {
		return err
	}

	r.Seek(3, 1) // skip counter sample decoding

	if err = read(r, &fs.SamplingRate); err != nil {
		return err
	}

	if err = read(r, &fs.SamplePool); err != nil {
		return err
	}

	if err = read(r, &fs.Drops); err != nil {
		return err
	}

	if err = read(r, &fs.Input); err != nil {
		return err
	}

	if err = read(r, &fs.Output); err != nil {
		return err
	}

	fmt.Printf("%#v\n", fs)

	return nil
}
