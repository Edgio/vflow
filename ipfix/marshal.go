// Package ipfix decodes IPFIX packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    marshal.go
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
package ipfix

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

var errUknownMarshalDataType = errors.New("unknown data type to marshal")

// JSONMarshal encodes IPFIX message
func (m *Message) JSONMarshal(e *bytes.Buffer) ([]byte, error) {
	e.WriteString("{")

	// encode agent id
	m.encodeAgent(e)

	// encode header
	m.encodeHeader(e)

	// encode data sets
	if err := m.encodeDataSet(e); err != nil {
		return nil, err
	}

	e.WriteString("}")

	return e.Bytes(), nil
}

func (m *Message) encodeDataSet(e *bytes.Buffer) error {
	var length, dsLength int

	e.WriteString("\"DataSets\":")
	dsLength = len(m.DataSets)

	e.WriteString("[")
	for i := range m.DataSets {
		length = len(m.DataSets[i])

		e.WriteString("[")
		for j := range m.DataSets[i] {
			e.WriteString("{\"ID\":")
			e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].ID), 10))
			e.WriteString(",\"Value\":")

			switch m.DataSets[i][j].Value.(type) {
			case uint:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint)), 10))
			case uint8:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint8)), 10))
			case uint16:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint16)), 10))
			case uint32:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint32)), 10))
			case uint64:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint64)), 10))
			case int:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int)), 10))
			case int8:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int8)), 10))
			case int16:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int16)), 10))
			case int32:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int32)), 10))
			case int64:
				e.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int64)), 10))
			case float32:
				e.WriteString(strconv.FormatFloat(float64(m.DataSets[i][j].Value.(float32)), 'E', -1, 32))
			case float64:
				e.WriteString(strconv.FormatFloat(m.DataSets[i][j].Value.(float64), 'E', -1, 64))
			case string:
				e.WriteString("\"")
				e.WriteString(m.DataSets[i][j].Value.(string))
				e.WriteString("\"")
			case net.IP:
				e.WriteString("\"")
				e.WriteString(m.DataSets[i][j].Value.(net.IP).String())
				e.WriteString("\"")
			case net.HardwareAddr:
				e.WriteString("\"")
				e.WriteString(m.DataSets[i][j].Value.(net.HardwareAddr).String())
				e.WriteString("\"")
			case []uint8:
				e.WriteString("\"")
				e.WriteString(fmt.Sprintf("0x%x", m.DataSets[i][j].Value.([]uint8)))
				e.WriteString("\"")
			default:
				return errUknownMarshalDataType
			}
			if j < length-1 {
				e.WriteString("},")
			} else {
				e.WriteString("}")
			}
		}
		if i < dsLength-1 {
			e.WriteString("],")
		} else {
			e.WriteString("]")
		}
	}
	e.WriteString("]")

	return nil
}

func (m *Message) encodeHeader(e *bytes.Buffer) {
	e.WriteString("\"Header\":{\"Version\":")
	e.WriteString(strconv.FormatInt(int64(m.Header.Version), 10))
	e.WriteString(",\"Length\":")
	e.WriteString(strconv.FormatInt(int64(m.Header.Length), 10))
	e.WriteString(",\"ExportTime\":")
	e.WriteString(strconv.FormatInt(int64(m.Header.ExportTime), 10))
	e.WriteString(",\"SequenceNo\":")
	e.WriteString(strconv.FormatInt(int64(m.Header.SequenceNo), 10))
	e.WriteString(",\"DomainID\":")
	e.WriteString(strconv.FormatInt(int64(m.Header.DomainID), 10))
	e.WriteString("},")
}

func (m *Message) encodeAgent(e *bytes.Buffer) {
	e.WriteString("\"AgentID\":\"")
	e.WriteString(m.AgentID)
	e.WriteString("\",")
}
