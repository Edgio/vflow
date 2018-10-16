//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    kafka.go.no
//: details: vflow kafka producer plugin
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
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/segmentio/kafka-go.v0"
	"gopkg.in/segmentio/kafka-go.v0/gzip"
	"gopkg.in/segmentio/kafka-go.v0/lz4"
	"gopkg.in/segmentio/kafka-go.v0/snappy"
	"gopkg.in/yaml.v2"
)

// Kafka represents kafka producer
type Kafka struct {
	producer *kafka.Writer
	config   kafka.WriterConfig
	fileconf FileConfig
	logger   *log.Logger
}

// KafkaConfig represents kafka configuration
type FileConfig struct {
	Brokers         []string `yaml:"brokers" env:"BROKERS"`
	BootstrapServer string   `yaml:"bootstrap_server" env:"BOOTSTRAP_SERVER"`
	Compression     string   `yaml:"compression" env:"COMPRESSION"`

	ConnectTimeout int   `yaml:"connect-timeout" env:"CONNECT_TIMEOUT"`
	RetryMax       int   `yaml:"retry-max" env:"RETRY_MAX"`
	RequestSizeMax int32 `yaml:"request-size-max" env:"REQUEST_SIZE_MAX"`
	RetryBackoff   int   `yaml:"retry-backoff" env:"RETRY_BACKOFF"`
	RequiredAcks   int   `yaml:"required-acks" env:"REQUIRED_ACKS"`

	TLSCertFile string `yaml:"tls-cert" env:"TLS_CERT"`
	TLSKeyFile  string `yaml:"tls-key" env:"TLS_KEY"`
	CAFile      string `yaml:"ca-file" env:"CA_FILE"`
	VerifySSL   bool   `yaml:"verify-ssl" env:"VERIFY_SSL"`
}

func (k *Kafka) setup(configFile string, logger *log.Logger) error {
	var err error

	// set default values
	k.fileconf = FileConfig{
		Brokers:        []string{"localhost:9092"},
		ConnectTimeout: 10,
		RequiredAcks:   1,
		RetryMax:       2,
		RequestSizeMax: 104857600,
		RetryBackoff:   10,
		VerifySSL:      true,
	}

	// setup logger
	k.logger = logger

	// load configuration file if available
	if err = k.load(configFile); err != nil {
		logger.Println(err)
	}

	// get env config
	k.loadEnv("VFLOW_KAFKA")

	// init kafka configuration
	k.config = kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Dialer: &kafka.Dialer{
			ClientID:  "vFlow.Kafka",
			Timeout:   10,
			DualStack: true,
		},
		Balancer:          &kafka.Hash{},
		MaxAttempts:       2,
		QueueCapacity:     1024,
		BatchSize:         512,
		RebalanceInterval: 10,
		RequiredAcks:      1,
	}

	if tlsConfig := k.tlsConfig(); tlsConfig != nil {
		k.config.Dialer.TLS = tlsConfig
		k.logger.Println("Kafka client TLS enabled")
	}

	switch k.fileconf.Compression {
	case "gzip":
		k.config.CompressionCodec = gzip.NewCompressionCodec()
	case "lz4":
		k.config.CompressionCodec = lz4.NewCompressionCodec()
	case "snappy":
		k.config.CompressionCodec = snappy.NewCompressionCodec()
	}

	return err
}

func (k *Kafka) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		ok  bool
	)

	k.logger.Printf("start producer: Kafka, brokers: %+v, topic: %s\n",
		k.config.Brokers, topic)
	k.config.Topic = topic
	k.producer = kafka.NewWriter(k.config)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		err := k.producer.WriteMessages(context.Background(), kafka.Message{
			Value: msg,
		})

		k.logger.Println(err.Error())
	}

	k.producer.Close()
}

func (k *Kafka) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &k.fileconf)
	if err != nil {
		return err
	}

	return nil
}

func (k Kafka) tlsConfig() *tls.Config {
	var t *tls.Config

	if k.fileconf.TLSCertFile != "" && k.fileconf.TLSKeyFile != "" && k.fileconf.CAFile != "" {
		cert, err := tls.LoadX509KeyPair(k.fileconf.TLSCertFile, k.fileconf.TLSKeyFile)
		if err != nil {
			k.logger.Fatal("Kafka TLS error: ", err)
		}

		caCert, err := ioutil.ReadFile(k.fileconf.CAFile)
		if err != nil {
			k.logger.Fatal("Kafka TLS error: ", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: k.fileconf.VerifySSL,
		}
	}

	return t
}

func (k *Kafka) loadEnv(prefix string) {
	v := reflect.ValueOf(&k.config).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		env := t.Field(i).Tag.Get("env")
		if env == "" {
			continue
		}

		val, ok := os.LookupEnv(prefix + "_" + env)
		if !ok {
			continue
		}

		switch f.Kind() {
		case reflect.Int:
			valInt, err := strconv.Atoi(val)
			if err != nil {
				k.logger.Println(err)
				continue
			}
			f.SetInt(int64(valInt))
		case reflect.String:
			f.SetString(val)
		case reflect.Slice:
			for _, elm := range strings.Split(val, ";") {
				f.Index(0).SetString(elm)
			}
		case reflect.Bool:
			valBool, err := strconv.ParseBool(val)
			if err != nil {
				k.logger.Println(err)
				continue
			}
			f.SetBool(valBool)
		}
	}
}