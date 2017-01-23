package packet

import "errors"

type Packet struct {
	L2   Datalink
	L3   interface{}
	L4   interface{}
	data []byte
}

var (
	errUnknownEtherType = errors.New("unknown ether type")
)

func NewPacket() Packet {
	return Packet{}
}

func (p *Packet) Decoder(data []byte) (*Packet, error) {
	var (
		err error
	)

	p.data = data
	err = p.decodeEthernet()
	if err != nil {
		return p, err
	}

	switch p.L2.EtherType {
	case EtherTypeIPv4:

		err = p.decodeIPv4Header()
		if err != nil {
			return p, err
		}

		err = p.decodeNextLayer()
		if err != nil {
			return p, err
		}

	case EtherTypeIPv6:

		err = p.decodeIPv6Header()
		if err != nil {
			return p, err
		}

		err = p.decodeNextLayer()
		if err != nil {
			return p, err
		}

	default:
		return p, errUnknownEtherType
	}

	return p, nil
}
