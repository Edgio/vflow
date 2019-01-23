//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    marshal.go
//: details: encoding of each decoded IPFIX data sets
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

package ipfix

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
)

var errUknownMarshalDataType = errors.New("unknown data type to marshal")

// JSONMarshal encodes IPFIX message
func (m *Message) JSONMarshal(b *bytes.Buffer) ([]byte, error) {
	b.WriteString("{")

	// encode agent id
	m.encodeAgent(b)

	// encode header
	m.encodeHeader(b)

	// encode data sets
	if err := m.encodeDataSet(b); err != nil {
		return nil, err
	}

	b.WriteString("}")

	return b.Bytes(), nil
}

func (m *Message) encodeDataSet(b *bytes.Buffer) error {
	var (
		length   int
		dsLength int
		err      error
	)

	b.WriteString("\"DataSets\":")
	dsLength = len(m.DataSets)

	b.WriteByte('[')

	for i := range m.DataSets {
		length = len(m.DataSets[i])

		b.WriteByte('[')
		for j := range m.DataSets[i] {
			b.WriteString("{\"I\":")
			b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].ID), 10))
			b.WriteString(",\"V\":")
			err = m.writeValue(b, i, j)

			if m.DataSets[i][j].EnterpriseNo != 0 {
				b.WriteString(",\"E\":")
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].EnterpriseNo), 10))
			}

			if j < length-1 {
				b.WriteString("},")
			} else {
				b.WriteByte('}')
			}
		}

		if i < dsLength-1 {
			b.WriteString("],")
		} else {
			b.WriteByte(']')
		}
	}

	b.WriteByte(']')

	return err
}

func (m *Message) encodeDataSetFlat(b *bytes.Buffer) error {
	var (
		length   int
		dsLength int
		err      error
	)

	b.WriteString("\"DataSets\":")
	dsLength = len(m.DataSets)

	b.WriteByte('[')

	for i := range m.DataSets {
		length = len(m.DataSets[i])

		b.WriteByte('{')
		for j := range m.DataSets[i] {
			b.WriteByte('"')
			b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].ID), 10))
			b.WriteString("\":")
			err = m.writeValue(b, i, j)

			if j < length-1 {
				b.WriteByte(',')
			} else {
				b.WriteByte('}')
			}
		}

		if i < dsLength-1 {
			b.WriteString(",")
		}
	}

	b.WriteByte(']')

	return err
}

func (m *Message) encodeHeader(b *bytes.Buffer) {
	b.WriteString("\"Header\":{\"Version\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Version), 10))
	b.WriteString(",\"Length\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Length), 10))
	b.WriteString(",\"ExportTime\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.ExportTime), 10))
	b.WriteString(",\"SequenceNo\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.SequenceNo), 10))
	b.WriteString(",\"DomainID\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.DomainID), 10))
	b.WriteString("},")
}

func (m *Message) encodeAgent(b *bytes.Buffer) {
	b.WriteString("\"AgentID\":\"")
	b.WriteString(m.AgentID)
	b.WriteString("\",")
}

func (m *Message) writeValue(b *bytes.Buffer, i, j int) error {
	switch m.DataSets[i][j].Value.(type) {
	case uint:
		b.WriteString(strconv.FormatUint(uint64(m.DataSets[i][j].Value.(uint)), 10))
	case uint8:
		b.WriteString(strconv.FormatUint(uint64(m.DataSets[i][j].Value.(uint8)), 10))
	case uint16:
		b.WriteString(strconv.FormatUint(uint64(m.DataSets[i][j].Value.(uint16)), 10))
	case uint32:
		b.WriteString(strconv.FormatUint(uint64(m.DataSets[i][j].Value.(uint32)), 10))
	case uint64:
		b.WriteString(strconv.FormatUint(m.DataSets[i][j].Value.(uint64), 10))
	case int:
		b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int)), 10))
	case int8:
		b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int8)), 10))
	case int16:
		b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int16)), 10))
	case int32:
		b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int32)), 10))
	case int64:
		b.WriteString(strconv.FormatInt(m.DataSets[i][j].Value.(int64), 10))
	case float32:
		b.WriteString(strconv.FormatFloat(float64(m.DataSets[i][j].Value.(float32)), 'E', -1, 32))
	case float64:
		b.WriteString(strconv.FormatFloat(m.DataSets[i][j].Value.(float64), 'E', -1, 64))
	case string:
		b.WriteByte('"')
		b.WriteString(m.DataSets[i][j].Value.(string))
		b.WriteByte('"')
	case net.IP:
		b.WriteByte('"')
		b.WriteString(m.DataSets[i][j].Value.(net.IP).String())
		b.WriteByte('"')
	case net.HardwareAddr:
		b.WriteByte('"')
		b.WriteString(m.DataSets[i][j].Value.(net.HardwareAddr).String())
		b.WriteByte('"')
	case []uint8:
		b.WriteByte('"')
		b.WriteString("0x" + hex.EncodeToString(m.DataSets[i][j].Value.([]uint8)))
		b.WriteByte('"')
	default:
		return errUknownMarshalDataType
	}

	return nil
}
