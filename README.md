# Unmanned Construction System - CORE services & middlewares

## Framework
```
web     |       web services for sys admin
conn    |       connection services
coder   |       (not used) general coder for messages 
data    |       database clients 
    | buff      redis & membuff
    | file      minio & large file services
    | orm       mysql & mongodb
    | rtdb      agilor & taos
    | olap      (not used) clickhouse 
    | stats     (not used) tools for data aggregation
dd      |       data distribution services
    | Fastdds   Agilor boosted Fast-DDS with TSN support
    | MQTT 
model   |       data models
dao     |       
ulog    |       logging tools
utils          
test           
mock
scripts         for DevOps 
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

## CGO compling
In case of C function calling, perform the following commands:
```bash
# gcc compling
gcc -c ${SOURCEFILE} -o ${OUTPUT}
gcc -shared -o lib${LIBNAME} ${OUTPUT}

# go build 
go build -o ${EXEC_NAME}
```

Note that the ```go build``` must explicitly declare the output excutable file (```-o ${OUTPUT}```) so the library paths defined in the source file effective.


## Licence
All rights reserved.

## Contact
- lingfliu@gmail.com
- lingfeng.liu@163.com