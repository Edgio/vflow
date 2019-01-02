//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder.go
//: details: decodes netflow version 5 packets
//: author:  Christopher Noel
//: date:    12/10/2018
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

package netflow5

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/VerizonDigital/vflow/reader"
)

type nonfatalError error

// PacketHeader represents Netflow v5 packet header
// Based on docs at https://www.plixer.com/support/netflow-v5/
// 24 bytes long
type PacketHeader struct {
	Version        uint16 // Version of Flow Record format exported in this packet
	Count          uint16 // The total number of flows in the Export Packet
	SysUpTimeMSecs uint32 // Time in milliseconds since this device was first booted
	UNIXSecs       uint32 // Time in seconds since 0000 UTC 1970
	UNIXNSecs      uint32 // Residual nanoseconds since 0000 UTC 1970
	SeqNum         uint32 // Incremental sequence counter of total flows
	EngType        uint8  // An 8-bit value that identifies the type of flow-switching engine
	EngID          uint8  // An 8-bit value that identifies the Slot number of the flow-switching engine
	SmpInt         uint16 // A 16-bit value that identifies the Sampling Interval
	// The Sampling Interval - first 2 bits are the sampling mode, the last 14 bits hold the sampling interval
}

// FlowRecord represents Netflow v5 flow
// Based on docs at https://www.plixer.com/support/netflow-v5/
// 48 bytes long
type FlowRecord struct {
	SrcAddr   uint32 // Source IP Address
	DstAddr   uint32 // Destination IP Address
	NextHop   uint32 // IP Address of the next hop router
	Input     uint16 // SNMP index of input interface
	Output    uint16 // SNMP index of output interface
	PktCount  uint32 // Number of packets in the flow
	L3Octets  uint32 // Total number of Layer 3 bytes in the packets of the flow
	StartTime uint32 // SysUptime at start of flow in ms since last boot
	EndTime   uint32 // SysUptime at end of the flow in ms since last boot
	SrcPort   uint16 // TCP/UDP source port number or equivalent
	DstPort   uint16 // TCP/UDP destination port number or equivalent
	Padding1  uint8  // Unused (zero) bytes
	TCPFlags  uint8  // Cumulative OR of TCP flags
	ProtType  uint8  // IP protocol type (for example, TCP = 6; UDP = 17)
	Tos       uint8  // IP type of service (ToS)
	SrcAsNum  uint16 // Autonomous system number of the source, either origin or peer
	DstAsNum  uint16 // Autonomous system number of the destination, either origin or peer
	SrcMask   uint8  // Source address prefix mask bits
	DstMask   uint8  // Destination address prefix mask bits
	Padding2  uint16 // Unused (zero) bytes
}

// Decoder represents Netflow payload and remote address
type Decoder struct {
	raddr  net.IP
	reader *reader.Reader
}

// Message represents Netflow v5 decoded data
type Message struct {
	AgentID string
	Header  PacketHeader
	Flows   []FlowRecord
}

//   The Packet Header format is specified as:
//
//    0                   1                   2                   3
//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |       Version Number          |            Count              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                           sysUpTime                           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                           UNIX Secs                           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                           UNIX NSecs                          |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                        Sequence Counter                       |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |  Engine Type  |   Engine ID   |       Sampling Interval       |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (h *PacketHeader) unmarshal(r *reader.Reader) error {
	var err error

	if h.Version, err = r.Uint16(); err != nil {
		return err
	}

	if h.Count, err = r.Uint16(); err != nil {
		return err
	}

	if h.SysUpTimeMSecs, err = r.Uint32(); err != nil {
		return err
	}

	if h.UNIXSecs, err = r.Uint32(); err != nil {
		return err
	}

	if h.UNIXNSecs, err = r.Uint32(); err != nil {
		return err
	}

	if h.SeqNum, err = r.Uint32(); err != nil {
		return err
	}

	if h.EngType, err = r.Uint8(); err != nil {
		return err
	}

	if h.EngID, err = r.Uint8(); err != nil {
		return err
	}

	if h.SmpInt, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

func (h *PacketHeader) validate() error {

	if h.Version != 5 {
		return fmt.Errorf("invalid netflow version, (expected: 5) (received: %d)", h.Version)
	} else if h.Count < 1 || h.Count > 30 {
		return fmt.Errorf("flow count out of bounds, (expected: [1...30]) (received: %d)", h.Count)
	}

	return nil
}

