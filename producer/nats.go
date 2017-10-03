//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    nats.go
//: details: vflow nats producer plugin
//: author:  Jeremy Rossi
//: date:    06/19/2017
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

	"github.com/nats-io/go-nats"
	"gopkg.in/yaml.v2"
)

// NATS represents nats producer
type NATS struct {
	connection *nats.Conn
	config     NATSConfig
	logger     *log.Logger
}

// NATSConfig is the struct that holds all configuation for NATS connections
type NATSConfig struct {
	URL string `yaml:"url"`
}

func (n *NATS) setup(configFile string, logger *log.Logger) error {
	var err error
	n.config = NATSConfig{
		URL: nats.DefaultURL,
	}

	if err = n.load(configFile); err != nil {
		logger.Println(err)
		return err
	}

	n.connection, err = nats.Connect(n.config.URL)
	if err != nil {
		logger.Println(err)
		return err
	}

	n.logger = logger

	return nil
}

func (n *NATS) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		err error
		ok  bool
	)

	n.logger.Printf("start producer: NATS, server: %+v, topic: %s\n",
		n.config.URL, topic)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		err = n.connection.Publish(topic, msg)
		if err != nil {
			n.logger.Println(err)
			*ec++
		}
	}
}

func (n *NATS) load(f string) error {
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
