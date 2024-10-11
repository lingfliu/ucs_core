#include <stdio.h>
#include <cstdint>
#include <stdbool.h>
#include <cstddef>
#include "../../../../无人工地测试/agilor/include/agilor_defs.h"

typedef struct agilor_point_t {
    char tag[64];         // 测点标签 *
    char descriptor[32];  // 测点描述 #
    char engunit[16];     // 测点数据单位（安培、摄氏度等） #
    int32 id;                        // 测点编号，由系统配置
    uint8 type;                      // 菜单类型(R浮点数/S字符串/B开关/L整形/E枚举) *
    uint8 scan;             // 测点扫描标识(0或>=0x80："禁止"， 1："输入", 2："输出" *
    float32 typical_value;  //典型值 #
    union {                 // 点值 #
        // float64 dval;                  // 双精度浮点//[agilor_point_t has not]
        // int64 llval;                   // 64位长整//[agilor_point_t has not]
        float32 rval;         // 浮点
        int32 lval;           // 长整
        agibool bval;         // 开关
        char sval[128];  // 字符串
        struct {
            short type;  // 0x0001：使用key, 0x0002：使用name，0x0003表示同时使用key,name
            short key;   // 枚举(值)
            char name[128];  // 枚举(字符串)
        } eval;                          // 枚举
    };
    char enum_desc[128];  // 枚举描述 （"2:1,2,on:0,3,off"），暂时无用，[hp has not]
    int64 timedate;                  // 时间戳
    int32 state;                     // 点状态（点的质量、实时点值、缓冲的点值）
    // 由系统配置，覆盖添加时=old.state
    char point_source[32];  // 测点的数据源站(设备名) *
    char source_group[32];   // 测点的数据源结点组 #
    char source_tag[128];     // 测点的源标签 *

    float32 upper_limit;                  //数据上限，用于压缩
    float32 lower_limit;                  //数据下限，用于压缩

    uint16 push_ref1;                     //实时推理规则标志 #
    uint16 rule_ref1;                     //实时推理规则标志 #

    // Exception reporting
    // Exception reporting ensures that a Agilor interface only sends meaningful
    // data, rather than sending unnecessary data that taxes the system.
    // 异常报告可确保Agilor接口只发送有意义的数据，而不是发送不必要的数据，从而加重系统的负担。

    // Exception reporting uses a simple deadband algorithm to determine whether
    // to send events to Agilor Data Archive. For each point, you can set
    // exception reporting specifications that create the deadband. The interface
    // ignores values that fall inside the deadband.
    // 异常报告使用一个简单的死区算法来确定是否将事件发送到PI数据存档。对于每一点，可以设置创建死区的异常报告规范。该接口忽略死区内的值。
    // TODO: exc_xxx这3个参数，只对接口有效，还是内核中也使用这3个参数？
    int64 exc_min;       // 实时数据处理最短间隔（接口参数）
    int64 exc_max;       //实时数据处理最大间隔（内核参数）
    //不管是否压缩，点值是否变化，当timedate-last_timedate >= exc_max时强制存储数据
    float32 exc_dev;     // 实时数据处理偏差（接口参数）：
    // 当fabs(tagvalue.rval) < fabs(lptag->rval) * (1 - lptag->exc_dev)
    // 或(fabs(tagvalue.rval) > fabs(lptag->rval) * (1 - lptag->exc_dev)
    // 表示点值变化超过偏差。这时当点值变化超过偏差且与上次发送的点值时间戳之差>=exc_min
    // 时，即使是过滤发生，也会将点值发送到内核。

    uint16 alarm_type;   // 报警类型
    uint16 alarm_state;  // 状态报警
    float32 alarm_hi;    // 上限报警
    float32 alarm_lo;    // 下限报警
    float32 alarm_hihi;
    float32 alarm_lolo;

    uint16 priority_hi;  // 报警优先级，暂时不处理
    uint16 priority_lo;
    uint16 priority_hihi;
    uint16 priority_lolo;

    agibool archive;   // 是否存储历史数据
    agibool compress;  // 是否进行历史压缩 *，但type=O时，compress=agifalse
    uint8 step;        // 历史数据的插值形式（线形，台阶），compress=agitrue时有效
    int32 his_idx;     // 历史记录索引号，由系统配置

    // Compression testing
    int64 comp_min;       // 压缩最短间隔(压缩最小时间), compress minimum time，暂时无用
    int64 comp_max;       // 压缩最长间隔(压缩最大时间), compress maximum time，暂时无用
    float32 comp_dev;     // 压缩灵敏度（压缩偏差）， compress deviation
    // 归档时压缩灵敏度=(upper_limit - lower_limit) * comp_dev

    float32 last_val;     // 上次数据存档的值 #  // TODO: 应该使用内存中的？
    int64 last_timedate;  // 上次数据存档的时间 # // TODO: 应该使用内存中的？
    int64 create_date;    // 采集点创建日期，由系统配置
} agilor_point_t;



typedef union {
    float32 rval; // 浮点
    int32_t lval; // 长整
    bool bval;    // 开关
    char sval[128];  // 字符串
    byte* blob_data; // blob字节数组
}Anydata;


typedef struct agilor_value_t {
    int64_t timedate;                  // 时间戳
    int32_t state;                     // 状态 (Agpt_SetPointValue不需要设置state)
    uint8_t type;                      // 点值类型
    int32_t blob_size;
    Anydata data
} agilor_value_t;

// 定义 DataMeta 结构体
typedef struct DataMeta{
    int Dimen;//维度
    int ByteLen;//字节长度
    char* Alias; // 代号
    char* Unit;//单位
    int DpClass;
    int DnClass;

    int DataClass;//数据类型
    bool Msb;
} DataMeta;

// 定义 DPoint 结构体


/*
c_Dpoint 为与go端Dpoint对应的结构
*/
typedef struct c_DPoint {
    int64_t Ts;       // 时间戳
    char tag[32];
    char ParentTag[32];
    //+DataType DataType;
    Anydata Data[3];
    DataMeta* Meta;   // 元数据
} c_DPoint;


//将字符串指针复制给字符数组
void charToArr(char* source,char arr[32] ) {
    if (source == NULL) {
        // 如果源字符串为 NULL，则清空目标数组
        arr[0] = '\0';
        return;
    }
    // 使用 strncpy 复制字符串，最多复制 31 个字符（留出一个位置给终止符）
    strncpy(arr, source, sizeof(arr) - 1);

    // 确保目标数组以 null 字符结尾
    arr[sizeof(arr) - 1] = '\0';
}

agilor_point_t DPointToAgipoint(agilor_point_t point,c_DPoint* dp){
    
    point.timedate = dp->Ts;
    // 赋值 source_tag 字段
    strncpy(point.source_tag, dp->tag, sizeof(point.source_tag) - 1);
    point.source_tag[sizeof(point.source_tag) - 1] = '\0'; // 确保字符串以 null 结尾
    // 赋值 point_source 字段
    strncpy(point.point_source, dp->ParentTag, sizeof(point.point_source) - 1);
    point.point_source[sizeof(point.point_source) - 1] = '\0'; 
    charToArr(dp->Meta->Alias, point.descriptor);
    
 
}

bool Insert(const char* server, agilor_point_t point,c_DPoint* dp, agibool IsOverwrite) {
    DPointToAgipoint(point, dp);
    if (Agpt_PointExist(server,point.tag)) {
        Agpt_SetPointValue(server,tag,point);
    }
    else
    {
        Agpt_AddPoint(server, point, IsOverwrite);
    }
    
}