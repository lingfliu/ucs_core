# Unmanned Construction System - CORE services & middlewares

## Framework
```
conn           connection engine
    | conn      
    | coder    data protocol coder
    | srv      workflow handling
data           
    | model    data modeling
    | orm      orm accessor
    | olap     space-time olap
    | stats    stats aggregate
    | file     object-relationship file service
    | export   data2file export
etl
    | membuff
    | st_slice  spatial-temporal data formatting
    | flow      Flink based stream computation
dd              data distribution
    | packet    no-RT msg dispatching
    | stream    data stream distribution (RT DD)
    | mapping   convert stream to another
mq              message queue
alg
    alg         algorithm wrapper
    algNode     alg node server
    node_mgr    alg node server manager (load balance, reverse proxy)
utils           tools
cfg             service, log, tmp output config
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