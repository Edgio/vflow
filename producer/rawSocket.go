//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    rawSocket.go
//: details: vflow tcp/udp producer plugin
//: author:  Joe Percivall
//: date:    12/18/2017
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
	"strings"

	"fmt"
	"gopkg.in/yaml.v2"
	"net"
)

// RawSocket represents RawSocket producer
type RawSocket struct {
	connection net.Conn
	config     RawSocketConfig
	logger     *log.Logger
}

// RawSocketConfig is the struct that holds all configuation for RawSocketConfig connections
type RawSocketConfig struct {
	URL      string `yaml:"url"`
	Protocol string `yaml:"protocol"`
	MaxRetry int    `yaml:"retry-max"`
}

func (rs *RawSocket) setup(configFile string, logger *log.Logger) error {
	var err error
	rs.config = RawSocketConfig{
		URL:      "localhost:9555",
		Protocol: "tcp",
		MaxRetry: 2,
	}

	if err = rs.load(configFile); err != nil {
		logger.Println(err)
		return err
	}

	rs.connection, err = net.Dial(rs.config.Protocol, rs.config.URL)
	if err != nil {
		logger.Println(err)
		return err
	}

	rs.logger = logger

	return nil
}

func (rs *RawSocket) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		err error
		ok  bool
	)

	rs.logger.Printf("start producer: RawSocket, server: %+v, Protocol: %s\n",
		rs.config.URL, rs.config.Protocol)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		for i := 0; ; i++ {
			_, err = fmt.Fprintf(rs.connection, string(msg)+"\n")
			if err == nil {
				break
			}

			*ec++

			if strings.HasSuffix(err.Error(), "broken pipe") {
				var newConnection, err = net.Dial(rs.config.Protocol, rs.config.URL)
				if err != nil {
					rs.logger.Println("Error when attempting to fix the broken pipe", err)
				} else {
					rs.logger.Println("Successfully reconnected")
					rs.connection = newConnection
				}
			}

			if i >= (rs.config.MaxRetry) {
				rs.logger.Println("message failed after the configured retry limit:", err)
				break
			} else {
				rs.logger.Println("retrying after error:", err)
			}
		}
	}
}

func (rs *RawSocket) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &rs.config)
	if err != nil {
		return err
	}

	return nil
}
