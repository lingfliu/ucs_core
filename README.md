# Unmanned Construction System - CORE services & middlewares

## Framework
```
admin: web pages for administration
conn: connection services
    | conn      collection of different connection protocols 
    | coder     data protocol coder
data           
    | file      minio
    | orm       mysql & mongodb
    | buff      redis
    | olap      space-time olap
    | stats     data aggregation
    | export    data2file export (to be merged)
etl: data processing
    | flink     flink computing 
    | membuff
    | st_op     spatial-temporal data handling (format, aggregate)
    | mapping   stream data mapping
dd              data distribution services
    | rtdb      agilor
    | rtps      dds / rtps
    | tunneling redirect
mq              
cfg             service, log, tmp output config
types:          data models
ulog:           logging tools
utils           tools
test            
func_test:      function tests
scripts         script for database, mq deployment
```
## How to install
To run directly on the system:
```sh
.\scripts\install.sh
.\scripts\start_services.sh
```

To run as a image:
```sh
docker make
docker run -p 8800:8800 -p 13600:13600 -p 25433:25433
```


## Licence
All rights reserved

## Contact
- lingfliu.github.com
- lingfeng.liu@163.com