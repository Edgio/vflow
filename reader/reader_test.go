//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    reader_test.go
//: details: unit testing for reader.go
//: author:  Mehrdad Arshad Rad
//: date:    03/22/2017
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

package reader

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUint8(t *testing.T) {
	b := []byte{0x05, 0x11}

	r := NewReader(b)
	i, err := r.Uint8()
	if err != nil {
		t.Error("unexpected error happened, got", err)
	}

	if i != 5 {
		t.Error("expect read 5, got", i)
	}
}

func TestUint16(t *testing.T) {
	b := []byte{0x05, 0x11}

	r := NewReader(b)
	i, err := r.Uint16()
	if err != nil {
		t.Error("unexpected error happened, got", err)
	}

	if i != 1297 {
		t.Error("expect read 1297, got", i)
	}
}

func TestUint32(t *testing.T) {
	b := []byte{0x05, 0x11, 0x01, 0x16}

	r := NewReader(b)
	i, err := r.Uint32()
	if err != nil {
		t.Error("unexpected error happened, got", err)
	}

	if i != 85000470 {
		t.Error("expect read 85000470, got", i)
	}
}

func TestUint64(t *testing.T) {
	b := []byte{0x05, 0x11, 0x01, 0x16, 0x05, 0x01, 0x21, 0x26}

	r := NewReader(b)
	i, err := r.Uint64()
	if err != nil {
		t.Error("unexpected error happened, got", err)
	}

	if i != 365074238878589222 {
		t.Error("expect read 365074238878589222, got", i)
	}
}

func TestReadN(t *testing.T) {
	b := []byte{0x05, 0x11, 0x01, 0x16}

	r := NewReader(b)
	i, err := r.Read(2)
	if err != nil {
		t.Error("unexpected error happened, got", err)
	}

	if !reflect.DeepEqual(i, []byte{0x05, 0x11}) {
		t.Error("expect read [5 17], got", i)
	}
}

func TestReadCount(t *testing.T) {
	b := make([]byte, 18)
	for i := range b {
		b[i] = byte(i)
	}
	r := NewReader(b)
	check := func(expected int) {
		count := r.ReadCount()
		if count != expected {
			t.Error("Unexpected ReadCount(). Expected", expected, "got", count)
		}
	}

	check(0)
	r.Uint8()
	check(1)
	r.Uint16()
	check(3)
	r.Uint32()
	check(7)
	r.Uint64()
	check(15)
	r.Read(3)
	check(18)
}

func TestReadString(t *testing.T) {
	x := "TetrationAnalytics"
	b := []byte{byte(len(x))}
	xb := []byte(x)
	b = append(b, xb...)
	fmt.Printf("byte: %+v\n", b)
	r := NewReader(b)
	d, err := r.Read(65535)
	if err != nil {
		t.Errorf("read failed: %s", err.Error())
	}
	dx := string(d)
	fmt.Printf("read string: %s", dx)
	if x != dx {
		t.Errorf("string differs: found %s, expected: %s", x, dx)
	}
}

func TestReadTwoStrings(t *testing.T) {
	x := []string{"TetrationAnalytics", "NetFlowSensor"}
	b := []byte{}
	for i := range x {
		xi := x[i]
		b = append(b, byte(len(xi)))
		xb := []byte(xi)
		b = append(b, xb...)
	}
	fmt.Printf("byte: %+v\n", b)
	r := NewReader(b)
	for i := range x {
		xi := x[i]
		d, err := r.Read(65535)
		if err != nil {
			t.Errorf("read failed: %s", err.Error())
		}
		dx := string(d)
		fmt.Printf("read string: %s", dx)
		if xi != dx {
			t.Errorf("string differs: found %s, expected: %s", x, dx)
		}
	}
}

func TestReadStringError(t *testing.T) {
	r := NewReader([]byte{})
	_, err := r.readString(10)
	if err == nil {
		t.Errorf("we should have errored out")
	}
	x := "TetrationAnalytics"
	b := []byte{5}
	xb := []byte(x)
	b = append(b, xb...)
	r = NewReader(b)
	fmt.Printf("byte: %+v\n", b)
	_, err = r.readString(65535)
	if err != nil {
		t.Errorf("no error expected even though expected len is %d, found %d", len(x), b[0])
	}
	b[0] = byte(len(x))
	r = NewReader(b)
	_, err = r.readString(5)
	if err == nil {
		t.Errorf("to read only %d, string length is %d", 5, b[0])
	}
	fmt.Printf("byte: %+v\n", b)
	r = NewReader(b)
	d, err := r.readString(65535)
	if err != nil {
		t.Errorf("no error expected")
	}
	if string(d) != x {
		t.Errorf("expected string %s, found %s", x, string(d))
	}
}
