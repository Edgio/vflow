//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    producer_test.go
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

package producer

import (
	"log"
	"sync"
	"testing"
)

type MQMock struct{}

func (k *MQMock) setup(configFile string, logger *log.Logger) error {
	return nil
}

func (k *MQMock) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	for {
		msg, ok := <-mCh
		if !ok {
			break
		}
		mCh <- msg
	}
}

func TestProducerChan(t *testing.T) {
	var (
		ch = make(chan []byte, 1)
		wg sync.WaitGroup
	)

	p := Producer{MQ: new(MQMock)}
	p.Chan = ch

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.Run(); err != nil {
			t.Error("unexpected error", err)
		}
	}()

	ch <- []byte("test")
	m := <-ch
	if string(m) != "test" {
		t.Error("expect to get test, got", string(m))
	}

	close(ch)

	wg.Wait()
}
