//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    vflow.go
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

// Package main is the vflow binary
package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	opts   *Options
	logger *log.Logger
)

func main() {
	var (
		wg       sync.WaitGroup
		signalCh = make(chan os.Signal, 1)
	)

	opts = GetOptions()

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	sFlow := NewSFlow()
	ipfix := NewIPFIX()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sFlow.run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ipfix.run()
	}()

	go statsHTTPServer(ipfix, sFlow)

	<-signalCh

	wg.Add(1)
	go func() {
		defer wg.Done()
		sFlow.shutdown()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ipfix.shutdown()
	}()

	wg.Wait()
}
