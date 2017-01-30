package ipfix

import (
	"encoding/binary"
	"io"
)

type IPFIXDecoder struct {
	reader io.Reader
}

type MessageHeader struct {
	Version    uint16
	Length     uint16
	ExportTime uint32
	SequenceNo uint32
	DomainID   uint32
}

func NewDecoder(r io.Reader) IPFIXDecoder {
	return IFDecoder{
		reader: r,
	}
}

func (d *IPFIXDecoder) Decode() error {
	var (
		h   MessageHeader
		err error
	)

	if err = h.unmarshal(d.reader); err != nil {
		return err
	}

	return nil
}

func (h *MessageHeader) unmarshal(d io.Reader) error {
	var err error

	if err = read(d, &h.Version); err != nil {
		return err
	}

	if err = read(d, &h.Length); err != nil {
		return err
	}

	if err = read(d, &h.ExportTime); err != nil {
		return err
	}

	if err = read(d, &h.SequenceNo); err != nil {
		return err
	}

	if err = read(d, &h.DomainID); err != nil {
		return err
	}

	return nil
}

func read(r io.Reader, v interface{}) error {
	return binary.Read(r, binary.BigEndian, v)
}
