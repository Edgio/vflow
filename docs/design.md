# Overview
The vFlow project is an enterprise IPFIX and sFlow collector. it produces the decoded samples to a message bus like Kafka
or NSQ. The vFlow is high performance and scaleable, It can be able to grow horizontally (each node can talk through RPC
to find out any missed IPFIX template). there is cloning IPFIX UDP packet feature with spoofing in case you need to have
the IPFIX raw data somewhere else.

# Architecture

![Architecture](/docs/imgs/architecture.gif)

# Discovery

Each vFlow uses multicasting on all interfaces to discover nodes to communicate in regard to get any IPFIX new template from other nodes. The multicast IP address is 224.0.0.55 and each node sends hello packet every second. You do not need to enable multicast communications across routers.
