//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    decoder_test.go
//: details: netflow v9 decoder tests and benchmarks
//: author:  Mehrdad Arshad Rad
//: date:    05/05/2017
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
	"net"
	"testing"
)

func TestDecodeNoData(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")
	body := []byte{}
	d := NewDecoder(ip, body)
	if _, err := d.Decode(mCache); err == nil {
		t.Error("expected err but nothing")
	}
}
