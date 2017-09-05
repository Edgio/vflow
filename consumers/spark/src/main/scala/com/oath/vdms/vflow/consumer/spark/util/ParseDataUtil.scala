//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    ParseDataUtil.scala
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

import org.apache.kafka.clients.consumer.ConsumerRecord
import java.lang.reflect.Field
import java.sql.Timestamp

import com.oath.vdms.vflow.consumer.spark.model.IPFix
import org.apache.log4j.Logger

import scala.collection.immutable.{List, Map}
import scala.util.Try
import scala.util.parsing.json._
import java.text.SimpleDateFormat
import java.util.{Date, TimeZone}

object ParseDataUtil {

  val logger = Logger.getLogger(getClass.getName)


  //Parse data from stream and convert it to json object
  def parseRec(record: ConsumerRecord[String, String]): Seq[IPFix] = {
    val resultMap = JSON.parseFull(record.value()).getOrElse(Map).asInstanceOf[Map[String, String]]
    println(resultMap.get("Header").toString)
    val epochTime:Long = Try(resultMap.get("Header").getOrElse(Map()).asInstanceOf[Map[String, _]].get("ExportTime").getOrElse(0).asInstanceOf[Double].toLong).getOrElse(0)
    println(epochTime)
    val time:Timestamp = generateTimeObj(epochTime)
    val dataSetList = resultMap.get("DataSets").getOrElse(List.empty).asInstanceOf[List[List[Map[String, _]]]]
    for (dataSet <- dataSetList) yield setFields(dataSet, time)
  }

  def generateTimeObj(epochTime: Long):Timestamp = {
    val sdf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss")
    sdf.setTimeZone(TimeZone.getTimeZone("UTC"))
    Timestamp.valueOf(sdf.format(new Date((epochTime * 1000))))
  }

  //Iterate through the json to extract data and set it in the right object field
  def setFields(dateSet: List[Map[String, _]], time: Timestamp): IPFix = {
    val ipFix = new IPFix(time, "", "", "", 0, 0, 0, 0, 0, 0, 0, "", 0)
    val c: Class[_] = ipFix.getClass
    for (entry <- dateSet) {
      try {
        val index = entry.get("I").getOrElse(-1.00).asInstanceOf[Double].toInt
        val fieldName = FieldMappings.indexMap.getOrElse(index, "")
        fieldName match {
          case "" =>
          case _ => {
            val field: Field = c.getDeclaredField(FieldMappings.indexMap(index))
            field.setAccessible(true)
            field.getType.getName match {
              case "java.lang.String" => field.set(ipFix, entry.get("V").getOrElse("").toString)
              case "long" => field.set(ipFix, Try(entry.get("V").getOrElse(0L).asInstanceOf[Double].toLong).getOrElse(0))
              case "int" => field.set(ipFix, Try(entry.get("V").getOrElse(0).asInstanceOf[Double].toInt).getOrElse(0))
              case _ =>
            }
          }
        }
      }
      catch {
        case nse: NoSuchFieldException => logger.error("Unknown field " + entry.get("I").toString + "->" + entry.get("V").toString + " " + nse.toString)
        case nfe: NumberFormatException => logger.error("Unknown field " + entry.get("I").toString + "->" + entry.get("V").toString + " " + nfe.toString)
        case iae: IllegalArgumentException => logger.error("Unknown field " + entry.get("I").toString + "->" + entry.get("V").toString + " " + iae.toString)
        case cce: ClassCastException => logger.error("Unknown field " + entry.get("I").toString + "->" + entry.get("V").toString + " " + cce.toString)
      }
    }
    ipFix
  }


}
