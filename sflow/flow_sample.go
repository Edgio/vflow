// Package sflow decodes sFlow packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    flow_sample.go
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
package sflow

import (
	"errors"
	"io"

	"github.com/VerizonDigital/vflow/packet"
)

const (
	// SFDataRawHeader is sFlow Raw Packet Header number
	SFDataRawHeader = 1

	// SFDataExtSwitch is sFlow Extended Switch Data number
	SFDataExtSwitch = 1001
)

// FlowSample represents single flow sample
type FlowSample struct {
	SequenceNo   uint32 // Incremented with each flow sample
	SourceID     byte   // sfSourceID
	SamplingRate uint32 // sfPacketSamplingRate
	SamplePool   uint32 // Total number of packets that could have been sampled
	Drops        uint32 // Number of times a packet was dropped due to lack of resources
	Input        uint32 // SNMP ifIndex of input interface
	Output       uint32 // SNMP ifIndex of input interface
	RecordsNo    uint32 // Number of records to follow
}

// SampledHeader represents sampled header
type SampledHeader struct {
	Protocol     uint32 // (enum SFLHeader_protocol)
	FrameLength  uint32 // Original length of packet before sampling
	Stripped     uint32 // Header/trailer bytes stripped by sender
	HeaderLength uint32 // Length of sampled header bytes to follow
	Header       []byte // Header bytes
}

// ExtSwitchData represents Extended Switch Data
type ExtSwitchData struct {
	SrcVlan     uint32 // The 802.1Q VLAN id of incoming frame
	SrcPriority uint32 // The 802.1p priority of incoming frame
	DstVlan     uint32 // The 802.1Q VLAN id of outgoing frame
	DstPriority uint32 // The 802.1p priority of outgoing frame
}

var (
	errMaxOutEthernetLength = errors.New("the ethernet lenght is greater than 1500")
)

func (fs *FlowSample) unmarshal(r io.ReadSeeker) error {
	var err error

	if err = read(r, &fs.SequenceNo); err != nil {
		return err
	}

	if err = read(r, &fs.SourceID); err != nil {
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

	return nil
}

func (sh *SampledHeader) unmarshal(r io.Reader) error {
	var err error

	if err = read(r, &sh.Protocol); err != nil {
		return err
	}

	if err = read(r, &sh.FrameLength); err != nil {
		return err
	}

	if err = read(r, &sh.Stripped); err != nil {
		return err
	}

	if err = read(r, &sh.HeaderLength); err != nil {
		return err
	}

	if sh.HeaderLength > 1500 {
		return errMaxOutEthernetLength
	}

	// cut off a header length mod 4 == 0 number of bytes
	tmp := (4 - sh.HeaderLength) % 4
	if tmp < 0 {
		tmp += 4
	}

	sh.Header = make([]byte, sh.HeaderLength+tmp)
	if _, err = r.Read(sh.Header); err != nil {
		return err
	}

	sh.Header = sh.Header[:sh.HeaderLength]

	return nil
}

func (es *ExtSwitchData) unmarshal(r io.Reader) error {
	var err error

	if err = read(r, &es.SrcVlan); err != nil {
		return err
	}

	if err = read(r, &es.SrcPriority); err != nil {
		return err
	}

	if err = read(r, &es.DstVlan); err != nil {
		return err
	}

	if err = read(r, &es.SrcPriority); err != nil {
		return err
	}

	return nil
}

func decodeFlowSample(r io.ReadSeeker) ([]interface{}, error) {
	var (
		fs          = new(FlowSample)
		data        []interface{}
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = fs.unmarshal(r); err != nil {
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
	var (
		h   = new(SampledHeader)
		err error
	)

	if err = h.unmarshal(r); err != nil {
		return nil, err
	}

	p := packet.NewPacket()
	d, err := p.Decoder(h.Header)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func decodeExtSwitchData(r io.Reader) (*ExtSwitchData, error) {
	var es = new(ExtSwitchData)

	if err := es.unmarshal(r); err != nil {
		return nil, err
	}

	return es, nil
}
