# vFlow configuration

## Introduction

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

|Key                  | Default                        | Description                                 |
|---------------------| -------------------------------|---------------------------------------------|
|log-file             | stdError                       | |
|verbose              | false                          | |
|pid-file             | /var/run/vflow.pid             | |
|ipfix-enabled        | true                           | |
|ipfix-port           | 4739                           | |
|ipfix-workers        | 200                            | |
|ipfix-udp-size       | 1500                           | |
|ipfix-mirror-addr    | 
|ipfix-mirror-port    | 4172
|ipfix-mirror-workers | 5
|ipfix-tpl-cache-file | /tmp/vflow.templates
|ipfix-rpc-enabled    | true
|sflow-enabled        | true
|sflow-port           | 6343
|sflow-workers        | 200
|sflow-udp-size       | 1500
|stats-enabled        | true
|stats-http-addr      | *
|stats-http-port      | 8081
|mq-name              | kafka
|mq-config-file       | /usr/local/vflow/etc/kafka.conf

The default configuration path is /usr/local/vflow/etc/vflow.conf but you can change it as below:
```
vflow -config /etc/vflow.conf
```
The vFlow version shows as below:
```
vflow -version
```
