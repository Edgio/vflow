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
