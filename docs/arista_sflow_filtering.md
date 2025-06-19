# Arista sFlow Filtering for OCI Object Store

## Overview

vFlow supports selective filtering of sFlow traffic when uploading to OCI Object Store. This feature allows you to send only Arista-specific sFlow data to object storage while continuing to send all sFlow traffic to the message queue. This reduces storage costs and processing overhead for object store operations.

## Architecture

```
sFlow UDP → sFlowWorker() → JSON → ┬→ sFlowMQCh → Message Queue (All Traffic)
                                   └→ sFlowOCICh → OCI Object Store (Arista Only)
```

## Arista Extensions Detected

The filtering system detects the following Arista-specific sFlow extensions:

### Flow Sample Extensions
- **ExtAristaBGP** (Type 1013): BGP Route Information
  - BGP next hop, AS path, communities
  - Local preference, MED, origin attributes
  - Source, destination, and peer AS numbers

- **ExtAristaVPLS** (Type 1014): VPLS Extension
  - VPLS instance name and pseudowire ID
  - Virtual circuit ID and type
  - MPLS/VPLS network visibility

- **ExtAristaDSCP** (Type 1015): DSCP Extension
  - Original DSCP value before rewriting
  - Rewritten DSCP value (if applicable)
  - Flag indicating if DSCP was rewritten

## Configuration

### Enable Arista-Only Filtering

Add the following to your `vflow.conf`:

```yaml
# sFlow OCI Object Store configuration
sflow-oci-enabled: true
sflow-oci-config-file: "oci.conf"
sflow-oci-arista-only: true  # Only send Arista sFlow to OCI
```

### Disable Filtering (Send All sFlow to OCI)

```yaml
sflow-oci-arista-only: false  # Send all sFlow to OCI
```

### Command Line Options

```bash
# Enable Arista-only filtering
vflow -sflow-oci-enabled=true -sflow-oci-arista-only=true

# Disable filtering (send all sFlow to OCI)
vflow -sflow-oci-enabled=true -sflow-oci-arista-only=false
```

### Environment Variables

```bash
export VFLOW_SFLOW_OCI_ENABLED=true
export VFLOW_SFLOW_OCI_ARISTA_ONLY=true
```

## Example JSON Output

### Arista sFlow Record (Sent to OCI)
```json
{
  "Version": 5,
  "AgentSubID": 0,
  "IPAddress": "192.168.1.100",
  "Samples": [{
    "SequenceNo": 12345,
    "SamplingRate": 1000,
    "Records": {
      "RawHeader": {...},
      "ExtAristaBGP": {
        "NextHop": "192.168.1.1",
        "ASPath": [65001, 65002],
        "Communities": [100:200],
        "LocalPref": 100,
        "SourceAS": 65001,
        "DestAS": 65002,
        "MED": 50,
        "Origin": 1
      },
      "ExtAristaDSCP": {
        "OriginalDSCP": 46,
        "RewrittenDSCP": 0,
        "DSCPRewritten": true
      }
    }
  }]
}
```

### Standard sFlow Record (Not Sent to OCI)
```json
{
  "Version": 5,
  "AgentSubID": 0,
  "IPAddress": "10.0.0.1",
  "Samples": [{
    "SequenceNo": 12346,
    "SamplingRate": 1000,
    "Records": {
      "RawHeader": {...},
      "ExtSwitch": {
        "SrcVlan": 100,
        "DstVlan": 200
      }
    }
  }]
}
```

## Benefits

### Cost Optimization
- **Reduced Storage**: Only relevant Arista traffic stored in object store
- **Lower API Calls**: Fewer upload operations to OCI
- **Efficient Processing**: Smaller data volumes for batch processing

### Performance Benefits
- **Faster Uploads**: Smaller file sizes complete faster
- **Reduced Memory**: Lower buffer usage for OCI operations
- **Better Throughput**: More efficient use of network bandwidth

### Operational Benefits
- **Targeted Analysis**: Focus on Arista-specific network intelligence
- **Compliance**: Separate Arista data for regulatory requirements
- **Flexibility**: Can disable filtering if needed

## Monitoring

### Log Messages

When verbose logging is enabled, you'll see:

```
[vflow] sFlow data with Arista extensions sent to OCI
```

### Statistics

The system tracks:
- Total sFlow packets received
- sFlow packets sent to message queue (all traffic)
- sFlow packets sent to OCI (filtered traffic)

## Troubleshooting

### No Data in OCI Object Store

1. **Check if filtering is enabled**: `sflow-oci-arista-only: true`
2. **Verify Arista extensions**: Ensure your Arista switches send BGP/VPLS/DSCP extensions
3. **Enable verbose logging**: Add `-verbose=true` to see filtering decisions

### All Data Going to OCI

1. **Check filtering setting**: Set `sflow-oci-arista-only: true`
2. **Restart vFlow**: Configuration changes require restart

### Mixed Environment Support

For networks with both Arista and non-Arista devices:
- **Recommended**: Enable filtering (`sflow-oci-arista-only: true`)
- **All traffic still goes to message queue** for standard processing
- **Only Arista traffic goes to object store** for specialized analysis

## Best Practices

1. **Enable Filtering by Default**: Reduces costs and improves performance
2. **Monitor Both Channels**: Ensure message queue gets all traffic, OCI gets filtered traffic
3. **Test Before Production**: Verify your Arista switches send expected extensions
4. **Use Verbose Logging**: During initial setup to verify filtering behavior
5. **Regular Monitoring**: Check OCI upload statistics to ensure expected traffic volume

## Integration Examples

### With Kafka
```yaml
# All sFlow → Kafka
mq-name: kafka
sflow-topic: vflow.sflow

# Arista sFlow → OCI
sflow-oci-enabled: true
sflow-oci-arista-only: true
```

### With NSQ
```yaml
# All sFlow → NSQ  
mq-name: nsq
sflow-topic: vflow.sflow

# Arista sFlow → OCI
sflow-oci-enabled: true
sflow-oci-arista-only: true
```

This filtering approach provides the flexibility to optimize storage costs while maintaining full visibility in your primary message queue system.