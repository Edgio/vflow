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
	"fmt"
	"io"
	"net"

	"github.com/guardicore/vflow/packet"
)

const (
	// SFDataRawHeader is sFlow Raw Packet Header number
	SFDataRawHeader = 1

	// SFDataEthernetFrame is sFlow Ethernet Frame Data
	// (parsing is not implemented yet)
	SFDataEthernetFrame = 2

	// SFDataPacketIPV4 is sFlow IP Version 4 Data number
	SFDataPacketIPV4 = 3

	// SFDataPacketIPV6 is sFlow IP Version 6 Data number
	SFDataPacketIPV6 = 4

	// SFDataExtSwitch is sFlow Extended Switch Data number
	SFDataExtSwitch = 1001

	// SFDataExtRouter is sFlow Extended Router Data number
	SFDataExtRouter = 1002
)

// FlowSample represents single flow sample
type FlowSample struct {
	SequenceNo   uint32 // Incremented with each flow sample
	SourceIDType uint32 // sfSourceID type
	SourceIDIdx  uint32 // sfSourceID index
	SamplingRate uint32 // sfPacketSamplingRate
	SamplePool   uint32 // Total number of packets that could have been sampled
	Drops        uint32 // Number of times a packet was dropped due to lack of resources
	InputType    uint32 // SNMP ifType of input interface
	InputIdx     uint32 // SNMP ifIndex of input interface
	OutputType   uint32 // SNMP ifType of output interface
	OutputIdx    uint32 // SNMP ifIndex of output interface
	RecordsNo    uint32 // Number of records to follow
	Records      map[string]Record
}

// SampledHeader represents sampled header
type SampledHeader struct {
	Protocol     uint32 // (enum SFLHeader_protocol)
	FrameLength  uint32 // Original length of packet before sampling
	Stripped     uint32 // Header/trailer bytes stripped by sender
	HeaderLength uint32 // Length of sampled header bytes to follow
	Header       []byte // Header bytes
}

// IPV4Data represents Packet IP Version 4 Data
type IPV4Data struct {
	Len      uint32 // The length of the IP packet excluding lower layer encapsulations
	Protocol uint32 // IP Protocol type: TCP = 6, UDP = 17)
	SrcIP    net.IP // Source IP Address
	DstIP    net.IP // Destination IP Address
	SrcPort  uint32 // TCP/UDP source port number or equivalent
	DstPort  uint32 // TCP/UDP destination port number or equivalent
	TcpFlags uint32 // TCP flags
	TOS      uint32 // IP type of service
}

// IPV6Data represents Packet IP Version 6 Data
type IPV6Data struct {
	Len      uint32 // The length of the IP packet excluding lower layer encapsulations
	Protocol uint32 // IP Protocol type: TCP = 6, UDP = 17)
	SrcIP    net.IP // Source IP Address
	DstIP    net.IP // Destination IP Address
	SrcPort  uint32 // TCP/UDP source port number or equivalent
	DstPort  uint32 // TCP/UDP destination port number or equivalent
	TcpFlags uint32 // TCP flags
	Priority uint32 // IP priority
}

// ExtSwitchData represents Extended Switch Data
type ExtSwitchData struct {
	SrcVlan     uint32 // The 802.1Q VLAN id of incoming frame
	SrcPriority uint32 // The 802.1p priority of incoming frame
	DstVlan     uint32 // The 802.1Q VLAN id of outgoing frame
	DstPriority uint32 // The 802.1p priority of outgoing frame
}

// ExtRouterData represents extended router data
type ExtRouterData struct {
	NextHop net.IP
	SrcMask uint32
	DstMask uint32
}

var (
	errMaxOutEthernetLength = errors.New("the ethernet length is greater than 1500")
)

