//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    FieldMappings.scala
//: details: vflow spark consumer
//: author:  Satheesh Ravi
//: date:    09/01/2017

//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at

//:     http://www.apache.org/licenses/LICENSE-2.0

//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package com.oath.vdms.vflow.consumer.spark.util

object FieldMappings {
  val indexMap = Map(
    8 -> "sourceIPAddress",
    27 -> "sourceIPAddress",
    12 -> "destinationIPAddress",
    28 -> "destinationIPAddress",
    15 -> "ipNextHopIPAddress",
    62 -> "ipNextHopIPAddress",
    16 -> "bgpSourceAsNumber",
    17 -> "bgpDestinationAsNumber",
    14 -> "ingressInterface",
    10 -> "egressInterface",
    7 -> "sourceTransportPort",
    11 -> "destinationTransportPort",
    4 -> "protocolIdentifier",
    6 -> "tcpControlBits",
    1 -> "octetDeltaCount"
  )
}