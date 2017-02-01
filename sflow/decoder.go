//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder.go
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
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	DataFlowSample    = 1 // Packet Flow Sampling
	DataCounterSample = 2 // Counter Sampling
)

type SFDecoder struct {
	reader io.ReadSeeker
	filter []uint32 // Filter data format(s)
}

type SFDatagram struct {
	Version    uint32 // Datagram version
	IPVersion  uint32 // Data gram sFlow version
	AgentSubId uint32 // Identifies a source of sFlow data
	SequenceNo uint32 // Sequence of sFlow Datagrams
	SysUpTime  uint32 // Current time (in milliseconds since device last booted
	SamplesNo  uint32 // Number of samples

	IPAddress net.IP // Agent IP address
}

type SFSampledHeader struct {
	HeaderProtocol uint32 // (enum SFHeaderProtocol)
	FrameLength    uint32 // Original length of packet before sampling
	Stripped       uint32 // Header/trailer bytes stripped by sender
	HeaderLength   uint32 // Length of sampled header bytes to follow
	HeaderBytes    []byte // Header bytes
}

type SFSample interface{}

var (
	nonEnterpriseStandard = errors.New("the enterprise is not standard sflow data")
	dataLengthUnknown     = errors.New("the sflow data length is unknown")
	sfVersionNotSupport   = errors.New("the sflow version doesn't support")
)

func NewSFDecoder(r io.ReadSeeker, f []uint32) SFDecoder {
	return SFDecoder{
		reader: r,
		filter: f,
	}
}

func (d *SFDecoder) SFDecode() ([]interface{}, error) {
	datagram, err := d.sfHeaderDecode()
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < datagram.SamplesNo; i++ {
		sfTypeFormat, sfDataLength, err := d.getSampleInfo()
		if err != nil {
			return nil, err
		}

		if m := d.isFilterMatch(sfTypeFormat); m {
			d.reader.Seek(int64(sfDataLength), 1)
			continue
		}

		switch sfTypeFormat {
		case DataFlowSample:
			h, err := decodeFlowSample(d.reader)
			return h, err
		case DataCounterSample:
			d.reader.Seek(int64(sfDataLength), 1)
		default:
			d.reader.Seek(int64(sfDataLength), 1)

		}

	}

	return nil, nil
}

func (d *SFDecoder) sfHeaderDecode() (*SFDatagram, error) {
	var (
		datagram     = &SFDatagram{}
		ipLen    int = 4
		err      error
	)

	if err = read(d.reader, &datagram.Version); err != nil {
		return nil, err
	}

	if datagram.Version != 5 {
		return nil, sfVersionNotSupport
	}

	if err = read(d.reader, &datagram.IPVersion); err != nil {
		return nil, err
	}

	// read the agent ip address
	if datagram.IPVersion == 2 {
		ipLen = 16
	}
	buff := make([]byte, ipLen)
	if _, err = d.reader.Read(buff); err != nil {
		return nil, err
	}
	datagram.IPAddress = buff

	if err = read(d.reader, &datagram.AgentSubId); err != nil {
		return nil, err
	}
	if err = read(d.reader, &datagram.SequenceNo); err != nil {
		return nil, err
	}
	if err = read(d.reader, &datagram.SysUpTime); err != nil {
		return nil, err
	}
	if err = read(d.reader, &datagram.SamplesNo); err != nil {
		return nil, err
	}

	return datagram, nil
}

func (d *SFDecoder) getSampleInfo() (uint32, uint32, error) {
	var (
		sfType           uint32
		sfTypeFormat     uint32
		sfTypeEnterprise uint32
		sfDataLength     uint32

		err error
	)

	if err = read(d.reader, &sfType); err != nil {
		return 0, 0, err
	}

	sfTypeEnterprise = sfType >> 12 // 20 bytes enterprise
	sfTypeFormat = sfType & 0xfff   // 12 bytes format

	// supports standard sflow data
	if sfTypeEnterprise != 0 {
		d.reader.Seek(int64(sfDataLength), 1)
		return 0, 0, nonEnterpriseStandard
	}

	if err = read(d.reader, &sfDataLength); err != nil {
		return 0, 0, dataLengthUnknown
	}

	return sfTypeFormat, sfDataLength, nil
}

func (d *SFDecoder) isFilterMatch(f uint32) bool {
	for _, v := range d.filter {
		if v == f {
			return true
		}
	}
	return false
}

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}
