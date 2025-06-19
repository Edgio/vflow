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
	"net"

	"github.com/EdgeCast/vflow/packet"
)

const (
	// SFDataRawHeader is sFlow Raw Packet Header number
	SFDataRawHeader = 1

	// SFDataExtSwitch is sFlow Extended Switch Data number
	SFDataExtSwitch = 1001

	// SFDataExtRouter is sFlow Extended Router Data number
	SFDataExtRouter = 1002

	// SFDataExtGateway is sFlow Extended Gateway Data number
	SFDataExtGateway = 1003

	// SFDataExtUser is sFlow Extended User Data number
	SFDataExtUser = 1004

	// SFDataExtURL is sFlow Extended URL Data number
	SFDataExtURL = 1005

	// SFDataExtMPLS is sFlow Extended MPLS Data number
	SFDataExtMPLS = 1006

	// SFDataExtNAT is sFlow Extended NAT Data number
	SFDataExtNAT = 1007

	// SFDataExtMPLSTunnel is sFlow Extended MPLS Tunnel number
	SFDataExtMPLSTunnel = 1008

	// SFDataExtMPLSVC is sFlow Extended MPLS VC number
	SFDataExtMPLSVC = 1009

	// SFDataExtMPLSFTN is sFlow Extended MPLS FTN number
	SFDataExtMPLSFTN = 1010

	// SFDataExtMPLSLDP_FEC is sFlow Extended MPLS LDP FEC number
	SFDataExtMPLSLDP_FEC = 1011

	// SFDataExtVLANTunnel is sFlow Extended VLAN Tunnel number
	SFDataExtVLANTunnel = 1012

	// Arista-specific enterprise extensions
	// SFDataExtAristaBGP is Arista BGP Route Information extension
	SFDataExtAristaBGP = 1013

	// SFDataExtAristaVPLS is Arista VPLS extension
	SFDataExtAristaVPLS = 1014

	// SFDataExtAristaDSCP is Arista DSCP extension
	SFDataExtAristaDSCP = 1015
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

// DSCPInfo represents DSCP field information for Arista extensions
type DSCPInfo struct {
	OriginalDSCP uint8 // Original DSCP value before rewriting
	RewrittenDSCP uint8 // DSCP value after rewriting (if applicable)
	DSCPRewritten bool  // Flag indicating if DSCP was rewritten
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

// ExtAristaBGPData represents Arista BGP Route Information extension
type ExtAristaBGPData struct {
	NextHop       net.IP   // BGP next hop IP address
	ASPath        []uint32 // AS path sequence
	Communities   []uint32 // BGP communities
	LocalPref     uint32   // Local preference
	SourceAS      uint32   // Source AS number
	DestAS        uint32   // Destination AS number
	PeerAS        uint32   // Peer AS number
	MED           uint32   // Multi-exit discriminator
	Origin        uint32   // BGP origin attribute
}

// ExtAristaVPLSData represents Arista VPLS extension
type ExtAristaVPLSData struct {
	InstanceName string // VPLS instance name
	PseudowireID uint32 // Pseudowire ID
	VCID         uint32 // VC ID
	VCType       uint32 // VC Type
}

var (
	errMaxOutEthernetLength = errors.New("the ethernet length is greater than 1500")
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

func (eab *ExtAristaBGPData) unmarshal(r io.Reader, l uint32) error {
	var err error
	var ipVersion uint32
	var pathLen uint32
	var commLen uint32

	// Read IP version for next hop
	if err = read(r, &ipVersion); err != nil {
		return err
	}

	// Read next hop IP address
	ipLen := 4
	if ipVersion == 2 {
		ipLen = 16
	}
	nextHopBuff := make([]byte, ipLen)
	if _, err = r.Read(nextHopBuff); err != nil {
		return err
	}
	eab.NextHop = nextHopBuff

	// Read AS path length and AS path
	if err = read(r, &pathLen); err != nil {
		return err
	}
	eab.ASPath = make([]uint32, pathLen)
	for i := uint32(0); i < pathLen; i++ {
		if err = read(r, &eab.ASPath[i]); err != nil {
			return err
		}
	}

	// Read communities length and communities
	if err = read(r, &commLen); err != nil {
		return err
	}
	eab.Communities = make([]uint32, commLen)
	for i := uint32(0); i < commLen; i++ {
		if err = read(r, &eab.Communities[i]); err != nil {
			return err
		}
	}

	// Read remaining BGP attributes
	if err = read(r, &eab.LocalPref); err != nil {
		return err
	}
	if err = read(r, &eab.SourceAS); err != nil {
		return err
	}
	if err = read(r, &eab.DestAS); err != nil {
		return err
	}
	if err = read(r, &eab.PeerAS); err != nil {
		return err
	}
	if err = read(r, &eab.MED); err != nil {
		return err
	}
	err = read(r, &eab.Origin)

	return err
}

func (dscp *DSCPInfo) unmarshal(r io.Reader) error {
	var err error

	if err = read(r, &dscp.OriginalDSCP); err != nil {
		return err
	}

	if err = read(r, &dscp.RewrittenDSCP); err != nil {
		return err
	}

	var rewrittenFlag uint8
	if err = read(r, &rewrittenFlag); err != nil {
		return err
	}
	dscp.DSCPRewritten = rewrittenFlag != 0

	// Skip padding byte to align to 4-byte boundary
	var padding uint8
	err = read(r, &padding)

	return err
}

func (eav *ExtAristaVPLSData) unmarshal(r io.Reader, l uint32) error {
	var err error
	var nameLen uint32

	// Read instance name length and name
	if err = read(r, &nameLen); err != nil {
		return err
	}
	
	nameBuff := make([]byte, nameLen)
	if _, err = r.Read(nameBuff); err != nil {
		return err
	}
	eav.InstanceName = string(nameBuff)

	// Read padding to align to 4-byte boundary
	padding := (4 - nameLen%4) % 4
	if padding > 0 {
		paddingBuff := make([]byte, padding)
		if _, err = r.Read(paddingBuff); err != nil {
			return err
		}
	}

	// Read VPLS identifiers
	if err = read(r, &eav.PseudowireID); err != nil {
		return err
	}
	if err = read(r, &eav.VCID); err != nil {
		return err
	}
	err = read(r, &eav.VCType)

	return err
}

func decodeFlowSample(r io.ReadSeeker) (*FlowSample, error) {
	var (
		fs          = new(FlowSample)
		rTypeFormat uint32
		rTypeLength uint32
		err         error
	)

	if err = fs.unmarshal(r); err != nil {
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
		case SFDataExtAristaBGP:
			d, err := decodeExtAristaBGPData(r, rTypeLength)
			if err != nil {
				return fs, err
			}

			fs.Records["ExtAristaBGP"] = d
		case SFDataExtAristaVPLS:
			d, err := decodeExtAristaVPLSData(r, rTypeLength)
			if err != nil {
				return fs, err
			}

			fs.Records["ExtAristaVPLS"] = d
		case SFDataExtAristaDSCP:
			d, err := decodeExtAristaDSCPData(r)
			if err != nil {
				return fs, err
			}

			fs.Records["ExtAristaDSCP"] = d
		default:
			r.Seek(int64(rTypeLength), 1)
		}
	}

	return fs, nil
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

func decodeExtAristaBGPData(r io.Reader, l uint32) (*ExtAristaBGPData, error) {
	var eab = new(ExtAristaBGPData)

	if err := eab.unmarshal(r, l); err != nil {
		return nil, err
	}

	return eab, nil
}

func decodeExtAristaVPLSData(r io.Reader, l uint32) (*ExtAristaVPLSData, error) {
	var eav = new(ExtAristaVPLSData)

	if err := eav.unmarshal(r, l); err != nil {
		return nil, err
	}

	return eav, nil
}

func decodeExtAristaDSCPData(r io.Reader) (*DSCPInfo, error) {
	var dscp = new(DSCPInfo)

	if err := dscp.unmarshal(r); err != nil {
		return nil, err
	}

	return dscp, nil
}
