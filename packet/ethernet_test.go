package packet

import "testing"

func TestDecodeIEEE802(t *testing.T) {
	b := []byte{
		0xd4, 0x4, 0xff, 0x1,
		0x1d, 0x9e, 0x30, 0x7c,
		0x5e, 0xe5, 0x59, 0xef,
		0x8, 0x0, 0x45, 0x0, 0x0,
	}

	d, err := decodeIEEE802(b)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if d.SrcMAC != "d4:04:ff:01:1d:9e" {
		t.Error("expected d4:04:ff:01:1d:9e, got", d.SrcMAC)
	}

	if d.DstMAC != "30:7c:5e:e5:59:ef" {
		t.Error("expected 30:7c:5e:e5:59:ef, got", d.DstMAC)
	}

	if d.Vlan != 0 {
		t.Error("expected 0, got", d.Vlan)
	}

	if d.EtherType != 0x800 {
		t.Error("expected 0x800, got", d.EtherType)
	}
}
