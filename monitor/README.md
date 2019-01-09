# vFlow Monitoring

### vFlow monitoring collects all statistical and diagnostic information about the vFlow itself.

|Metric          | Type   |API Key        | Protocol  | Description                                  |
|----------------| -------|-------------  |:---------:| ---------------------------------------------|
|udp.queue       | Gauge  |UDPQueue       | IPFIX     | UDP packets in queue                         |
|udp.rate        | Gauge  |UDPCount       | IPFIX     | UDP packets per second                       |
|decode.rate     | Gauge  |DecodedCount   | IPFIX     | Decoded packets per second                   |
|udp.mirror.queue| Gauge  |UDPMirrorQueue | IPFIX     | UDP packets in mirror's queue                |
|mq.queue        | Gauge  |MessageQueue   | IPFIX     | Decoded message in mq's queue                |
|mq.error.rate   | Gauge  |MQErrorCount   | IPFIX     | Message queue errors per second              |
|udp.queue       | Gauge  |UDPQueue       | NetflowV9 | UDP packets in queue                         |
|udp.rate        | Gauge  |UDPCount       | NetflowV9 | UDP packets per second                       |
|decode.rate     | Gauge  |DecodedCount   | NetflowV9 | Decoded packets per second                   |
|udp.mirror.queue| Gauge  |UDPMirrorQueue | NetflowV9 | UDP packets in mirror's queue                |
|mq.queue        | Gauge  |MessageQueue   | NetflowV9 | Decoded message in mq's queue                |
|mq.error.rate   | Gauge  |MQErrorCount   | NetflowV9 | Message queue errors per second              |
|udp.queue       | Gauge  |UDPQueue       | SFLOW     | UDP packets in queue                         |
|udp.rate        | Gauge  |UDPCount       | SFLOW     | UDP packets per second                       |
|decode.rate     | Gauge  |DecodedCount   | SFLOW     | Decoded packets per second                   |
|mq.queue        | Gauge  |MessageQueue   | SFLOW     | Decoded message in mq's queue                |
|mq.error.rate   | Gauge  |MQErrorCount   | SFLOW     | Message queue errors per second              |
|mem.heap.alloc  | Gauge  |MemHeapAlloc   | SYSTEM    | HeapAlloc is bytes of allocated heap objects |
|mem.alloc       | Gauge  |MemAlloc       | SYSTEM    | Bytes allocated and not yet freed            |
|mcache.inuse    | Gauge  |MCacheInuse    | SYSTEM    | Bytes used by mcache structures              |
|mem.total.alloc | Counter|MemTotalAlloc  | SYSTEM    | Bytes allocated                              |
|mem.heap.sys    | Gauge  |MemHeapSys     | SYSTEM    | Bytes obtained from system                   |
|num.goroutine   | Gauge  |NumGoRoutine   | SYSTEM    | The number of goroutines that currently exist|

## Grafana sample dashboard

![Alt text](/docs/imgs/grafana.gif?raw=true "vFlow")

## vFlow API

You can hit vFlow stats API directy to create your own monitoring

Flow API : http://localhost:8081/flow 

```
{
   "IPFIX" : {
      "MessageQueue" : 0,
      "DecodedCount" : 733,
      "MQErrorCount" : 0,
      "UDPCount" : 733,
      "UDPMirrorQueue" : 0,
      "UDPQueue" : 0,
      "Workers" : 100
   },
   "NetflowV5" : {
      "MessageQueue" : 0,
      "DecodedCount" : 322,
      "MQErrorCount" : 0,
      "UDPCount" : 322,
      "UDPMirrorQueue" : 0,
      "UDPQueue" : 0,
      "Workers" : 50
   },
   "NetflowV9" : {
      "MessageQueue" : 0,
      "DecodedCount" : 562,
      "MQErrorCount" : 0,
      "UDPCount" : 562,
      "UDPMirrorQueue" : 0,
      "UDPQueue" : 0,
      "Workers" : 80
   },   
   "SFlow" : {
      "MessageQueue" : 0,
      "UDPCount" : 268,
      "MQErrorCount" : 0,
      "DecodedCount" : 253,
      "UDPQueue" : 0,
      "Workers" : 100
   },
   "StartTime" : 1490134512
}
```

System API : http://localhost:8081/sys

```
{
   "GCSys" : 450560,
   "MemTotalAlloc" : 11435376,
   "MCacheInuse" : 4800,
   "GCNext" : 5053510,
   "MemHeapReleased" : 0,
   "NumGoroutine" : 237,
   "GoVersion" : "go1.7.4",
   "NumLogicalCPU" : 4,
   "MaxProcs" : 4, 
   "GCLast" : "2017-03-21 22:17:50.923246779 +0000 UTC",
   "MemHeapAlloc" : 4151416,
   "MemAlloc" : 4151416,
   "StartTime" : 1490134512,
   "MemHeapSys" : 5734400
}
```

## Configuration Keys - Command line
The monitor command line configuration contains the following keys

|Key                | Default         | 
|-------------------| ----------------|
|db-type            | influxdb        |     
|vflow-host         | localhost:8081  | 
|influxdb-api-addr  | localhost:8086  |
|influxdb-db-name   | vflow           |
|tsdb-api-addr      | localhost:4242  |
|hostname           | system hostname |

crontab every 1 minute example:
```
* * * * * monitor -vflow-host 192.168.0.7 -influxdb-api-addr 192.168.0.15
```
