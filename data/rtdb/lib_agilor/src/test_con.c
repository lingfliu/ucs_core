#include"agilor_wrap.h"
#include <stdio.h>
#include <stdlib.h>
#include <time.h>       // 包含 time() 函数的头文件
#include <sys/time.h>

int main() {

	   //  初始化库
	   int32_t startupResult =c_Agcn_Startup(0, agifalse);
	   if (startupResult != 0) {
	       printf("Startup failed with error code %d\n", startupResult);
	       return 0;	   
	  }

	   //数据库连接参数
	   const char* server = "Agilor";
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
    void PointNumber(){
                agirecordset  recordset =  c_Agpt_DeviceInfo(server);
	int32_t device_id=0;
	agilor_deviceinfo_t device_info;
	while(Agpt_NextDeviceInfo(recordset, &device_id, &device_info)){
	 printf("返回的设备信息：设备id:%d,设备名：%s，设备测点数量：%d\n",device_id,device_info.device_name,device_info.point_count);
	}
}
/*
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


**/




/////////////////////////////////////
////////测试addPoint////////
////////////////////////////////////

/*
printf("*******addPoint测试********\n");

    agilor_point_t pt = {};
    
    pt.type = 'L';
    strncpy(pt.tag, "testPoint_14", sizeof(pt.tag) - 1);
    pt.tag[sizeof(pt.tag) - 1] = '\0'; 
    strncpy(pt.point_source, "DV3",sizeof(pt.point_source) - 1);
    pt.point_source[sizeof(pt.point_source) - 1] = '\0'; 
    strncpy(pt.source_tag, "testPoint_5", sizeof(pt.source_tag) - 1);
    pt.source_tag[sizeof(pt.source_tag) - 1] = '\0';  

    strncpy(pt.descriptor, "测试点", sizeof(pt.descriptor) - 1);
    pt.descriptor[sizeof(pt.descriptor) - 1] = '\0';
    pt.archive = agitrue;
    pt.compress = agifalse;
    pt.scan = SCAN_OUTPUT;
    c_Agpt_AddPoint(server, &pt,agitrue);
    PointNumber();
**/
 

//PS:tag相同的点位会被覆盖
printf("*******createPoint测试********\n");
ucs_pt_t p = {};
strncpy(p.tag, "testPoint_144", sizeof(p.tag) - 1);
    p.tag[sizeof(p.tag) - 1] = '\0'; 
    time_t create_now;
    time(&create_now);
    p.ts =(int64_t)create_now*1000;
    agilor_ucs_pt_create( &p);
    PointNumber();


/*
//根据点位ID删除
printf("*******dropPoint测试********\n");
    ucs_pt_t p = {};
    p.id = 3000;
    agilor_ucs_pt_drop(&p);
    PointNumber();
**/
/*
printf("*******insertPoint测试********\n");
//agilor_value_t  p = {};
ucs_pt_t p = {};
strncpy(p.tag, "testPoint_14", sizeof(p.tag) - 1);
   p.tag[sizeof(p.tag) - 1] = '\0'; 
    double a=999;
    p.pt_value =&a;  
    time_t insrt_now;
    time(&insrt_now);
    p.ts =(int64_t)insrt_now*1000;
    agilor_ucs_pt_insert(&p);
**/

printf("*******QuaryByTime测试********\n");

    int64_t start_time =1708581463414;
    time_t now;
    time(&now);
    int64_t end_time =(int64_t)now*1000;; 
    int64_t step = 0;            
    char tag[64] = "testPoint_144";
    agirecordset result = c_Agda_TimedValue(server, tag, start_time,end_time,step);
    if (result > 0) {
    char tag1[64];
    agilor_value_t value;
        while(c_Agda_NextValue(result,tag1,&value,agitrue)){
	printf("tagName:%s,time:%lld,数据值value:%lf\n",tag1,value.timedate,value.rval);
        }
    } else {
        printf("c_Agda_TimedValue 请求失败。\n");
}


printf("********Snapshot测试***********\n");
    char tags[32] = "testPoint_144";
     int32_t count =1;
     agirecordset  res1 = c_Agda_Snapshot(server,tags,count);
     if (res1 > 0) {
     char tag[32];
     agilor_value_t value;
     while(c_Agda_NextValue(res1,tag,&value,agitrue)){
	printf("tagName:%s,time:%lld,数据值value:%lf\n",tag,value.timedate,value.rval);
        }
    } else {
        printf("Snapshot请求失败。\n");
}


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