func (fs *FlowSample) unmarshal(r io.ReadSeeker, expanded bool) error {
	var (
		err error
		val uint32
	)

	if err = read(r, &fs.SequenceNo); err != nil {
		return err
	}

	if !expanded {
		if err = read(r, &val); err != nil {
			return err
		}
		fs.SourceIDType = (val >> 24) & 0xFF
		fs.SourceIDIdx = val & 0xFFFFFF
	} else {
		if err = read(r, &fs.SourceIDType); err != nil {
			return err
		}
		if err = read(r, &fs.SourceIDIdx); err != nil {
			return err
		}
	}

	if err = read(r, &fs.SamplingRate); err != nil {
		return err
	}

	if err = read(r, &fs.SamplePool); err != nil {
		return err
	}

	if err = read(r, &fs.Drops); err != nil {
		return err
	}

	if !expanded {
		if err = read(r, &val); err != nil {
			return err
		}
		fs.InputType = (val >> 16) & 0xFFFF // 2 most significant bytes
		fs.InputIdx = val & 0xFFFF          // 2 least significant bytes
	} else {
		if err = read(r, &fs.InputType); err != nil {
			return err
		}
		if err = read(r, &fs.InputIdx); err != nil {
			return err
		}
	}

	if !expanded {
		if err = read(r, &val); err != nil {
			return err
		}
		fs.OutputType = (val >> 16) & 0xFFFF // 2 most significant bytes
		fs.OutputIdx = val & 0xFFFF          // 2 least significant bytes
	} else {
		if err = read(r, &fs.OutputType); err != nil {
			return err
		}
		if err = read(r, &fs.OutputIdx); err != nil {
			return err
		}
	}

	err = read(r, &fs.RecordsNo)

	return err
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

func (d *IPV4Data) unmarshal(r io.Reader) error {
	var err error

	if err = read(r, &d.Len); err != nil {
		return err
	}

	if err = read(r, &d.Protocol); err != nil {
		return err
	}

	buff := make([]byte, 4)
	if err = read(r, &buff); err != nil {
		return err
	}
	d.SrcIP = net.IPv4(buff[0], buff[1], buff[2], buff[3])

	if err = read(r, &buff); err != nil {
		return err
	}
	d.DstIP = net.IPv4(buff[0], buff[1], buff[2], buff[3])

	if err = read(r, &d.SrcPort); err != nil {
		return err
	}

	if err = read(r, &d.DstPort); err != nil {
		return err
	}

	if err = read(r, &d.TcpFlags); err != nil {
		return err
	}

	if err = read(r, &d.TOS); err != nil {
		return err
	}

	return nil
}

func (d *IPV6Data) unmarshal(r io.Reader) error {
	var err error

	if err = read(r, &d.Len); err != nil {
		return err
	}

	if err = read(r, &d.Protocol); err != nil {
		return err
	}

	buff := make([]byte, 16)
	if err = read(r, &buff); err != nil {
		return err
	}
	d.SrcIP = buff

	buff = make([]byte, 16)
	if err = read(r, &buff); err != nil {
		return err
	}
	d.DstIP = buff

	if err = read(r, &d.SrcPort); err != nil {
		return err
	}

	if err = read(r, &d.DstPort); err != nil {
		return err
	}

	if err = read(r, &d.TcpFlags); err != nil {
		return err
	}

	if err = read(r, &d.Priority); err != nil {
		return err
	}

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

	err = read(r, &es.SrcPriority)

	return err
}

func (er *ExtRouterData) unmarshal(r io.Reader, l uint32) error {
	var err error

	buff := make([]byte, l-8)
	if err = read(r, &buff); err != nil {
		return err
	}
	er.NextHop = buff[4:]

	if err = read(r, &er.SrcMask); err != nil {
		return err
	}

	err = read(r, &er.DstMask)

	return err
}

func decodeFlowSample(r io.ReadSeeker, expanded bool) (*FlowSample, error) {
	var (
		fs          = new(FlowSample)
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = fs.unmarshal(r, expanded); err != nil {
		return nil, err
	}

	fs.Records = make(map[string]Record)

	for i := uint32(0); i < fs.RecordsNo; i++ {
		if err = read(r, &rTypeFormat); err != nil {
			return nil, err
		}
		if err = read(r, &rTypeLength); err != nil {
			return nil, err
		}

		switch rTypeFormat {
		case SFDataRawHeader:
			d, err := decodeSampledHeader(r)
			if err != nil {
				return fs, err
			}
			fs.Records["RawHeader"] = d
		case SFDataPacketIPV4:
			d, err := decodeIPV4Data(r)
			if err != nil {
				return fs, err
			}
			fs.Records["IPV4"] = d
		case SFDataPacketIPV6:
			d, err := decodeIPV6Data(r)
			if err != nil {
				return fs, err
			}
			fs.Records["IPV6"] = d
		case SFDataExtSwitch:
			d, err := decodeExtSwitchData(r)
			if err != nil {
				return fs, err
			}

			fs.Records["ExtSwitch"] = d
		case SFDataExtRouter:
			d, err := decodeExtRouterData(r, rTypeLength)
			if err != nil {
				return fs, err
			}
			fs.Records["ExtRouter"] = d
		default:
			err = fmt.Errorf("sflow: unknown rTypeFormat: %d/0x%x, size: %d", rTypeFormat, rTypeFormat, rTypeLength)
			r.Seek(int64(rTypeLength), 1)
		}
	}

	return fs, err
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
	d, err := p.Decoder(h.Header, h.Protocol)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func decodeIPV4Data(r io.Reader) (*IPV4Data, error) {
	var d = new(IPV4Data)

	if err := d.unmarshal(r); err != nil {
		return nil, err
	}

	return d, nil
}

func decodeIPV6Data(r io.Reader) (*IPV6Data, error) {
	var d = new(IPV6Data)

	if err := d.unmarshal(r); err != nil {
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

func decodeExtRouterData(r io.Reader, l uint32) (*ExtRouterData, error) {
	var er = new(ExtRouterData)

	if err := er.unmarshal(r, l); err != nil {
		return nil, err
	}

	return er, nil
}
