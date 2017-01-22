package packet

import "errors"

type ICMP struct {
	Type int
	Code int
}

var errICMPHLenTooSHort = errors.New("ICMP header length is too short")

func decodeICMP(b []byte) (ICMP, error) {
	if len(b) < 4 {
		return ICMP{}, errICMPHLenTooSHort
	}

	return ICMP{
		Type: int(b[0]),
		Code: int(b[1]),
	}, nil
}
