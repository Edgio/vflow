// Package producer push decoded messages to messaging queue
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    kafka.go
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
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bitly/go-nsq"
)

// NSQ represents nsq producer
type NSQ struct {
	producer *nsq.Producer
	config   NSQConfig
	logger   *log.Logger
}

type NSQConfig struct {
	Broker string `json:"broker"`
}

func (n *NSQ) setup(configFile string, logger *log.Logger) error {
	n.producer, _ = nsq.NewProducer("127.0.0.1:4150", nil)
	// TODO

	return nil
}

func (n *NSQ) inputMsg(topic string, mCh chan string) {
	// TODO
}

func (n *NSQ) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &n.config)
	if err != nil {
		return err
	}

	return nil
}
