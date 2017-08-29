## Clickhouse, Apache Kafka
ClickHouse is an open source column-oriented database management system capable of real time generation of analytical data reports using SQL queries. ClickHouse's performance exceeds comparable column-oriented DBMS currently available on the market. It processes hundreds of millions to more than a billion rows and tens of gigabytes of data per single server per second. ClickHouse uses all available hardware to it's full potential to process each query as fast as possible. The peak processing performance for a single query (after decompression, only used columns) stands at more than 2 terabytes per second. (https://clickhouse.yandex/)
![Alt text](/docs/imgs/clickhouse.jpeg?raw=true "vFlow")
The below clickhouse setup needs a zookeeper server, replica server is optional.

### Configuration (/etc/clickhouse-server/config.xml)
Configure at least a zookeeper host (replica server is optional)

```
<zookeeper>
    <node index="1">
        <host>zk001</host>
        <port>2181</port>
    </node>
    <session_timeout_ms>1000</session_timeout_ms>
</zookeeper>

<remote_servers>
    <logs>
        <shard>
            <weight>1</weight>
            <internal_replication>false</internal_replication>
            <replica>
                <host>CLICKHOUSE_SRV1</host>
                <port>9000</port>
            </replica>
        </shard>
    </logs>
</remote_servers>    
```

### Create Database
```
CREATE DATABASE vflow
```

### Create Table
```
CREATE TABLE vflow.samples
(
    date Date,
    time DateTime,
    device String,
    src String,
    dst String,
    srcASN UInt64,
    dstASN UInt64,
    proto UInt8
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/1/', '1', date, (src, time), 8192);
```
### Build Kafka Consumer
```
go get -d ./...
go build main.go
```

### Benchmark Details
I tried it with two clickhouse servers, one for ingest and one for query. they had below hardware information and you can see the below results based on the above database configuration.

Hardware
- CPU Intel Core Processor (Haswell, no TSX) cores = 8, 2.6GHz, x86_64
- Memory 16G
- Drive SSD in software RAID

![Alt text](/docs/imgs/clickhouse_s1.png?raw=true "vFlow")
![Alt text](/docs/imgs/clickhouse_s2.png?raw=true "vFlow")
