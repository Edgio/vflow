//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    stress.go
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

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/EdgeCast/vflow/stress/hammer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var opts = struct {
	vflowAddr      string
	ipfixPort      int
	sflowPort      int
	prometheusPort int
	ipfixTick      string
	ipfixRateLimit int
	sflowRateLimit int
}{
	"127.0.0.1",
	4739,
	6343,
	8090,
	"10s",
	25000,
	25000,
}

func init() {
	flag.IntVar(&opts.ipfixPort, "ipfix-port", opts.ipfixPort, "ipfix port number")
	flag.IntVar(&opts.sflowPort, "sflow-port", opts.sflowPort, "sflow port number")
	flag.IntVar(&opts.prometheusPort, "prometheus-port", opts.prometheusPort, "prometheus port number")
	flag.StringVar(&opts.ipfixTick, "ipfix-interval", opts.ipfixTick, "ipfix template interval in seconds")
	flag.IntVar(&opts.ipfixRateLimit, "ipfix-rate-limit", opts.ipfixRateLimit, "ipfix rate limit packets per second")
	flag.IntVar(&opts.sflowRateLimit, "sflow-rate-limit", opts.sflowRateLimit, "sflow rate limit packets per second")
	flag.StringVar(&opts.vflowAddr, "vflow-addr", opts.vflowAddr, "vflow ip address")

	flag.Parse()
}

func main() {
	var (
		wg    sync.WaitGroup
		vflow = net.ParseIP(opts.vflowAddr)
	)

	prometheus.Unregister(collectors.NewGoCollector())
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opts.prometheusPort), nil))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ipfix, err := hammer.NewIPFIX(vflow)
		if err != nil {
			log.Fatalf("got error while NewIPFIX, %v", err)
		}
		ipfix.Port = opts.ipfixPort
		ipfix.Tick, err = time.ParseDuration(opts.ipfixTick)
		ipfix.RateLimit = opts.ipfixRateLimit
		if err != nil {
			log.Fatal(err)
		}
		ipfix.Run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sflow, err := hammer.NewSFlow(vflow)
		if err != nil {
			log.Fatalf("got error while NewSFlow, %v", err)
		}
		sflow.Port = opts.sflowPort
		sflow.RateLimit = opts.sflowRateLimit
		sflow.Run()
	}()

	log.Printf("Stress is attacking %s target ...", opts.vflowAddr)

	wg.Wait()
}
