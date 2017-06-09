## Clickhouse, Apache Kafka

### Create Database
```
CREATE DATABASE vflow
```

### Create Table
```
CREATE TABLE samples
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
Hardware
- CPU Intel Core Processor (Haswell, no TSX) cores = 8, 2.6GHz, x86_64
- Memory 16G
- Drive SSD in software RAID

![Alt text](/docs/imgs/clickhouse_s1.png?raw=true "vFlow")
![Alt text](/docs/imgs/clickhouse_s2.png?raw=true "vFlow")
