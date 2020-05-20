//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    stats.go
//: details: exposes flow status
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
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime = time.Now().Unix()

type rest struct {
	StartTime int64
	IPFIX     *IPFIXStats
	SFlow     *SFlowStats
	NetflowV5 *NetflowV5Stats
	NetflowV9 *NetflowV9Stats
}

func statsSysHandler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)
	var data = &struct {
		MemAlloc        uint64
		MemTotalAlloc   uint64
		MemHeapAlloc    uint64
		MemHeapSys      uint64
		MemHeapReleased uint64
		MCacheInuse     uint64
		GCSys           uint64
		GCNext          uint64
		GCLast          string
		NumLogicalCPU   int
		NumGoroutine    int
		MaxProcs        int
		GoVersion       string
		StartTime       int64
	}{
		mem.Alloc,
		mem.TotalAlloc,
		mem.HeapAlloc,
		mem.HeapSys,
		mem.HeapReleased,
		mem.MCacheInuse,
		mem.GCSys,
		mem.NextGC,
		time.Unix(0, int64(mem.LastGC)).String(),
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		runtime.GOMAXPROCS(-1),
		runtime.Version(),
		startTime,
	}

	j, err := json.Marshal(data)
	if err != nil {
		logger.Println(err)
	}

	if _, err = w.Write(j); err != nil {
		logger.Println(err)
	}
}

func statsFlowHandler(protos []proto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rd := &rest{StartTime: startTime}
		for _, p := range protos {
			switch p.(type) {
			case *IPFIX:
				ipfix, _ := p.(*IPFIX)
				rd.IPFIX = ipfix.status()
			case *SFlow:
				sflow, _ := p.(*SFlow)
				rd.SFlow = sflow.status()
			case *NetflowV5:
				netflowv5, _ := p.(*NetflowV5)
				rd.NetflowV5 = netflowv5.status()
			case *NetflowV9:
				netflowv9, _ := p.(*NetflowV9)
				rd.NetflowV9 = netflowv9.status()
			}
		}

		j, err := json.Marshal(rd)
		if err != nil {
			logger.Println(err)
		}

		if _, err = w.Write(j); err != nil {
			logger.Println(err)
		}
	}
}

func statsExpose(protos []proto) {
	if opts.StatsFormat != "prometheus" {
		statsRest(protos)
	} else {
		statsPrometheus(protos)
	}
}

func statsRest(protos []proto) {
	if !opts.StatsEnabled {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sys", statsSysHandler)
	mux.HandleFunc("/flow", statsFlowHandler(protos))

	logger.Println("starting stats http server ...")

	addr := net.JoinHostPort(opts.StatsHTTPAddr, opts.StatsHTTPPort)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		logger.Fatal(err)
	}
}

func statsPrometheus(protos []proto) {
	for _, p := range protos {
		promCounterDecoded(p)
		promCounterMQError(p)
		promCounterUDP(p)
		promGaugeMessageQueue(p)
		promGaugeUDPQueue(p)
		promGaugeWorkers(p)
		promGaugeUDPMirrorQueue(p)
	}

	logger.Println("starting prometheus http server ...")

	addr := net.JoinHostPort(opts.StatsHTTPAddr, opts.StatsHTTPPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)
}

func promCounterDecoded(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_ipfix_decoded_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().DecodedCount)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_decoded_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().DecodedCount)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_decoded_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().DecodedCount)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_decoded_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().DecodedCount)
			})
	}
}

func promCounterMQError(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_ipfix_mq_error",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MQErrorCount)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_mq_error",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MQErrorCount)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_mq_error",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MQErrorCount)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_mq_error",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MQErrorCount)
			})
	}
}

func promCounterUDP(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_ipfix_udp_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPCount)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_udp_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPCount)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_udp_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPCount)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_udp_packets",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPCount)
			})
	}
}

func promGaugeMessageQueue(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "vflow_ipfix_message_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MessageQueue)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_message_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MessageQueue)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_message_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MessageQueue)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_message_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().MessageQueue)
			})
	}
}

func promGaugeUDPQueue(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "vflow_ipfix_udp_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPQueue)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_udp_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPQueue)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_udp_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPQueue)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_udp_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPQueue)
			})
	}
}

func promGaugeWorkers(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "vflow_ipfix_workers",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().Workers)
			})
	case *SFlow:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_sflow_workers",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().Workers)
			})
	case *NetflowV5:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv5_workers",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().Workers)
			})
	case *NetflowV9:
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "vflow_netflowv9_workers",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().Workers)
			})
	}
}

func promGaugeUDPMirrorQueue(p interface{}) {
	switch flow := p.(type) {
	case *IPFIX:
		promauto.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "vflow_ipfix_udp_mirror_queue",
			Help: "",
		},
			func() float64 {
				return float64(flow.status().UDPMirrorQueue)
			})
	}
}
