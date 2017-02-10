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

type Producer struct {
	mq           MQueue
	mqConfigFile string

	iChan chan string
	sChan chan string

	logger *log.Logger
}

type MQueue interface {
	setup(string, *log.Logger) error
	inputMsg(string, chan string)
}

var mqRegistered = map[string]MQueue{"kafka": new(Kafka)}

func NewProducer(mqName, mqConfigFile string, logger *log.Logger) *Producer {
	return &Producer{
		mq:           mqRegistered[mqName],
		mqConfigFile: mqConfigFile,
		logger:       logger,
	}
}

func (p *Producer) RegSFlowChan(sCh chan string) {
	p.sChan = sCh
}

func (p *Producer) RegIPFIXChan(iCh chan string) {
	p.iChan = iCh
}

func (p *Producer) Run() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	err = p.mq.setup(p.mqConfigFile, p.logger)
	if err != nil {
		return err
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		p.mq.inputMsg("vflow.ipfix", p.iChan)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		p.mq.inputMsg("vflow.sflow", p.sChan)
	}()

	wg.Wait()

	return nil
}
