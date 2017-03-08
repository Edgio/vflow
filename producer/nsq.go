// Package producer push decoded messages to messaging queue
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    nsq.go
//: details: vflow nsq producer plugin
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
	"io/ioutil"
	"log"

	"github.com/bitly/go-nsq"
	"gopkg.in/yaml.v2"
)

// NSQ represents nsq producer
type NSQ struct {
	producer *nsq.Producer
	config   NSQConfig
	logger   *log.Logger
}

// NSQConfig represents NSQ configuration
type NSQConfig struct {
	Broker string `json:"broker"`
}

func (n *NSQ) setup(configFile string, logger *log.Logger) error {
	var err error
	// set default values
	n.config = NSQConfig{
		Broker: "localhost:4150",
	}

	// load configuration if available
	if err = n.load(configFile); err != nil {
		logger.Println(err)
	}

	n.producer, _ = nsq.NewProducer(n.config.Broker, nil)

	return nil
}

func (n *NSQ) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		err error
		ok  bool
	)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		err = n.producer.Publish(topic, msg)
		if err != nil {
			n.logger.Println(err)
			*ec++
		}
	}
}

func (n *NSQ) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &n.config)
	if err != nil {
		return err
	}

	return nil
}
