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
|cpu-cap                 | the number of available CPUs   | sets the maximum number of CPUs                  |
|ipfix-enabled           | true                           | enable/disable IPFIX decoders                    |
|ipfix-port              | 4739                           | server IPFIX UDP port                            |
|ipfix-addr              | -                              | server IPFIX UDP IP address to bind to           |
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
|sflow-type-filter       | -                              | filter sflow type(s)                             |
|netflow5-enabled        | true                           | enable/disable netflow v5 decoders               |
|netflow5-port           | 9996                           | server netflow v5 UDP port                       |
|netflow5-workers        | 50                             | netflow v5 concurrent decoders                   |
|netflow5-topic          | vflow.netflow5                 | netflow v5 message queue topic name              |
|netflow5-udp-size       | 1500                           | maximum netflow v9 UDP packet size 
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
|mq-name                 | kafka                          | [message queues](#message-queues)        |
|mq-config-file          | /etc/vflow/mq.conf             | message queue config file                        |

The default configuration path is /etc/vflow/vflow.conf but you can change it as below:
```
vflow -config /usr/local/etc/vflow.conf
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
## Message Queues 
The vFlow supports these message queuing 
- kafka
- kafka.segmentio
- nsq
- nat
- rawSocket

Note: there are two kafka drivers: [Kafka Sarama](https://github.com/Shopify/sarama) (Default) and [Kafka Segmentio](https://github.com/segmentio/kafka-go) (Kafka-Go)


# Kafka Configuration

## Format
A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```
key: value
```

The default configuration file is /etc/vflow/mq.conf, you can be able to change it through vFlow configuration.

## Configuration Keys (Default / Sarama)
The Kafka configuration contains the following key

|Key                  | Default     |  Environment variable        | Description                                                                        |
|---------------------| ------------|------------------------------|------------------------------------------------------------------------------------|
|brokers              | -           | VFLOW_KAFKA_BROKERS          | kafka broker addresses                                                             |
|compression          | none        | VFLOW_KAFKA_COMPRESSION      | compression codecs: gzip, snappy, lz4                                              |
|retry-max            | 2           | VFLOW_KAFKA_RETRY_MAX        | the total number of times to retry                                                 |
|request-size-max     | 104857600   | VFLOW_KAFKA_REQUEST_SIZE_MAX | the maximum size (in bytes) of any request that will be attempted to send to Kafka |
|retry-backoff        | 10          | VFLOW_KAFKA_RETRY_BACKOFF    | wait for leader election to occur before retrying in milliseconds              |
|tls-enabled          | false       | VFLOW_KAFKA_TLS_ENABLED      | connect using TLS                                                                   |
|tls-cert             | none        | VFLOW_KAFKA_TLS_CERT         | certificate file for client authentication                                         |
|tls-key              | none        | VFLOW_KAFKA_TLS_KEY          | key file for client authentication                                                 |
|ca-file              | none        | VFLOW_KAFKA_CA_FILE          | certificate authority file for TLS client authentication                           |
|tls-skip-verify      | true        | VFLOW_KAFKA_TLS_SKIP_VERIFY  | if true, the server's certificate will not validate                                                      |

## Example
```
brokers:
    - 192.16.1.25:9092
retry-max: 1
retry-backoff: 30
```

## Configuration Keys (Segmentio)
The Kafka configuration contains the following key

|Key                  | Default     |  Environment variable        | Description                                                                        |
|---------------------| ------------|------------------------------|------------------------------------------------------------------------------------|
|brokers              |             | VFLOW_KAFKA_BROKERS          |                                                                                    |
|bootstrap-server     |             | VFLOW_KAFKA_BOOTSTRAP_SERVER |                                                                                    |
|client-id            |             | VFLOW_KAFKA_CLIENT_ID        |                                                                                    |
|compression          |             | VFLOW_KAFKA_COMPRESSION      |                                                                                    |
|max-attempts         |             | VFLOW_KAFKA_MAX_ATTEMPTS     |                                                                                    |
|queue-size           |             | VFLOW_KAFKA_QUEUE_SIZE       |                                                                                    |
|batch-size           |             | VFLOW_KAFKA_BATCH_SIZE       |                                                                                    |
|keepalive            |             | VFLOW_KAFKA_KEEPALIVE        |                                                                                    |
|connect-timeout      |             | VFLOW_KAFKA_CONNECT_TIMEOUT  |                                                                                    |
|required-acks        |             | VFLOW_KAFKA_REQUIRED_ACKS    |                                                                                    |
|pflush               |             | VFLOW_KAFKA_PERIODIC_FLUSH   |                                                                                    |
|tls-cert             |             | VFLOW_KAFKA_TLS_CERT         |                                                                                    |
|tls-key              |             | VFLOW_KAFKA_TLS_KEY          |                                                                                    |
|ca-file              |             | VFLOW_KAFKA_CA_FILE          |                                                                                    |
|verify-ssl           |             | VFLOW_KAFKA_VERIFY_SSL       |                                                                                    |


# NSQ Configuration

## Format
A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```
key: value
```

The default configuration file is /etc/vflow/mq.conf, you can be able to change it through vFlow configuration.

## Configuration Keys
The NSQ configuration contains the following key

|Key                  | Default        |  Environment variable    | Description                                                      |
|---------------------| ---------------|--------------------------|------------------------------------------------------------------|
|server               | localhost:4150 | NA                       | NSQ server addresse and port

# NATS Configuration

## Format
A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```
key: value
```

The default configuration file is /etc/vflow/mq.conf, you can be able to change it through vFlow configuration.

## Configuration Keys
The NATS configuration contains the following key

|Key                  | Default               |  Environment variable    | Description                                                      |
|---------------------| ----------------------|--------------------------|------------------------------------------------------------------|
|url                  | nats://localhost:4222 | NA                       | URL addresse

# Raw Socket Configuration

Note that for messages sent over TCP and UDP using this producer, the message deliminator is a new line character ("\n").

## Format
A config file is a plain text file in [YAML](https://en.wikipedia.org/wiki/YAML) format.

```
key: value
```

The default configuration file is /etc/vflow/mq.conf, you can be able to change it through vFlow configuration.


## Configuration Keys
The rawSocket configuration contains the following key

|Key                  | Default               |  Environment variable    | Description                                                          |
|---------------------| ----------------------|--------------------------|----------------------------------------------------------------------|
|url                  | localhost:9555        | NA                       | URL address to send to. Includes the hostname and port.              |
|protocol             | tcp                   | NA                       | Protocol to use to send. Can be either "tcp" or "udp"                |
|retry-max            | 2                     | NA                       | The number of times a message will be retried before giving up on it |
