// Package producer push decoded messages to messaging queue
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    producer.go
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
)

// Producer represents messaging queue
type Producer struct {
	MQ           MQueue
	MQConfigFile string

	Topic string
	Chan  chan []byte

	Logger *log.Logger
}

// MQueue represents messaging queue methods
type MQueue interface {
	setup(string, *log.Logger) error
	inputMsg(string, chan []byte)
}

// register messaging queues
var mqRegistered = map[string]MQueue{
	"kafka": new(Kafka),
	"nsq":   new(NSQ),
}

// NewProducer constructs new Messaging Queue
func NewProducer(mqName string) *Producer {
	return &Producer{
		MQ: mqRegistered[mqName],
	}
}

// Run configs and tries to be ready to produce
func (p *Producer) Run() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	err = p.MQ.setup(p.MQConfigFile, p.Logger)
	if err != nil {
		return err
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		p.MQ.inputMsg("vflow."+p.Topic, p.Chan)
	}()

	wg.Wait()

	return nil
}

// Shutdown stops the producer
func (p *Producer) Shutdown() {
	close(p.Chan)
}
