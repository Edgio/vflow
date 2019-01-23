//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    marshal.go
//: details: encoding of each decoded netflow v5 flow set
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
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

var errUknownMarshalDataType = errors.New("unknown data type to marshal")

// JSONMarshal encodes netflow v9 message
func (m *Message) JSONMarshal(b *bytes.Buffer) ([]byte, error) {
	b.WriteString("{")

	// encode agent id
	m.encodeAgent(b)

	// encode header
	m.encodeHeader(b)

	// encode flows
	// encode data sets
	if err := m.encodeFlows(b); err != nil {
		return nil, err
	}

	b.WriteString("}")

	return b.Bytes(), nil
}

func (m *Message) encodeHeader(b *bytes.Buffer) {
	b.WriteString("\"Header\":{\"Version\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Version), 10))
	b.WriteString(",\"Count\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Count), 10))
	b.WriteString(",\"SysUpTimeMSecs\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.SysUpTimeMSecs), 10))
	b.WriteString(",\"UNIXSecs\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.UNIXSecs), 10))
	b.WriteString(",\"UNIXNSecs\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.UNIXNSecs), 10))
	b.WriteString(",\"SeqNum\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.SeqNum), 10))
	b.WriteString(",\"EngType\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.EngType), 10))
	b.WriteString(",\"EngID\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.EngID), 10))
	b.WriteString(",\"SmpInt\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.SmpInt), 10))
	b.WriteString("},")
}

func (m *Message) encodeAgent(b *bytes.Buffer) {
	b.WriteString("\"AgentID\":\"")
	b.WriteString(m.AgentID)
	b.WriteString("\",")
}

func (m *Message) encodeFlow(r FlowRecord, b *bytes.Buffer) {

	ip := make(net.IP, 4)

	b.WriteString("\"SrcAddr\":\"")
	binary.BigEndian.PutUint32(ip, r.SrcAddr)
	b.WriteString(ip.String())

	b.WriteString("\",\"DstAddr\":\"")
	binary.BigEndian.PutUint32(ip, r.DstAddr)
	b.WriteString(ip.String())

	b.WriteString("\",\"NextHop\":\"")
	binary.BigEndian.PutUint32(ip, r.NextHop)
	b.WriteString(ip.String())

	b.WriteString("\",\"Input\":")
	b.WriteString(strconv.FormatInt(int64(r.Input), 10))
	b.WriteString(",\"Output\":")
	b.WriteString(strconv.FormatInt(int64(r.Output), 10))
	b.WriteString(",\"PktCount\":")
	b.WriteString(strconv.FormatInt(int64(r.PktCount), 10))
	b.WriteString(",\"L3Octets\":")
	b.WriteString(strconv.FormatInt(int64(r.L3Octets), 10))

	// if these ever need to be translated actual time
	// then this will require knowing some information from
	// the header. I believe the basic formula is
	// UnixSecs - SysUpTime + StartTime

	b.WriteString(",\"StartTime\":")
	b.WriteString(strconv.FormatInt(int64(r.StartTime), 10))
	b.WriteString(",\"EndTime\":")
	b.WriteString(strconv.FormatInt(int64(r.EndTime), 10))

	b.WriteString(",\"SrcPort\":")
	b.WriteString(strconv.FormatInt(int64(r.SrcPort), 10))
	b.WriteString(",\"DstPort\":")
	b.WriteString(strconv.FormatInt(int64(r.DstPort), 10))
	b.WriteString(",\"Padding1\":")
	b.WriteString(strconv.FormatInt(int64(r.Padding1), 10))
	b.WriteString(",\"TCPFlags\":")
	b.WriteString(strconv.FormatInt(int64(r.TCPFlags), 10))
	b.WriteString(",\"ProtType\":")
	b.WriteString(strconv.FormatInt(int64(r.ProtType), 10))
	b.WriteString(",\"Tos\":")
	b.WriteString(strconv.FormatInt(int64(r.Tos), 10))
	b.WriteString(",\"SrcAsNum\":")
	b.WriteString(strconv.FormatInt(int64(r.SrcAsNum), 10))
	b.WriteString(",\"DstAsNum\":")
	b.WriteString(strconv.FormatInt(int64(r.DstAsNum), 10))
	b.WriteString(",\"SrcMask\":")
	b.WriteString(strconv.FormatInt(int64(r.SrcMask), 10))
	b.WriteString(",\"DstMask\":")
	b.WriteString(strconv.FormatInt(int64(r.DstMask), 10))
	b.WriteString(",\"Padding2\":")
	b.WriteString(strconv.FormatInt(int64(r.Padding2), 10))
}

func (m *Message) encodeFlows(b *bytes.Buffer) error {
	var (
		fLength int
		err     error
	)

	b.WriteString("\"Flows\":")
	fLength = len(m.Flows)

	b.WriteByte('[')

	for i := range m.Flows {
		b.WriteString("{")
		m.encodeFlow(m.Flows[i], b)
		b.WriteString("}")
		if i < fLength-1 {
			b.WriteString(",")
		}
	}

	b.WriteByte(']')

	return err
}
