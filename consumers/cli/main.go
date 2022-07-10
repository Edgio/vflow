//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    main.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    05/25/2017
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

// Package main is the command line vflow IPFIX consumer with simple filter feature
package main

import (
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	cluster "github.com/bsm/sarama-cluster"
)

type options struct {
	Broker  string
	Topic   string
	Id      int
	Value   string
	Debug   bool
	Workers int
}

type dataField struct {
	I int
	V interface{}
}

type ipfix struct {
	AgentID  string
	DataSets [][]dataField
}

var opts options

func init() {
	flag.StringVar(&opts.Broker, "broker", "127.0.0.1:9092", "broker ipaddress:port")
	flag.StringVar(&opts.Topic, "topic", "vflow.ipfix", "kafka topic")
	flag.StringVar(&opts.Value, "value", "8.8.8.8", "element value - string")
	flag.BoolVar(&opts.Debug, "debug", false, "enabled/disabled debug")
	flag.IntVar(&opts.Id, "id", 12, "IPFIX element ID")
	flag.IntVar(&opts.Workers, "workers", 16, "workers number / partition number")

	flag.Parse()
}

func main() {
	var wg sync.WaitGroup

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	wg.Add(opts.Workers)

	for i := 0; i < opts.Workers; i++ {
		go func(ti int) {
			var objmap ipfix

			brokers := []string{opts.Broker}
			topics := []string{opts.Topic}
			consumer, err := cluster.NewConsumer(brokers, "mygroup", topics, config)

			if err != nil {
				panic(err)
			}
			defer consumer.Close()

			pCount := 0
			count := 0
			tik := time.Tick(10 * time.Second)

			for {
				select {
				case <-tik:
					if opts.Debug {
						log.Printf("partition GroupId#%d,  rate=%d\n", ti, (count-pCount)/10)
					}
					pCount = count
				case msg, more := <-consumer.Messages():
					if more {
						if err := json.Unmarshal(msg.Value, &objmap); err != nil {
							log.Println(err)
						} else {
							for _, data := range objmap.DataSets {
								for _, dd := range data {
									if dd.I == opts.Id && dd.V == opts.Value {
										log.Printf("%#v\n", data)
									}
								}
							}
						}

						consumer.MarkOffset(msg, "")
						count++
					}
				}
			}
		}(i)
	}

	wg.Wait()
}
