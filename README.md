![vFlow](docs/imgs/vflow_logo.png?raw=true "vFlow logo")
##
[![Build Status](https://travis-ci.org/VerizonDigital/vflow.svg?branch=master)](https://travis-ci.org/VerizonDigital/vflow) 
[![Go Report Card](https://goreportcard.com/badge/github.com/VerizonDigital/vflow)](https://goreportcard.com/report/github.com/VerizonDigital/vflow)
[![GoDoc](https://godoc.org/github.com/VerizonDigital/vflow?status.svg)](https://godoc.org/github.com/VerizonDigital/vflow)
[![Join Slack](https://img.shields.io/badge/join-slack-9B69A0.svg)](https://join.slack.com/t/vflowworkspace/shared_invite/enQtNTAwNjMyNzg0MDY5LWJlNGExZDNiYThmYjkyNmM1NDAyZGY4NmMyZjYwYmE0ZjAzZjA2MTZlZjRkYjY3Njc1MDJjYTlhZDU2OTk2MGE)

High-performance, scalable and reliable IPFIX, sFlow and Netflow collector (written in pure Golang).

## Features
- IPFIX RFC7011 collector
- sFLow v5 raw header / counters collector
- Netflow v5 collector
- Netflow v9 collector
- Decoding sFlow raw header L2/L3/L4 
- Produce to Apache Kafka, NSQ, NATS
- Replicate IPFIX to 3rd party collector
- Supports IPv4 and IPv6
- Monitoring with InfluxDB and OpenTSDB backend

![Alt text](/docs/imgs/vflow.gif?raw=true "vFlow")

## Documentation
- [Architecture](/docs/design.md).
- [Configuration](/docs/config.md).
- [Quick Start](/docs/quick_start_nsq.md).
- [JUNOS Integration](/docs/junos_integration.md).
- [Monitoring](/monitor/README.md).
- [Stress / Load Generator](/stress/README.md).
- [Kafka consumer examples](https://github.com/VerizonDigital/vflow/tree/master/consumers).

## Decoded IPFIX data
The IPFIX data decodes to JSON format and IDs are [IANA IPFIX element ID](http://www.iana.org/assignments/ipfix/ipfix.xhtml)
```json
{"AgentID":"192.168.21.15","Header":{"Version":10,"Length":420,"ExportTime":1483484642,"SequenceNo":1434533677,"DomainID":32771},"DataSets":[[{"I":8,"V":"192.16.28.217"},{"I":12,"V":"180.10.210.240"},{"I":5,"V":2},{"I":4,"V":6},{"I":7,"V":443},{"I":11,"V":64381},{"I":32,"V":0},{"I":10,"V":811},{"I":58,"V":0},{"I":9,"V":24},{"I":13,"V":20},{"I":16,"V":4200000000},{"I":17,"V":27747},{"I":15,"V":"180.105.10.210"},{"I":6,"V":"0x10"},{"I":14,"V":1113},{"I":1,"V":22500},{"I":2,"V":15},{"I":52,"V":63},{"I":53,"V":63},{"I":152,"V":1483484581770},{"I":153,"V":1483484622384},{"I":136,"V":2},{"I":243,"V":0},{"I":245,"V":0}]]}
```

## Decoded sFlow data
```json
{"Version":5,"IPVersion":1,"AgentSubID":5,"SequenceNo":37591,"SysUpTime":3287084017,"SamplesNo":1,"Samples":[{"SequenceNo":1530345639,"SourceID":0,"SamplingRate":4096,"SamplePool":1938456576,"Drops":0,"Input":536,"Output":728,"RecordsNo":3,"Records":{"ExtRouter":{"NextHop":"115.131.251.90","SrcMask":24,"DstMask":14},"ExtSwitch":{"SrcVlan":0,"SrcPriority":0,"DstVlan":0,"DstPriority":0},"RawHeader":{"L2":{"SrcMAC":"58:00:bb:e7:57:6f","DstMAC":"f4:a7:39:44:a8:27","Vlan":0,"EtherType":2048},"L3":{"Version":4,"TOS":0,"TotalLen":1452,"ID":13515,"Flags":0,"FragOff":0,"TTL":62,"Protocol":6,"Checksum":8564,"Src":"10.1.8.5","Dst":"161.140.24.181"},"L4":{"SrcPort":443,"DstPort":56521,"DataOffset":5,"Reserved":0,"Flags":16}}}}],"IPAddress":"192.168.10.0"}
```
## Decoded Netflow v5 data
``` json
{"AgentID":"114.23.3.231","Header":{"Version":5,"Count":3,"SysUpTimeMSecs":51469784,"UNIXSecs":1544476581,"UNIXNSecs":0,"SeqNum":873873830,"EngType":0,"EngID":0,"SmpInt":1000},"Flows":[{"SrcAddr":"125.238.46.48","DstAddr":"114.23.236.96","NextHop":"114.23.3.231","Input":791,"Output":817,"PktCount":4,"L3Octets":1708,"StartTime":51402145,"EndTime":51433264,"SrcPort":49233,"DstPort":443,"Padding1":0,"TCPFlags":16,"ProtType":6,"Tos":0,"SrcAsNum":4771,"DstAsNum":56030,"SrcMask":20,"DstMask":22,"Padding2":0},{"SrcAddr":"125.238.46.48","DstAddr":"114.23.236.96","NextHop":"114.23.3.231","Input":791,"Output":817,"PktCount":1,"L3Octets":441,"StartTime":51425137,"EndTime":51425137,"SrcPort":49233,"DstPort":443,"Padding1":0,"TCPFlags":24,"ProtType":6,"Tos":0,"SrcAsNum":4771,"DstAsNum":56030,"SrcMask":20,"DstMask":22,"Padding2":0},{"SrcAddr":"210.5.53.48","DstAddr":"103.22.200.210","NextHop":"122.56.118.157","Input":564,"Output":802,"PktCount":1,"L3Octets":1500,"StartTime":51420072,"EndTime":51420072,"SrcPort":80,"DstPort":56108,"Padding1":0,"TCPFlags":16,"ProtType":6,"Tos":0,"SrcAsNum":56030,"DstAsNum":13335,"SrcMask":24,"DstMask":23,"Padding2":0}]}
```
## Decoded Netflow v9 data
```json
{"AgentID":"10.81.70.56","Header":{"Version":9,"Count":1,"SysUpTime":357280,"UNIXSecs":1493918653,"SeqNum":14,"SrcID":87},"DataSets":[[{"I":1,"V":"0x00000050"},{"I":2,"V":"0x00000002"},{"I":4,"V":2},{"I":5,"V":192},{"I":6,"V":"0x00"},{"I":7,"V":0},{"I":8,"V":"10.81.70.56"},{"I":9,"V":0},{"I":10,"V":0},{"I":11,"V":0},{"I":12,"V":"224.0.0.22"},{"I":13,"V":0},{"I":14,"V":0},{"I":15,"V":"0.0.0.0"},{"I":16,"V":0},{"I":17,"V":0},{"I":21,"V":300044},{"I":22,"V":299144}]]}
```

## Supported platform
- Linux
- Windows

## Build
Given that the Go Language compiler (version 1.11 preferred) is installed, you can build it with:
```
go get github.com/VerizonDigital/vflow/vflow
cd $GOPATH/src/github.com/VerizonDigital/vflow

make build
or
go get -d ./...
cd vflow; go build 
```
## Installation
You can download and install pre-built debian package as below ([RPM and Linux binary are available](https://github.com/VerizonDigital/vflow/releases/tag/v0.7.0)). 

dpkg -i [vflow-0.7.0-x86_64.deb](https://github.com/VerizonDigital/vflow/releases/download/v0.7.0/vflow-0.7.0-x86_64.deb)

Once you installed you need to configure the below files, for more information check [configuration guide](/docs/config.md):
```
/etc/vflow/vflow.conf
/etc/vflow/mq.conf
```
You can start the service by the below:
```
service vflow start
```

## Docker
1. Install [Docker](https://www.docker.com/).
2. Download vFlow and Kafka images from public [Docker Hub ](https://hub.docker.com/): 
```
docker pull mehrdadrad/vflow
docker pull spotify/kafka
```
3. You can run them like below:
```
docker run -d -p 2181:2181 -p 9092:9092 spotify/kafka
docker run -d -p 4739:4739 -p 4729:4729 -p 6343:6343 -p 8081:8081 -e VFLOW_KAFKA_BROKERS="172.17.0.1:9092" mehrdadrad/vflow
```

## Slack

You can also join the vFlow Team on Slack [https://vflowworkspace.slack.com](https://join.slack.com/t/vflowworkspace/shared_invite/enQtNTAwNjMyNzg0MDY5LWJlNGExZDNiYThmYjkyNmM1NDAyZGY4NmMyZjYwYmE0ZjAzZjA2MTZlZjRkYjY3Njc1MDJjYTlhZDU2OTk2MGE) and chat with developers.

## License
Licensed under the Apache License, Version 2.0 (the "License")

## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
