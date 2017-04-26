//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ipfix.go
//: details: IP Flow Information Export (IPFIX) entities model - https://www.iana.org/assignments/ipfix/ipfix.xhtml
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
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

package ipfix

// FieldType is IPFIX Abstract Data Types RFC5102#section-3.1
type FieldType int

// ElementKey represents field specifier format ids
type ElementKey struct {
	EnterpriseNo uint32
	ElementID    uint16
}

// InfoElementEntry represents standard name and
// type for a field - RFC5102
type InfoElementEntry struct {
	FieldID uint16
	Name    string
	Type    FieldType
}

// IANAInfoModel represents IPFIX field's name, identification and type
type IANAInfoModel map[ElementKey]InfoElementEntry

const (
	// Unknown data type
	Unknown FieldType = iota

	// Uint8 represents a non-negative integer value in the
	// range of 0 to 255.
	Uint8

	// Uint16 represents a non-negative integer value in the
	// range of 0 to 65535.
	Uint16

	// Uint32 represents a non-negative integer value in the
	// range of 0 to 4294967295.
	Uint32

	// Uint64 represents a non-negative integer value in the
	// range of 0 to 18446744073709551615.
	Uint64

	// Int8 represents an integer value in the range of -128
	// to 127.
	Int8

	// Int16 represents an integer value in the range of
	// -32768 to 32767.
	Int16

	// Int32 represents an integer value in the range of
	// -2147483648 to 2147483647.
	Int32

	// Int64 represents an integer value in the range of
	// -9223372036854775808 to 9223372036854775807.
	Int64

	// Float32 corresponds to an IEEE single-precision 32-bit
	// floating point type as defined in [IEEE.754.1985].
	Float32

	// Float64 corresponds to an IEEE double-precision 64-bit
	// floating point type as defined in [IEEE.754.1985].
	Float64

	// Boolean represents a binary value.  The only allowed
	// values are "true" and "false".
	Boolean

	// MacAddress represents a string of 6 octets.
	MacAddress

	// OctetArray represents a finite-length string of octets.
	OctetArray

	// String represents a finite-length string of valid
	String

	// DateTimeSeconds represents a time value in units of
	// seconds based on coordinated universal time (UTC).
	DateTimeSeconds

	// DateTimeMilliseconds represents a time value in units of
	// milliseconds based on coordinated universal time (UTC).
	DateTimeMilliseconds

	// DateTimeMicroseconds represents a time value in units of
	// microseconds based on coordinated universal time (UTC).
	DateTimeMicroseconds

	// DateTimeNanoseconds represents a time value in units of
	// nanoseconds based on coordinated universal time (UTC).
	DateTimeNanoseconds

	// Ipv4Address represents a value of an IPv4 address.
	Ipv4Address

	// Ipv6Address represents a value of an IPv6 address.
	Ipv6Address
)

// FieldTypes represents data types
var FieldTypes = map[string]FieldType{
	"unsigned8":            Uint8,
	"unsigned16":           Uint16,
	"unsigned32":           Uint32,
	"unsigned64":           Uint64,
	"signed8":              Int8,
	"signed16":             Int16,
	"signed32":             Int32,
	"signed64":             Int64,
	"float32":              Float32,
	"float64":              Float64,
	"boolean":              Boolean,
	"macAddress":           MacAddress,
	"octetArray":           OctetArray,
	"string":               String,
	"dateTimeSeconds":      DateTimeSeconds,
	"dateTimeMilliseconds": DateTimeMilliseconds,
	"dateTimeMicroseconds": DateTimeMicroseconds,
	"dateTimeNanoseconds":  DateTimeNanoseconds,
	"ipv4Address":          Ipv4Address,
	"ipv6Address":          Ipv6Address,
}

