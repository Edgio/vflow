## vFlow [![Build Status](https://travis-ci.org/VerizonDigital/vflow.svg?branch=master)](https://travis-ci.org/VerizonDigital/vflow) [![Go Report Card](https://goreportcard.com/badge/github.com/VerizonDigital/vflow)](https://goreportcard.com/report/github.com/VerizonDigital/vflow)


## Features
- IPFIX RFC7011 collector
- sFLow v5 raw header packet collector
- Decoding sFlow raw header L2/L3/L4 
- Producer to Apache Kafka, NSQ
- Cloning IPFIX stream to other source
- Supports IPv4 and IPv6
- Monitoring with InfluxDB and OpenTSDB backend

![Alt text](/docs/imgs/vflow.gif?raw=true "vFlow")

## Documentation
- [Architecture](/docs/design.md).
- [Monitoring](/monitor/README.md).
- [Stress / Load Generator](/stress/README.md).

## Decoded IPFIX data
The IPFIX data decodes to JSON format and IDs are [IANA IPFIX element ID](http://www.iana.org/assignments/ipfix/ipfix.xhtml)
```json
{"AgentID":"192.168.21.15","Header":{"Version":10,"Length":420,"ExportTime":1483484642,"SequenceNo":1434533677,"DomainID":32771},"DataSets":[[{"ID":8,"Value":"192.16.28.217"},{"ID":12,"Value":"180.10.210.240"},{"ID":5,"Value":2},{"ID":4,"Value":6},{"ID":7,"Value":443},{"ID":11,"Value":64381},{"ID":32,"Value":0},{"ID":10,"Value":811},{"ID":58,"Value":0},{"ID":9,"Value":24},{"ID":13,"Value":20},{"ID":16,"Value":4200000000},{"ID":17,"Value":27747},{"ID":15,"Value":"180.105.10.210"},{"ID":6,"Value":"0x10"},{"ID":14,"Value":1113},{"ID":1,"Value":22500},{"ID":2,"Value":15},{"ID":52,"Value":63},{"ID":53,"Value":63},{"ID":152,"Value":1483484581770},{"ID":153,"Value":1483484622384},{"ID":136,"Value":2},{"ID":243,"Value":0},{"ID":245,"Value":0}]]}
```
## License
Licensed under the Apache License, Version 2.0 (the "License")

## Contribute
Welcomes any kind of contribution, please follow the next steps:

- Fork the project on github.com.
- Create a new branch.
- Commit changes to the new branch.
- Send a pull request.
