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
	Server string `yaml:"server"`
}

func (n *NSQ) setup(configFile string, logger *log.Logger) error {
	var (
		err error
		cfg = nsq.NewConfig()
	)

	// set default values
	n.config = NSQConfig{
		Server: "localhost:4150",
	}

	// load configuration if available
	if err = n.load(configFile); err != nil {
		logger.Println(err)
	}

	cfg.ClientID = "vflow.nsq"

	n.producer, err = nsq.NewProducer(n.config.Server, cfg)
	if err != nil {
		logger.Println(err)
		return err
	}

	n.logger = logger

	return nil
}

func (n *NSQ) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		err error
		ok  bool
	)

	n.logger.Printf("start producer: NSQ, server: %+v, topic: %s\n",
		n.config.Server, topic)

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
