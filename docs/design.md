# Overview

The vFlow project is an enterprise IPFIX and sFlow collector. it produces the decoded samples to a message bus like Kafka
or NSQ. The vFlow is high performance and scalable, It can be able to grow horizontally (each node can talk through RPC
to find out any missed IPFIX template). there is cloning IPFIX UDP packet feature with spoofing in case you need to have
the IPFIX raw data somewhere else.

# Architecture

![Architecture](/docs/imgs/architecture.gif)

# Dynamic pool

The number of worker processes can be changed at runtime automated based on the incoming load. the minimum workers can be able to configure then vFlow adjusts it at runtime gradually.  

# Discovery

Each vFlow uses multicasting on all interfaces to discover nodes to communicate in regard to get any IPFIX new template from other nodes. The multicast IP address is 224.0.0.55 and each node sends hello packet every second. You do not need to enable multicast communications across routers.

# Pluggable architecture

The vFlow accepts message queue plugin. for the time being it has Kafka and NSQ plugins but you can write for a message queue like RabbitMQ quick and easy.

# Hardware requirements

|Load|IPFIX PPS|CPU|RAM|
|----|---------|---|---|
|low| < 1K |2-4|64M|
|moderate| < 10K| 8+| 256M|
|high| < 50K| 12+| 512M|
|x-high| < 100K | 24+ | 1G|
