package main

import (
	"net"
	"strconv"
	"sync"
	"time"

	"git.edgecastcdn.net/vflow/ipfix"
)

type IPFIX struct {
	port    int
	addr    string
	udpSize int
	workers int
	stop    bool
}

type IPFIXUDPMsg struct {
	raddr *net.UDPAddr
	body  []byte
}

var (
	ipfixUdpCh = make(chan IPFIXUDPMsg, 1000)
)

func NewIPFIX(opts *Options) *IPFIX {
	return &IPFIX{
		port:    opts.IPFIXPort,
		udpSize: opts.IPFIXUDPSize,
		workers: opts.IPFIXWorkers,
	}
}

func (i *IPFIX) run() {
	var wg sync.WaitGroup

	hostPort := net.JoinHostPort(i.addr, strconv.Itoa(i.port))
	udpAddr, _ := net.ResolveUDPAddr("udp", hostPort)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {

	}

	for n := 0; n < i.workers; n++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			ipfixWorker()

		}()
	}

	logger.Printf("ipfix is running (workers#: %d)", i.workers)

	for !i.stop {
		b := make([]byte, i.udpSize)
		conn.SetReadDeadline(time.Now().Add(1e9))
		n, raddr, err := conn.ReadFromUDP(b)
		if err != nil {
			continue
		}
		ipfixUdpCh <- IPFIXUDPMsg{raddr, b[:n]}
	}

	wg.Wait()
}

func (i *IPFIX) shutdown() {
	i.stop = true
	logger.Println("stopped ipfix service gracefully ...")
	time.Sleep(1 * time.Second)
	logger.Println("ipfix has been shutdown")
	close(ipfixUdpCh)
}

func ipfixWorker() {
	var (
		msg IPFIXUDPMsg
		ok  bool
	)

	for {
		if msg, ok = <-ipfixUdpCh; !ok {
			break
		}

		if verbose {
			logger.Printf("rcvd ipfix data from: %s, size: %d bytes",
				msg.raddr, len(msg.body))
		}

		d := ipfix.NewDecoder(msg.body)
		d.Decode()
	}
}
