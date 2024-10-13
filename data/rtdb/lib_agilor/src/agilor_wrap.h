#ifndef _AGILOR_WRAPPER_C_H_
#define _AGILOR_WRAPPER_C_H_

#ifdef __cplusplus
extern "C" {
#endif
    #include <stdint.h>
    #include <stdbool.h>
    #include<agilor.h>
    #include<agilor_defs.h>
    #include<stdio.h>
    #include <string.h>
	
    const int32_t Step = 0;
    const agibool removed = agitrue;//是否移除记录集
    const agibool overwrite = agifalse;//添加点位遇到tag相同的点是否覆盖

    typedef struct {
        long long id;
        int64_t node_id;
        char descrip[32];
        char tag[64];
        void* pt_value;
        int64_t ts;
    } ucs_pt_t;

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
////////////Add New Funtion //////////////////
//////////////////////////////////////////////
    agilor_point_t ucsptToAgilorPt(ucs_pt_t* p);//ucs_pt_t转换
	
    void agilor_ucs_pt_create(ucs_pt_t* p,const char* server); //创建一个ucs点位
    void agilor_ucs_pt_drop(ucs_pt_t* p,const char* server); //删除一个ucs点位
    void agilor_ucs_pt_insert(ucs_pt_t* p); //插入一个ucs点位值
    int agilor_ucs_pt_query(char* tag, int64_t start_time, int64_t end_time, int64_t step, ucs_pt_t* p_list); //查询范围内点位值
    void agilor_ucs_pt_query_now(char* tag, ucs_pt_t* pt); //查询范围内点位值
    void agilor_ucs_pt_remove_before(char* tag, int64_t before_time); //删除范围内点位值
#ifdef __cplusplus
}
#endif

#endif // _AGILOR_WRAPPER_C_H_