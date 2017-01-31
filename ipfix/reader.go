package ipfix

import (
	"encoding/binary"
	"errors"
)

type Reader struct {
	data []byte
}

var errReader = errors.New("can not read the data")

func NewReader(b []byte) *Reader {
	return &Reader{
		data: b,
	}
}

func (r *Reader) Uint8() (uint8, error) {
	if len(r.data) < 1 {
		return 0, errReader
	}

	r.data = r.data[1:]

	return r.data[0], nil
}

func (r *Reader) Uint16() (uint16, error) {
	if len(r.data) < 2 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]

	return d, nil
}

func (r *Reader) Uint32() (uint32, error) {
	if len(r.data) < 4 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]

	return d, nil
}

func (r *Reader) Uint64() (uint64, error) {
	if len(r.data) < 8 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint64(r.data)
	r.data = r.data[8:]

	return d, nil
}

func (r *Reader) Read(n int) ([]byte, error) {
	if len(r.data) < n {
		return []byte{}, errReader
	}

	d := r.data[:n]
	r.data = r.data[n:]

	return d, nil
}

func (r *Reader) Len() int {
	return len(r.data)
}
