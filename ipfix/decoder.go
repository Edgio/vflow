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
package ipfix

import (
	"errors"
	"io"
	"net"
	"sync"
)

type IPFIXDecoder struct {
	raddr  net.IP
	reader *Reader
}

type MessageHeader struct {
	Version    uint16 // Version of IPFIX to which this Message conforms
	Length     uint16 // Total length of the IPFIX Message, measured in octets
	ExportTime uint32 // Time at which the IPFIX Message Header leaves the Exporter
	SequenceNo uint32 // Incremental sequence counter modulo 2^32
	DomainID   uint32 // A 32-bit id that is locally unique to the Exporting Process
}

type TemplateHeader struct {
	TemplateID uint16
	FieldCount uint16
}

type TemplateRecords struct {
	TemplateID      uint16
	FieldCount      uint16
	FieldSpecifiers []TemplateFieldSpecifier
}

type TemplateFieldSpecifier struct {
	ElementID    uint16
	Length       uint16
	EnterpriseNo uint32
}

type OptsTemplateHeader struct {
	TemplateID      uint16
	FieldCount      uint16
	ScopeFieldCount uint16
}

type OptsTemplateRecords struct {
	TemplateID      uint16
	FieldCount      uint16
	ScopeFieldCount uint16
}

type Message struct {
	Header   MessageHeader
	DataSets []DataSet
}

type DataSet struct {
}

type Session struct {
	buff *sync.Pool
}

type SetHeader struct {
	SetID  uint16
	Length uint16
}

var (
	errInvalidVersion = errors.New("invalid ipfix version")
)

func NewDecoder(raddr net.IP, b []byte) *IPFIXDecoder {
	return &IPFIXDecoder{raddr, NewReader(b)}
}

func (d *IPFIXDecoder) Decode(mem MemCache) (*Message, error) {
	var (
		msg = new(Message)
		err error
	)

	// IPFIX Message Header decoding
	if err = msg.Header.unmarshal(d.reader); err != nil {
		return nil, err
	}
	// IPFIX Message Header validation
	if err = msg.Header.validate(); err != nil {
		return nil, err
	}

	for d.reader.Len() > 0 {

		setHeader := new(SetHeader)
		setHeader.unmarshal(d.reader)

		if setHeader.Length < 4 {
			return nil, io.ErrUnexpectedEOF
		}

		switch {
		case setHeader.SetID == 2:
			// Template set
			tr := TemplateRecords{}
			tr.unmarshal(d.reader)
			mem.insert(tr.TemplateID, d.raddr, tr)
		case setHeader.SetID == 3:
			// Option set
			tr := OptsTemplateRecords{}
			tr.unmarshal(d.reader)
			mem.insert(tr.TemplateID, d.raddr, tr)
		case setHeader.SetID >= 4 && setHeader.SetID <= 255:
			// Reserved
		default:
			// data
			mem.retrieve(setHeader.SetID, d.raddr)
		}

		break

	}

	return msg, nil
}

// RFC 7011 - part 3.1. Message Header Format
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Version Number          |            Length             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           Export Time                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       Sequence Number                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Observation Domain ID                      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (h *MessageHeader) unmarshal(r *Reader) error {
	var err error

	if h.Version, err = r.Uint16(); err != nil {
		return err
	}

	if h.Length, err = r.Uint16(); err != nil {
		return err
	}

	if h.ExportTime, err = r.Uint32(); err != nil {
		return err
	}

	if h.SequenceNo, err = r.Uint32(); err != nil {
		return err
	}

	if h.DomainID, err = r.Uint32(); err != nil {
		return err
	}

	return nil
}

func (h *MessageHeader) validate() error {
	if h.Version != 0x000a {
		return errInvalidVersion
	}

	// TODO: needs more validation

	return nil
}

// RFC 7011 - part 3.3.2 Set Header Format
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID               |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (h *SetHeader) unmarshal(r *Reader) error {
	var err error

	if h.SetID, err = r.Uint16(); err != nil {
		return err
	}

	if h.Length, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

// RFC 7011
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       Set ID = (2 or 3)       |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (t *TemplateHeader) unmarshal(r *Reader) error {
	var err error

	if t.TemplateID, err = r.Uint16(); err != nil {
		return err
	}

	if t.FieldCount, err = r.Uint16(); err != nil {
		return err
	}

	return nil

}

// RFC 7011 3.4.2.2.  Options Template Record Format
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID = 3           |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |         Field Count = N + M   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope Field Count = N     |0|  Scope 1 Infor. Element id. |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (t *OptsTemplateHeader) unmarshal(r *Reader) error {
	var err error

	if t.TemplateID, err = r.Uint16(); err != nil {
		return err
	}

	if t.FieldCount, err = r.Uint16(); err != nil {
		return err
	}

	if t.ScopeFieldCount, err = r.Uint16(); err != nil {
		return err
	}

	return nil

}

// RFC 7011
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |E|  Information Element ident. |        Field Length           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                      Enterprise Number                        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (f *TemplateFieldSpecifier) unmarshal(r *Reader) error {
	var err error

	if f.ElementID, err = r.Uint16(); err != nil {
		return err
	}

	if f.Length, err = r.Uint16(); err != nil {
		return err
	}

	if f.ElementID > 0x8000 {
		f.ElementID = f.ElementID & 0x7fff
		if f.EnterpriseNo, err = r.Uint32(); err != nil {
			return err
		}
	}

	return nil
}

func (tr *TemplateRecords) unmarshal(r *Reader) {
	var (
		th = TemplateHeader{}
		tf = TemplateFieldSpecifier{}
	)

	th.unmarshal(r)
	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount

	for i := th.FieldCount; i > 0; i-- {
		tf.unmarshal(r)
		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}
}

func (tr *OptsTemplateRecords) unmarshal(r *Reader) {
	var (
		th = OptsTemplateHeader{}
	)

	th.unmarshal(r)
	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount
	tr.ScopeFieldCount = th.ScopeFieldCount

	// TODO
}
