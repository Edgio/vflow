package main

import (
	"encoding/json"
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

func statsHTTPServer(opts *Options) {
	var logger = opts.Logger

	if !opts.StatsEnabled {
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/sys", StatsSysHandler)

	logger.Println("starting stats web server ...")
	logger.Println(http.ListenAndServe(":8080", mux))
}
