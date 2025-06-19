# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Arista sFlow Extensions Support**: Enhanced sFlow implementation to support Arista-specific extensions based on Arista EOS sFlow documentation
  - Added BGP Route Information Extension (record type 1013)
    - BGP next hop IP address parsing
    - AS path sequence decoding with variable length support
    - BGP communities extraction (standard 32-bit format)
    - Local preference, MED, and origin attribute parsing
    - Source, destination, and peer AS number fields
  - Added VPLS Extension (record type 1014)
    - VPLS instance name with variable length string support
    - Pseudowire ID and Virtual Circuit (VC) identifiers
    - VC type classification for MPLS/VPLS networks
  - Added DSCP information structure for traffic class detection
  - Extended record type constants for standard sFlow extensions:
    - Gateway (1003), User (1004), URL (1005)
    - MPLS (1006), NAT (1007), MPLS Tunnel (1008)
    - MPLS VC (1009), MPLS FTN (1010), MPLS LDP FEC (1011)
    - VLAN Tunnel (1012)

### Enhanced
- **sFlow Decoder Improvements**:
  - Updated `decodeFlowSample` function to handle new Arista record types
  - Added proper binary unmarshaling with 4-byte alignment for VPLS strings
  - Enhanced error handling for malformed extension data
  - Improved memory management for variable-length fields

### Added - Testing
- **Comprehensive Test Coverage**:
  - Added unit tests for `ExtAristaBGPData` unmarshaling with mock data
  - Added unit tests for `ExtAristaVPLSData` unmarshaling with mock data
  - Added integration tests for decoder functions
  - Added benchmark tests for performance validation
  - All existing tests continue to pass ensuring backward compatibility

### Added - Documentation
- **Enhanced Documentation**:
  - Updated CLAUDE.md with Arista sFlow extension details
  - Added JSON output examples for new record types
  - Documented record type constants and their purposes
  - Included architecture notes for extension integration

### Technical Details
- **Data Structures**:
  ```go
  type ExtAristaBGPData struct {
      NextHop       net.IP   // BGP next hop IP address
      ASPath        []uint32 // AS path sequence
      Communities   []uint32 // BGP communities
      LocalPref     uint32   // Local preference
      SourceAS      uint32   // Source AS number
      DestAS        uint32   // Destination AS number
      PeerAS        uint32   // Peer AS number
      MED           uint32   // Multi-exit discriminator
      Origin        uint32   // BGP origin attribute
  }

  type ExtAristaVPLSData struct {
      InstanceName string // VPLS instance name
      PseudowireID uint32 // Pseudowire ID
      VCID         uint32 // VC ID
      VCType       uint32 // VC Type
  }
  ```

- **JSON Output Format**:
  - BGP extensions appear as `"ExtAristaBGP"` in Records map
  - VPLS extensions appear as `"ExtAristaVPLS"` in Records map
  - Maintains backward compatibility with existing record types
  - Follows existing vFlow JSON structure patterns

### Files Modified
- `sflow/flow_sample.go`: Added new structures, constants, and decoder functions
- `sflow/decoder_test.go`: Added comprehensive test cases and benchmarks
- `CLAUDE.md`: Updated with Arista extension documentation

### Compatibility
- **Backward Compatible**: All existing sFlow functionality preserved
- **Standard Compliant**: Follows sFlow v5 specification for enterprise extensions
- **Performance**: Minimal overhead for non-Arista sFlow data
- **Memory Efficient**: Proper cleanup and reuse of byte buffers

## [0.9.0] - Previous Release
- Existing vFlow functionality
- Support for standard sFlow v5, IPFIX, and Netflow
- Message queue integration (Kafka, NSQ, NATS)
- Dynamic worker scaling
- Packet mirroring capabilities

---

**Note**: This changelog was created to document the Arista sFlow enhancements. For historical changes prior to these enhancements, please refer to the git commit history.