//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    reader.go
//: details: decodes a variable from buffer
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

// Package reader decodes a variable from buffer
package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Reader represents the data bytes for reading
type Reader struct {
	data  []byte
	count int
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
	r.advance(1)

	return d, nil
}

// Uint16 reads two bytes as big-endian
func (r *Reader) Uint16() (uint16, error) {
	if len(r.data) < 2 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint16(r.data)
	r.advance(2)

	return d, nil
}

// Uint32 reads four bytes as big-endian
func (r *Reader) Uint32() (uint32, error) {
	if len(r.data) < 4 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint32(r.data)
	r.advance(4)

	return d, nil
}

// Uint64 reads eight bytes as big-endian
func (r *Reader) Uint64() (uint64, error) {
	if len(r.data) < 8 {
		return 0, errReader
	}

	d := binary.BigEndian.Uint64(r.data)
	r.advance(8)

	return d, nil
}

// Read reads n bytes and returns it
func (r *Reader) Read(n int) ([]byte, error) {
	if n == 65535 {
		return r.readString(n)
	}
	if len(r.data) < n {
		return []byte{}, errReader
	}

	d := r.data[:n]
	r.advance(n)

	return d, nil
}

func (r *Reader) readString(n int) ([]byte, error) {
	if len(r.data) < 1 {
		return []byte(""), errReader
	}
	penlen := int(r.data[0])
	r.advance(1)
	if penlen == 255 {
		x, err := r.Uint16()
		if err != nil {
			return []byte(""), fmt.Errorf("not enought data available to read the length of the variable length element: %d", len(r.data))
		}
		penlen = int(x)
	}
	if len(r.data) < penlen || penlen > n {
		return []byte(""), fmt.Errorf("not enough data available to read %d length of the variable length element, available: %d", penlen, len(r.data))
	}
	if penlen != 0 {
		d := r.data[:penlen]
		r.advance(penlen)
		return d, nil
	}
	return []byte(""), nil
}

// PeekUint16 peeks the next two bytes interpreted as big-endian two-byte integer
func (r *Reader) PeekUint16() (res uint16, err error) {
	var b []byte
	if b, err = r.Peek(2); err == nil {
		res = binary.BigEndian.Uint16(b)
	}
	return
}

// Peek returns the next n bytes in the reader without advancing in the stream
func (r *Reader) Peek(n int) ([]byte, error) {
	if len(r.data) < n {
		return []byte{}, errReader
	}
	return r.data[:n], nil
}

// Len returns the current length of the reader's data
func (r *Reader) Len() int {
	return len(r.data)
}

func (r *Reader) advance(num int) {
	r.data = r.data[num:]
	r.count += num
}

// ReadCount returns the number of bytes that have been read from this Reader in total
func (r *Reader) ReadCount() int {
	return r.count
}
