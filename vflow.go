package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
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

	for {
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Println(err)
			continue
		}
		udpChn <- UDPMsg{raddr, bytes.NewReader(b[:n])}
	}

	wg.Wait()
}

func sFlowWorker() {
	filter := []uint32{sflow.DataCounterSample}

	for {
		msg := <-udpChn
		println("rcvd", msg.body.Size())
		d := sflow.NewSFDecoder(msg.body, filter)
		d.SFDecode()
	}
}

func main() {
	fmt.Println("start listening")

	sFlow := SFServer{
		port:    "6343",
		udpSize: 1500,
		workers: 10,
	}
	sFlow.run()
}
