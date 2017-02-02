//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    stats.go
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
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"time"
)

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
	}

	j, _ := json.Marshal(data)

	w.Write(j)
}

func statsHTTPServer() {
	if !opts.StatsEnabled {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sys", StatsSysHandler)

	addr := net.JoinHostPort(opts.StatsHTTPAddr, opts.StatsHTTPPort)

	logger.Println("starting stats web server ...")
	logger.Println(http.ListenAndServe(addr, mux))
}
