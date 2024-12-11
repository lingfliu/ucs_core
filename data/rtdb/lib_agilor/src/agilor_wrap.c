#include "agilor_wrap.h"


////////////////////////////////////////////////////////////////////////
/////////////////////////// Agilor Connect C API ///////////////////////////
//////////////////////////////////////////////////////////////////////////
//
int32_t c_Agcn_Startup(uint64_t thread_id, agibool through_firewall) {
    int32_t Startup_result = Agcn_Startup(thread_id, through_firewall);
    switch(Startup_result){
        case 0:
            printf("启动成功。\n");
            break;
        case -1:
            printf("启动失败。\n");
            break;
        default:
            printf("未知错误代码：%d\n", Startup_result);
            break;
    }
    return Startup_result;
}

int32_t c_Agcn_Connect(const char* server, const char* host_addr, const char* username,
        const char* password, uint32_t port)
{
	int32_t connectResult = Agcn_Connect(server, host_addr, username, password, port);
	switch (connectResult) {
		case 0:
		printf("数据库连接成功。\n");
		break;
		case -1:
		printf( "已经处于连接状态。\n");
		break;
		case -2:
		printf("创建连接点失败。\n");
		break;
		case -3:
		printf("连接点状态错误。\n");
		break;
		case -5:
		printf("已经连接到实时数据库。\n");
		break;
		case -6:
		printf( "创建socket失败。\n");
		break;
		case -7:
		case -8:
		printf("网络错误。\n");
		break;
		case -501:
		printf( "用户名和密码错误。\n");
		break;
		default:
		printf("未知错误： %d.\n", connectResult);
		break;
	}

	return connectResult;
	   
}


int32_t c_Agcn_Disconnect(const char* server) {
	int32_t Discon_result = Agcn_Disconnect(server);
	switch(Discon_result){
		case 0:
		printf("从数据库断开成功。\n");
		break;
		case -2:
		printf("移除连接点失败。\n");
		break;
		case -3:
		printf("连接点状态错误。\n");
		break;
		case -4:
		printf("没有连接到实时数据库。\n");	
	       	break;
	}
	return  Discon_result;

}

int32_t c_Agcn_Cleanup() {
	int32_t Clearnup_result = Agcn_Cleanup();
	switch(Clearnup_result){
		case 0:
		printf("成功清理资源。\n");
		break;
		case -1:
		printf("清理资源失败。\n");
		break;	
	}
	return Clearnup_result;

}

bool c_Agcn_ServerInfo(int32_t* server_id, agilor_serverinfo_t* server_info) {
	
	bool ServerInfo_result = Agcn_ServerInfo(server_id, server_info);
	if(ServerInfo_result ){
		printf("查询实时数据库信息成功。\n");
	}else{
		printf("查询实时数据库信息错误。\n");
	}
	return ServerInfo_result;
}
/////////////////////////////////////////////////
///////////////////订阅/////////////////////////
int32_t c_Agda_Subscribe(const char* server, const char* tags, int32_t count) {

    int32_t result = Agda_Subscribe(server, tags, count);
    switch(result) {
        case -3:
            printf("连接点状态错误\n");
            break;
        case -4:
            printf("没有连接到实时数据库\n");
            break;
        case -201:
            printf("发送请求失败\n");
            break;
        case -211:
            printf("等待超时\n");
            break;
        case -101:
            printf("标签tag错误\n");
            break;
        case -110:
            printf("订阅测点数量超过了接口单次订阅数量限制\n");
            break;
        case -111:
            printf("订阅测点数量超过了实时数据库单次订阅数量的限制\n");
            break;
        case -112:
            printf("订阅测点数量超过了实时数据库剩余订阅数量的限制\n");
            break;
        case -502:
            printf("没有查看测点的权限\n");
            break;
        case -506:
            printf("没有订阅测点的权限\n");
            break;
        default:
            if (result > 0) {
                printf("成功，成功订阅的测点数量: %d\n", result);
            } else {
                printf("未知错误 %d.\n", result);
            }
            break;
    }

    return result;

}


