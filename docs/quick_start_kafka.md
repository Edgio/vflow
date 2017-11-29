# Install vFlow with Kafka - Linux

## Download and install vFlow 
### Debian Package
``` 
wget https://github.com/VerizonDigital/vflow/releases/download/v0.4.1/vflow-0.4.1-amd64.deb
dpkg -i vflow-0.4.1-amd64.deb
```
### RPM Package
```
wget https://github.com/VerizonDigital/vflow/releases/download/v0.4.1/vflow-0.4.1.amd64.rpm
rpm -ivh vflow-0.4.1.amd64.rpm 
or
yum localinstall vflow-0.4.1.amd64.rpm
```
## Download Kafka
```
wget https://www.apache.org/dyn/closer.cgi?path=/kafka/0.11.0.0/kafka_2.11-0.11.0.0.tgz
tar -xzf kafka_2.11-0.11.0.0.tgz
cd kafka_2.11-0.11.0.0
```
Kafka uses ZooKeeper so you need to first start a ZooKeeper server if already you don't have one
```
bin/zookeeper-server-start.sh config/zookeeper.properties
```
start the Kafka server
```
bin/kafka-server-start.sh config/server.properties
```
## vFlow - start service
```
service vflow start
```

## vFlow - load generator
```
vflow_stress -sflow-rate-limit 1 0ipfix-rate-limit 1 &
```

## Consume IPFIX topic from NSQ
```
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic vflow.ipfix
```
