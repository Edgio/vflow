### Juniper MX Series routers (MX5, MX10, MX40, MX80, MX104, MX120, MX240, MX480, MX960)

```
set services flow-monitoring version-ipfix template vflow flow-active-timeout 10
set services flow-monitoring version-ipfix template vflow flow-inactive-timeout 10
set services flow-monitoring version-ipfix template vflow template-refresh-rate packets 1000
set services flow-monitoring version-ipfix template vflow template-refresh-rate seconds 10
set services flow-monitoring version-ipfix template vflow option-refresh-rate packets 1000
set services flow-monitoring version-ipfix template vflow option-refresh-rate seconds 10
set services flow-monitoring version-ipfix template vflow ipv4-template

set forwarding-options sampling instance ipfix input rate 100
set forwarding-options sampling instance ipfix family inet output flow-server 192.168.0.10 port 4739
set forwarding-options sampling instance ipfix family inet output flow-server 192.168.0.10 version-ipfix template vflow
set forwarding-options sampling instance ipfix family inet output inline-jflow source-address 192.168.0.1
```
