# vFlow configuration

## Format

A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format also you can be able to configure
 through the command line. in case you configure a key at config file and command line, the command line would be preferred.

### config file
```
key: value
```
### command line
```
-key value
```
## Configuration Keys
The vFlow configuration contains the following keys

|Key                     | Default                        | Description                                      |
|------------------------| -------------------------------|--------------------------------------------------|
|log-file                | stdError                       | name of log file to send logging output to       |
|verbose                 | false                          | enable the full logging                          |
|pid-file                | /var/run/vflow.pid             | file in which server should write its process ID |
|ipfix-enabled           | true                           | enable/disable IPFIX decoders                    |
|ipfix-port              | 4739                           | server IPFIX UDP port                            |
|ipfix-workers           | 200                            | IPFIX concurrent decoders                        |
|ipfix-topic             | vflow.ipfix                    | ipfix message queue topic name                   |
|ipfix-udp-size          | 1500                           | maximum IPFIX UDP packet size                    |
|ipfix-mirror-addr       | -                              | IPFIX 3rd party collector address                |
|ipfix-mirror-port       | 4172                           | IPFIX 3rd party collector port                   |
|ipfix-mirror-workers    | 5                              | IPFIX replicator concurrent packet generator     |
|ipfix-tpl-cache-file    | /tmp/vflow.templates           | IPFIX templates cache file                       |
|ipfix-rpc-enabled       | true                           | enable/disable IPFIX RPC                         |
|sflow-enabled           | true                           | enable/disable sFlow decoders                    |
|sflow-port              | 6343                           | server sFlow UDP port                            |
|sflow-workers           | 200                            | sFlow concurrent decoders                        |
|sflow-udp-size          | 1500                           | maximum sFlow UDP packet size                    |
|sflow-topic             | vflow.sflow                    | sFlow message queue topic name                   |
|netflow9-enabled        | true                           | enable/disable netflow v9 decoders               |
|netflow9-port           | 4729                           | server netflow v9 UDP port                       |
|netflow9-workers        | 50                             | netflow v9 concurrent decoders                   |
|netflow9-topic          | vflow.netflow9                 | netflow v9 message queue topic name              |
|netflow9-udp-size       | 1500                           | maximum netflow v9 UDP packet size               |
|netflow9-tpl-cache-file | /tmp/netflow9.templates        | netflow v9 templates cache file                  |
|dynamic-workers         | true                           | enable/disable dynamic workers feature           |
|stats-enabled           | true                           | enable/disable web stats listener                |
|stats-http-addr         | *                              | web stats address option at server startup       |
|stats-http-port         | 8081                           | web stats TCP port                               |
|mq-name                 | kafka                          | message queueing name (kafka or nsq)             |
|mq-config-file          | /usr/local/vflow/etc/kafka.conf| message queue config file                        |

The default configuration path is /usr/local/vflow/etc/vflow.conf but you can change it as below:
```
vflow -config /etc/vflow.conf
```
The vFlow version shows as below:
```
vflow -version
```

## Example
```
ipfix-workers: 600
sflow-workers: 300
log-file: /var/log/vflow.log
```

# Kafka Configuration

## Format
A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```
key: value
```

The default configuration file is /usr/local/vflow/etc/kafka.conf, you can be able to change it through vFlow configuration.

## Configuration Keys
The Kafka configuration contains the following key

|Key                  | Default |  Environment variable    | Description                                                    |
|---------------------| --------|--------------------------|----------------------------------------------------------------|
|brokers              | -       | VFLOW_KAFKA_BROKERS      | kafka broker addresses                                         |
|compression          | none    | VFLOW_KAFKA_COMPRESSION  | compression codecs: gzip, snappy, lz4                          |
|retry-max            | 0       | VFLOW_KAFKA_RETRY_MAX    | the total number of times to retry                             |
|retry-backoff        | 0       | VFLOW_KAFKA_RETRY_BACKOFF| wait for leader election to occur before retrying in milliseconds|

## Example
```
brokers: 
    - 192.16.1.25:9092
retry-max: 2
retry-backoff: 10
```
