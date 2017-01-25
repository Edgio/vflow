package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"git.edgecastcdn.net/vflow/packet"
	"git.edgecastcdn.net/vflow/sflow"
)

type SFUDPMsg struct {
	raddr *net.UDPAddr
	body  *bytes.Reader
}

var (
	sFlowUdpCh = make(chan SFUDPMsg, 1000)
)

type SFServer struct {
	port        int
	addr        string
	laddr       *net.UDPAddr
	readTimeout time.Duration
	udpSize     int
	workers     int
	stop        bool
}

func NewSFlow(opts *Options) *SFServer {
	return &SFServer{
		port:    opts.SFlowPort,
		udpSize: opts.SFlowUDPSize,
		workers: opts.SFlowWorkers,
	}
}

func (s *SFServer) run() {
	var (
		b  = make([]byte, s.udpSize)
		wg sync.WaitGroup
	)

	hostPort := net.JoinHostPort(s.addr, strconv.Itoa(s.port))
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
		sFlowUdpCh <- SFUDPMsg{raddr, bytes.NewReader(b[:n])}
	}

	wg.Wait()
}

func (s *SFServer) shutdown() {
	s.stop = true
	log.Println("stopped sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	log.Println("vFlow has been shutdown")
	close(sFlowUdpCh)
}

func sFlowWorker() {
	var (
		msg    SFUDPMsg
		ok     bool
		filter = []uint32{sflow.DataCounterSample}
	)

	for {
		if msg, ok = <-sFlowUdpCh; !ok {
			break
		}

		log.Println("rcvd sflow data", msg.body.Size())
		d := sflow.NewSFDecoder(msg.body, filter)
		data, err := d.SFDecode()
		if err != nil {
			log.Println(err)
		}

		switch data.(type) {
		case *packet.Packet:
			log.Printf("%#v\n", data)
		}
	}
}
