//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    main.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    06/01/2017
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

// Package main is the vflow IPFIX consumer for the ClickHouse database (https://clickhouse.com)
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Shopify/sarama"
)

type options struct {
	Broker  string
	Topic   string
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

type dIPFIXSample struct {
	device string
	src    string
	dst    string
	srcASN uint64
	dstASN uint64
	proto  uint8
}

type consumerGroupHandler struct {
	ti    int
	ch    chan ipfix
	debug bool
}

var opts options

func init() {
	flag.StringVar(&opts.Broker, "broker", "127.0.0.1:9092", "broker ipaddress:port")
	flag.StringVar(&opts.Topic, "topic", "vflow.ipfix", "kafka topic")
	flag.BoolVar(&opts.Debug, "debug", false, "enabled/disabled debug")
	flag.IntVar(&opts.Workers, "workers", 16, "workers number / partition number")

	flag.Parse()
}

func (h consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	pCount := 0
	count := 0
	tik := time.Tick(10 * time.Second)

	for msg := range claim.Messages() {
		select {
		case <-tik:
			if h.debug {
				log.Printf("partition GroupId#%d,  rate=%d\n", h.ti, (count-pCount)/10)
			}
			pCount = count
		default:
			objmap := ipfix{}
			if err := json.Unmarshal(msg.Value, &objmap); err != nil {
				log.Println(err)
			} else {
				h.ch <- objmap
			}
			sess.MarkMessage(msg, "")
			count++
		}
	}

	return nil
}

func main() {
	var (
		wg sync.WaitGroup
		ch = make(chan ipfix, 10000)
	)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_1_0_0

	for i := 0; i < 5; i++ {
		go ingestClickHouse(ch)
	}

	wg.Add(opts.Workers)

	for i := 0; i < opts.Workers; i++ {
		go func(ti int) {
			brokers := []string{opts.Broker}
			topics := []string{opts.Topic}
			consumerGroup, err := sarama.NewConsumerGroup(brokers, "mygroup", config)
			if err != nil {
				log.Fatalf("Failed to create consumer group: %s", err)
			}
			defer consumerGroup.Close()

			for {
				err := consumerGroup.Consume(context.Background(), topics, consumerGroupHandler{ti: ti, ch: ch, debug: opts.Debug})
				if err != nil {
					log.Printf("Error from consumer: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
}

func ingestClickHouse(ch chan ipfix) {
	var sample ipfix

	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=false")
	if err != nil {
		log.Println(err)
		return
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Println(err)
		}
		return
	}

	tx, err := connect.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO vflow.samples (date,time,device,src,dst,srcASN,dstASN, proto) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}

	for {
		for i := 0; i < 10000; i++ {
			sample = <-ch
			for _, data := range sample.DataSets {
				s := dIPFIXSample{}
				for _, dd := range data {
					switch dd.I {
					case 8, 27:
						s.src = dd.V.(string)
					case 12, 28:
						s.dst = dd.V.(string)
					case 16:
						s.srcASN = uint64(dd.V.(float64))
					case 17:
						s.dstASN = uint64(dd.V.(float64))
					case 4:
						s.proto = uint8(dd.V.(float64))
					}
				}
				if _, err := stmt.Exec(
					time.Now(),
					time.Now(),
					sample.AgentID,
					s.src,
					s.dst,
					s.srcASN,
					s.dstASN,
					s.proto,
				); err != nil {
					log.Println(err)
					return
				}

			}
		}

		go func(tx *sql.Tx) {
			if err := tx.Commit(); err != nil {
				log.Println(err)
			}
		}(tx)
	}
}