int32_t c_Agda_Unsubscribe(const char* server, const char* tags, int32_t count) {
    int32_t result = Agda_Unsubscribe(server, tags, count);
    switch(result) {
        case -3:
            printf("连接点状态错误\n");
            break;
        case -4:
            printf("没有连接到实时数据库\n");
            break;
        case -201:
            printf("发送请求失败\n");
            break;
        case -211:
            printf("等待超时\n");
            break;
        case -101:
            printf("标签tag错误\n");
            break;
        default:
            if (result > 0) {
                printf("成功，发送给实时数据库的合法tag的数量:%d\n", result);
            } else {
                printf("未知错误 %d.\n", result);
            }
            break;
    }
    return result;
}


int32_t c_Agda_UnsubscribeAll(const char* server) {
    int32_t result = Agda_UnsubscribeAll(server);
    switch(result) {
        case -3:
            printf("连接点状态错误\n");
            break;
        case -4:
            printf("没有连接到实时数据库\n");
            break;
        case -201:
            printf("发送请求失败\n");
            break;
        case -211:
            printf("等待超时\n");
            break;	
	default:
	 if (result > 0) {		
	  printf("成功，实时数据库成功取消订阅的测点数量:%d\n", result);
	} else {
	 printf("未知错误 %d.\n", result);
	 }
	break;
    }
    return result;
}


agibool c_Agda_GetSubscribeValue(char* server, char* tag, agilor_value_t* value) {
    agibool result = Agda_GetSubscribeValue(server, tag, value);
    if (result == true) {
        printf("枚举订阅测点实时数据成功\n");
    } else {
        printf("枚举订阅测点实时数据失败\n");
    }
    return result;
}

//removed从该记录集中取出所有数据后，是否清除该记录集
agibool c_Agda_NextValue(agirecordset recordset, char* tag, agilor_value_t* value, agibool removed) {
    agibool result = Agda_NextValue(recordset, tag, value, removed);
    if (result == true) {
        printf("成功枚举数据\n");
    } else {
        printf("枚举结束\n");
    }
    return result;
}

