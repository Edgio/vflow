//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    flow_counter.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    08/08/2018
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package sflow

import (
	"io"
)

// GenericInterfaceCounters represents Generic Interface Counters RFC2233
type GenericInterfaceCounters struct {
	Index               uint32
	Type                uint32
	Speed               uint64
	Direction           uint32
	Status              uint32
	InOctets            uint64
	InUnicastPackets    uint32
	InMulticastPackets  uint32
	InBroadcastPackets  uint32
	InDiscards          uint32
	InErrors            uint32
	InUnknownProtocols  uint32
	OutOctets           uint64
	OutUnicastPackets   uint32
	OutMulticastPackets uint32
	OutBroadcastPackets uint32
	OutDiscards         uint32
	OutErrors           uint32
	PromiscuousMode     uint32
}

// EthernetInterfaceCounters represents Ethernet Interface Counters RFC2358
type EthernetInterfaceCounters struct {
	AlignmentErrors           uint32
	FCSErrors                 uint32
	SingleCollisionFrames     uint32
	MultipleCollisionFrames   uint32
	SQETestErrors             uint32
	DeferredTransmissions     uint32
	LateCollisions            uint32
	ExcessiveCollisions       uint32
	InternalMACTransmitErrors uint32
	CarrierSenseErrors        uint32
	FrameTooLongs             uint32
	InternalMACReceiveErrors  uint32
	SymbolErrors              uint32
}

// VlanCounters represents VLAN Counters
type VlanCounters struct {
	ID               uint32
	Octets           uint64
	UnicastPackets   uint32
	MulticastPackets uint32
	BroadcastPackets uint32
	Discards         uint32
}

// ProcessorCounters represents Processor Information
type ProcessorCounters struct {
	CPU5s       uint32
	CPU1m       uint32
	CPU5m       uint32
	TotalMemory uint64
	FreeMemory  uint64
}

func (gic *GenericInterfaceCounters) unmarshal(r io.ReadSeeker) error {
	var err error

	fields := []interface{}{
		&gic.Index,
		&gic.Type,
		&gic.Speed,
		&gic.Direction,
		&gic.Status,
		&gic.InOctets,
		&gic.InUnicastPackets,
		&gic.InMulticastPackets,
		&gic.InBroadcastPackets,
		&gic.InDiscards,
		&gic.InErrors,
		&gic.InUnknownProtocols,
		&gic.OutOctets,
		&gic.OutUnicastPackets,
		&gic.OutMulticastPackets,
		&gic.OutBroadcastPackets,
		&gic.OutDiscards,
		&gic.OutErrors,
		&gic.PromiscuousMode,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func (eic *EthernetInterfaceCounters) unmarshal(r io.ReadSeeker) error {
	var err error

	fields := []interface{}{
		&eic.AlignmentErrors,
		&eic.FCSErrors,
		&eic.SingleCollisionFrames,
		&eic.MultipleCollisionFrames,
		&eic.SQETestErrors,
		&eic.DeferredTransmissions,
		&eic.LateCollisions,
		&eic.ExcessiveCollisions,
		&eic.InternalMACTransmitErrors,
		&eic.CarrierSenseErrors,
		&eic.FrameTooLongs,
		&eic.InternalMACReceiveErrors,
		&eic.SymbolErrors,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func (vc *VlanCounters) unmarshal(r io.ReadSeeker) error {
	var err error
	fields := []interface{}{
		&vc.ID,
		&vc.Octets,
		&vc.UnicastPackets,
		&vc.MulticastPackets,
		&vc.BroadcastPackets,
		&vc.Discards,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}

func (pc *ProcessorCounters) unmarshal(r io.ReadSeeker) error {
	var err error
	fields := []interface{}{
		&pc.CPU5s,
		&pc.CPU1m,
		&pc.CPU5m,
		&pc.TotalMemory,
		&pc.FreeMemory,
	}

	for _, field := range fields {
		if err = read(r, field); err != nil {
			return err
		}
	}

	return nil
}
