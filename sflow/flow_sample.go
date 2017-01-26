package sflow

import (
	"errors"
	"io"

	"git.edgecastcdn.net/vflow/packet"
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
	SourceId     byte   // fsSourceId
	SamplingRate uint32 // fsPacketSamplingRate
	SamplePool   uint32 // Total number of packets that could have been sampled
	Drops        uint32 // Number of times a packet was dropped due to lack of resources
	Input        uint32 // SNMP ifIndex of input interface
	Output       uint32 // SNMP ifIndex of input interface
	RecordsNo    uint32 // Number of records to follow
}

type SampledHeader struct {
	Protocol     uint32 // (enum SFLHeader_protocol)
	FrameLength  uint32 // Original length of packet before sampling
	Stripped     uint32 // Header/trailer bytes stripped by sender
	HeaderLength uint32 // Length of sampled header bytes to follow
	Header       []byte // Header bytes
}

type ExtSwitchData struct {
	SrcVlan     uint32
	SrcPriority uint32
	DstVlan     uint32
	DstPriority uint32
}

var (
	maxOutEthernetLength = errors.New("the ethernet lenght is greater than 1500")
)

func decodeFlowSample(r io.ReadSeeker) ([]interface{}, error) {
	var (
		fs          = new(FlowSample)
		data        []interface{}
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = read(r, &fs.SequenceNo); err != nil {
		return nil, err
	}

	if err = read(r, &fs.SourceId); err != nil {
		return nil, err
	}

	r.Seek(3, 1) // skip counter sample decoding

	if err = read(r, &fs.SamplingRate); err != nil {
		return nil, err
	}

	if err = read(r, &fs.SamplePool); err != nil {
		return nil, err
	}

	if err = read(r, &fs.Drops); err != nil {
		return nil, err
	}

	if err = read(r, &fs.Input); err != nil {
		return nil, err
	}

	if err = read(r, &fs.Output); err != nil {
		return nil, err
	}

	if err = read(r, &fs.RecordsNo); err != nil {
		return nil, err
	}

	data = append(data, fs)

	for i := uint32(0); i < fs.RecordsNo; i++ {
		if err = read(r, &rTypeFormat); err != nil {
			return nil, err
		}
		if err = read(r, &rTypeLength); err != nil {
			return nil, err
		}

		switch rTypeFormat {
		case SFDataRawHeader:
			h, err := decodeSampledHeader(r)
			if err != nil {
				return data, err
			}
			data = append(data, h)
		case SFDataExtSwitch:
			d, err := decodeExtSwitchData(r)
			if err != nil {
				return data, err
			}
			data = append(data, d)
		default:
			r.Seek(int64(rTypeLength), 1)
		}
	}

	return data, nil
}

func decodeSampledHeader(r io.Reader) (*packet.Packet, error) {
	// TODO
	var (
		h   = new(SampledHeader)
		err error
	)
	if err = read(r, &h.Protocol); err != nil {
		return nil, err
	}

	if err = read(r, &h.FrameLength); err != nil {
		return nil, err
	}

	if err = read(r, &h.Stripped); err != nil {
		return nil, err
	}

	if err = read(r, &h.HeaderLength); err != nil {
		return nil, err
	}

	if h.HeaderLength > 1500 {
		return nil, maxOutEthernetLength
	}

	// TODO: make sure the padding works!!
	// cut off a header length mod 4 == 0 number of bytes
	tmp := (4 - h.HeaderLength) % 4
	if tmp < 0 {
		tmp += 4
	}

	h.Header = make([]byte, h.HeaderLength+tmp)
	if _, err = r.Read(h.Header); err != nil {
		return nil, err
	}

	h.Header = h.Header[:h.HeaderLength]

	p := packet.NewPacket()
	d, err := p.Decoder(h.Header)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func decodeExtSwitchData(r io.Reader) (*ExtSwitchData, error) {
	var (
		d   = new(ExtSwitchData)
		err error
	)

	if err = read(r, &d.SrcVlan); err != nil {
		return nil, err
	}

	if err = read(r, &d.SrcPriority); err != nil {
		return nil, err
	}

	if err = read(r, &d.DstVlan); err != nil {
		return nil, err
	}

	if err = read(r, &d.SrcPriority); err != nil {
		return nil, err
	}

	return d, nil
}
