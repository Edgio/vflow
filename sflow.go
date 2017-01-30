package main

import (
	"bytes"
	"log"
	"net"
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
	logger     *log.Logger
	verbose    bool
)

type SFlow struct {
	port        int
	addr        string
	laddr       *net.UDPAddr
	readTimeout time.Duration
	udpSize     int
	workers     int
	stop        bool
}

func NewSFlow(opts *Options) *SFlow {
	logger = opts.Logger
	verbose = opts.Verbose

	return &SFlow{
		port:    opts.SFlowPort,
		udpSize: opts.SFlowUDPSize,
		workers: opts.SFlowWorkers,
	}
}

func (s *SFlow) run() {
	var (
		b  = make([]byte, s.udpSize)
		wg sync.WaitGroup
	)

	hostPort := net.JoinHostPort(s.addr, strconv.Itoa(s.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal(err)
	}

	for i := 0; i < s.workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			sFlowWorker()

		}()
	}

	logger.Printf("sFlow is running (workers#: %d)", s.workers)

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

func (s *SFlow) shutdown() {
	s.stop = true
	logger.Println("stopped sflow service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("vFlow has been shutdown")
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

		if verbose {
			logger.Printf("rcvd sflow data from: %s, size: %d bytes",
				msg.raddr, msg.body.Size())
		}

		d := sflow.NewSFDecoder(msg.body, filter)
		records, err := d.SFDecode()
		if err != nil {
			logger.Println(err)
		}
		for _, data := range records {
			switch data.(type) {
			case *packet.Packet:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			case *sflow.ExtSwitchData:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			case *sflow.FlowSample:
				if verbose {
					logger.Printf("%#v\n", data)
				}
			}
		}
	}
}