//InfoModel maps element to name and type based on the field id and enterprise id
var InfoModel = IANAInfoModel{
	ElementKey{0, 1}:   InfoElementEntry{FieldID: 1, Name: "octetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 2}:   InfoElementEntry{FieldID: 2, Name: "packetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 3}:   InfoElementEntry{FieldID: 3, Name: "deltaFlowCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 4}:   InfoElementEntry{FieldID: 4, Name: "protocolIdentifier", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 5}:   InfoElementEntry{FieldID: 5, Name: "ipClassOfService", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 6}:   InfoElementEntry{FieldID: 6, Name: "tcpControlBits", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 7}:   InfoElementEntry{FieldID: 7, Name: "sourceTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 8}:   InfoElementEntry{FieldID: 8, Name: "sourceIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 9}:   InfoElementEntry{FieldID: 9, Name: "sourceIPv4PrefixLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 10}:  InfoElementEntry{FieldID: 10, Name: "ingressInterface", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 11}:  InfoElementEntry{FieldID: 11, Name: "destinationTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 12}:  InfoElementEntry{FieldID: 12, Name: "destinationIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 13}:  InfoElementEntry{FieldID: 13, Name: "destinationIPv4PrefixLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 14}:  InfoElementEntry{FieldID: 14, Name: "egressInterface", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 15}:  InfoElementEntry{FieldID: 15, Name: "ipNextHopIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 16}:  InfoElementEntry{FieldID: 16, Name: "bgpSourceAsNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 17}:  InfoElementEntry{FieldID: 17, Name: "bgpDestinationAsNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 18}:  InfoElementEntry{FieldID: 18, Name: "bgpNextHopIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 19}:  InfoElementEntry{FieldID: 19, Name: "postMCastPacketDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 20}:  InfoElementEntry{FieldID: 20, Name: "postMCastOctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 21}:  InfoElementEntry{FieldID: 21, Name: "flowEndSysUpTime", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 22}:  InfoElementEntry{FieldID: 22, Name: "flowStartSysUpTime", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 23}:  InfoElementEntry{FieldID: 23, Name: "postOctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 24}:  InfoElementEntry{FieldID: 24, Name: "postPacketDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 25}:  InfoElementEntry{FieldID: 25, Name: "minimumIpTotalLength", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 26}:  InfoElementEntry{FieldID: 26, Name: "maximumIpTotalLength", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 27}:  InfoElementEntry{FieldID: 27, Name: "sourceIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 28}:  InfoElementEntry{FieldID: 28, Name: "destinationIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 29}:  InfoElementEntry{FieldID: 29, Name: "sourceIPv6PrefixLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 30}:  InfoElementEntry{FieldID: 30, Name: "destinationIPv6PrefixLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 31}:  InfoElementEntry{FieldID: 31, Name: "flowLabelIPv6", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 32}:  InfoElementEntry{FieldID: 32, Name: "icmpTypeCodeIPv4", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 33}:  InfoElementEntry{FieldID: 33, Name: "igmpType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 34}:  InfoElementEntry{FieldID: 34, Name: "samplingInterval", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 35}:  InfoElementEntry{FieldID: 35, Name: "samplingAlgorithm", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 36}:  InfoElementEntry{FieldID: 36, Name: "flowActiveTimeout", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 37}:  InfoElementEntry{FieldID: 37, Name: "flowIdleTimeout", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 38}:  InfoElementEntry{FieldID: 38, Name: "engineType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 39}:  InfoElementEntry{FieldID: 39, Name: "engineId", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 40}:  InfoElementEntry{FieldID: 40, Name: "exportedOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 41}:  InfoElementEntry{FieldID: 41, Name: "exportedMessageTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 42}:  InfoElementEntry{FieldID: 42, Name: "exportedFlowRecordTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 43}:  InfoElementEntry{FieldID: 43, Name: "ipv4RouterSc", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 44}:  InfoElementEntry{FieldID: 44, Name: "sourceIPv4Prefix", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 45}:  InfoElementEntry{FieldID: 45, Name: "destinationIPv4Prefix", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 46}:  InfoElementEntry{FieldID: 46, Name: "mplsTopLabelType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 47}:  InfoElementEntry{FieldID: 47, Name: "mplsTopLabelIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 48}:  InfoElementEntry{FieldID: 48, Name: "samplerId", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 49}:  InfoElementEntry{FieldID: 49, Name: "samplerMode", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 50}:  InfoElementEntry{FieldID: 50, Name: "samplerRandomInterval", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 51}:  InfoElementEntry{FieldID: 51, Name: "classId", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 52}:  InfoElementEntry{FieldID: 52, Name: "minimumTTL", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 53}:  InfoElementEntry{FieldID: 53, Name: "maximumTTL", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 54}:  InfoElementEntry{FieldID: 54, Name: "fragmentIdentification", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 55}:  InfoElementEntry{FieldID: 55, Name: "postIpClassOfService", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 56}:  InfoElementEntry{FieldID: 56, Name: "sourceMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 57}:  InfoElementEntry{FieldID: 57, Name: "postDestinationMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 58}:  InfoElementEntry{FieldID: 58, Name: "vlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 59}:  InfoElementEntry{FieldID: 59, Name: "postVlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 60}:  InfoElementEntry{FieldID: 60, Name: "ipVersion", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 61}:  InfoElementEntry{FieldID: 61, Name: "flowDirection", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 62}:  InfoElementEntry{FieldID: 62, Name: "ipNextHopIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 63}:  InfoElementEntry{FieldID: 63, Name: "bgpNextHopIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 64}:  InfoElementEntry{FieldID: 64, Name: "ipv6ExtensionHeaders", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 70}:  InfoElementEntry{FieldID: 70, Name: "mplsTopLabelStackSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 71}:  InfoElementEntry{FieldID: 71, Name: "mplsLabelStackSection2", Type: FieldTypes["octetArray"]},
	ElementKey{0, 72}:  InfoElementEntry{FieldID: 72, Name: "mplsLabelStackSection3", Type: FieldTypes["octetArray"]},
	ElementKey{0, 73}:  InfoElementEntry{FieldID: 73, Name: "mplsLabelStackSection4", Type: FieldTypes["octetArray"]},
	ElementKey{0, 74}:  InfoElementEntry{FieldID: 74, Name: "mplsLabelStackSection5", Type: FieldTypes["octetArray"]},
	ElementKey{0, 75}:  InfoElementEntry{FieldID: 75, Name: "mplsLabelStackSection6", Type: FieldTypes["octetArray"]},
	ElementKey{0, 76}:  InfoElementEntry{FieldID: 76, Name: "mplsLabelStackSection7", Type: FieldTypes["octetArray"]},
	ElementKey{0, 77}:  InfoElementEntry{FieldID: 77, Name: "mplsLabelStackSection8", Type: FieldTypes["octetArray"]},
	ElementKey{0, 78}:  InfoElementEntry{FieldID: 78, Name: "mplsLabelStackSection9", Type: FieldTypes["octetArray"]},
	ElementKey{0, 79}:  InfoElementEntry{FieldID: 79, Name: "mplsLabelStackSection10", Type: FieldTypes["octetArray"]},
	ElementKey{0, 80}:  InfoElementEntry{FieldID: 80, Name: "destinationMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 81}:  InfoElementEntry{FieldID: 81, Name: "postSourceMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 82}:  InfoElementEntry{FieldID: 82, Name: "interfaceName", Type: FieldTypes["string"]},
	ElementKey{0, 83}:  InfoElementEntry{FieldID: 83, Name: "interfaceDescription", Type: FieldTypes["string"]},
	ElementKey{0, 84}:  InfoElementEntry{FieldID: 84, Name: "samplerName", Type: FieldTypes["string"]},
	ElementKey{0, 85}:  InfoElementEntry{FieldID: 85, Name: "octetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 86}:  InfoElementEntry{FieldID: 86, Name: "packetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 87}:  InfoElementEntry{FieldID: 87, Name: "flagsAndSamplerId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 88}:  InfoElementEntry{FieldID: 88, Name: "fragmentOffset", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 89}:  InfoElementEntry{FieldID: 89, Name: "forwardingStatus", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 90}:  InfoElementEntry{FieldID: 90, Name: "mplsVpnRouteDistinguisher", Type: FieldTypes["octetArray"]},
	ElementKey{0, 91}:  InfoElementEntry{FieldID: 91, Name: "mplsTopLabelPrefixLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 92}:  InfoElementEntry{FieldID: 92, Name: "srcTrafficIndex", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 93}:  InfoElementEntry{FieldID: 93, Name: "dstTrafficIndex", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 94}:  InfoElementEntry{FieldID: 94, Name: "applicationDescription", Type: FieldTypes["string"]},
	ElementKey{0, 95}:  InfoElementEntry{FieldID: 95, Name: "applicationId", Type: FieldTypes["octetArray"]},
	ElementKey{0, 96}:  InfoElementEntry{FieldID: 96, Name: "applicationName", Type: FieldTypes["string"]},
	ElementKey{0, 98}:  InfoElementEntry{FieldID: 98, Name: "postIpDiffServCodePoint", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 99}:  InfoElementEntry{FieldID: 99, Name: "multicastReplicationFactor", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 100}: InfoElementEntry{FieldID: 100, Name: "className", Type: FieldTypes["string"]},
	ElementKey{0, 101}: InfoElementEntry{FieldID: 101, Name: "classificationEngineId", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 102}: InfoElementEntry{FieldID: 102, Name: "layer2packetSectionOffset", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 103}: InfoElementEntry{FieldID: 103, Name: "layer2packetSectionSize", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 104}: InfoElementEntry{FieldID: 104, Name: "layer2packetSectionData", Type: FieldTypes["octetArray"]},
	ElementKey{0, 128}: InfoElementEntry{FieldID: 128, Name: "bgpNextAdjacentAsNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 129}: InfoElementEntry{FieldID: 129, Name: "bgpPrevAdjacentAsNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 130}: InfoElementEntry{FieldID: 130, Name: "exporterIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 131}: InfoElementEntry{FieldID: 131, Name: "exporterIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 132}: InfoElementEntry{FieldID: 132, Name: "droppedOctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 133}: InfoElementEntry{FieldID: 133, Name: "droppedPacketDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 134}: InfoElementEntry{FieldID: 134, Name: "droppedOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 135}: InfoElementEntry{FieldID: 135, Name: "droppedPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 136}: InfoElementEntry{FieldID: 136, Name: "flowEndReason", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 137}: InfoElementEntry{FieldID: 137, Name: "commonPropertiesId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 138}: InfoElementEntry{FieldID: 138, Name: "observationPointId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 139}: InfoElementEntry{FieldID: 139, Name: "icmpTypeCodeIPv6", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 140}: InfoElementEntry{FieldID: 140, Name: "mplsTopLabelIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 141}: InfoElementEntry{FieldID: 141, Name: "lineCardId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 142}: InfoElementEntry{FieldID: 142, Name: "portId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 143}: InfoElementEntry{FieldID: 143, Name: "meteringProcessId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 144}: InfoElementEntry{FieldID: 144, Name: "exportingProcessId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 145}: InfoElementEntry{FieldID: 145, Name: "templateId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 146}: InfoElementEntry{FieldID: 146, Name: "wlanChannelId", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 147}: InfoElementEntry{FieldID: 147, Name: "wlanSSID", Type: FieldTypes["string"]},
	ElementKey{0, 148}: InfoElementEntry{FieldID: 148, Name: "flowId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 149}: InfoElementEntry{FieldID: 149, Name: "observationDomainId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 150}: InfoElementEntry{FieldID: 150, Name: "flowStartSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 151}: InfoElementEntry{FieldID: 151, Name: "flowEndSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 152}: InfoElementEntry{FieldID: 152, Name: "flowStartMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 153}: InfoElementEntry{FieldID: 153, Name: "flowEndMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 154}: InfoElementEntry{FieldID: 154, Name: "flowStartMicroseconds", Type: FieldTypes["dateTimeMicroseconds"]},
	ElementKey{0, 155}: InfoElementEntry{FieldID: 155, Name: "flowEndMicroseconds", Type: FieldTypes["dateTimeMicroseconds"]},
	ElementKey{0, 156}: InfoElementEntry{FieldID: 156, Name: "flowStartNanoseconds", Type: FieldTypes["dateTimeNanoseconds"]},
	ElementKey{0, 157}: InfoElementEntry{FieldID: 157, Name: "flowEndNanoseconds", Type: FieldTypes["dateTimeNanoseconds"]},
	ElementKey{0, 158}: InfoElementEntry{FieldID: 158, Name: "flowStartDeltaMicroseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 159}: InfoElementEntry{FieldID: 159, Name: "flowEndDeltaMicroseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 160}: InfoElementEntry{FieldID: 160, Name: "systemInitTimeMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 161}: InfoElementEntry{FieldID: 161, Name: "flowDurationMilliseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 162}: InfoElementEntry{FieldID: 162, Name: "flowDurationMicroseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 163}: InfoElementEntry{FieldID: 163, Name: "observedFlowTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 164}: InfoElementEntry{FieldID: 164, Name: "ignoredPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 165}: InfoElementEntry{FieldID: 165, Name: "ignoredOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 166}: InfoElementEntry{FieldID: 166, Name: "notSentFlowTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 167}: InfoElementEntry{FieldID: 167, Name: "notSentPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 168}: InfoElementEntry{FieldID: 168, Name: "notSentOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 169}: InfoElementEntry{FieldID: 169, Name: "destinationIPv6Prefix", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 170}: InfoElementEntry{FieldID: 170, Name: "sourceIPv6Prefix", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 171}: InfoElementEntry{FieldID: 171, Name: "postOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 172}: InfoElementEntry{FieldID: 172, Name: "postPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 173}: InfoElementEntry{FieldID: 173, Name: "flowKeyIndicator", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 174}: InfoElementEntry{FieldID: 174, Name: "postMCastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 175}: InfoElementEntry{FieldID: 175, Name: "postMCastOctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 176}: InfoElementEntry{FieldID: 176, Name: "icmpTypeIPv4", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 177}: InfoElementEntry{FieldID: 177, Name: "icmpCodeIPv4", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 178}: InfoElementEntry{FieldID: 178, Name: "icmpTypeIPv6", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 179}: InfoElementEntry{FieldID: 179, Name: "icmpCodeIPv6", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 180}: InfoElementEntry{FieldID: 180, Name: "udpSourcePort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 181}: InfoElementEntry{FieldID: 181, Name: "udpDestinationPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 182}: InfoElementEntry{FieldID: 182, Name: "tcpSourcePort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 183}: InfoElementEntry{FieldID: 183, Name: "tcpDestinationPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 184}: InfoElementEntry{FieldID: 184, Name: "tcpSequenceNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 185}: InfoElementEntry{FieldID: 185, Name: "tcpAcknowledgementNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 186}: InfoElementEntry{FieldID: 186, Name: "tcpWindowSize", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 187}: InfoElementEntry{FieldID: 187, Name: "tcpUrgentPointer", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 188}: InfoElementEntry{FieldID: 188, Name: "tcpHeaderLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 189}: InfoElementEntry{FieldID: 189, Name: "ipHeaderLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 190}: InfoElementEntry{FieldID: 190, Name: "totalLengthIPv4", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 191}: InfoElementEntry{FieldID: 191, Name: "payloadLengthIPv6", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 192}: InfoElementEntry{FieldID: 192, Name: "ipTTL", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 193}: InfoElementEntry{FieldID: 193, Name: "nextHeaderIPv6", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 194}: InfoElementEntry{FieldID: 194, Name: "mplsPayloadLength", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 195}: InfoElementEntry{FieldID: 195, Name: "ipDiffServCodePoint", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 196}: InfoElementEntry{FieldID: 196, Name: "ipPrecedence", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 197}: InfoElementEntry{FieldID: 197, Name: "fragmentFlags", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 198}: InfoElementEntry{FieldID: 198, Name: "octetDeltaSumOfSquares", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 199}: InfoElementEntry{FieldID: 199, Name: "octetTotalSumOfSquares", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 200}: InfoElementEntry{FieldID: 200, Name: "mplsTopLabelTTL", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 201}: InfoElementEntry{FieldID: 201, Name: "mplsLabelStackLength", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 202}: InfoElementEntry{FieldID: 202, Name: "mplsLabelStackDepth", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 203}: InfoElementEntry{FieldID: 203, Name: "mplsTopLabelExp", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 204}: InfoElementEntry{FieldID: 204, Name: "ipPayloadLength", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 205}: InfoElementEntry{FieldID: 205, Name: "udpMessageLength", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 206}: InfoElementEntry{FieldID: 206, Name: "isMulticast", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 207}: InfoElementEntry{FieldID: 207, Name: "ipv4IHL", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 208}: InfoElementEntry{FieldID: 208, Name: "ipv4Options", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 209}: InfoElementEntry{FieldID: 209, Name: "tcpOptions", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 210}: InfoElementEntry{FieldID: 210, Name: "paddingOctets", Type: FieldTypes["octetArray"]},
	ElementKey{0, 211}: InfoElementEntry{FieldID: 211, Name: "collectorIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 212}: InfoElementEntry{FieldID: 212, Name: "collectorIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 213}: InfoElementEntry{FieldID: 213, Name: "exportInterface", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 214}: InfoElementEntry{FieldID: 214, Name: "exportProtocolVersion", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 215}: InfoElementEntry{FieldID: 215, Name: "exportTransportProtocol", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 216}: InfoElementEntry{FieldID: 216, Name: "collectorTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 217}: InfoElementEntry{FieldID: 217, Name: "exporterTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 218}: InfoElementEntry{FieldID: 218, Name: "tcpSynTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 219}: InfoElementEntry{FieldID: 219, Name: "tcpFinTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 220}: InfoElementEntry{FieldID: 220, Name: "tcpRstTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 221}: InfoElementEntry{FieldID: 221, Name: "tcpPshTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 222}: InfoElementEntry{FieldID: 222, Name: "tcpAckTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 223}: InfoElementEntry{FieldID: 223, Name: "tcpUrgTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 224}: InfoElementEntry{FieldID: 224, Name: "ipTotalLength", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 225}: InfoElementEntry{FieldID: 225, Name: "postNATSourceIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 226}: InfoElementEntry{FieldID: 226, Name: "postNATDestinationIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 227}: InfoElementEntry{FieldID: 227, Name: "postNAPTSourceTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 228}: InfoElementEntry{FieldID: 228, Name: "postNAPTDestinationTransportPort", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 229}: InfoElementEntry{FieldID: 229, Name: "natOriginatingAddressRealm", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 230}: InfoElementEntry{FieldID: 230, Name: "natEvent", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 231}: InfoElementEntry{FieldID: 231, Name: "initiatorOctets", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 232}: InfoElementEntry{FieldID: 232, Name: "responderOctets", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 233}: InfoElementEntry{FieldID: 233, Name: "firewallEvent", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 234}: InfoElementEntry{FieldID: 234, Name: "ingressVRFID", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 235}: InfoElementEntry{FieldID: 235, Name: "egressVRFID", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 236}: InfoElementEntry{FieldID: 236, Name: "VRFname", Type: FieldTypes["string"]},
	ElementKey{0, 237}: InfoElementEntry{FieldID: 237, Name: "postMplsTopLabelExp", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 238}: InfoElementEntry{FieldID: 238, Name: "tcpWindowScale", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 239}: InfoElementEntry{FieldID: 239, Name: "biflowDirection", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 240}: InfoElementEntry{FieldID: 240, Name: "ethernetHeaderLength", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 241}: InfoElementEntry{FieldID: 241, Name: "ethernetPayloadLength", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 242}: InfoElementEntry{FieldID: 242, Name: "ethernetTotalLength", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 243}: InfoElementEntry{FieldID: 243, Name: "dot1qVlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 244}: InfoElementEntry{FieldID: 244, Name: "dot1qPriority", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 245}: InfoElementEntry{FieldID: 245, Name: "dot1qCustomerVlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 246}: InfoElementEntry{FieldID: 246, Name: "dot1qCustomerPriority", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 247}: InfoElementEntry{FieldID: 247, Name: "metroEvcId", Type: FieldTypes["string"]},
	ElementKey{0, 248}: InfoElementEntry{FieldID: 248, Name: "metroEvcType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 249}: InfoElementEntry{FieldID: 249, Name: "pseudoWireId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 250}: InfoElementEntry{FieldID: 250, Name: "pseudoWireType", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 251}: InfoElementEntry{FieldID: 251, Name: "pseudoWireControlWord", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 252}: InfoElementEntry{FieldID: 252, Name: "ingressPhysicalInterface", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 253}: InfoElementEntry{FieldID: 253, Name: "egressPhysicalInterface", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 254}: InfoElementEntry{FieldID: 254, Name: "postDot1qVlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 255}: InfoElementEntry{FieldID: 255, Name: "postDot1qCustomerVlanId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 256}: InfoElementEntry{FieldID: 256, Name: "ethernetType", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 257}: InfoElementEntry{FieldID: 257, Name: "postIpPrecedence", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 258}: InfoElementEntry{FieldID: 258, Name: "collectionTimeMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 259}: InfoElementEntry{FieldID: 259, Name: "exportSctpStreamId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 260}: InfoElementEntry{FieldID: 260, Name: "maxExportSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 261}: InfoElementEntry{FieldID: 261, Name: "maxFlowEndSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 262}: InfoElementEntry{FieldID: 262, Name: "messageMD5Checksum", Type: FieldTypes["octetArray"]},
	ElementKey{0, 263}: InfoElementEntry{FieldID: 263, Name: "messageScope", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 264}: InfoElementEntry{FieldID: 264, Name: "minExportSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 265}: InfoElementEntry{FieldID: 265, Name: "minFlowStartSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 266}: InfoElementEntry{FieldID: 266, Name: "opaqueOctets", Type: FieldTypes["octetArray"]},
	ElementKey{0, 267}: InfoElementEntry{FieldID: 267, Name: "sessionScope", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 268}: InfoElementEntry{FieldID: 268, Name: "maxFlowEndMicroseconds", Type: FieldTypes["dateTimeMicroseconds"]},
	ElementKey{0, 269}: InfoElementEntry{FieldID: 269, Name: "maxFlowEndMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 270}: InfoElementEntry{FieldID: 270, Name: "maxFlowEndNanoseconds", Type: FieldTypes["dateTimeNanoseconds"]},
	ElementKey{0, 271}: InfoElementEntry{FieldID: 271, Name: "minFlowStartMicroseconds", Type: FieldTypes["dateTimeMicroseconds"]},
	ElementKey{0, 272}: InfoElementEntry{FieldID: 272, Name: "minFlowStartMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 273}: InfoElementEntry{FieldID: 273, Name: "minFlowStartNanoseconds", Type: FieldTypes["dateTimeNanoseconds"]},
	ElementKey{0, 274}: InfoElementEntry{FieldID: 274, Name: "collectorCertificate", Type: FieldTypes["octetArray"]},
	ElementKey{0, 275}: InfoElementEntry{FieldID: 275, Name: "exporterCertificate", Type: FieldTypes["octetArray"]},
	ElementKey{0, 276}: InfoElementEntry{FieldID: 276, Name: "dataRecordsReliability", Type: FieldTypes["boolean"]},
	ElementKey{0, 277}: InfoElementEntry{FieldID: 277, Name: "observationPointType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 278}: InfoElementEntry{FieldID: 278, Name: "newConnectionDeltaCount", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 279}: InfoElementEntry{FieldID: 279, Name: "connectionSumDurationSeconds", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 280}: InfoElementEntry{FieldID: 280, Name: "connectionTransactionId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 281}: InfoElementEntry{FieldID: 281, Name: "postNATSourceIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 282}: InfoElementEntry{FieldID: 282, Name: "postNATDestinationIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 283}: InfoElementEntry{FieldID: 283, Name: "natPoolId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 284}: InfoElementEntry{FieldID: 284, Name: "natPoolName", Type: FieldTypes["string"]},
	ElementKey{0, 285}: InfoElementEntry{FieldID: 285, Name: "anonymizationFlags", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 286}: InfoElementEntry{FieldID: 286, Name: "anonymizationTechnique", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 287}: InfoElementEntry{FieldID: 287, Name: "informationElementIndex", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 288}: InfoElementEntry{FieldID: 288, Name: "p2pTechnology", Type: FieldTypes["string"]},
	ElementKey{0, 289}: InfoElementEntry{FieldID: 289, Name: "tunnelTechnology", Type: FieldTypes["string"]},
	ElementKey{0, 290}: InfoElementEntry{FieldID: 290, Name: "encryptedTechnology", Type: FieldTypes["string"]},
	ElementKey{0, 291}: InfoElementEntry{FieldID: 291, Name: "basicList", Type: FieldTypes["basicList"]},
	ElementKey{0, 292}: InfoElementEntry{FieldID: 292, Name: "subTemplateList", Type: FieldTypes["subTemplateList"]},
	ElementKey{0, 293}: InfoElementEntry{FieldID: 293, Name: "subTemplateMultiList", Type: FieldTypes["subTemplateMultiList"]},
	ElementKey{0, 294}: InfoElementEntry{FieldID: 294, Name: "bgpValidityState", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 295}: InfoElementEntry{FieldID: 295, Name: "IPSecSPI", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 296}: InfoElementEntry{FieldID: 296, Name: "greKey", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 297}: InfoElementEntry{FieldID: 297, Name: "natType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 298}: InfoElementEntry{FieldID: 298, Name: "initiatorPackets", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 299}: InfoElementEntry{FieldID: 299, Name: "responderPackets", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 300}: InfoElementEntry{FieldID: 300, Name: "observationDomainName", Type: FieldTypes["string"]},
	ElementKey{0, 301}: InfoElementEntry{FieldID: 301, Name: "selectionSequenceId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 302}: InfoElementEntry{FieldID: 302, Name: "selectorId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 303}: InfoElementEntry{FieldID: 303, Name: "informationElementId", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 304}: InfoElementEntry{FieldID: 304, Name: "selectorAlgorithm", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 305}: InfoElementEntry{FieldID: 305, Name: "samplingPacketInterval", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 306}: InfoElementEntry{FieldID: 306, Name: "samplingPacketSpace", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 307}: InfoElementEntry{FieldID: 307, Name: "samplingTimeInterval", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 308}: InfoElementEntry{FieldID: 308, Name: "samplingTimeSpace", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 309}: InfoElementEntry{FieldID: 309, Name: "samplingSize", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 310}: InfoElementEntry{FieldID: 310, Name: "samplingPopulation", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 311}: InfoElementEntry{FieldID: 311, Name: "samplingProbability", Type: FieldTypes["float64"]},
	ElementKey{0, 312}: InfoElementEntry{FieldID: 312, Name: "dataLinkFrameSize", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 313}: InfoElementEntry{FieldID: 313, Name: "ipHeaderPacketSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 314}: InfoElementEntry{FieldID: 314, Name: "ipPayloadPacketSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 315}: InfoElementEntry{FieldID: 315, Name: "dataLinkFrameSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 316}: InfoElementEntry{FieldID: 316, Name: "mplsLabelStackSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 317}: InfoElementEntry{FieldID: 317, Name: "mplsPayloadPacketSection", Type: FieldTypes["octetArray"]},
	ElementKey{0, 318}: InfoElementEntry{FieldID: 318, Name: "selectorIdTotalPktsObserved", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 319}: InfoElementEntry{FieldID: 319, Name: "selectorIdTotalPktsSelected", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 320}: InfoElementEntry{FieldID: 320, Name: "absoluteError", Type: FieldTypes["float64"]},
	ElementKey{0, 321}: InfoElementEntry{FieldID: 321, Name: "relativeError", Type: FieldTypes["float64"]},
	ElementKey{0, 322}: InfoElementEntry{FieldID: 322, Name: "observationTimeSeconds", Type: FieldTypes["dateTimeSeconds"]},
	ElementKey{0, 323}: InfoElementEntry{FieldID: 323, Name: "observationTimeMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 324}: InfoElementEntry{FieldID: 324, Name: "observationTimeMicroseconds", Type: FieldTypes["dateTimeMicroseconds"]},
	ElementKey{0, 325}: InfoElementEntry{FieldID: 325, Name: "observationTimeNanoseconds", Type: FieldTypes["dateTimeNanoseconds"]},
	ElementKey{0, 326}: InfoElementEntry{FieldID: 326, Name: "digestHashValue", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 327}: InfoElementEntry{FieldID: 327, Name: "hashIPPayloadOffset", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 328}: InfoElementEntry{FieldID: 328, Name: "hashIPPayloadSize", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 329}: InfoElementEntry{FieldID: 329, Name: "hashOutputRangeMin", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 330}: InfoElementEntry{FieldID: 330, Name: "hashOutputRangeMax", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 331}: InfoElementEntry{FieldID: 331, Name: "hashSelectedRangeMin", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 332}: InfoElementEntry{FieldID: 332, Name: "hashSelectedRangeMax", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 333}: InfoElementEntry{FieldID: 333, Name: "hashDigestOutput", Type: FieldTypes["boolean"]},
	ElementKey{0, 334}: InfoElementEntry{FieldID: 334, Name: "hashInitialiserValue", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 335}: InfoElementEntry{FieldID: 335, Name: "selectorName", Type: FieldTypes["string"]},
	ElementKey{0, 336}: InfoElementEntry{FieldID: 336, Name: "upperCILimit", Type: FieldTypes["float64"]},
	ElementKey{0, 337}: InfoElementEntry{FieldID: 337, Name: "lowerCILimit", Type: FieldTypes["float64"]},
	ElementKey{0, 338}: InfoElementEntry{FieldID: 338, Name: "confidenceLevel", Type: FieldTypes["float64"]},
	ElementKey{0, 339}: InfoElementEntry{FieldID: 339, Name: "informationElementDataType", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 340}: InfoElementEntry{FieldID: 340, Name: "informationElementDescription", Type: FieldTypes["string"]},
	ElementKey{0, 341}: InfoElementEntry{FieldID: 341, Name: "informationElementName", Type: FieldTypes["string"]},
	ElementKey{0, 342}: InfoElementEntry{FieldID: 342, Name: "informationElementRangeBegin", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 343}: InfoElementEntry{FieldID: 343, Name: "informationElementRangeEnd", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 344}: InfoElementEntry{FieldID: 344, Name: "informationElementSemantics", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 345}: InfoElementEntry{FieldID: 345, Name: "informationElementUnits", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 346}: InfoElementEntry{FieldID: 346, Name: "privateEnterpriseNumber", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 347}: InfoElementEntry{FieldID: 347, Name: "virtualStationInterfaceId", Type: FieldTypes["octetArray"]},
	ElementKey{0, 348}: InfoElementEntry{FieldID: 348, Name: "virtualStationInterfaceName", Type: FieldTypes["string"]},
	ElementKey{0, 349}: InfoElementEntry{FieldID: 349, Name: "virtualStationUUID", Type: FieldTypes["octetArray"]},
	ElementKey{0, 350}: InfoElementEntry{FieldID: 350, Name: "virtualStationName", Type: FieldTypes["string"]},
	ElementKey{0, 351}: InfoElementEntry{FieldID: 351, Name: "layer2SegmentId", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 352}: InfoElementEntry{FieldID: 352, Name: "layer2OctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 353}: InfoElementEntry{FieldID: 353, Name: "layer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 354}: InfoElementEntry{FieldID: 354, Name: "ingressUnicastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 355}: InfoElementEntry{FieldID: 355, Name: "ingressMulticastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 356}: InfoElementEntry{FieldID: 356, Name: "ingressBroadcastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 357}: InfoElementEntry{FieldID: 357, Name: "egressUnicastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 358}: InfoElementEntry{FieldID: 358, Name: "egressBroadcastPacketTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 359}: InfoElementEntry{FieldID: 359, Name: "monitoringIntervalStartMilliSeconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 360}: InfoElementEntry{FieldID: 360, Name: "monitoringIntervalEndMilliSeconds", Type: FieldTypes["dateTimeMilliseconds"]},
	ElementKey{0, 361}: InfoElementEntry{FieldID: 361, Name: "portRangeStart", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 362}: InfoElementEntry{FieldID: 362, Name: "portRangeEnd", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 363}: InfoElementEntry{FieldID: 363, Name: "portRangeStepSize", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 364}: InfoElementEntry{FieldID: 364, Name: "portRangeNumPorts", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 365}: InfoElementEntry{FieldID: 365, Name: "staMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 366}: InfoElementEntry{FieldID: 366, Name: "staIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 367}: InfoElementEntry{FieldID: 367, Name: "wtpMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 368}: InfoElementEntry{FieldID: 368, Name: "ingressInterfaceType", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 369}: InfoElementEntry{FieldID: 369, Name: "egressInterfaceType", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 370}: InfoElementEntry{FieldID: 370, Name: "rtpSequenceNumber", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 371}: InfoElementEntry{FieldID: 371, Name: "userName", Type: FieldTypes["string"]},
	ElementKey{0, 372}: InfoElementEntry{FieldID: 372, Name: "applicationCategoryName", Type: FieldTypes["string"]},
	ElementKey{0, 373}: InfoElementEntry{FieldID: 373, Name: "applicationSubCategoryName", Type: FieldTypes["string"]},
	ElementKey{0, 374}: InfoElementEntry{FieldID: 374, Name: "applicationGroupName", Type: FieldTypes["string"]},
	ElementKey{0, 375}: InfoElementEntry{FieldID: 375, Name: "originalFlowsPresent", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 376}: InfoElementEntry{FieldID: 376, Name: "originalFlowsInitiated", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 377}: InfoElementEntry{FieldID: 377, Name: "originalFlowsCompleted", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 378}: InfoElementEntry{FieldID: 378, Name: "distinctCountOfSourceIPAddress", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 379}: InfoElementEntry{FieldID: 379, Name: "distinctCountOfDestinationIPAddress", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 380}: InfoElementEntry{FieldID: 380, Name: "distinctCountOfSourceIPv4Address", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 381}: InfoElementEntry{FieldID: 381, Name: "distinctCountOfDestinationIPv4Address", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 382}: InfoElementEntry{FieldID: 382, Name: "distinctCountOfSourceIPv6Address", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 383}: InfoElementEntry{FieldID: 383, Name: "distinctCountOfDestinationIPv6Address", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 384}: InfoElementEntry{FieldID: 384, Name: "valueDistributionMethod", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 385}: InfoElementEntry{FieldID: 385, Name: "rfc3550JitterMilliseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 386}: InfoElementEntry{FieldID: 386, Name: "rfc3550JitterMicroseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 387}: InfoElementEntry{FieldID: 387, Name: "rfc3550JitterNanoseconds", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 388}: InfoElementEntry{FieldID: 388, Name: "dot1qDEI", Type: FieldTypes["boolean"]},
	ElementKey{0, 389}: InfoElementEntry{FieldID: 389, Name: "dot1qCustomerDEI", Type: FieldTypes["boolean"]},
	ElementKey{0, 390}: InfoElementEntry{FieldID: 390, Name: "flowSelectorAlgorithm", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 391}: InfoElementEntry{FieldID: 391, Name: "flowSelectedOctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 392}: InfoElementEntry{FieldID: 392, Name: "flowSelectedPacketDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 393}: InfoElementEntry{FieldID: 393, Name: "flowSelectedFlowDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 394}: InfoElementEntry{FieldID: 394, Name: "selectorIDTotalFlowsObserved", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 395}: InfoElementEntry{FieldID: 395, Name: "selectorIDTotalFlowsSelected", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 396}: InfoElementEntry{FieldID: 396, Name: "samplingFlowInterval", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 397}: InfoElementEntry{FieldID: 397, Name: "samplingFlowSpacing", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 398}: InfoElementEntry{FieldID: 398, Name: "flowSamplingTimeInterval", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 399}: InfoElementEntry{FieldID: 399, Name: "flowSamplingTimeSpacing", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 400}: InfoElementEntry{FieldID: 400, Name: "hashFlowDomain", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 401}: InfoElementEntry{FieldID: 401, Name: "transportOctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 402}: InfoElementEntry{FieldID: 402, Name: "transportPacketDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 403}: InfoElementEntry{FieldID: 403, Name: "originalExporterIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 404}: InfoElementEntry{FieldID: 404, Name: "originalExporterIPv6Address", Type: FieldTypes["ipv6Address"]},
	ElementKey{0, 405}: InfoElementEntry{FieldID: 405, Name: "originalObservationDomainId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 406}: InfoElementEntry{FieldID: 406, Name: "intermediateProcessId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 407}: InfoElementEntry{FieldID: 407, Name: "ignoredDataRecordTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 408}: InfoElementEntry{FieldID: 408, Name: "dataLinkFrameType", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 409}: InfoElementEntry{FieldID: 409, Name: "sectionOffset", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 410}: InfoElementEntry{FieldID: 410, Name: "sectionExportedOctets", Type: FieldTypes["unsigned16"]},
	ElementKey{0, 411}: InfoElementEntry{FieldID: 411, Name: "dot1qServiceInstanceTag", Type: FieldTypes["octetArray"]},
	ElementKey{0, 412}: InfoElementEntry{FieldID: 412, Name: "dot1qServiceInstanceId", Type: FieldTypes["unsigned32"]},
	ElementKey{0, 413}: InfoElementEntry{FieldID: 413, Name: "dot1qServiceInstancePriority", Type: FieldTypes["unsigned8"]},
	ElementKey{0, 414}: InfoElementEntry{FieldID: 414, Name: "dot1qCustomerSourceMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 415}: InfoElementEntry{FieldID: 415, Name: "dot1qCustomerDestinationMacAddress", Type: FieldTypes["macAddress"]},
	ElementKey{0, 417}: InfoElementEntry{FieldID: 417, Name: "postLayer2OctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 418}: InfoElementEntry{FieldID: 418, Name: "postMCastLayer2OctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 420}: InfoElementEntry{FieldID: 420, Name: "postLayer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 421}: InfoElementEntry{FieldID: 421, Name: "postMCastLayer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 422}: InfoElementEntry{FieldID: 422, Name: "minimumLayer2TotalLength", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 423}: InfoElementEntry{FieldID: 423, Name: "maximumLayer2TotalLength", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 424}: InfoElementEntry{FieldID: 424, Name: "droppedLayer2OctetDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 425}: InfoElementEntry{FieldID: 425, Name: "droppedLayer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 426}: InfoElementEntry{FieldID: 426, Name: "ignoredLayer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 427}: InfoElementEntry{FieldID: 427, Name: "notSentLayer2OctetTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 428}: InfoElementEntry{FieldID: 428, Name: "layer2OctetDeltaSumOfSquares", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 429}: InfoElementEntry{FieldID: 429, Name: "layer2OctetTotalSumOfSquares", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 430}: InfoElementEntry{FieldID: 430, Name: "layer2FrameDeltaCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 431}: InfoElementEntry{FieldID: 431, Name: "layer2FrameTotalCount", Type: FieldTypes["unsigned64"]},
	ElementKey{0, 432}: InfoElementEntry{FieldID: 432, Name: "pseudoWireDestinationIPv4Address", Type: FieldTypes["ipv4Address"]},
	ElementKey{0, 433}: InfoElementEntry{FieldID: 433, Name: "ignoredLayer2FrameTotalCount", Type: FieldTypes["unsigned64"]},
}
