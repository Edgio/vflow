package sflow

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	SampleDataFormat  = 1 // Packet Flow Sampling
	CounterDataFormat = 2 // Counter Sampling
)

type SFDecoder struct {
	reader io.ReadSeeker
	filter []uint32 // Filter data format(s)
}

type SFDatagram struct {
	Version    uint32 // Datagram version
	IPVersion  uint32 // Data gram sFlow version
	AgentSubId uint32 // Identifies a source of sFlow data
	SequenceNo uint32 // Sequence of sFlow Datagrams
	SysUpTime  uint32 // Current time (in milliseconds since device last booted
	SamplesNo  uint32 // Number of samples

	IPAddress net.IP // Agent IP address
}

type SFSampledHeader struct {
	HeaderProtocol uint32 // (enum SFHeaderProtocol)
	FrameLength    uint32 // Original length of packet before sampling
	Stripped       uint32 // Header/trailer bytes stripped by sender
	HeaderLength   uint32 // Length of sampled header bytes to follow
	HeaderBytes    []byte // Header bytes
}

type SFSample interface {
}

func NewSFDecoder(r io.ReadSeeker, f []uint32) SFDecoder {
	return SFDecoder{
		reader: r,
		filter: f,
	}
}

func (d *SFDecoder) SFDecode() (*SFDatagram, error) {
	var (
		datagram     = &SFDatagram{}
		formatSample uint32
		lengthSample uint32
		ipLen        int = 4
		err          error
	)

	if err = binary.Read(d.reader, binary.BigEndian, &datagram.Version); err != nil {
		return nil, err
	}

	if datagram.Version != 5 {
		return nil, fmt.Errorf("sflow version doesn't support")
	}

	if err = binary.Read(d.reader, binary.BigEndian, &datagram.IPVersion); err != nil {
		return nil, err
	}

	// read the agent ip address
	if datagram.IPVersion == 2 {
		ipLen = 16
	}
	buff := make([]byte, ipLen)
	if _, err = d.reader.Read(buff); err != nil {
		return nil, err
	}
	datagram.IPAddress = buff

	if err = binary.Read(d.reader, binary.BigEndian, &datagram.AgentSubId); err != nil {
		return nil, err
	}
	if err = binary.Read(d.reader, binary.BigEndian, &datagram.SequenceNo); err != nil {
		return nil, err
	}
	if err = binary.Read(d.reader, binary.BigEndian, &datagram.SysUpTime); err != nil {
		return nil, err
	}
	if err = binary.Read(d.reader, binary.BigEndian, &datagram.SamplesNo); err != nil {
		return nil, err
	}

	// decode sample(s)
	for i := uint32(0); i < datagram.SamplesNo; i++ {
		if err = binary.Read(d.reader, binary.BigEndian, &formatSample); err != nil {
			return nil, err
		}
		if err = binary.Read(d.reader, binary.BigEndian, &lengthSample); err != nil {
			return nil, err
		}

		if m := d.isFilterMatch(formatSample); m {
			d.reader.Seek(int64(lengthSample), 1)
			continue
		}

		switch formatSample {
		case SampleDataFormat:
			// TODO
			d.reader.Seek(int64(lengthSample), 1)
		case CounterDataFormat:
			// TODO
		default:
			d.reader.Seek(int64(lengthSample), 1)

		}
	}

	return datagram, nil
}

func (d *SFDecoder) isFilterMatch(f uint32) bool {
	for _, v := range d.filter {
		if v == f {
			return true
		}
	}
	return false
}
