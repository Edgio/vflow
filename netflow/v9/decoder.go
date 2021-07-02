//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder.go
//: details: decodes netflow version 9 packets
//: author:  Mehrdad Arshad Rad
//: date:    04/10/2017
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

package netflow9

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/EdgeCast/vflow/ipfix"
	"github.com/EdgeCast/vflow/reader"
)

type nonfatalError error

// PacketHeader represents Netflow v9  packet header
type PacketHeader struct {
	Version   uint16 // Version of Flow Record format exported in this packet
	Count     uint16 // The total number of records in the Export Packet
	SysUpTime uint32 // Time in milliseconds since this device was first booted
	UNIXSecs  uint32 // Time in seconds since 0000 UTC 197
	SeqNum    uint32 // Incremental sequence counter of all Export Packets
	SrcID     uint32 // A 32-bit value that identifies the Exporter
}

// SetHeader represents netflow v9 data flowset id and length
type SetHeader struct {
	FlowSetID uint16 // FlowSet ID value 0:: template, 1:: options template, 255< :: data
	Length    uint16 // Total length of this FlowSet
}

// TemplateHeader represents netflow v9 data template id and field count
type TemplateHeader struct {
	TemplateID     uint16 // Template ID
	FieldCount     uint16 // Number of fields in this Template Record
	OptionLen      uint16 // The length in bytes of any Scope field definition (Option)
	OptionScopeLen uint16 // The length in bytes of any options field definitions (Option)
}

// TemplateFieldSpecifier represents field properties
type TemplateFieldSpecifier struct {
	ElementID uint16
	Length    uint16
}

// TemplateRecord represents template fields
type TemplateRecord struct {
	TemplateID           uint16
	FieldCount           uint16
	FieldSpecifiers      []TemplateFieldSpecifier
	ScopeFieldCount      uint16
	ScopeFieldSpecifiers []TemplateFieldSpecifier
}

// DecodedField represents a decoded field
type DecodedField struct {
	ID    uint16
	Value interface{}
}

// Decoder represents Netflow payload and remote address
type Decoder struct {
	raddr  net.IP
	reader *reader.Reader
}

// Message represents Netflow decoded data
type Message struct {
	AgentID  string
	Header   PacketHeader
	DataSets [][]DecodedField
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
//   |                       Sequence Number                         |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//   |                        Source ID                              |
//   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (h *PacketHeader) unmarshal(r *reader.Reader) error {
	var err error

	if h.Version, err = r.Uint16(); err != nil {
		return err
	}

	if h.Count, err = r.Uint16(); err != nil {
		return err
	}

	if h.SysUpTime, err = r.Uint32(); err != nil {
		return err
	}

	if h.UNIXSecs, err = r.Uint32(); err != nil {
		return err
	}

	if h.SeqNum, err = r.Uint32(); err != nil {
		return err
	}

	if h.SrcID, err = r.Uint32(); err != nil {
		return err
	}

	return nil
}

