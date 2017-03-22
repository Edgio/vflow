// Package ipfix decodes IPFIX packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    reader.go
//: details: bytes reader
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
	"encoding/binary"
	"errors"
)

// Reader represents the data bytes for reading
type Reader struct {
	data []byte
}

var errReader = errors.New("can not read the data")

// NewReader constructs a reader
func NewReader(b []byte) *Reader {
	return &Reader{
		data: b,
	}
}

// Uint8 reads a byte
func (r *Reader) Uint8() (uint8, error) {
	if len(r.data) < 1 {
		return 0, errReader
	}

	d := r.data[0]
	r.data = r.data[1:]

	return d, nil
}

// Uint16 reads two bytes as big-endian
func (r *Reader) Uint16() (uint16, error) {
	if len(r.data) < 2 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]

	return d, nil
}

// Uint32 reads four bytes as big-endian
func (r *Reader) Uint32() (uint32, error) {
	if len(r.data) < 4 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]

	return d, nil
}

// Uint64 reads eight bytes as big-endian
func (r *Reader) Uint64() (uint64, error) {
	if len(r.data) < 8 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint64(r.data)
	r.data = r.data[8:]

	return d, nil
}

// Read reads n bytes and returns it
func (r *Reader) Read(n int) ([]byte, error) {
	if len(r.data) < n {
		return []byte{}, errReader
	}

	d := r.data[:n]
	r.data = r.data[n:]

	return d, nil
}

// Len returns the current length of the reader's data
func (r *Reader) Len() int {
	return len(r.data)
}
