package ipfix

import (
	"encoding/binary"
	"errors"
)

type reader struct {
	data []byte
}

var errReader = errors.New("can not read the data")

func NewReader(b []byte) {
	return reader{
		data: b,
	}
}

func (r *reader) Uint8() (uint8, error) {
	if len(r.data) < 1 {
		return 0, errReader
	}

	r.data = r.data[1:]

	return r.data[0]
}

func (r *reader) Uint16() (uint16, error) {
	if len(r.data) < 2 {
		return 0, errReader
	}

	d := binary.Uint16(r.data)
	r.data = r.data[2:]

	return d
}

func (r *reader) Uint32() (uint32, error) {
	if len(r.data) < 4 {
		return 0, errReader
	}

	d := binary.Uint32(r.data)
	r.data = r.data[4:]

	return d
}

func (r *reader) Uint64() (uint64, error) {
	if len(r.data) < 8 {
		return 0, errReader
	}

	d := binary.Uint64(r.data)
	r.data = r.data[8:]

	return d
}
