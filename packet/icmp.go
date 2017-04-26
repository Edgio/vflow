//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    icmp.go
//: details: TODO
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

package packet

import "errors"

// ICMP represents ICMP header
type ICMP struct {
	// Type is ICMP type
	Type int

	// Code is ICMP subtype
	Code int
}

var errICMPHLenTooSHort = errors.New("ICMP header length is too short")

func decodeICMP(b []byte) (ICMP, error) {
	if len(b) < 4 {
		return ICMP{}, errICMPHLenTooSHort
	}

	return ICMP{
		Type: int(b[0]),
		Code: int(b[1]),
	}, nil
}
