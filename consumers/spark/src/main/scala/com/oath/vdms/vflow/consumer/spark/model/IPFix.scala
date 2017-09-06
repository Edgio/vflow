//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    IPFix.scala
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


package com.oath.vdms.vflow.consumer.spark.model

import java.sql.Timestamp

@SerialVersionUID(19900528L)
case class IPFix(
                  var exportTime: Timestamp,
                  var sourceIPAddress: String,
                  var destinationIPAddress: String,
                  var ipNextHopIPAddress: String,
                  var bgpSourceAsNumber: Long,
                  var bgpDestinationAsNumber: Long,
                  var ingressInterface: Long,
                  var egressInterface: Long,
                  var sourceTransportPort: Int,
                  var destinationTransportPort: Int,
                  var protocolIdentifier: Int,
                  var tcpControlBits: String,
                  var octetDeltaCount: Long) {

  import java.sql.Timestamp

}
