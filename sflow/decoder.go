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
	"time"
)

const (
	// DataFlowSample defines packet flow sampling
	DataFlowSample = 1

	// DataCounterSample defines counter sampling
	DataCounterSample = 2
)

// SFDecoder represents sFlow decoder
type SFDecoder struct {
	reader io.ReadSeeker
	filter []uint32 // Filter data format(s)
}

// SFDatagram represents sFlow datagram
type SFDatagram struct {
	Version    uint32 // Datagram version
	IPVersion  uint32 // Data gram sFlow version
	AgentSubID uint32 // Identifies a source of sFlow data
	SequenceNo uint32 // Sequence of sFlow Datagrams
	SysUpTime  uint32 // Current time (in milliseconds since device last booted
	SamplesNo  uint32 // Number of samples
	Samples    []Sample
	Counters   []Counter

	IPAddress net.IP // Agent IP address
	ColTime   int64  // Collected time
}

// SFSampledHeader represents sFlow sample header
type SFSampledHeader struct {
	HeaderProtocol uint32 // (enum SFHeaderProtocol)
	FrameLength    uint32 // Original length of packet before sampling
	Stripped       uint32 // Header/trailer bytes stripped by sender
	HeaderLength   uint32 // Length of sampled header bytes to follow
	HeaderBytes    []byte // Header bytes
}

// Sample represents sFlow sample flow
type Sample interface{}

// Counter represents sFlow counters
type Counter interface{}

// Record represents sFlow sample record record
type Record interface{}

var (
	errNoneEnterpriseStandard = errors.New("the enterprise is not standard sflow data")
	errDataLengthUnknown      = errors.New("the sflow data length is unknown")
	errSFVersionNotSupport    = errors.New("the sflow version doesn't support")
)

// NewSFDecoder constructs new sflow decoder
func NewSFDecoder(r io.ReadSeeker, f []uint32) SFDecoder {
	return SFDecoder{
		reader: r,
		filter: f,
	}
}

// SFDecode decodes sFlow data
func (d *SFDecoder) SFDecode() (*SFDatagram, error) {
	datagram, err := d.sfHeaderDecode()
	if err != nil {
		return nil, err
	}

	datagram.Samples = []Sample{}
	datagram.Counters = []Counter{}

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
			d, err := decodeFlowSample(d.reader)
			if err != nil {
				return datagram, err
			}
			datagram.Samples = append(datagram.Samples, d)
		case DataCounterSample:
			d, err := decodeFlowCounter(d.reader)
			if err != nil {
				return datagram, err
			}
			datagram.Counters = append(datagram.Counters, d)
		default:
			d.reader.Seek(int64(sfDataLength), 1)
		}

	}

	return datagram, nil
}

func (d *SFDecoder) sfHeaderDecode() (*SFDatagram, error) {
	var (
		datagram = &SFDatagram{}
		ipLen    = 4
		err      error
	)

	if err = read(d.reader, &datagram.Version); err != nil {
		return nil, err
	}

	if datagram.Version != 5 {
		return nil, errSFVersionNotSupport
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

	if err = read(d.reader, &datagram.AgentSubID); err != nil {
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

	datagram.ColTime = time.Now().Unix()

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
		return 0, 0, errNoneEnterpriseStandard
	}

	if err = read(d.reader, &sfDataLength); err != nil {
		return 0, 0, errDataLengthUnknown
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
