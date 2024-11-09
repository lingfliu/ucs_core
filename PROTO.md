@ -1,132 +0,0 @@
# 1 Data modeling
We define the UCS with three level components:
1. Project
2. Complex (static) and Machs (mobile)
3. Nodes and Cams

## 1.1 Nodes and Points (节点与点位)
考虑对复杂监测与控制场景，节点可能包含多种类型的点位，而点位可以是多维的，因此，对节点-点位的设定如下：  

1. 将点位分为两类：  
    1. DNode-Dpoint 监测点位
    2. CtlNode-CtlPoint 控制点位

2. 节点是对来自同一设备的点位的集合，同一节点对点位数据的采集是同步的，且每次采样长度各个点位一样。

3. 每一个点位的数值是同一类型数据

## 1.2 Data and Messages (数据与消息)

数据同样分为DData （监测数据）与CtlData （控制数据）

数据传递的载体为Message，为减少数据解析与存储的复杂度，数据内应当存储数据的属性信息，包括：

1. Node ID / Addr / Point Offset ID与地址
2. DataMeta 数据属性
3. 值
4. Ts & Idx 时间索引
5. Pos (Optional) 位置
6. Session (Optional) 会话ID，该属性主要用于节点存在不同业务流程的时候进行分辨

Data以Node为单元，一次Data传输应当包括该Node上所有的有效DPoint。  

Message作为对上述数据的编码，可以划分为两种：
1. 基于Json的消息
2. 基于二进制的报文

### 1.2.1 Json数据格式
1. 监测数据样本
```json
{
    "node_id": 1,
    "addr": "1.2.3.4:8080",
    "point": [
        {
            "offset": 0,
            "data_meta": 
            {
                "alias": "前臂关节角",
                "type": "int32",
                "byte_len": 4,
                "dimen": 3,
            }
            "data": [1,2,3]
        },
        {
            "offset": 1,
            "data_meta": 
            {
                "alias": "横臂关节角",
                "type": "int32",
                "byte_len": 4,
                "dimen": 3,
            }
            "data": [4,5,6]
        }
    ],
    "ts": 17800000201,
    "idx": 4, //索引，optional
    "session": "default", //会话信息, optional
    "pos": [1,2,3], //位置信息，optional
    "mode": 1, //见源码定义
    "sps": 100, //samplelen > 1时有效
    "sample_len": 1 
}
```

Json消息格式与上述一致，区别为"data"字段为字节流

### 1.2.2 二进制数据与消息格式 (未启用)
二进制按上述Json数据格式进行二进制编码，并在报头部分填写荷载长度，各点位维度与尺度等信息


| ||||||
|---|---|---|---|---|---|
|HEADER（4）|TYPE (2)|META_LEN (2)|TS(8)|IDX(2)|
|MODE(1)|SPS（2）|SAMPLE_LEN(2)|NUM_PT(1)|NUM_PT(1)|META_PT（VAR）|
|DATA_PT(VAR)||||  

META_PT格式:  
|||||
|---|---|---|---|
|OFFSET|TYPE(1)|BYTE_LEN(1)|DIMEN(1)|

# 2 CRUD 操作
## 2.1 规则
1. 点位为非独立数据，DPoint与CtlPoint必须基于Dnode与CtlNode进行操作
2. 优先基于模板创建

## 2.2 模板编辑
支持两种模板：
1. Mach 模板

可编辑部分：

    1. ID
    2. Name
    3. Addr
    4. Cam
    5. Node

2. Node / CAM 模板  
可编辑部分：  

    1. Name
    2. Addr

# 3 TAOS数据库操作规范

TAOS数据库采用单点位表存储，每个表存储一个点位数据，超表命名规则：
```
节点类型_点位类型
```

如果为固定节点，则column与tag为
```
(TS, V1, V2, ..., V_K) Tag(node_id, point_offset, alias, unit, x, y, z)
```

如果为移动节点，则column与tag为
```
(TS, X, Y, Z, V1, V2, ..., V_K) Tag(node_id, point_offset, alias, unit)
```