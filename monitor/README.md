# vFlow Monitoring

### vFlow monitoring collects all statistical and diagnostic information about the vFlow itself.

|Metric          | Type   |API Key        | Protocol| Description                                  |
|----------------| -------|-------------  |:-------:| ---------------------------------------------|
|udp.queue       | Gauge  |UDPQueue       | IPFIX   | UDP packets in queue                         |
|udp.rate        | Gauge  |UDPCount       | IPFIX   | UDP packets per second                       |
|decode.rate     | Gauge  |DecodedCount   | IPFIX   | Decoded packets per second                   |
|udp.mirror.queue| Gauge  |UDPMirrorQueue | IPFIX   | UDP packets in mirror's queue                |
|mq.queue        | Gauge  |MessageQueue   | IPFIX   | Decoded message in mq's queue                |
|mq.error.rate   | Gauge  |MQErrorCount   | IPFIX   | Message queue errors per second              |
|udp.queue       | Gauge  |UDPQueue       | SFLOW   | UDP packets in queue                         |
|udp.rate        | Gauge  |UDPCount       | SFLOW   | UDP packets per second                       |
|decode.rate     | Gauge  |DecodedCount   | SFLOW   | Decoded packets per second                   |
|mq.queue        | Gauge  |MessageQueue   | SFLOW   | Decoded message in mq's queue                |
|mq.error.rate   | Gauge  |MQErrorCount   | SFLOW   | Message queue errors per second              |
|mem.heap.alloc  | Gauge  |MemHeapAlloc   | SYSTEM  | HeapAlloc is bytes of allocated heap objects |
|mem.alloc       | Gauge  |MemAlloc       | SYSTEM  | Bytes allocated and not yet freed            |
|mcache.inuse    | Gauge  |MCacheInuse    | SYSTEM  | Bytes used by mcache structures              |
|mem.total.alloc | Counter|MemTotalAlloc  | SYSTEM  | Bytes allocated                              |
|mem.heap.sys    | Gauge  |MemHeapSys     | SYSTEM  | Bytes obtained from system                   |
|num.goroutine   | Gauge  |NumGoRoutine   | SYSTEM  | The number of goroutines that currently exist|

## Grafana sample dashboard

![Alt text](/docs/imgs/grafana.gif?raw=true "vFlow")
