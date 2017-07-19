//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    memcache_test.go
//: details: memory template cache testing
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
	"net"
	"reflect"
	"testing"
)

func TestMemCacheRetrieve(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")
	d := NewDecoder(ip, tpl)
	d.Decode(mCache)
	v, ok := mCache.retrieve(256, ip)
	if !ok {
		t.Error("expected mCache retrieve status true, got", ok)
	}
	if v.TemplateID != 256 {
		t.Error("expected template id#:256, got", v.TemplateID)
	}
}

func TestMemCacheInsert(t *testing.T) {
	var tpl TemplateRecord
	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")

	tpl.TemplateID = 310
	mCache.insert(310, ip, tpl)

	v, ok := mCache.retrieve(310, ip)
	if !ok {
		t.Error("expected mCache retrieve status true, got", ok)
	}
	if v.TemplateID != 310 {
		t.Error("expected template id#:310, got", v.TemplateID)
	}
}

func TestMemCacheAllSetIds(t *testing.T) {
	var tpl TemplateRecord
	ip := net.ParseIP("127.0.0.1")
	mCache := GetCache("cache.file")

	tpl.TemplateID = 310
	mCache.insert(tpl.TemplateID, ip, tpl)
	tpl.TemplateID = 410
	mCache.insert(tpl.TemplateID, ip, tpl)
	tpl.TemplateID = 210
	mCache.insert(tpl.TemplateID, ip, tpl)

	expected := []int{210, 310, 410}
	actual := mCache.allSetIds()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected set IDs %v, got %v", expected, actual)
	}
}
