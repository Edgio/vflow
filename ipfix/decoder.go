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
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/EdgeCast/vflow/reader"
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
	ID           uint16
	Value        interface{}
	EnterpriseNo uint32
}

// SetHeader represents set header fields
type SetHeader struct {
	SetID  uint16
	Length uint16
}

type nonfatalError struct {
	error
}

var rpcChan = make(chan RPCRequest, 1)

// NewDecoder constructs a decoder
func NewDecoder(raddr net.IP, b []byte) *Decoder {
	return &Decoder{raddr, reader.NewReader(b)}
}

// Decode decodes the IPFIX raw data
func (d *Decoder) Decode(mem MemCache) (*Message, error) {
	var msg = new(Message)

	// IPFIX Message Header decoding
	if err := msg.Header.unmarshal(d.reader); err != nil {
		return nil, err
	}
	// IPFIX Message Header validation
	if err := msg.Header.validate(); err != nil {
		return nil, err
	}

	// Add source IP address as Agent ID
	msg.AgentID = d.raddr.String()

	// In case there are multiple non-fatal errors, collect them and report all of them.
	// The rest of the received sets will still be interpreted, until a fatal error is encountered.
	// A non-fatal error is for example an illegal data record or unknown template id.
	var decodeErrors []error
	for d.reader.Len() > 4 {
		if err := d.decodeSet(mem, msg); err != nil {
			switch err.(type) {
			case nonfatalError:
				decodeErrors = append(decodeErrors, err)
			default:
				return nil, err
			}
		}
	}

	return msg, combineErrors(decodeErrors...)
}

// RFC 7011 - part 3.B IPFIX Message Format
// +----------------------------------------------------+
// | Message Header                                     |
// +----------------------------------------------------+
// | Set                                                |
// +----------------------------------------------------+
// | Set                                                |
// +----------------------------------------------------+
//   ...
// +----------------------------------------------------+
// | Set                                                |
// +----------------------------------------------------+

