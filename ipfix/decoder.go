package ipfix

import (
	"errors"
	"io"
	"sync"
)

type IPFIXDecoder struct {
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

type Message struct {
	Header       MessageHeader
	TemplateSets []TemplateSet
	DataSets     []DataSet
}

type TemplateSet struct {
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

func NewDecoder(r io.Reader) (*IPFIXDecoder, error) {
	data := make([]byte, 1500)
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &IPFIXDecoder{NewReader(data[:n])}, nil
}

func (d *IPFIXDecoder) Decode() error {
	var (
		msg Message
		err error
	)

	// IPFIX Message Header decoding
	if err = msg.Header.unmarshal(d.reader); err != nil {
		return err
	}
	// IPFIX Message Header validation
	if err = msg.Header.validate(); err != nil {
		return err
	}

	for d.reader.Len() > 0 {

		setHeader := new(SetHeader)
		setHeader.unmarshal(d.reader)

		if setHeader.Length < 4 {
			return io.ErrUnexpectedEOF
		}

		switch {
		case setHeader.SetID == 2:
			// Template set
			ts := new(TemplateSet)
			ts.unmarshal(d.reader)
		case setHeader.SetID == 3:
			// Option set
		case setHeader.SetID >= 4 && setHeader.SetID <= 255:
			// Reserved
		default:
			// data
		}

		break

	}

	return nil
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

func (t *TemplateSet) unmarshal(r *Reader) error {
	th := new(TemplateHeader)
	th.unmarshal(r)

	println(th.TemplateID, th.FieldCount)

	return nil
}
