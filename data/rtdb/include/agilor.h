// /////////////////////////////////////////////////////////////////////////////
// ///////////////////////////// AgilorDB API ///////////// ////////////////////
// /////////////////////////////////////////////////////////////////////////////
#ifndef AGILORDB_SRC_COMMON_AGILOR_H_
#define AGILORDB_SRC_COMMON_AGILOR_H_

#include "agilor_defs.h"

#ifdef __cplusplus
extern "C" {
#endif

//****** 数据点报警状态位 *********
const uint16 ALARM_TYPE_HILIMIT_MASK = 0x0001;    //高报警
const uint16 ALARM_TYPE_LOLIMIT_MASK = 0x0002;    //低报警
const uint16 ALARM_TYPE_HIHILIMIT_MASK = 0x0004;  //高高报警
const uint16 ALARM_TYPE_LOLOLIMIT_MASK = 0x0008;  //低低报警
const uint16 ALARM_TYPE_SWITCHON_MASK = 0x0010;   //开报警
const uint16 ALARM_TYPE_SWITCHOFF_MASK = 0x0020;  //关报警

//****** 启动数据采集标志 *********
const uint8 SCAN_INPUT = 0x01;   //输入允许
const uint8 SCAN_OUTPUT = 0x02;  //输出允许
// const uint8 SCAN_DISABLE		= 0x80;	//禁止I/O
const uint8 SCAN_DISABLE_MUSK = 0x80;  //禁止(128)

#ifndef _WIN32  // linux
#define WM_USER                         0x0400
#define MYSIG_SUBDATAARRIVAL SIGUSR1  // monitor the arrival of sub data
#define MYSIG_DISCONNECTED SIGUSR2    // monitor the connection with server
#else
#define MYSIG_SUBDATAARRIVAL WM_USER + 101  // monitor the arrival of sub data
#define MYSIG_DISCONNECTED WM_USER + 102    // monitor the connection with server
#endif

#define WM_DRTDB_STATE_INFORMATION WM_USER + 1000
#define WM_DRTDBAPI_STATE_CONNECTED WM_USER + 1001
#define WM_DRTDBAPI_STATE_CONNECTING WM_USER + 1002
#define WM_DRTDBAPI_STATE_RECONNECTING WM_USER + 1003
#define WM_DRTDBAPI_STATE_DISCONNECTED WM_USER + 1004
#define WM_DRTDBAPI_STATE_SENDBUFFERDATA WM_USER + 1005
#define WM_DRTDBAPI_STATE_SENDBUFFERDATA_END WM_USER + 1006
#define WM_DRTDBAPI_STATE_SENDREALTIMEDATA WM_USER + 1007
#define WM_DRTDBAPI_STATE_SENDDATATOBUFFER WM_USER + 1008
#define WM_DRTDBAPI_STATE_TAGINFORM WM_USER + 1009
#define WM_DRTDBAPI_STATE_SERVER_REQUEST_DISCONNECT WM_USER + 1010
#define WM_DRTDBAPI_STATE_SENDDATA WM_USER + 1011

enum AggregateFunction {
  AF_SUMMARY,
  AF_MINIMUM,
  AF_MAXIMUM,
  AF_AVERAGE,
  AF_COUNT,
  AF_SUMMARY_FOR_CONTINOUS,
  AF_AVERAGE_FOR_CONTINOUS,
};

typedef struct agilor_serverinfo_t {
  char server_name[C_SERVERNAME_LEN];
  char server_addr[C_SERVERADDR_LEN];
  char username[C_USERNAME_LEN];
  char password[C_PASSWORD_LEN];
  agibool is_connected;
} agilor_serverinfo_t;

typedef struct agilor_deviceinfo_t {
  char device_name[C_DEVICENAME_LEN];
  agibool is_online;
  int32 point_count;
} agilor_deviceinfo_t;

typedef struct agilor_value_t {
  int64 timedate;                  // 时间戳
  int32 state;                     // 状态 (Agpt_SetPointValue不需要设置state)
  uint8 type;                      // 点值类型
  int32 blob_size;
  union {                          // 点值
    float32 rval;                  // 浮点
    int32 lval;                    // 长整
    agibool bval;                  // 开关
    char sval[C_STRINGVALUE_LEN];  // 字符串
    byte* blob_data;               // blob字节数组
  };
} agilor_value_t;

typedef struct agilor_blob_t {
  int64 timedate;       // timestamp
  int32 state;          // 查询快照时，state是点的状态，查询历史数据时，blob是blob数据所在的文件索引
  int32 amount;         // number of bytes to be send to server
  byte* buffer;         // the buffer containing the data to be send to server
} agilor_blob_t;

typedef struct agilor_point_t {
  char tag[C_TAGNAME_LEN];         // 测点标签 *
  char descriptor[C_TAGDESC_LEN];  // 测点描述 #
  char engunit[C_TAGUNIT_LEN];     // 测点数据单位（安培、摄氏度等） #
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
    char sval[C_STRINGVALUE_LEN];  // 字符串
    struct {
      short type;  // 0x0001：使用key, 0x0002：使用name，0x0003表示同时使用key,name
      short key;   // 枚举(值)
      char name[C_STRINGVALUE_LEN];  // 枚举(字符串)
    } eval;                          // 枚举
  };
  char enum_desc[C_ENUMDESC_LEN];  // 枚举描述 （"2:1,2,on:0,3,off"），暂时无用，[hp has not]
  int64 timedate;                  // 时间戳
  int32 state;                     // 点状态（点的质量、实时点值、缓冲的点值）
                                   // 由系统配置，覆盖添加时=old.state
  char point_source[C_DEVICENAME_LEN];  // 测点的数据源站(设备名) *
  char source_group[C_GROUPNAME_LEN];   // 测点的数据源结点组 #
  char source_tag[C_SOURCETAG_LEN];     // 测点的源标签 *

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

typedef struct agilor_devicepoint_t {
  int32 local_id;                    // 本地重新分配的测点id
  int32 id;                          // 测点id
  char source_tag[C_SOURCETAG_LEN];  // 测点的源标签
  float32 exc_dev;
  int64 exc_min;
  int64 exc_max;
  uint16 type;
  uint16 scan;
  int64 timedate;
  int32 state;
  union {
    float32 rval;
    int32 lval;
    agibool bval;
    char sval[C_STRINGVALUE_LEN];
  };
} agilor_devicepoint_t;

// 设备连接时的初始化参数
typedef struct agilor_deviceconf_t {
  agibool throw_firewall;           // 是否穿过防火墙
  agibool save_pointtable_to_file;  // 是否把点表保存到本地
  int32 reconnect_interval;            // 重连间隔
  int32 reconnect_max_time;                // 最大重连次数
  int32 socket_resend_max_time;            // socket重发数据最大次数
  agibool save_raw_data;            // 是否保存原始数据
  int32 max_raw_data_count;                // 最大文件数量
  int32 max_raw_data_size;         // 文件最大容量
  agibool send_not_changed_data;    // 是否发送没有变化的数据
  int32 send_not_changed_data_interval;
  int32 data_buffer_size;
  // TODO:
  agibool is_data_need_buffer;       // 网络断开时，是否需要缓存数据
  agibool is_buffer_data_need_send;  // 是否需要发送缓存内容
  char device_pointtable_path[C_FILE_PATH_LEN];
  char raw_datafile_path[C_FILE_PATH_LEN];
} agilor_deviceconf_t;

////////////////////////////////////////////////////////////////////////
///////////////////////// Agilor Connect API ///////////////////////////
////////////////////////////////////////////////////////////////////////

// 0: startup success
// -1: startup failed
AGILOR_EXPORT int32 Agcn_Startup(uint64 thread_id = 0, agibool through_firewall = agifalse);

// 0: success
// -1: has connected yet
// -2: create node failed
// -3: invalid node state
// -5: already connected
// -6: create socket failed
// -7/-8: network error
// -501: username or password incorrect
AGILOR_EXPORT int32 Agcn_Connect(const char* server, const char* host_addr, const char* username,
                                 const char* password, uint32 port = 3955);

// 0: success
// -2: remove node failed
// -3: invalid node state
// -4: has not connected
AGILOR_EXPORT int32 Agcn_Disconnect(const char* server);

AGILOR_EXPORT int32 Agcn_Cleanup();

// enumerate the server nodes, which has been created
// if return ture, lpSvrInfo pointers to a valid struct containing related information
// nServerID = 0 indicates the start of a new iterative process
// 0: success
// -3: invalid node state
AGILOR_EXPORT bool Agcn_ServerInfo(int32* server_id, agilor_serverinfo_t* server_info);

////////////////////////////////////////////////////////////////////////
///////////////////////// Agilor Point API /////////////////////////////
////////////////////////////////////////////////////////////////////////

// 0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -503: has no permission to modify point
// 其他负数：内核添加失败
AGILOR_EXPORT int32 Agpt_AddPoint(const char* server, const agilor_point_t& point,
                                  agibool overwrite = agitrue);

// 0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -503: has no permission to modify point
AGILOR_EXPORT int32 Agpt_RemovePoint(const char* server, int32 point_id);

// >0: success, record set id
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -502: has no permission to view point
AGILOR_EXPORT agirecordset Agpt_DeviceInfo(const char* server);

// enumerate the device information returned by querydeviceinfo
// nDeviceID must better be set to zero at beginning
AGILOR_EXPORT agibool Agpt_NextDeviceInfo(agirecordset recordset, int32* device_id,
                                          agilor_deviceinfo_t* device_info);

// 0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -502: has no permission to view point
// 1019: device not exist
AGILOR_EXPORT int32 Agpt_PointCount(const char* server, const char* device_name, int32* count);

// >0: success, record set id
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -502: has no permission to view point
// 1019: device not exist
AGILOR_EXPORT agirecordset Agpt_GetPointByDevice(const char* server, const char* device_name);

// >0: success, record set id
// 0：success，but can't find matching point
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -502: has no permission to view point
AGILOR_EXPORT agirecordset Agpt_GetPointByTagMask(const char* server, const char* tag_mask);

// enumerate the tagname returned by querytags
AGILOR_EXPORT agibool Agpt_NextPoint(agirecordset recordset, int32* point_id, char* tag);

// 0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -502: has no permission to view point
AGILOR_EXPORT int32 Agpt_Point(const char* server, const char* tag, agilor_point_t* point);

// 0: success（存在一个点有修改权限，但内核不一定能settagvalue成功）
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -502: has no permission to view point
// -503: has no permission to modify point
AGILOR_EXPORT int32 Agpt_SetPointValue(const char* server, const char* tag,
                                       const agilor_value_t& value, agibool manual,
                                       const char* comment = NULL);

// >0: success, tag id
// -3: invalid node state
// -4: has not connected
// -5: tag not exist
// -201: send request failed
// -211: wait timed out
// -502: has no permission to view point
AGILOR_EXPORT int32 Agpt_PointExist(const char* server, const char* tag);

// 0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect id
// -502: has no permission to view point
AGILOR_EXPORT int32 Agpt_Tag(const char* server, int32 point_id, char* tag);

////////////////////////////////////////////////////////////////////////
///////////////////////// Agilor Data API //////////////////////////////
////////////////////////////////////////////////////////////////////////

// >0: success, nmber of points successfully subscribed
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag, or tag has been subscribed
// -110: exceeds api limits per time
// -111: exceeds server limits per time
// -112: exceeds server remaining limits per time
// -502: has no permission to view point
// -506: has no permission to subscribe point
// 接口中判断只要有一个点id获取失败，就返回，都获取成功才发给内核;内核中，只要存在没有权限的点，就不会在内核中订阅任何点，直接返回
AGILOR_EXPORT int32 Agda_Subscribe(const char* server, const char* tags, int32 count);

// >0: success, （client发送给server合法tag的数量）
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: all tags are incorrect
// 内核会将所有传入的点，都取消订阅（如：传入了Simu1_1, Simu1_2,
// Simu1_3,但是之前只订阅了Simu1_1，返回值是3不是1）
AGILOR_EXPORT int32 Agda_Unsubscribe(const char* server, const char* tags, int32 count);

// >0: success, （server成功取消订阅的点数）
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
AGILOR_EXPORT int32 Agda_UnsubscribeAll(const char* server);
// agitrue: get sub val success
// agifalse: get sub val failed
// when notified on sub data arrival, call this funtion to get the data
AGILOR_EXPORT agibool Agda_GetSubscribeValue(char* server, char* tag, agilor_value_t* value);

// >0: success, record set id
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -110: 超过单次查询数量上限（MAX_NUMBER_OF_SERVER_COMPUTE_PER_TIME = 2010）
// -502: has no permission to view point
// -504: has no permission to view snapshot
AGILOR_EXPORT agirecordset Agda_Snapshot(const char* server, const char* tags, int32 count);


// =0: success
// -1: 查询出错（查询到的blob快照数量不正确）
// -3: invalid node state
// -4: has not connected
// -111: 查询超过了内核单次处理上限
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -502: has no permission to view point
// -504: has no permission to view snapshot
AGILOR_EXPORT int32 Agda_BlobSnapshot(const char* server, const char* tags,
                                      int32 *count, agilor_blob_t** values);

// >0: success, record set id
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -502: has no permission to view point
// -505: has no permission to view history data
AGILOR_EXPORT agirecordset Agda_TimedValue(const char* server, const char* tag, int64 start_time,
                                           int64 end_time, int64 step = 0);

// >0: success, record set id
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -502: has no permission to view point（存在没有查看点权限的点）
// -505: has no permission to view history
// data（所有点都没有查看历史数据的权限）
AGILOR_EXPORT agirecordset Agda_TimedValues(const char* server, const char* tags, int32 count,
                                            int64 start_time, int64 end_time, int64 step = 0);

// query blob historical data
AGILOR_EXPORT int32 Agda_TimedBlob(const char* server, const char* tag, int64 start_time,
                                           int64 end_time, int64 step,
                                           int32 *count, agilor_blob_t** values);

AGILOR_EXPORT agibool Agda_NextValue(agirecordset recordset, char* tag, agilor_value_t* value,
                                     agibool removed = agitrue);

// 所有查询到的blob数据（快照、历史数据）都必须手动调用Agdb_FreeBlob进行清理，否则会出现内存泄漏
AGILOR_EXPORT void Agda_FreeBlob(int32 count, agilor_blob_t* values);

// server　端统计：内核直接返回统计结果
// 统计指定点（单个点，且type只能是'R'/'L'）指定时间段的数值信息
// 内核会返回统计结果的time、state、value
// 统计的数值结果根据type的不同，分别存储在pTagVal的rval, lval, bval中
// TODO: BOOL类型也可以统计？？？？

// =0: success
// -3: invalid node state
// -4: has not connected
// -201: send request failed
// -211: wait timed out
// -101: incorrect tag
// -103: incorrect statistic type
// -502: has no permission to view point
// -505: has no permission to view history data
// -801: server comput failed 内核统计失败（错误，没有统计到结果）
AGILOR_EXPORT int32 Agda_TimedValueStatistic(const char* server, const char* tag, int64 start_time,
                                             int64 end_time, int32 statistic_type,
                                             agilor_value_t* value);

// client 统计：先获取查询数据，再统计
// 统计指定点（单个点，且type只能是'R'/'L'）指定时间段的数值信息

// 统计类型为（MAX / MIN）时，输出参数：
// pTagVal->rval: 统计值
// lpTagVal->timedate: hRecordset中最后一条数据的time
// lpTagVal->state: hRecordset中最后一条数据的state

// 统计类型为（SUM / AVG）时，输出参数：
// pTagVal->rval: 统计值
// lpTagVal->timedate: hRecordset中最后一条数据的time
// lpTagVal->state: hRecordset中数据量

// 统计类型为（COUNT）时，输出参数：
// pTagVal->lval: 统计值
// lpTagVal->timedate: hRecordset中最后一条数据的time
// lpTagVal->state: hRecordset中数据量

// =0: success
// -103: incorrect statistic type
// -104: incorrect record set id
AGILOR_EXPORT int32 Agda_TimedValueAggregate(agirecordset recordset, char* tag,
                                             agilor_value_t* value,
                                             int32 statistic_type = AF_SUMMARY,
                                             agibool removed = agitrue);

// agitrue: success
// agifalse: close record set failed(结果集id不存在)
AGILOR_EXPORT agibool Agda_CloseRecordset(agirecordset recordset);

////////////////////////////////////////////////////////////////////////
///////////////////////// Agilor Archive API ///////////////////////////
////////////////////////////////////////////////////////////////////////

// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -22: device already connected
// -23: get finish event failed
// -201: send request failed
// -211: wait timed out
// -507: has no permission to update real time data
// 1019: device not exist
// 1002: sys not init
// 1004: sys not start
// 1020: device has linked by current user
// 1022: device link exceeds limits
AGILOR_EXPORT int32 Agar_Register(const char* server, const char* device_name,
                                  agibool time_sync = agifalse,
                                  const agilor_deviceconf_t* conf = NULL);

// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device already disconnected
// -23: get finish event failed
// -201: send request failed
// -211: wait timed out
AGILOR_EXPORT int32 Agar_Unregister(const char* server, const char* device_name);

// 发送单点值给Agilor服务器,可以设置在点信息没有变化时是否需要更新并发送点信息到Agilor
// 会通过本地点表设置传入的point_value的一些点属性(id, type)
// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
// -23: get finish event failed
// -31: local device pt compute failed(update local device pt failed)
// -35: send data to buffer failed
// -201: send request failed
// -211: wait timed out

// AGILOR_EXPORT int32 Agar_PutValue(agilor_pointvalue_t *point_value,
//                                   agibool filtered = agifalse);

AGILOR_EXPORT int32 Agar_PutValue(const char* server, const char* source_tag,
                                  agilor_value_t* value, agibool filtered = agifalse);

// 根据id直接发送新的点值，不过滤，且不会更改point_value中的点属性
// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
// -23: get finish event failed
// -35: send data to buffer failed
// -201: send request failed
// -211: wait timed out
// AGILOR_EXPORT int32 Agar_PutValueById(const agilor_pointvalue_t
// &point_value);
AGILOR_EXPORT int32 Agar_PutValueById(const char* server, int32 point_id,
                                      agilor_value_t* value);
// register
// device后，会在本地保存该device上的全部点信息，该接口用于获取本地保存的点信息
// local_point_id为开始查询的位置，输出设备点信息数组索引为"local_point_id ~
// total_device_point_count"之间的所有的点
// local_point_id：查询点信息时，设备点信息数组的起始索引，若local_point_id=0，则查找出该设备所有的点

// AGILOR_EXPORT int32 Agar_PutBlob(const char* server, int32 point_id, agilor_blob_t& data);

// 0：后续没有点了
// 1：后续仍有点，可继续查询
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
// -31: local device pt compute failed(get node from local device pt failed)
AGILOR_EXPORT int32 Agar_NextDevicePoint(const char* server, int32* local_point_id,
                                         agilor_devicepoint_t* device_point);

// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
// -31: local device pt compute failed(set device status to local device pt
// failed) -35: send data to buffer failed -201: send request failed -211: wait
// timed out
AGILOR_EXPORT int32 Agar_SetDeviceStatus(const char* server, agibool device_connected);

// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
// -30: invalid local device pt
// -31: local device pt compute failed(get point from local device pt failed)
// -32: local device pt invalid input
// -35: send data to buffer failed
// -201: send request failed
// -211: wait timed out
AGILOR_EXPORT int32 Agar_GetDevicePoint(const char* server, const char* source_tag,
                                        agilor_devicepoint_t* device_point);

// =0: success
// -3: invalid node state
// -4: has not connected
// -20: invalid device manager
// -21: device has not connect
AGILOR_EXPORT int32 Agar_SetCallbackFn(
    const char* server, void (*on_add_point)(const agilor_devicepoint_t& device_point),
    void (*on_remove_point)(const agilor_devicepoint_t& device_point),
    void (*on_set_devicepoint_value)(int32 point_id, const char* source_tag,
                                     const agilor_value_t& tag_val),
    void (*on_get_devicepoint_value)(agilor_devicepoint_t* device_point));

////////////////////////////////////////////////////////////////////////
///////////////////////// Agilor Time API //////////////////////////////
////////////////////////////////////////////////////////////////////////

AGILOR_EXPORT int64 Agtm_Time();

// 获取远程服务端当前时间戳
// =0: success
// -3: invalid node state
// -4: has not connected
AGILOR_EXPORT int32 Agtm_ServerTime(const char* server, int64* server_time);

// convert time to a int64 integer as second count since 1970.1.1
AGILOR_EXPORT int64 Agtm_Time2Long(int32 hour, int32 min, int32 sec, int32 millisec, int32 year,
                                   int32 mon, int32 day);

// convert unix timestamp to YYYY-MM-DD HH:MM:SS.mmm
AGILOR_EXPORT void Agtm_Long2String(int64 date_time, char* data_str);

#ifdef __cplusplus
} /* end extern "C" */
#endif

#endif
