package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"git.edgecastcdn.net/vflow/sflow"
)

type SFServer struct {
	port        string
	addr        string
	laddr       *net.UDPAddr
	readTimeout time.Duration
	udpSize     int
	workers     int
	stop        bool
}

type UDPMsg struct {
	raddr *net.UDPAddr
	body  *bytes.Reader
}

var (
	udpChn = make(chan UDPMsg, 1000)
)

func (s *SFServer) run() {
	var (
		b  = make([]byte, s.udpSize)
		wg sync.WaitGroup
	)

	hostPort := net.JoinHostPort(s.addr, s.port)
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for i := 0; i < s.workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			sFlowWorker()

		}()
	}

	for !s.stop {
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		udpChn <- UDPMsg{raddr, bytes.NewReader(b[:n])}
	}

	wg.Wait()
}

func (s *SFServer) shutdown() {
	s.stop = true
	log.Println("stopped sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	log.Println("vFlow has been shutdown")
	close(udpChn)
}

func sFlowWorker() {
	var (
		msg    UDPMsg
		ok     bool
		filter = []uint32{sflow.DataCounterSample}
	)

	for {
		if msg, ok = <-udpChn; !ok {
			break
		}
		log.Println("rcvd", msg.body.Size())
		d := sflow.NewSFDecoder(msg.body, filter)
		d.SFDecode()
	}
}

func main() {
	var wg sync.WaitGroup
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	sFlow := SFServer{
		port:    "6343",
		udpSize: 1500,
		workers: 10,
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		sFlow.run()
	}()

	<-signalCh
	sFlow.shutdown()
	wg.Wait()
}
