package mirror

import "encoding/binary"

type UDP struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
}

func (u *UDP) Marshal() []byte {
	b := make([]byte, UDPHLen)

	binary.BigEndian.PutUint16(b[0:], uint16(u.SrcPort))
	binary.BigEndian.PutUint16(b[2:], uint16(u.DstPort))
	binary.BigEndian.PutUint16(b[4:], uint16(UDPHLen+u.Length))
	binary.BigEndian.PutUint16(b[6:], uint16(u.Checksum))

	return b
}
