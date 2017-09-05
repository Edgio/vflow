//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    IngestStream.scala
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

package com.oath.vdms.vflow.consumer.spark.driver


import com.oath.vdms.vflow.consumer.spark.model.IPFix
import com.oath.vdms.vflow.consumer.spark.util.ParseDataUtil
import org.apache.log4j.Logger
import org.apache.spark.{SparkConf, SparkContext}
import org.apache.spark.streaming.{Seconds, StreamingContext}
import org.apache.kafka.common.serialization.StringDeserializer
import org.apache.spark.streaming.kafka010._
import org.apache.spark.streaming.kafka010.LocationStrategies.PreferConsistent
import org.apache.spark.streaming.kafka010.ConsumerStrategies.Subscribe
import org.apache.spark.sql.{SparkSession}

/**
  * Created by sravi on 9/1/17.
  */
object IngestStream {
  val appName = "Vflow Ingestion"
  val sparkConf = new SparkConf().setAppName(appName)
  val logger = Logger.getLogger(getClass.getName)

  def generateSparkSession(sparkConf: SparkConf) = {
    SparkSession
      .builder()
      .config(sparkConf)
      .enableHiveSupport()
      .getOrCreate()
  }


  def main(args: Array[String]): Unit = {

    //Initilize context and stream
    val sc = new SparkContext(sparkConf)
    val streamingContext = new StreamingContext(sc, Seconds(20))
    System.setProperty("spark.hadoop.dfs.replication", "1")

    val argsLen = args.length

    //Get valid number of arguments
    argsLen match {
      case 5 => {
        val topics = Array(args(0))
        val bootstrap = args(1)
        val groupID = args(2)
        val writeFormat: String = args(3)
        val tablename: String = args(4)
        val sparkSession = generateSparkSession(sparkConf)
        val kafkaParams = Map[String, Object](
          "bootstrap.servers" -> bootstrap,
          "key.deserializer" -> classOf[StringDeserializer],
          "value.deserializer" -> classOf[StringDeserializer],
          "group.id" -> groupID,
          "auto.offset.reset" -> "latest",
          "enable.auto.commit" -> (true: java.lang.Boolean)
        )
        val stream = KafkaUtils.createDirectStream[String, String](
          streamingContext,
          PreferConsistent,
          Subscribe[String, String](topics, kafkaParams)
        )
        import sparkSession.implicits._
        val tf = stream.flatMap(ParseDataUtil.parseRec)
        tf.foreachRDD(record => {
          record.toDS().write.format(writeFormat).mode(org.apache.spark.sql.SaveMode.Append).saveAsTable(tablename)
        })
        streamingContext.start()
        streamingContext.awaitTermination()
      }
      case _ => logger.error("Invalid Argument Count. Please give arguments in this order: topic bootstrap_server group_id format tablename")
    }
  }
}
