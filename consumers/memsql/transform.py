#!/usr/bin/python
#: ----------------------------------------------------------------------------
#: Copyright (C) 2017 Verizon.  All Rights Reserved.
#: All Rights Reserved
#:
#: file:    transform.py
#: details: memsql pipline transform python script
#: author:  Mehrdad Arshad Rad
#: date:    04/27/2017
#:
#: Licensed under the Apache License, Version 2.0 (the "License");
#: you may not use this file except in compliance with the License.
#: You may obtain a copy of the License at
#:
#:     http://www.apache.org/licenses/LICENSE-2.0
#:
#: Unless required by applicable law or agreed to in writing, software
#: distributed under the License is distributed on an "AS IS" BASIS,
#: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#: See the License for the specific language governing permissions and
#: limitations under the License.
#: ----------------------------------------------------------------------------

import json
import struct
import sys
import time


def transform_records():
    while True:
        byte_len = sys.stdin.read(8)
        if len(byte_len) == 8:
            byte_len = struct.unpack("L", byte_len)[0]
            result = sys.stdin.read(byte_len)
            yield result
        else:
            assert len(byte_len) == 0, byte_len
            return

for records in transform_records():
    flows = json.loads(records)
    exported_time = time.strftime('%Y-%m-%d %H:%M:%S',
                                  time.localtime(flows["Header"]["ExportTime"]))

    try:
        for flow in flows["DataSets"]:
            sourceIPAddress = "unknown"
            destinationIPAddress = "unknown"
            bgpSourceAsNumber = "unknown"
            bgpDestinationAsNumber = "unknown"
            protocolIdentifier = 0
            sourceTransportPort = 0
            destinationTransportPort = 0
            tcpControlBits = "unknown"
            ipNextHopIPAddress = "unknown"
            octetDeltaCount = 0
            ingressInterface = 0
            egressInterface = 0

            for field in flow:
                if field["I"] in [214]:
                    raise
                elif field["I"] in [8, 27]:
                    sourceIPAddress = field["V"]
                elif field["I"] in [12, 28]:
                    destinationIPAddress = field["V"]
                elif field["I"] in [15, 62]:
                    ipNextHopIPAddress = field["V"]
                elif field["I"] == 16:
                    bgpSourceAsNumber = field["V"]
                elif field["I"] == 17:
                    bgpDestinationAsNumber = field["V"]
                elif field["I"] == 14:
                    ingressInterface = field["V"]
                elif field["I"] == 10:
                    egressInterface = field["V"]
                elif field["I"] == 7:
                    sourceTransportPort = field["V"]
                elif field["I"] == 11:
                    destinationTransportPort = field["V"]
                elif field["I"] == 4:
                    protocolIdentifier = field["V"]
                elif field["I"] == 6:
                    tcpControlBits = field["V"]
                elif field["I"] == 1:
                    octetDeltaCount = field["V"]

            out = b"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n" \
                  % (
                        flows["AgentID"],
                        sourceIPAddress,
                        destinationIPAddress,
                        ipNextHopIPAddress,
                        bgpSourceAsNumber,
                        bgpDestinationAsNumber,
                        protocolIdentifier,
                        sourceTransportPort,
                        destinationTransportPort,
                        tcpControlBits,
                        ingressInterface,
                        egressInterface,
                        octetDeltaCount,
                        exported_time,
                    )

            sys.stdout.write(out)
    except:
        continue
