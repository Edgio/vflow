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

The default format is prometheus: http://localhost:8081/metrics

```
# HELP vflow_ipfix_decoded_packets 
# TYPE vflow_ipfix_decoded_packets counter
vflow_ipfix_decoded_packets 0
# HELP vflow_ipfix_message_queue 
# TYPE vflow_ipfix_message_queue gauge
vflow_ipfix_message_queue 0
# HELP vflow_ipfix_mq_error 
# TYPE vflow_ipfix_mq_error counter
vflow_ipfix_mq_error 0
# HELP vflow_ipfix_udp_mirror_queue 
# TYPE vflow_ipfix_udp_mirror_queue gauge
vflow_ipfix_udp_mirror_queue 0
# HELP vflow_ipfix_udp_packets 
# TYPE vflow_ipfix_udp_packets counter
vflow_ipfix_udp_packets 0
# HELP vflow_ipfix_udp_queue 
# TYPE vflow_ipfix_udp_queue gauge
vflow_ipfix_udp_queue 0
# HELP vflow_ipfix_workers 
# TYPE vflow_ipfix_workers gauge
vflow_ipfix_workers 200
# HELP vflow_netflowv5_decoded_packets 
# TYPE vflow_netflowv5_decoded_packets counter
vflow_netflowv5_decoded_packets 0
# HELP vflow_netflowv5_message_queue 
# TYPE vflow_netflowv5_message_queue counter
vflow_netflowv5_message_queue 0
# HELP vflow_netflowv5_mq_error 
# TYPE vflow_netflowv5_mq_error counter
vflow_netflowv5_mq_error 0
# HELP vflow_netflowv5_udp_packets 
# TYPE vflow_netflowv5_udp_packets counter
vflow_netflowv5_udp_packets 0
# HELP vflow_netflowv5_udp_queue 
# TYPE vflow_netflowv5_udp_queue counter
vflow_netflowv5_udp_queue 0
# HELP vflow_netflowv5_workers 
# TYPE vflow_netflowv5_workers counter
vflow_netflowv5_workers 200
# HELP vflow_netflowv9_decoded_packets 
# TYPE vflow_netflowv9_decoded_packets counter
vflow_netflowv9_decoded_packets 0
# HELP vflow_netflowv9_message_queue 
# TYPE vflow_netflowv9_message_queue counter
vflow_netflowv9_message_queue 0
# HELP vflow_netflowv9_mq_error 
# TYPE vflow_netflowv9_mq_error counter
vflow_netflowv9_mq_error 0
# HELP vflow_netflowv9_udp_packets 
# TYPE vflow_netflowv9_udp_packets counter
vflow_netflowv9_udp_packets 0
# HELP vflow_netflowv9_udp_queue 
# TYPE vflow_netflowv9_udp_queue counter
vflow_netflowv9_udp_queue 0
# HELP vflow_netflowv9_workers 
# TYPE vflow_netflowv9_workers counter
vflow_netflowv9_workers 200
# HELP vflow_sflow_decoded_packets 
# TYPE vflow_sflow_decoded_packets counter
vflow_sflow_decoded_packets 0
# HELP vflow_sflow_message_queue 
# TYPE vflow_sflow_message_queue counter
vflow_sflow_message_queue 0
# HELP vflow_sflow_mq_error 
# TYPE vflow_sflow_mq_error counter
vflow_sflow_mq_error 0
# HELP vflow_sflow_udp_packets 
# TYPE vflow_sflow_udp_packets counter
vflow_sflow_udp_packets 0
# HELP vflow_sflow_udp_queue 
# TYPE vflow_sflow_udp_queue counter
vflow_sflow_udp_queue 0
# HELP vflow_sflow_workers 
# TYPE vflow_sflow_workers counter
vflow_sflow_workers 200
```

If you configured the [stats-format](https://github.com/EdgeCast/vflow/blob/master/docs/config.md#Configuration-Keys) to restful then the metrics will be available at http://localhost:8081/flow for flow and system at http://localhost:8081/sys

```json
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

System API

```json
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
