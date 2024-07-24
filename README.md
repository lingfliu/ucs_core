# Unmanned Construction System - CORE services & middlewares

## Framework
```
web: web services for sys admin
conn: connection services
    | conn      collection of different connection protocols 
    | coder     data protocol coder
data:           
    | buff      redis & membuff
    | file      minio & large file services
    | orm       mysql & mongodb
    | rtdb      agilor
    | olap      space-time olap using clickhouse
    | stats     tools for data aggregation
etl: data processing
    | st_op     spatial-temporal data handling (merge, formatting, aggregate)
dd:              data distribution services
    | fast-dds  dds with TSN support 
    | inmem dd
mq  
    | mqtt
    | zeromq
types:          data models
ulog:           logging tools
utils           tools
test:           unit tests           
mock:           mock system
scripts         script for database, mq deployment
```
## Install
Runs under go 2.16.1

To install
```bash
.\scripts\install.sh
```

To build
```bash
.\scripts\build.sh
```

To run by docker:
```bash
docker make
docker run -p 8800:8800 -p 13600:13600 -p 25433:25433
```


## Licence
All rights reserved.

## Contact
- lingfliu@gmail.com
- lingfeng.liu@163.com