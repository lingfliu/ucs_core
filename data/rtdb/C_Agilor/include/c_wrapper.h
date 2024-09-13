#include<stdint.h>
#include <stdbool.h>
#include<agilor.h>
#include<agilor_defs.h>

//const int C_SERVERNAME_LEN = 16;
//const int C_SERVERADDR_LEN = 16;
//const int C_USERNAME_LEN = 32;
//const int C_PASSWORD_LEN = 16;

//typedef bool agibool;
//#define agitrue   ((agibool)1)
//#define agifalse   ((agibool)0)

//typedef struct agilor_serverinfo_t {
//	char server_name[C_SERVERNAME_LEN];
//	char server_addr[C_SERVERADDR_LEN];
//	char username[C_USERNAME_LEN];
//	char password[C_PASSWORD_LEN];
//	agibool is_connected;
//} agilor_serverinfo_t;

//typedef int agibool;
//#define agitrue           1
//#define agifalse          0

//////////////////////////////////////////////
//////////////Agilor Connect//////////////////
/////////////////////////////////////////////

int32_t c_Agcn_Startup(uint64_t thread_id, agibool through_firewall);
int32_t c_Agcn_Connect(const char* server, const char* host_addr, const char* username,const char* password, uint32_t port);
int32_t c_Agcn_Disconnect(const char* server);
int32_t c_Agcn_Cleanup();
bool c_Agcn_ServerInfo(int32_t* server_id, agilor_serverinfo_t* server_info);

//////////////////////////////////////////////
///////////////Agilor Data //////////////////
//////////////////////////////////////////////
int32_t c_Agda_Subscribe(const char* server, const char* tags, int32_t count);
int32_t c_Agda_Unsubscribe(const char* server, const char* tags, int32_t count);
int32_t c_Agda_UnsubscribeAll(const char* server);
agibool c_Agda_GetSubscribeValue(char* server,char* tag, agilor_value_t* value);
agibool c_Agda_NextValue(agirecordset recordset,  char* tag, agilor_value_t* value,agibool removed);
agirecordset c_Agda_TimedValue(const char* server, const char* tag, int64_t start_time,int64_t end_time, int64_t step);
agirecordset c_Agda_TimedValues(const char* server, const char* tags, int32_t count, int64_t start_time, int64_t end_time, int64_t step);

//////////////////////////////////////////////
////////////Agilor Point //// ////////////////
//////////////////////////////////////////////
int32_t c_Agpt_AddPoint(const char* server, const agilor_point_t point,agibool overwrite);
int32_t c_Agpt_RemovePoint(const char* server, int32_t point_id);
agirecordset c_Agpt_DeviceInfo(const char* server);
agibool c_Agpt_NextDeviceInfo(agirecordset recordset, int32_t* device_id,agilor_deviceinfo_t* device_info);
int32_t c_Agpt_PointCount(const char* server, const char* device_name, int32_t* count);
agirecordset c_Agpt_GetPointByDevice(const char* server, const char* device_name);
agibool c_Agpt_NextPoint(agirecordset recordset, int32_t* point_id, char* tag);
int32_t c_Agpt_Point(const char* server, const char* tag, agilor_point_t* point);
int32_t c_Agpt_PointExist(const char* server, const char* tag);
int32_t c_Agpt_Tag(const char* server, int32_t point_id, char* tag);

//////////////////////////////////////////////
////////////Agilor Archive //////////////////
//////////////////////////////////////////////
//int32 Agar_Register(const char* server, const char* device_name,agibool time_sync = agifalse, const agilor_deviceconf_t* conf = NULL);
//int32 Agar_Unregister(const char* server, const char* device_name);








