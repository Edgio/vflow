# CLAUDE.md
## Instruction 
1. First think through the problem, read the codebase for relevant files, and write a plan to todo.md.
2. The plan should have a list of todo items that you can check off as you complete them
3. Before you begin working, check in with me and I will verify the plan.
4. Then, begin working on the todo items, marking them as complete as you go.
5. Please every step of the way just give me a high level explanation of what changes you made
6. Make every task and code change you do as simple as possible. We want to avoid making any massive or complex changes. Every change should impact as little code as possible. Everything is about simplicity.
7. Finally, add a review section to the todo.md file with a summary of the changes you made and any other relevant information.



This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vFlow is a high-performance, scalable IPFIX, sFlow and Netflow collector written in Go. It decodes network flow data and produces JSON messages to message queues like Kafka, NSQ, or NATS. The project supports horizontal scaling through RPC communication between nodes and includes packet mirroring capabilities.

## Common Commands

### Building and Running
- `make build` - Build the main vflow binary and stress testing tool
- `make test` - Run all tests with 1-minute timeout
- `make bench` - Run benchmarks with 2-minute timeout
- `make run` - Build and run vflow with default worker configuration
- `make debug` - Build and run vflow with verbose logging enabled

### Development Tools
- `make lint` - Run golint on all packages
- `make cyclo` - Check cyclomatic complexity (threshold: 15)
- `make errcheck` - Run error checking analysis
- `make tools` - Install development tools (golint, errcheck, gocyclo)

### Direct Go Commands
- `cd vflow; go build` - Build just the vflow binary
- `go test -v ./... -timeout 1m` - Run tests directly

## Architecture

### Core Components
- **vflow/**: Main application entry point and protocol handlers
- **ipfix/**: IPFIX RFC7011 decoder with template caching and RPC memcache
- **sflow/**: sFlow v5 decoder for raw headers and counters
- **netflow/**: Netflow v5 and v9 decoders
- **packet/**: Low-level packet parsing (Ethernet, IP, TCP/UDP, ICMP)
- **producer/**: Message queue producers (Kafka, NSQ, NATS, raw socket)
- **mirror/**: UDP packet mirroring with IP spoofing
- **monitor/**: Prometheus metrics and monitoring storage
- **reader/**: Configuration file reader with YAML support

### Key Features
- Multi-protocol support: IPFIX, sFlow v5, Netflow v5/v9
- Pluggable message queue architecture
- Horizontal scaling with multicast discovery (224.0.0.55)
- Dynamic worker pool adjustment based on load
- Template caching with cross-node RPC sharing
- Packet mirroring and replication
- RESTful API and Prometheus monitoring

### Configuration
Main config: `scripts/vflow.conf` - Worker counts and log paths
Message queue config: `scripts/kafka.conf` - Producer settings
IPFIX elements: `scripts/ipfix.elements` - Field definitions

### Data Flow
1. Raw network packets received on UDP ports (IPFIX: 4739, sFlow: 6343, Netflow: 4729)
2. Protocol-specific decoders parse and extract flow records
3. Templates cached locally and shared via RPC between nodes
4. Decoded data converted to JSON with IANA field IDs
5. JSON messages sent to configured message queue producers
6. Optional packet mirroring to third-party collectors

### Worker Model
Each protocol uses configurable worker pools:
- IPFIX workers handle template management and decoding
- sFlow workers process sampling data and raw headers
- Netflow workers decode v5/v9 flow records
- Dynamic scaling adjusts worker count based on incoming load

### sFlow Arista Extensions
The sFlow implementation includes support for Arista-specific extensions:

**BGP Route Information Extension (SFDataExtAristaBGP: 1013)**
- BGP next hop, AS path, communities
- Local preference, MED, origin attributes
- Source, destination, and peer AS numbers
- Provides comprehensive BGP routing context

**VPLS Extension (SFDataExtAristaVPLS: 1014)**
- VPLS instance name and pseudowire ID
- Virtual circuit ID and type
- Support for MPLS/VPLS network visibility

**Extended Record Types Supported:**
- Standard: Raw Header (1), Extended Switch (1001), Extended Router (1002)
- Arista: BGP Route Info (1013), VPLS (1014)
- Additional: Gateway, User, URL, MPLS, NAT, Tunnel extensions

**JSON Output Examples:**
```json
{
  "Records": {
    "ExtAristaBGP": {
      "NextHop": "192.168.1.1",
      "ASPath": [256, 512, 768],
      "Communities": [16777316, 33554632],
      "LocalPref": 100,
      "SourceAS": 256,
      "DestAS": 512,
      "PeerAS": 768,
      "MED": 50,
      "Origin": 1
    },
    "ExtAristaVPLS": {
      "InstanceName": "vpls1234",
      "PseudowireID": 291,
      "VCID": 1110,
      "VCType": 5
    }
  }
}
```

### sFlow OCI Object Store Integration
Added support for uploading sFlow records to OCI Object Store with configurable batching and chunking:

**Features:**
- Parallel processing: sFlow data sent to both message queue and OCI simultaneously
- Batching: Configurable file size limits and flush intervals for cost optimization
- File chunking: Automatic splitting of large files to prevent timeouts
- Timestamp-based naming: `sflow_YYYYMMDD_HHMMSS_001.json`

**Configuration:**
```yaml
# Enable OCI upload
sflow-oci-enabled: true
sflow-oci-config-file: "oci.conf"
```

**OCI Configuration (oci.conf):**
```yaml
region: "us-phoenix-1"
tenancy_ocid: "ocid1.tenancy.oc1..example"
user_ocid: "ocid1.user.oc1..example"
fingerprint: "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"
private_key_path: "/path/to/oci_private_key.pem"
bucket_name: "sflow-data"
namespace: "mycompany"
chunk_size_mb: 10
flush_interval: "60s"
```

**Architecture:**
```
sFlow UDP → sFlowWorker() → JSON → ┬→ sFlowMQCh → Message Queue
                                   └→ sFlowOCICh → OCI Batch Worker → OCI Object Store
```