agirecordset c_Agda_TimedValue(const char* server, const char* tag, int64_t start_time, int64_t end_time, int64_t step)
{
    agirecordset result = Agda_TimedValue(server, tag, start_time, end_time, step);
    
    if (result > 0) {
        printf("查询成功，记录集 ID: %lld\n", result);
    } else {
        switch (result) {
            case -3:
                printf("连接状态错误。\n");
                break;
            case -4:
                printf("尚未连接实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -101:
                printf("标签错误。\n");
                break;
            case -502:
                printf("没有查看测点的权限。\n");
                break;
            case -505:
                printf("没有查看测点历史数据的权限。\n");
                break;
            default:
                printf("未知错误：%lld。\n", result);
                break;
        }
    }

    return result;
}

agirecordset c_Agda_TimedValues(const char* server, const char* tags, int32_t count, int64_t start_time, int64_t end_time, int64_t step)
{
    agirecordset result = Agda_TimedValues(server, tags, count, start_time, end_time, step);
    
    if (result > 0) {
        printf("查询成功，记录集 ID: %lld\n", result);
    } else {
        switch (result) {
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("尚未连接到实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -101:
                printf("标签tag错误。\n");
                break;
            case -502:
                printf("存在测点没有查看该测点的权限。\n");
                break;
            case -505:
                printf("对于所有测点都没有权限查看历史数据。\n");
                break;
            default:
                printf("未知错误：%lld。\n", result);
                break;
        }
    }

    return result;
}

agirecordset c_Agda_Snapshot(const char* server, const char* tags, int32_t count){
    agirecordset result = Agda_Snapshot(server,tags,count);
    if (result > 0) {
        printf("查询成功，记录集 ID: %lld\n", result);
    } else {
        switch (result) {
            case -3:
                printf("连接点状态错误\n");
                break;
            case -4:
                printf("没有连接到实时数据库\n");
                break;
            case -101:
                printf("标签tag错误\n");
                break;
            case -201:
                printf("发送请求失败\n");
                break;
            case -211:
                printf("等待超时\n");
                break;
            case -502:
                printf("没有查看测点的权限\n");
                break;
            case -504:
                printf("没有查看测点快照的权限\n");
                break;
            default:
                printf("未知错误：%lld\n", result);
                break;
        }
    }
    return result;
}


//////////////////////////////////////////////
////////Agilor Point Function ////////////////
//////////////////////////////////////////////

int32_t c_Agpt_AddPoint(const char* server, const agilor_point_t *point,agibool overwrite) {
    int32_t AddPoint_result = Agpt_AddPoint(server, point, overwrite);
    switch(AddPoint_result){
        case 0:
            printf("成功添加点。\n");
            break;
        case -3:
            printf("连接点状态错误。\n");
            break;
        case -4:
            printf("未连接到实时数据库。\n");
            break;
        case -201:
            printf("发送请求失败。\n");
            break;
        case -211:
            printf("等待超时。\n");
            break;
        case -503:
            printf("无权限修改点。\n");
            break;
        default:
            if (AddPoint_result < 0)
                printf("服务端错误。%d\n",AddPoint_result);
            else
                printf("未知错误：%d\n", AddPoint_result);
            break;
    }
    return AddPoint_result;
}

int32_t c_Agpt_RemovePoint(const char* server, int32 point_id) {
    int32_t RemovePoint_result = Agpt_RemovePoint(server, point_id);
    switch(RemovePoint_result){
        case 0:
            printf("成功移除id：%d的点。\n",point_id);
            break;
        case -3:
            printf("连接点状态出错。\n");
            break;
        case -4:
            printf("未连接到实时数据库。\n");
            break;
        case -201:
            printf("发送请求失败。\n");
            break;
        case -211:
            printf("等待超时。\n");
            break;
        case -503:
            printf("无权限修改点。\n");
            break;
        default:
            printf("未知错误：%d\n", RemovePoint_result);
            break;
    }
    return RemovePoint_result;
}

    agirecordset c_Agpt_DeviceInfo(const char* server) {
         agirecordset DeviceInfo_result = Agpt_DeviceInfo(server);
         agirecordset recordset = -1; // 防止出错时该值无结果
         if(DeviceInfo_result > 0){
            printf("设备信息获取成功，记录集ID：%lld\n", DeviceInfo_result);
            recordset = DeviceInfo_result;
    } else {
        switch(DeviceInfo_result){
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("未连接到实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -502:
                printf("无权限查看点。\n");
                break;
            default:
                printf("未知错误：%lld\n",DeviceInfo_result);
                break;
        }
    }
    return recordset;
}

agibool c_Agpt_NextDeviceInfo(agirecordset recordset, int32_t* device_id,agilor_deviceinfo_t* device_info) {
    agibool NextDeviceInfo_result = Agpt_NextDeviceInfo(recordset, device_id, device_info);
    if(NextDeviceInfo_result == agitrue){
        printf("成功获取设备信息。\n");
    } else {
        printf("获取设备信息失败。\n");
    }
    return NextDeviceInfo_result;
}

int32_t c_Agpt_Tag(const char* server, int32_t point_id,char* tag) {
    int32_t Tag_result = Agpt_Tag(server, point_id, tag);
    switch (Tag_result) {
        case 0:
            printf("标签查询成功。\n");
            break;
        case -3:
            printf("连接点状态错误。\n");
            break;
        case -4:
            printf("未连接到实时数据库。\n");
            break;
        case -201:
            printf("发送请求失败。\n");
            break;
        case -211:
            printf("等待超时。\n");
            break;
        case -101:
            printf("测点ID不正确。\n");
            break;
        case -502:
            printf("无权查看该点。\n");
            break;
        default:
            printf("未知错误：%d\n", Tag_result);
            break;
    }

    return Tag_result;
}


int32_t c_Agpt_PointCount(const char* server, const char* device_name, int32_t* count) {
    int32_t PointCount_result = Agpt_PointCount(server, device_name, count);
    if (PointCount_result== 0) {
        printf("点计数获取成功，点数：%d\n", *count);
     } else {
        switch (PointCount_result) {
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("未连接到实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -502:
                printf("无权限查看点。\n");
                break;
            case 1019:
                printf("设备不存在。\n");
                break;
            default:
                printf("未知错误：%d\n", PointCount_result );
                break;
        }
    }
    return PointCount_result;
}

agirecordset c_Agpt_GetPointByDevice(const char* server, const char* device_name) {
    agirecordset recordset =-1;//防止出错时记录集无结果
    agirecordset GetPointByDevice_result = Agpt_GetPointByDevice(server, device_name);
    if (GetPointByDevice_result > 0) {
        recordset =GetPointByDevice_result;
        printf("设备点获取成功，记录集ID：%lld\n", recordset);
    } else {
        switch (GetPointByDevice_result ) {
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("未连接到实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -502:
                printf("无权限查看点。\n");
                break;
            case 1019:
                printf("设备不存在。\n");
                break;
            default:
                printf("未知错误：%lld\n",GetPointByDevice_result);
                break;
        }
    }
    return recordset;
}
agibool c_Agpt_NextPoint(agirecordset recordset, int32_t* point_id, char* tag) {
    agibool NextPoint_result= Agpt_NextPoint(recordset, point_id, tag);
    if (NextPoint_result == agitrue) {
        printf("成功获取测点标签，测点ID：%d，标签：%s\n", *point_id, tag);
    } else {
        printf("获取测点标签失败。\n");
    }
    return NextPoint_result;
}

int32_t c_Agpt_Point(const char* server, const char* tag, agilor_point_t* point) {
    int32_t Point_result = Agpt_Point(server, tag, point);
    if (Point_result == 0) {
        printf("测点信息获取成功。\n");
    } else {
        switch (Point_result) {
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("未连接到实时数据库。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -101:
                printf("标签tag错误。\n");
                break;
            case -502:
                printf("无权限查看测点。\n");
                break;
            default:
                printf("未知错误：%d\n", Point_result );
                break;
        }
    }
    return Point_result;
}

int32_t c_Agpt_PointExist(const char* server, const char* tag) {
    int32_t PointExist = Agpt_PointExist(server, tag);
    if (PointExist > 0) {
        printf("测点存在，标签ID：%d\n",PointExist);
    } else {
        switch (PointExist ) {
            case -3:
                printf("连接点状态错误。\n");
                break;
            case -4:
                printf("未连接到实时数据库。\n");
                break;
            case -5:
                printf("测点不存在。\n");
                break;
            case -201:
                printf("发送请求失败。\n");
                break;
            case -211:
                printf("等待超时。\n");
                break;
            case -502:
                printf("无权限查看测点。\n");
                break;
            default:
                printf("未知错误：%d\n", PointExist );
                break;
        }
    }
    return PointExist;
}

agirecordset c_Agpt_GetPointByTagMask(const char* server,char*  tag_mask){
	agirecordset result = Agpt_GetPointByTagMask(server,tag_mask);
	switch(result) {
        		case 0:
            		printf("成功，但没有找到匹配的数据点\n");
            		break;
        		case -3:
            		printf("连接点状态错误\n");
            		break;
        		case -4:
            		printf("没有连接到实时数据库\n");
            		break;
        		case -201:
            		printf("发送请求失败\n");
            		break;
        		case -211:
            		printf("等待超时\n");
           		break;
		case -502:
            		printf("没有查看测点权限\n");
           		break;
        		default:
            		if (result > 0) {
                	printf("查询成功，返回记录集为:%lld\n", result);
		} else {
                	           printf("未知错误 %lld\n", result);
            			    }
           		 	break;
		}
    	return result;

}

int32_t c_Agpt_SetPointValue(const char* server, const char* tag,const agilor_value_t *value, agibool manual,const char* comment){
    int32_t result = Agpt_SetPointValue(server,tag,value,manual,comment);
    switch(result) {
        case 0:
        printf("成功插入数据\n");
        break;
        case -3:
        printf("连接点状态错误\n");
        break;
        case -4:
        printf("没有连接到实时数据库\n");
        break;
        case -201:
        printf("发送请求失败\n");
        break;
        case -211:
        printf("等待超时\n");
        break;
        case -101:
        printf("标签tag错误\n");
        break;
        case -502:
        printf("没有查看测点的权限\n");
        break;
        case -503:
        printf("没有修改测点的权限\n");
        break;
        default:
        if (result > 0) {
        printf("查询成功，返回记录集为:%d\n", result);
	} else {
                    printf("未知错误 %d\n", result);
          }
        break;
    }
    return result;
}
/////////////////////////////////////////////
/////////////Agar Function/////////////////
////////////////////////////////////////////
int32_t c_Agar_Register(const char* server, const char* device_name,agibool time_sync,const agilor_deviceconf_t* conf){
    int32_t result = Agar_Register(server,device_name,time_sync,conf);
     switch(result) {
        case 0:
        printf("设备注册成功\n");
        break;
        case -3:
        printf("连接点状态错误\n");
        break;
        case -4:
        printf("没有连接到实时数据库\n");
        break;
        case -20:
        printf("设备管理器无效\n");
        break;
        case -22:
        printf("设备已经注册\n");
        break;
        case -23:
        printf("获取完成事件失败\n");
        break;
        case -201:
        printf("发送请求失败\n");
        break;
        case -211:
        printf("等待超时\n");
        break;
        case -507:
        printf("没有更新实时数据的权限\n");
        break;
        case 1002:
        printf("系统内部错误1002：sys not init\n");
        break;
        case 1004:
        printf("系统内部错误1004：sys not start\n");
        break;
        case 1019:
        printf("系统内部错误1019：device not exist\n");
        break;
        case 1020:
        printf("系统内部错误1020：device has linked by current user\n");
        break;
        case 1022:
        printf("系统内部错误1022：device linkage exceeds limits\n");
        break;
        default:   
         printf("未知错误 %d\n", result);
         break;
    }
    return result;
  }



int32_t c_Agar_Unregister(const char* server, const char* device_name){
     int32_t result = Agar_Unregister(server,device_name);
     switch(result) {
        case 0:
        printf("设备注册断开成功\n");
        break;
        case -3:
        printf("连接点状态错误\n");
        break;
        case -4:
        printf("没有连接到实时数据库\n");
        break;
        case -20:
        printf("设备管理器无效\n");
        break;
        case -21:
        printf("设备已断开注册\n");
        break;
        case -23:
        printf("获取完成事件失败\n");
        break;
        case -201:
        printf("发送请求失败\n");
        break;
        case -211:
        printf("等待超时\n");
        break;
        default:   
         printf("未知错误 %d\n", result);
         break;
    }
    return result;
}


///////////////////////////////////////////////
/////////////New function//////////////////
/////////////////////////////////////////////

char ucsTypeToagiType(int32_t v_type){
    switch (v_type) {
        case 9:
            return 'R'; // 浮点数
        case 0:
            return 'S'; // 字符串
        case 11:
            return 'B'; // 开关(bool)
        case 5:
            return 'L'; // 整形
        default:
            return '\0';
    }
}

int32_t agiTypeToucsType( char type){
    switch (type) {
        case 'R':
            return 9; // 浮点数
        case'S':
            return 0; // 字符串
        case 'B':
            return 11; // 开关(bool)
        case 'L':
            return 5; // 整形
        default:
            return '\0';
    }
}

ucs_pt_t agilorPtToucspt(agilor_value_t* value){

//这里由于value没有tag字段tag需要在后面补上
         ucs_pt_t  pt;
         pt.ts = value->timedate;
         pt.v_type = agiTypeToucsType(value->type);
	switch(value->type){
	case 'R':
	pt.rval = value->rval;
	break;
	case 'S':
                strncpy(pt.sval, value->sval, sizeof(pt.sval)-1);
                pt.sval[sizeof(pt.sval) - 1] = '\0';  
	break;
	case 'B':
	pt.bval= value->bval;
	break;
	case 'L':
	pt.lval= value->lval;
	break;
	}    
       return pt;    
}

agilor_point_t ucsptToAgilorPt(ucs_pt_t* p){
    agilor_point_t pt = {};
    strncpy(pt.tag, p->tag, sizeof(pt.tag)-1);
    pt.tag[sizeof(pt.tag) - 1] = '\0';  
    strncpy(pt.point_source, device_name,sizeof(pt.point_source) - 1);
    pt.point_source[sizeof(pt.point_source) - 1] = '\0'; 
    strncpy(pt.source_tag, p->tag, sizeof(pt.source_tag) - 1);
    pt.source_tag[sizeof(pt.source_tag) - 1] = '\0';
    pt.timedate = p->ts;
    pt.scan = 1;
    pt.archive = agitrue;
    pt.type = ucsTypeToagiType(p->v_type);
    return pt;
}
   
agilor_value_t ucsptToAgilorValue(ucs_pt_t* p){
    agilor_value_t value = {};
    value.type = ucsTypeToagiType(p->v_type);
    value.timedate = p->ts;
    switch(value.type){
        case 'R':
        value.rval =p->rval;
        break;
        case 'S':
        strncpy(value.sval, p->sval, sizeof(value.sval)-1);
        value.sval[sizeof(value.sval) - 1] = '\0';  
        break;
        case 'B':
        value.bval =p->bval;
        break;
        case 'L':
        value.lval  = p->lval;
        break;
    }
    return value;
}

void agilor_ucs_pt_create(ucs_pt_t* p) {
    agilor_point_t  pt = ucsptToAgilorPt(p);
//    const char* server = "Agilor";
    c_Agpt_AddPoint(server, &pt,Isoverwrite);
}

void agilor_ucs_pt_drop(ucs_pt_t* p) {
 //   const char* server = "Agilor";
    agilor_point_t point;
    char tag[64];
    strncpy(tag,p->tag, sizeof(tag) - 1);
    tag[sizeof(tag) - 1] = '\0'; 
    c_Agpt_Point(server,tag,&point);
    c_Agpt_RemovePoint(server, point.id);
}

void agilor_ucs_pt_insert(ucs_pt_t* p) {
//  const char* server = "Agilor";  
    agilor_value_t value = {};
    value = ucsptToAgilorValue(p);
// char* device_name="DV1";
    c_Agar_Register(server,device_name,agifalse,NULL);
    c_Agpt_SetPointValue(server,p->tag,&value,agifalse,NULL);
    c_Agar_Unregister(server,device_name);
}

int agilor_ucs_pt_query(char* tag, int64_t start_time, int64_t end_time, int64_t step, ucs_pt_t* pt_list) {
// const char* server = "Agilor";   
    int count = 0;
    agirecordset result = c_Agda_TimedValue(server, tag, start_time,end_time,step);
    if (result > 0) {
        char tag1[64];
        agilor_value_t value;
        while(c_Agda_NextValue(result,tag1,&value,agitrue)){
            pt_list[count] =agilorPtToucspt(&value);
            strncpy(pt_list[count].tag, tag1, sizeof(pt_list[count].tag)-1);
            pt_list[count].tag[sizeof(pt_list[count].tag) - 1] = '\0';  
            pt_list[count].v_type = agiTypeToucsType(value.type);
            count++;		
        }
     }
    return count;
}

void agilor_ucs_pt_query_now(char* tag, ucs_pt_t* pt) {
//  const char* server = "Agilor";  
     int32_t count =1;
     agirecordset  res = c_Agda_Snapshot(server,tag,count);
     if(res>0){
         char tag1[32];
         agilor_value_t value;
         while(c_Agda_NextValue(res,tag1,&value,agitrue)){
                     *pt=agilorPtToucspt(&value);
	      strncpy(pt->tag, tag1, sizeof(pt->tag)-1);
        	      pt->tag[sizeof(pt->tag) - 1] = '\0';  	
           }
     } 
}

void agilor_ucs_pt_remove_before(char* tag, int64_t before_time) {

} 