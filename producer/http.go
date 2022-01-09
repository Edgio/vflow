//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    http.go
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
	"context"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
	"net"
	"net/http"
)

// Http represents Http producer
type Http struct {
	connection http.Client
	config     HttpConfig
	logger     *log.Logger
}

// HttpConfig is the struct that holds all configuation for HttpConfig connections
type HttpConfig struct {
	Address  string `yaml:"address"`
	URL      string `yaml:"url"`
	Protocol string `yaml:"protocol"`
	MaxRetry int    `yaml:"retry-max"`
}

func (h *Http) setup(configFile string, logger *log.Logger) error {
	var err error
	h.config = HttpConfig{
		Address:  "localhost:9555",
		URL:      "",
		Protocol: "tcp",
		MaxRetry: 2,
	}

	if err = h.load(configFile); err != nil {
		logger.Println(err)
		return err
	}

	h.connection = http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial(h.config.Protocol, h.config.Address)
			},
		},
	}

	h.logger = logger

	return nil
}

func (h *Http) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		err error
		ok  bool
	)

	h.logger.Printf("start producer: Http, server: %+v, Protocol: %s\n",
		h.config.Address, h.config.Protocol)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		event := fmt.Sprintf("[%s]", string(msg))
		for i := 0; ; i++ {
			_, err = h.connection.Post("http://unix"+h.config.URL, "application/json", strings.NewReader(event))
			if err == nil {
				break
			}

			h.logger.Println("Error sending event via POST:", err)

			*ec++

			if strings.HasSuffix(err.Error(), "broken pipe") {
				var newConnection = http.Client{
					Timeout: time.Second * 60,
					Transport: &http.Transport{
						DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
							return net.Dial(h.config.Protocol, h.config.Address)
						},
					},
				}
				if err != nil {
					h.logger.Println("Error when attempting to fix the broken pipe", err)
				} else {
					h.logger.Println("Successfully reconnected")
					h.connection = newConnection
				}
			}

			if i >= (h.config.MaxRetry) {
				h.logger.Println("message failed after the configured retry limit:", err)
				break
			} else {
				h.logger.Println("retrying after error:", err)
			}
		}
	}
}

func (h *Http) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &h.config)
	if err != nil {
		return err
	}

	return nil
}
