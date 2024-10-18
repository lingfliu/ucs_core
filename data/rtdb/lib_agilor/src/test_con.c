#include"agilor_wrap.h"
#include <stdio.h>
#include <stdlib.h>


int main() {

	   //  初始化库
	   int32_t startupResult =c_Agcn_Startup(0, agifalse);
	   if (startupResult != 0) {
	       printf("Startup failed with error code %d\n", startupResult);
	       return 0;	   
	  }

	   //数据库连接参数
	   const char* server = "issc-dev-lingfliu.lan";
	   const char* host_addr = "192.168.66.27";
	   const char* username = "test";
	   const char* password = "123";
	   int port = 3955; 

	   //尝试连接数据库
	 int32_t connectResult = c_Agcn_Connect(server, host_addr, username, password, port);


	/////////////////////////////////////////////
	/////////////////Agda//////////////////////
	////////////////////////////////////////////
/*	
	//订阅 查询订阅数据
	c_Agda_Subscribe(server,"R_point1",1);
	char server1[]="issc-dev-lingfliu.lan";
	agilor_value_t test_val = {0};
	char tag_Agda[64];
	c_Agda_GetSubscribeValue(server1,tag_Agda, &test_val);
	printf("timedate:%lld,state:%d,type:%d\n",test_val.timedate,test_val.state,test_val.type);
	//c_Agda_Unsubscribe(server, "R_point2", 3);

	//c_Agda_NextValue(125,"R_point1",&test_val,false);
	
**/
//查询实时库中所有设备信息(其中device_info还包含其他的信息未输出)
	agirecordset  recordset =  c_Agpt_DeviceInfo(server);
	int32_t device_id=0;
	agilor_deviceinfo_t device_info;
	while(Agpt_NextDeviceInfo(recordset, &device_id, &device_info)){
	 printf("返回的设备信息：设备id:%d,设备名：%s，设备测点数量：%d\n",device_id,device_info.device_name,device_info.point_count);
	}

//查寻指定设备上的测点数量
	int32_t count;
	char name[5]="DV2";
	const char* device_name=name;
	c_Agpt_PointCount(server,device_name, &count);


//查询指定设备上所有的测点信息	
	char tag[128];
	int32_t point_id ;
	const char* device_name1=name;
	agirecordset  recordset1 = Agpt_GetPointByDevice(server, device_name1);
	c_Agpt_NextPoint(recordset1, &point_id, tag);

	agilor_point_t point;
	c_Agpt_Point(server,tag, &point);
	printf("测点大小：%lu\n",sizeof(point));
	printf("测点的详细信息：%s,描述：%s,数据单位：%s,测点编号：%d.\n",point.tag,point.descriptor,point.engunit,point.id);

//查询实时数据库服务端信息
	agilor_serverinfo_t server_info;
	int32_t server_id=0;
	if(c_Agcn_ServerInfo(&server_id,&server_info)){
		printf("查询到实时数据库信息：server_id：%d,server_name:%s,server_addr:%s,username:%s,password:%s.\n",server_id,server_info.server_name,server_info.server_addr,server_info.username,server_info.password);
	}  



/*

int64_t start_time = 1672531200; // 2023-01-01 00:00:00 UTC

    int64_t end_time = 1672534800;   // 2023-01-01 01:00:00 UTC

    int64_t step = 3600;             // 1 hour


    agirecordset result = c_Agda_TimedValue(server, tag, start_time, end_time, step);

    if (result > 0) {

        printf("c_Agda_TimedValue 请求成功，记录集 ID: %lld\n",result);

    } else {

        printf("c_Agda_TimedValue 请求失败。\n");

    }


    // 测试 c_Agda_TimedValues

    const char* tags = "sensor1,sensor2,sensor3";

    int32_t count1 = 3;

    result = c_Agda_TimedValues(server, tags, count1, start_time, end_time, step);

    if (result > 0) {

        printf("c_Agda_TimedValues 请求成功，记录集 ID: %lld\n",result);

    } else {

        printf("c_Agda_TimedValues 请求失败。\n");

    }
**/

/////////////////////////////////////
////////10.13测试addPoint////////
////////////////////////////////////
/*
ucs_pt_t p = {};
int value = 2; 
strncpy(p.tag, "testPoint_1", sizeof(p.tag) - 1);
    p.tag[sizeof(p.tag) - 1] = '\0'; 
    strncpy(p.descrip, "测试点1", sizeof(p.descrip) - 1);
    p.descrip[sizeof(p.descrip) - 1] = '\0';  // 确保字符串以 null 结尾
    // 赋值其他成员
    p.id = 3;
    p.node_id = 123;
    p.pt_value =&value;  // 初始化为 NULL 或者指向有效的数据
    p.ts =1728798010;
    
   // agilor_ucs_pt_create( &p);

agilor_ucs_pt_drop(&p);
**/

/////////////////////////////////////
///////////////insert///////////////
///////////////////////////////////

/*
printf("11111111");
ucs_pt_t p = {};
uint8_t v = 66;
strncpy(p.tag, "R_point1001", sizeof(p.tag) - 1);
    p.tag[sizeof(p.tag) - 1] = '\0'; 
    strncpy(p.descrip, "测试点1", sizeof(p.descrip) - 1);
    p.descrip[sizeof(p.descrip) - 1] = '\0';  // 确保字符串以 null 结尾
    // 赋值其他成员
    //p.id = 3;
    p.node_id = 123;
    p.pt_value = &v;  
    p.ts =1728798017;
    agilor_ucs_pt_insert(&p);
**/
//若连接成功则断开连接
	   if (connectResult == 0) {
		int disconnectResult =c_Agcn_Disconnect(server);
		if (disconnectResult != 0) {
		printf("Failed to disconnect from the database.\n");
	 }
	   }
	   //清理资源
	   c_Agcn_Cleanup();

return 0;

}




