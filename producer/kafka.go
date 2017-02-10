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

import "github.com/Shopify/sarama"

type Kafka struct {
	producer sarama.AsyncProducer
}

func (k *Kafka) setup(configFile string) error {
	var err error

	k.producer, err = sarama.NewAsyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) inputMsg(topic string, mCh chan string) {
	var msg string
	for {
		msg = <-mCh
		select {
		case k.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(msg),
		}:
		case err := <-k.producer.Errors():
			println(err.Error())
		}
	}
}
