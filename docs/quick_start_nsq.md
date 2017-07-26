# Install vFlow with NSQ - Linux Debian

## Download vFlow
``` 
wget https://github.com/VerizonDigital/vflow/releases/download/v0.4.1/vflow-0.4.1-amd64.deb
dpkg -i vflow-0.4.1-amd64.deb
```

## Download NSQ
```
wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-1.0.0-compat.linux-amd64.go1.8.tar.gz
tar -xvf nsq-1.0.0-compat.linux-amd64.go1.8.tar.gz
cp nsq-1.0.0-compat.linux-amd64.go1.8/bin/* /usr/bin
```
## NSQ - start service

```
nsqd &
```

## vFlow - NSQ config
```
echo "mq-name: nsq" >> /etc/vflow/vflow.conf
```

## vFlow - start service
```
service vflow start
```

## vFlow - Load generator
```
vflow_stress -sflow-rate-limit 1 0ipfix-rate-limit 1 &
```
## Consume IPFIX topic from NSQ
```
nsq_tail --topic vflow.ipfix -nsqd-tcp-address localhost:4150
```




