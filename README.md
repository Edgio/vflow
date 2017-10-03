![vFlow](docs/imgs/vflow_logo.png?raw=true "vFlow logo")
##
[![Build Status](https://travis-ci.org/VerizonDigital/vflow.svg?branch=master)](https://travis-ci.org/VerizonDigital/vflow) 
[![Go Report Card](https://goreportcard.com/badge/github.com/VerizonDigital/vflow)](https://goreportcard.com/report/github.com/VerizonDigital/vflow)
[![GoDoc](https://godoc.org/github.com/VerizonDigital/vflow?status.svg)](https://godoc.org/github.com/VerizonDigital/vflow)

High-performance, scalable and reliable IPFIX, sFlow and Netflow collector (written in pure Golang).

## Features
- IPFIX RFC7011 collector
- sFLow v5 raw header packet collector
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
{"Header":{"Version":5,"IPVersion":1,"AgentSubID":0,"SequenceNo":24324,"SysUpTime":766903208,"SamplesNo":1,"IPAddress":"192.16.14.0"},"ExtSWData":{"SrcVlan":0,"SrcPriority":0,"DstVlan":12,"DstPriority":0},"Sample":{"SequenceNo":0,"SourceID":0,"SamplingRate":2000,"SamplePool":0,"Drops":0,"Input":552,"Output":0,"RecordsNo":2},"Packet":{"L2":{"SrcMAC":"d4:04:ff:01:1d:9e","DstMAC":"30:7c:5e:e5:59:ef","Vlan":12,"EtherType":34525},"L3":{"Version":6,"TrafficClass":0,"FlowLabel":0,"PayloadLen":265,"NextHeader":17,"HopLimit":57,"Src":"2600:8000:5207:6f00::1","Dst":"2606:2800:404e:2:1663:6fe:2cc6:100a"},"L4":{"SrcPort":53,"DstPort":34234}}}
```
## Decoded Netflow v9 data
```json
{"AgentID":"10.81.70.56","Header":{"Version":9,"Count":1,"SysUpTime":357280,"UNIXSecs":1493918653,"SeqNum":14,"SrcID":87},"DataSets":[[{"I":1,"V":"0x00000050"},{"I":2,"V":"0x00000002"},{"I":4,"V":2},{"I":5,"V":192},{"I":6,"V":"0x00"},{"I":7,"V":0},{"I":8,"V":"10.81.70.56"},{"I":9,"V":0},{"I":10,"V":0},{"I":11,"V":0},{"I":12,"V":"224.0.0.22"},{"I":13,"V":0},{"I":14,"V":0},{"I":15,"V":"0.0.0.0"},{"I":16,"V":0},{"I":17,"V":0},{"I":21,"V":300044},{"I":22,"V":299144}]]}
```

## Supported platform
- Linux
- Windows

## Build
Given that the Go Language compiler (version 1.8 preferred) is installed, you can build it with:
```
go get github.com/VerizonDigital/vflow/vflow
cd $GOPATH/src/github.com/VerizonDigital/vflow

make build
or
go get -d ./...
cd vflow; go build 
```
## Installation
You can download and install pre-built debian package as below ([RPM and Linux binary are available](https://github.com/VerizonDigital/vflow/releases/tag/v0.4.1)). 

dpkg -i [vflow-0.4.1-amd64.deb](https://github.com/VerizonDigital/vflow/releases/download/v0.4.1/vflow-0.4.1-amd64.deb)

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

## License
Licensed under the Apache License, Version 2.0 (the "License")

## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