func (h *PacketHeader) validate() error {
	if h.Version != 9 {
		return fmt.Errorf("invalid netflow version (%d)", h.Version)
	}

	// TODO: needs more validation

	return nil
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        FlowSet ID             |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (h *SetHeader) unmarshal(r *reader.Reader) error {
	var err error

	if h.FlowSetID, err = r.Uint16(); err != nil {
		return err
	}

	if h.Length, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Template ID            |         Field Count           |
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

func (t *TemplateHeader) unmarshalOpts(r *reader.Reader) error {
	var err error

	if t.TemplateID, err = r.Uint16(); err != nil {
		return err
	}

	if t.OptionScopeLen, err = r.Uint16(); err != nil {
		return err
	}

	if t.OptionLen, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type             |         Field Length          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (f *TemplateFieldSpecifier) unmarshal(r *reader.Reader) error {
	var err error

	if f.ElementID, err = r.Uint16(); err != nil {
		return err
	}

	if f.Length, err = r.Uint16(); err != nil {
		return err
	}

	return nil
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      Template ID 256          |         Field Count           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 1           |         Field Length 1        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type 2           |         Field Length 2        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |             ...               |              ...              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Field Type N           |         Field Length N        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (tr *TemplateRecord) unmarshal(r *reader.Reader) error {
	var (
		th  = TemplateHeader{}
		tf  = TemplateFieldSpecifier{}
		err error
	)

	if err = th.unmarshal(r); err != nil {
		return err
	}

	tr.TemplateID = th.TemplateID
	tr.FieldCount = th.FieldCount

	for i := th.FieldCount; i > 0; i-- {
		if err = tf.unmarshal(r); err != nil {
			return err
		}
		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}

	return nil
}

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |       FlowSet ID = 1          |          Length               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |         Template ID           |      Option Scope Length      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |        Option Length          |       Scope 1 Field Type      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope 1 Field Length      |               ...             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Scope N Field Length      |      Option 1 Field Type      |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Option 1 Field Length     |             ...               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Option M Field Length     |           Padding             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func (tr *TemplateRecord) unmarshalOpts(r *reader.Reader) error {
	var (
		th  = TemplateHeader{}
		tf  = TemplateFieldSpecifier{}
		err error
	)

	if err = th.unmarshalOpts(r); err != nil {
		return err
	}

	tr.TemplateID = th.TemplateID

	for i := th.OptionScopeLen / 4; i > 0; i-- {
		if err = tf.unmarshal(r); err != nil {
			return err
		}

		tr.ScopeFieldSpecifiers = append(tr.ScopeFieldSpecifiers, tf)
	}

	for i := th.OptionLen / 4; i > 0; i-- {
		if err = tf.unmarshal(r); err != nil {
			return err
		}

		tr.FieldSpecifiers = append(tr.FieldSpecifiers, tf)
	}

	return nil
}

func (d *Decoder) decodeData(tr TemplateRecord) ([]DecodedField, error) {
	var (
		fields []DecodedField
		err    error
		b      []byte
	)

	r := d.reader

	for i := 0; i < len(tr.ScopeFieldSpecifiers); i++ {
		b, err = r.Read(int(tr.ScopeFieldSpecifiers[i].Length))
		if err != nil {
			return nil, err
		}

		m, ok := ipfix.InfoModel[ipfix.ElementKey{
			0,
			tr.ScopeFieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, nonfatalError(fmt.Errorf("Netflow element key (%d) not exist (scope)",
				tr.ScopeFieldSpecifiers[i].ElementID))
		}

		fields = append(fields, DecodedField{
			ID:    m.FieldID,
			Value: ipfix.Interpret(&b, m.Type),
		})
	}

	for i := 0; i < len(tr.FieldSpecifiers); i++ {
		b, err = r.Read(int(tr.FieldSpecifiers[i].Length))
		if err != nil {
			return nil, err
		}

		m, ok := ipfix.InfoModel[ipfix.ElementKey{
			0,
			tr.FieldSpecifiers[i].ElementID,
		}]

		if !ok {
			return nil, nonfatalError(fmt.Errorf("Netflow element key (%d) not exist",
				tr.FieldSpecifiers[i].ElementID))
		}

		fields = append(fields, DecodedField{
			ID:    m.FieldID,
			Value: ipfix.Interpret(&b, m.Type),
		})
	}

	return fields, nil
}

// NewDecoder constructs a decoder
func NewDecoder(raddr net.IP, b []byte) *Decoder {
	return &Decoder{raddr, reader.NewReader(b)}
}

// Decode decodes the flow records
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
	if setHeader.FlowSetID > 255 {
		var ok bool
		tr, ok = mem.retrieve(setHeader.FlowSetID, d.raddr)
		if !ok {
			err = nonfatalError(fmt.Errorf("%s unknown netflow template id# %d",
				d.raddr.String(),
				setHeader.FlowSetID,
			))
		}
	}

	// the next set should be greater than 4 bytes otherwise that's padding
	for err == nil && (int(setHeader.Length)-(d.reader.ReadCount()-startCount) > 4) && d.reader.Len() > 4 {
		if setId := setHeader.FlowSetID; setId == 0 || setId == 1 {
			// Template record or template option record
			tr := TemplateRecord{}
			if setId == 0 {
				err = tr.unmarshal(d.reader)
			} else {
				err = tr.unmarshalOpts(d.reader)
			}
			if err == nil {
				mem.insert(tr.TemplateID, d.raddr, tr)
			}
		} else if setId >= 4 && setId <= 255 {
			// Reserved set, do not read any records
			break
		} else {
			// Data set
			var data []DecodedField
			data, err = d.decodeData(tr)
			if err == nil {
				msg.DataSets = append(msg.DataSets, data)
			}
		}
	}

	// Skip the rest of the set in order to properly continue with the next set
	// This is necessary if the set is padded, has a reserved set ID, or a nonfatal error occurred
	leftoverBytes := int(setHeader.Length) - (d.reader.ReadCount() - startCount)
	if leftoverBytes > 0 {
		_, skipErr := d.reader.Read(int(leftoverBytes))
		if skipErr != nil {
			err = skipErr
		}
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
