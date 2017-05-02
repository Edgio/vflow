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
)

var startTime = time.Now().Unix()

// StatsSysHandler handles /sys endpoint
func StatsSysHandler(w http.ResponseWriter, r *http.Request) {
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

// StatsFlowHandler handles /flow endpoint
func StatsFlowHandler(i *IPFIX, s *SFlow, n *NetflowV9) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data = &struct {
			StartTime int64
			IPFIX     *IPFIXStats
			SFlow     *SFlowStats
			NetflowV9 *NetflowV9Stats
		}{
			startTime,
			i.status(),
			s.status(),
			n.status(),
		}

		j, err := json.Marshal(data)
		if err != nil {
			logger.Println(err)
		}

		if _, err = w.Write(j); err != nil {
			logger.Println(err)
		}
	}
}

func statsHTTPServer(ipfix *IPFIX, sflow *SFlow, netflow9 *NetflowV9) {
	if !opts.StatsEnabled {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sys", StatsSysHandler)
	mux.HandleFunc("/flow", StatsFlowHandler(ipfix, sflow, netflow9))

	addr := net.JoinHostPort(opts.StatsHTTPAddr, opts.StatsHTTPPort)

	logger.Println("starting stats web server ...")
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		logger.Fatal(err)
	}
}
