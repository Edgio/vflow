![vFlow](docs/imgs/vflow_logo.png?raw=true "vFlow logo")
## [![Build Status](https://travis-ci.org/VerizonDigital/vflow.svg?branch=master)](https://travis-ci.org/VerizonDigital/vflow) [![Go Report Card](https://goreportcard.com/badge/github.com/VerizonDigital/vflow)](https://goreportcard.com/report/github.com/VerizonDigital/vflow)

High-performance, scalable and reliable IPFIX and sFlow collector. 

## Features
- IPFIX RFC7011 collector
- sFLow v5 raw header packet collector
- Decoding sFlow raw header L2/L3/L4 
- Produce to Apache Kafka, NSQ
- Replicate IPFIX to 3rd party collector
- Supports IPv4 and IPv6
- Monitoring with InfluxDB and OpenTSDB backend

![Alt text](/docs/imgs/vflow.gif?raw=true "vFlow")

## Documentation
- [Architecture](/docs/design.md).
- [Configuration](/docs/config.md).
- [Monitoring](/monitor/README.md).
- [Stress / Load Generator](/stress/README.md).

## Decoded IPFIX data
The IPFIX data decodes to JSON format and IDs are [IANA IPFIX element ID](http://www.iana.org/assignments/ipfix/ipfix.xhtml)
```json
{"AgentID":"192.168.21.15","Header":{"Version":10,"Length":420,"ExportTime":1483484642,"SequenceNo":1434533677,"DomainID":32771},"DataSets":[[{"ID":8,"Value":"192.16.28.217"},{"ID":12,"Value":"180.10.210.240"},{"ID":5,"Value":2},{"ID":4,"Value":6},{"ID":7,"Value":443},{"ID":11,"Value":64381},{"ID":32,"Value":0},{"ID":10,"Value":811},{"ID":58,"Value":0},{"ID":9,"Value":24},{"ID":13,"Value":20},{"ID":16,"Value":4200000000},{"ID":17,"Value":27747},{"ID":15,"Value":"180.105.10.210"},{"ID":6,"Value":"0x10"},{"ID":14,"Value":1113},{"ID":1,"Value":22500},{"ID":2,"Value":15},{"ID":52,"Value":63},{"ID":53,"Value":63},{"ID":152,"Value":1483484581770},{"ID":153,"Value":1483484622384},{"ID":136,"Value":2},{"ID":243,"Value":0},{"ID":245,"Value":0}]]}
```

## Decoded sFlow data
```json
{"Header":{"Version":5,"IPVersion":1,"AgentSubID":0,"SequenceNo":24324,"SysUpTime":766903208,"SamplesNo":1,"IPAddress":"192.16.14.0"},"ExtSWData":{"SrcVlan":0,"SrcPriority":0,"DstVlan":12,"DstPriority":0},"Sample":{"SequenceNo":0,"SourceID":0,"SamplingRate":2000,"SamplePool":0,"Drops":0,"Input":552,"Output":0,"RecordsNo":2},"Packet":{"L2":{"SrcMAC":"d4:04:ff:01:1d:9e","DstMAC":"30:7c:5e:e5:59:ef","Vlan":12,"EtherType":34525},"L3":{"Version":6,"TrafficClass":0,"FlowLabel":0,"PayloadLen":265,"NextHeader":17,"HopLimit":57,"Src":"2600:8000:5207:6f00::1","Dst":"2606:2800:404e:2:1663:6fe:2cc6:100a"},"L4":{"SrcPort":53,"DstPort":34234}}}
```

## Build
Given that the Go Language compiler (version 1.8 preferred) is installed, you can build it with:
```
go get github.com/VerizonDigital/vflow
cd $GOPATH/src/github.com/VerizonDigital/vflow

make build
or
go get -d ./...
cd vflow; go build 
```

## License
Licensed under the Apache License, Version 2.0 (the "License")

## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
