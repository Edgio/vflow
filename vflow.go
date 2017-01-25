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
		opts     = NewOptions()
	)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	sFlow := NewSFlow(opts)

	go func() {
		wg.Add(1)
		defer wg.Done()
		sFlow.run()
	}()

	<-signalCh
	sFlow.shutdown()
	wg.Wait()
}
