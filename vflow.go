package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var (
		wg       sync.WaitGroup
		signalCh = make(chan os.Signal, 1)
		opts     = GetOptions()
	)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	sFlow := NewSFlow(opts)
	ipfix := NewIPFIX(opts)

	go func() {
		wg.Add(1)
		defer wg.Done()
		sFlow.run()
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		ipfix.run()
	}()

	go statsHTTPServer(opts)

	<-signalCh
	go sFlow.shutdown()
	go ipfix.shutdown()
	wg.Wait()
}
