//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    sarama.go
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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"gopkg.in/yaml.v2"
)

// KafkaSarama represents kafka producer
type KafkaSarama struct {
	producer sarama.AsyncProducer
	config   KafkaSaramaConfig
	logger   *log.Logger
}

// KafkaSaramaConfig represents kafka configuration
type KafkaSaramaConfig struct {
	Brokers        []string `yaml:"brokers" env:"BROKERS"`
	Compression    string   `yaml:"compression" env:"COMPRESSION"`
	RetryMax       int      `yaml:"retry-max" env:"RETRY_MAX"`
	RequestSizeMax int32    `yaml:"request-size-max" env:"REQUEST_SIZE_MAX"`
	RetryBackoff   int      `yaml:"retry-backoff" env:"RETRY_BACKOFF"`
	TLSEnabled     bool     `yaml:"tls-enabled" env:"TLS_ENABLED"`
	TLSCertFile    string   `yaml:"tls-cert" env:"TLS_CERT"`
	TLSKeyFile     string   `yaml:"tls-key" env:"TLS_KEY"`
	CAFile         string   `yaml:"ca-file" env:"CA_FILE"`
	TLSSkipVerify  bool     `yaml:"tls-skip-verify" env:"TLS-SKIP-VERIFY"`
}

func (k *KafkaSarama) setup(configFile string, logger *log.Logger) error {
	var (
		config = sarama.NewConfig()
		err    error
	)

	// set default values
	k.config = KafkaSaramaConfig{
		Brokers:        []string{"localhost:9092"},
		RetryMax:       2,
		RequestSizeMax: 104857600,
		RetryBackoff:   10,
		TLSEnabled:     false,
		TLSSkipVerify:  true,
	}

	k.logger = logger

	// load configuration if available
	if err = k.load(configFile); err != nil {
		logger.Println(err)
	}

	// init kafka configuration
	config.ClientID = "vFlow.Kafka"
	config.Producer.Retry.Max = k.config.RetryMax
	config.Producer.Retry.Backoff = time.Duration(k.config.RetryBackoff) * time.Millisecond

	sarama.MaxRequestSize = k.config.RequestSizeMax

	switch k.config.Compression {
	case "gzip":
		config.Producer.Compression = sarama.CompressionGZIP
	case "lz4":
		config.Producer.Compression = sarama.CompressionLZ4
	case "snappy":
		config.Producer.Compression = sarama.CompressionSnappy
	default:
		config.Producer.Compression = sarama.CompressionNone
	}

	if tlsConfig := k.tlsConfig(); tlsConfig != nil || k.config.TLSEnabled {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
		if k.config.TLSSkipVerify {
			k.logger.Printf("kafka client TLS enabled (server certificate didn't validate)")
		} else {
			k.logger.Printf("kafka client TLS enabled")
		}
	}

	// get env config
	k.loadEnv("VFLOW_KAFKA")

	if err = config.Validate(); err != nil {
		logger.Fatal(err)
	}

	k.producer, err = sarama.NewAsyncProducer(k.config.Brokers, config)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaSarama) inputMsg(topic string, mCh chan []byte, ec *uint64) {
	var (
		msg []byte
		ok  bool
	)

	k.logger.Printf("start producer: Kafka, brokers: %+v, topic: %s\n",
		k.config.Brokers, topic)

	for {
		msg, ok = <-mCh
		if !ok {
			break
		}

		select {
		case k.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msg),
		}:
		case err := <-k.producer.Errors():
			k.logger.Println(err)
			*ec++
		}
	}

	k.producer.Close()
}

func (k *KafkaSarama) load(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &k.config)
	if err != nil {
		return err
	}

	return nil
}

func (k KafkaSarama) tlsConfig() *tls.Config {
	var t *tls.Config

	if k.config.TLSCertFile != "" || k.config.TLSKeyFile != "" || k.config.CAFile != "" {
		cert, err := tls.LoadX509KeyPair(k.config.TLSCertFile, k.config.TLSKeyFile)
		if err != nil {
			k.logger.Fatal("kafka TLS load X509 key pair error: ", err)
		}

		caCert, err := ioutil.ReadFile(k.config.CAFile)
		if err != nil {
			k.logger.Fatal("kafka TLS CA file error: ", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: k.config.TLSSkipVerify,
		}
	}

	return t
}

func (k *KafkaSarama) loadEnv(prefix string) {
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
