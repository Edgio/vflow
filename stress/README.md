# Stress 

## Stress tries to check correct behavior and robustness of vFlow

## Features
- Generate IPFIX data, template and template option
- Generate sFlow v5 sample header data
- Simulate different IPFIX agents

![Alt text](/docs/imgs/stress.gif?raw=true "vFlow")

## Usage Manual
````
-vflow-addr         vflow ip address (default 127.0.0.1)
-ipfix-port         ipfix port number (default 4739)
-sflow-port         sflow port number (default 6343)
-ipfix-interval     ipfix template interval (default 10s)
-ipfix-rate-limit   ipfix rate limit packets per second (default 25000 PPS)
-sflow-rate-limit   sflow rate limit packets per second (default 25000 PPS)
````
