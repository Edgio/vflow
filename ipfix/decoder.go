//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder.go
//: details: decodes IPFIX packets
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
	"fmt"
	"io"
	"net"

	"github.com/VerizonDigital/vflow/reader"
)

// Decoder represents IPFIX payload and remote address
type Decoder struct {
	raddr  net.IP
	reader *reader.Reader
}

// MessageHeader represents IPFIX message header
type MessageHeader struct {
	Version    uint16 // Version of IPFIX to which this Message conforms
	Length     uint16 // Total length of the IPFIX Message, measured in octets
	ExportTime uint32 // Time at which the IPFIX Message Header leaves the Exporter
	SequenceNo uint32 // Incremental sequence counter modulo 2^32
	DomainID   uint32 // A 32-bit id that is locally unique to the Exporting Process
}

// TemplateHeader represents template fields
type TemplateHeader struct {
	TemplateID      uint16
	FieldCount      uint16
	ScopeFieldCount uint16
}

// TemplateRecord represents template records
type TemplateRecord struct {
	TemplateID           uint16
	FieldCount           uint16
	FieldSpecifiers      []TemplateFieldSpecifier
	ScopeFieldCount      uint16
	ScopeFieldSpecifiers []TemplateFieldSpecifier
}

// TemplateFieldSpecifier represents field properties
type TemplateFieldSpecifier struct {
	ElementID    uint16
	Length       uint16
	EnterpriseNo uint32
}

// Message represents IPFIX decoded data
type Message struct {
	AgentID  string
	Header   MessageHeader
	DataSets [][]DecodedField
}

// DecodedField represents a decoded field
type DecodedField struct {
	ID    uint16
	Value interface{}
}

// SetHeader represents set header fields
type SetHeader struct {
	SetID  uint16
	Length uint16
}

var rpcChan = make(chan RPCRequest, 1)

// NewDecoder constructs a decoder
func NewDecoder(raddr net.IP, b []byte) *Decoder {
	return &Decoder{raddr, reader.NewReader(b)}
}

// Decode decodes the IPFIX raw data
func (d *Decoder) Decode(mem MemCache) (*Message, error) {
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

	// Add source IP address as Agent ID
	msg.AgentID = d.raddr.String()

	for d.reader.Len() > 4 {

		setHeader := new(SetHeader)
		setHeader.unmarshal(d.reader)

		if setHeader.Length < 4 {
			return nil, io.ErrUnexpectedEOF
		}

		switch {
		case setHeader.SetID == 2:
			// Template set
			tr := TemplateRecord{}
			tr.unmarshal(d.reader)
			mem.insert(tr.TemplateID, d.raddr, tr)
		case setHeader.SetID == 3:
			// Option set
			tr := TemplateRecord{}
			tr.unmarshalOpts(d.reader)
			mem.insert(tr.TemplateID, d.raddr, tr)
		case setHeader.SetID >= 4 && setHeader.SetID <= 255:
			// Reserved
		default:
			// data
			tr, ok := mem.retrieve(setHeader.SetID, d.raddr)
			if !ok {
				select {
				case rpcChan <- RPCRequest{
					ID: setHeader.SetID,
					IP: d.raddr,
				}:
				default:
				}
				return msg, fmt.Errorf("%s unknown ipfix template id# %d",
					d.raddr.String(),
					setHeader.SetID,
				)
			}

			// data records
			for d.reader.Len() > 2 {
				data, err := decodeData(d.reader, tr)
				if err != nil {
					return msg, err
				}

				msg.DataSets = append(msg.DataSets, data)
			}
		}
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

func (h *MessageHeader) unmarshal(r *reader.Reader) error {
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
		return fmt.Errorf("invalid ipfix version (%d)", h.Version)
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

func (h *SetHeader) unmarshal(r *reader.Reader) error {
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

func (t *TemplateHeader) unmarshal(r *reader.Reader) error {
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

func (t *TemplateHeader) unmarshalOpts(r *reader.Reader) error {
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

func (f *TemplateFieldSpecifier) unmarshal(r *reader.Reader) error {
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

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |          Set ID = 2           |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID = 256        |         Field Count = N       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |1| Information Element id. 1.1 |        Field Length 1.1       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                    Enterprise Number  1.1                     |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0| Information Element id. 1.2 |        Field Length 1.2       |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (tr *TemplateRecord) unmarshal(r *reader.Reader) {
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

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |          Set ID = 3           |          Length               |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |         Template ID = X       |         Field Count = N + M   |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |     Scope Field Count = N     |0|  Scope 1 Infor. Element id. |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |     Scope 1 Field Length      |0|  Scope 2 Infor. Element id. |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |     Scope 2 Field Length      |             ...               |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |            ...                |1|  Scope N Infor. Element id. |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |     Scope N Field Length      |   Scope N Enterprise Number  ...
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// ...  Scope N Enterprise Number   |1| Option 1 Infor. Element id. |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |    Option 1 Field Length      |  Option 1 Enterprise Number  ...
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// ... Option 1 Enterprise Number   |              ...              |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |             ...               |0| Option M Infor. Element id. |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |     Option M Field Length     |      Padding (optional)       |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (tr *TemplateRecord) unmarshalOpts(r *reader.Reader) {
	var (
		th = TemplateHeader{}
		tf = TemplateFieldSpecifier{}
	)

	th.unmarshalOpts(r)
	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount
	tr.ScopeFieldCount = th.ScopeFieldCount

	for i := th.ScopeFieldCount; i > 0; i-- {
		tf.unmarshal(r)
		tr.ScopeFieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}

	for i := th.FieldCount - th.ScopeFieldCount; i > 0; i-- {
		tf.unmarshal(r)
		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}
}

func decodeData(r *reader.Reader, tr TemplateRecord) ([]DecodedField, error) {
	var (
		fields []DecodedField
		err    error
		b      []byte
	)

	for i := 0; i < len(tr.FieldSpecifiers); i++ {
		b, err = r.Read(int(tr.FieldSpecifiers[i].Length))
		if err != nil {
			return nil, err
		}

		m, ok := InfoModel[ElementKey{
			tr.FieldSpecifiers[i].EnterpriseNo,
			tr.FieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, fmt.Errorf("IPFIX element key not exist")
		}

		fields = append(fields, DecodedField{
			ID:    m.FieldID,
			Value: Interpret(&b, m.Type),
		})
	}

	for i := 0; i < len(tr.ScopeFieldSpecifiers); i++ {
		b, err = r.Read(int(tr.ScopeFieldSpecifiers[i].Length))
		if err != nil {
			return nil, err
		}

		m, ok := InfoModel[ElementKey{
			tr.ScopeFieldSpecifiers[i].EnterpriseNo,
			tr.ScopeFieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, fmt.Errorf("IPFIX element key not exist (scope)")
		}

		fields = append(fields, DecodedField{
			ID:    m.FieldID,
			Value: Interpret(&b, m.Type),
		})
	}

	return fields, nil
}
