package sflow

import (
	"fmt"
	"io"
)

const (
	SFDataRawHeader      = 1
	SFDataEthernetHeader = 2
	SFDataIPV4Header     = 3
	SFDataIPV6Header     = 4

	SFDataExtSwitch     = 1001
	SFDataExtRouter     = 1002
	SFDataExtGateway    = 1003
	SFDataExtUser       = 1004
	SFDataExtURL        = 1005
	SFDataExtMPLS       = 1006
	SFDataExtNAT        = 1007
	SFDataExtMPLSTunnel = 1008
	SFDataExtMPLSVC     = 1009
	SFDataExtMPLSFEC    = 1010
	SFDataExtMPLSLVPFEC = 1011
	SFDataExtVLANTunnel = 1012
)

type FlowSample struct {
	SequenceNo   uint32 // Incremented with each flow sample
	SourceId     byte
	SamplingRate uint32
	SamplePool   uint32
	Drops        uint32
	Input        uint32
	Output       uint32
	RecordsNo    uint32
}

func decodeFlowSample(r io.ReadSeeker) error {
	var (
		fs          = new(FlowSample)
		rTypeFormat uint32
		rTypeLength uint32
		err         error
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

	if err = read(r, &fs.RecordsNo); err != nil {
		return err
	}

	fmt.Printf("%#v\n", fs) // just for test

	for i := uint32(0); i < fs.RecordsNo; i++ {
		if err = read(r, &rTypeFormat); err != nil {
			return err
		}
		if err = read(r, &rTypeLength); err != nil {
			return err
		}

		switch rTypeFormat {
		case SFDataRawHeader:
			decodeSampledHeader(r)
			r.Seek(int64(rTypeLength), 1)
		case SFDataExtSwitch:
			decodeSampledExtSwitch(r)
			r.Seek(int64(rTypeLength), 1)
		default:
			r.Seek(int64(rTypeLength), 1)
		}
	}

	return nil
}

func decodeSampledHeader(r io.Reader) {
	// TODO
}

func decodeSampledExtSwitch(r io.Reader) {
	// TODO
}

func decodeSampledIPv4() {
	// TODO
}

func decodeSampledIPv6() {
	// TODO
}