//   The Flow Record format is specified as:
//
//    0                   1                   2                   3
//    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                            Src Addr                           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                            Dst Addr                           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                            Next Hop                           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             Input             |             Output            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                          Packet Count                         |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                           Octet Count                         |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                         Flow Start Time                       |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                          Flow End Time                        |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |           Src Port            |           Dst Port            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |    Padding1   |   TCP Flags   |  Protocol     |     TOS       |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |             Src AS            |             Dst AS            |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |    Src Mask   |    Dst Mask   |            Padding2           |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (fr *FlowRecord) unmarshal(r *reader.Reader) error {
	var err error

	if fr.SrcAddr, err = r.Uint32(); err != nil {
		return err
	}

	if fr.DstAddr, err = r.Uint32(); err != nil {
		return err
	}

	if fr.NextHop, err = r.Uint32(); err != nil {
		return err
	}

	if fr.Input, err = r.Uint16(); err != nil {
		return err
	}

	if fr.Output, err = r.Uint16(); err != nil {
		return err
	}

	if fr.PktCount, err = r.Uint32(); err != nil {
		return err
	}

	if fr.L3Octets, err = r.Uint32(); err != nil {
		return err
	}

	if fr.StartTime, err = r.Uint32(); err != nil {
		return err
	}

	if fr.EndTime, err = r.Uint32(); err != nil {
		return err
	}

	if fr.SrcPort, err = r.Uint16(); err != nil {
		return err
	}

	if fr.DstPort, err = r.Uint16(); err != nil {
		return err
	}

	if fr.Padding1, err = r.Uint8(); err != nil {
		return err
	}

	if fr.TCPFlags, err = r.Uint8(); err != nil {
		return err
	}

	if fr.ProtType, err = r.Uint8(); err != nil {
		return err
	}

	if fr.Tos, err = r.Uint8(); err != nil {
		return err
	}

	if fr.SrcAsNum, err = r.Uint16(); err != nil {
		return err
	}

	if fr.DstAsNum, err = r.Uint16(); err != nil {
		return err
	}

	if fr.SrcMask, err = r.Uint8(); err != nil {
		return err
	}

	if fr.DstMask, err = r.Uint8(); err != nil {
		return err
	}

	if fr.Padding2, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

// NewDecoder constructs a decoder
func NewDecoder(raddr net.IP, b []byte) *Decoder {
	return &Decoder{raddr, reader.NewReader(b)}
}

// Decode decodes the flow records
func (d *Decoder) Decode() (*Message, error) {
	var msg = new(Message)

	// Decode the Packet Header
	if err := msg.Header.unmarshal(d.reader); err != nil {
		return nil, err
	}
	// Validate the Packet Header
	if err := msg.Header.validate(); err != nil {
		return nil, err
	}

	// Add source IP address as Agent ID
	msg.AgentID = d.raddr.String()

	// Decode the Flows
	var decodeErrors []error
	flowCount := int(msg.Header.Count)
	if err := d.decodeFlows(flowCount, msg); err != nil {
		switch err.(type) {
		case nonfatalError:
			decodeErrors = append(decodeErrors, err)
		default:
			return nil, err
		}
	}

	return msg, combineErrors(decodeErrors...)

}

func (d *Decoder) decodeFlows(flowCount int, msg *Message) error {
	remainingLen := d.reader.Len()
	expectedLen := flowCount * 48
	flowIndex := 0

	var err error

	if expectedLen > remainingLen {
		err = fmt.Errorf("Expect %v bytes to read, %v remaining bytes encountered", expectedLen, remainingLen)
	}

	// there should be *flowCount* number of flows in the message, each 48 bytes long
	for err == nil && flowIndex < flowCount {

		fr := FlowRecord{}
		err = fr.unmarshal(d.reader)
		if err == nil {
			msg.Flows = append(msg.Flows, fr)
		}
		flowIndex++
	}

	return err
}

func combineErrors(errorSlice ...error) (err error) {
	switch len(errorSlice) {
	case 0:
	case 1:
		err = errorSlice[0]
	default:
		var errMsg bytes.Buffer
		errMsg.WriteString("Multiple errors:")
		for _, subError := range errorSlice {
			errMsg.WriteString("\n- " + subError.Error())
		}
		err = errors.New(errMsg.String())
	}
	return
}