func (d *Decoder) decodeSet(mem MemCache, msg *Message) error {
	startCount := d.reader.ReadCount()

	setHeader := new(SetHeader)
	if err := setHeader.unmarshal(d.reader); err != nil {
		return err
	}
	if setHeader.Length < 4 {
		return io.ErrUnexpectedEOF
	}

	var tr TemplateRecord
	var err error
	// This check is somewhat redundant with the switch-clause below, but the retrieve() operation should not be executed inside the loop.
	if setHeader.SetID > 255 {
		var ok bool
		if tr, ok = mem.retrieve(setHeader.SetID, d.raddr, msg.Header.DomainID); !ok {
			select {
			case rpcChan <- RPCRequest{
				ID: setHeader.SetID,
				IP: d.raddr,
				SrcID: msg.Header.DomainID,
			}:
			default:
			}
			err = nonfatalError{fmt.Errorf("%s unknown ipfix template id# %d with domain ID %d",
				d.raddr.String(),
				setHeader.SetID,
				msg.Header.DomainID,
			)}
		}
	}

	// the next set should be greater than 4 bytes otherwise that's padding
	for err == nil && setHeader.Length > uint16(d.reader.ReadCount()-startCount) && d.reader.Len() > 4 && setHeader.Length-uint16(d.reader.ReadCount()-startCount) > 4 {
		if setID := setHeader.SetID; setID == 2 || setID == 3 {
			// Template record or template option record

			// Check if only padding is left in this set. A template id of zero indicates padding bytes, which MUST be zero.
			if templateID, err := d.reader.PeekUint16(); err == nil && templateID == 0 {
				break
			}

			tr := TemplateRecord{}
			if setID == 2 {
				err = tr.unmarshal(d.reader)
			} else {
				err = tr.unmarshalOpts(d.reader)
			}
			if err == nil {
				mem.insert(tr.TemplateID, d.raddr, tr, msg.Header.DomainID)
			}
		} else if setID >= 4 && setID <= 255 {
			// Reserved set, do not read any records
			break
		} else if setID == 0 {
			// Invalid set
			return fmt.Errorf("failed to decodeSet / invalid setID")
		} else {
			// Data set
			var data []DecodedField
			if data, err = d.decodeData(tr); err == nil {
				msg.DataSets = append(msg.DataSets, data)
			} else {
				switch err.(type) {
				case nonfatalError:
				default:
					return err
				}
			}
		}
	}

	// Skip the rest of the set in order to properly continue with the next set
	// This is necessary if the set is padded, has a reserved set ID, or a nonfatal error occurred
	leftoverBytes := setHeader.Length - uint16(d.reader.ReadCount()-startCount)
	if leftoverBytes > 0 {
		if _, skipErr := d.reader.Read(int(leftoverBytes)); skipErr != nil {
			err = skipErr
		}
	}
	return err
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
	} else {
		f.EnterpriseNo = 0
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

func (tr *TemplateRecord) unmarshal(r *reader.Reader) error {
	var (
		th = TemplateHeader{}
		tf = TemplateFieldSpecifier{}
	)

	if err := th.unmarshal(r); err != nil {
		return err
	}
	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount

	for i := th.FieldCount; i > 0; i-- {
		if err := tf.unmarshal(r); err != nil {
			return err
		}
		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}
	return nil
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

func (tr *TemplateRecord) unmarshalOpts(r *reader.Reader) error {
	var (
		th = TemplateHeader{}
		tf = TemplateFieldSpecifier{}
	)

	if err := th.unmarshalOpts(r); err != nil {
		return err
	}
	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount
	tr.ScopeFieldCount = th.ScopeFieldCount

	for i := th.ScopeFieldCount; i > 0; i-- {
		if err := tf.unmarshal(r); err != nil {
			return err
		}
		tr.ScopeFieldSpecifiers = append(tr.ScopeFieldSpecifiers, tf)
	}

	for i := th.FieldCount - th.ScopeFieldCount; i > 0; i-- {
		if err := tf.unmarshal(r); err != nil {
			return err
		}
		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}
	return nil
}

func (d *Decoder) getDataLength(fieldSpecifierLen uint16, t FieldType) (uint16, error) {
	var (
		err        error
		readLength uint16
	)

	r := d.reader

	if (t == String || t == OctetArray) && (fieldSpecifierLen == 65535) {
		var len8 uint8
		if len8, err = r.Uint8(); err != nil {
			return 0, err
		} else if len8 == 255 {
			if readLength, err = r.Uint16(); err != nil {
				return 0, err
			}
		} else {
			readLength = uint16(len8)
		}
	} else {
		readLength = fieldSpecifierLen
	}

	return readLength, nil
}

func (d *Decoder) decodeData(tr TemplateRecord) ([]DecodedField, error) {
	var (
		fields     []DecodedField
		err        error
		b          []byte
		readLength uint16
	)

	r := d.reader

	for i := 0; i < len(tr.ScopeFieldSpecifiers); i++ {
		m, ok := InfoModel[ElementKey{
			tr.ScopeFieldSpecifiers[i].EnterpriseNo,
			tr.ScopeFieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, nonfatalError{fmt.Errorf("IPFIX element key (%d) not exist (scope)",
				tr.ScopeFieldSpecifiers[i].ElementID)}
		}

		if readLength, err = d.getDataLength(tr.ScopeFieldSpecifiers[i].Length, m.Type); err != nil {
			return nil, err
		}

		if b, err = r.Read(int(readLength)); err != nil {
			return nil, err
		}

		fields = append(fields, DecodedField{
			ID:           m.FieldID,
			Value:        Interpret(&b, m.Type),
			EnterpriseNo: tr.ScopeFieldSpecifiers[i].EnterpriseNo,
		})
	}

	for i := 0; i < len(tr.FieldSpecifiers); i++ {
		m, ok := InfoModel[ElementKey{
			tr.FieldSpecifiers[i].EnterpriseNo,
			tr.FieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, nonfatalError{fmt.Errorf("IPFIX element key (%d) not exist",
				tr.FieldSpecifiers[i].ElementID)}
		}

		if readLength, err = d.getDataLength(tr.FieldSpecifiers[i].Length, m.Type); err != nil {
			return nil, err
		}

		if b, err = r.Read(int(readLength)); err != nil {
			return nil, err
		}

		fields = append(fields, DecodedField{
			ID:           m.FieldID,
			Value:        Interpret(&b, m.Type),
			EnterpriseNo: tr.FieldSpecifiers[i].EnterpriseNo,
		})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("failed to decodeData")
	}

	return fields, nil
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
