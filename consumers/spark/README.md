Spark consumer for vFlow data has the following components. This document assume users have basic knowledge about spark and hadoop insfrastructure

1. Spark Component (Yarn or Standalone cluster. Works with Pseudo mode also)
2. HDFS & Hive Component
3. Presto (Optional)

Vflow data from Kafka will be processed using the spark component and stored in HDFS component as Hive tables. Spark consumer is highly scalable and reliable. Consumer is tested against Spark 2.1. It uses scala 2.11.8. This consumer serves as the skeleton for processing vflow data using spark. Complex processing and analysis can be built on top of this code. With Presto (https://prestodb.io), billions of entries can be queried/joined in few minutes or seconds depending on the cluster configuration. Superset can be used for visualization (https://superset.incubator.apache.org). Consumer can also be easily modified to use other storage frameworks.

# Build
`sbt assembly` 

# Spark Submit  
`spark-submit --master <master> --class com.oath.vdms.vflow.consumer.spark.driver.IngestStream vflow_spark_consumer-assembly-1.0.jar <kafka_topic> <bootstrap_server> <consumer_group> <storage_format> <table_name>` 

# Example
`spark-submit --master spark://master:7077 --driver-memory 8G --executor-memory 4G --executor-cores 2 --conf "spark.driver.extraJavaOptions=-Dspark.hadoop.dfs.replication=1" --class com.oath.vdms.vflow.consumer.spark.driver.IngestStream vflow_spark_consumer-assembly-1.0.jar vflow.ipfix localhost:9092 consumer-group ORC ipfix-table` 
