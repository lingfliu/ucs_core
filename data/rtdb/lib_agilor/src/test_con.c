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


/////////////////////////////////////
////////测试addPoint////////
////////////////////////////////////

/*
tag、ts、Type必须传
**/

/*
//PS:tag相同的点位会被覆盖
printf("*******createPoint测试********\n");
ucs_pt_t p = {};
strncpy(p.tag, "testPoint_55", sizeof(p.tag) - 1);
    p.tag[sizeof(p.tag) - 1] = '\0'; 
    p.v_type=9;
    time_t create_now;
    time(&create_now);
    p.ts =(int64_t)create_now*1000;
    agilor_ucs_pt_create( &p);
PointNumber();
**/

/*
//根据点位Tag删除
printf("*******dropPoint测试********\n");
    ucs_pt_t p = {};
    strncpy(p.tag, "testPoint_55", sizeof(p.tag) - 1);
    p.tag[sizeof(p.tag) - 1] = '\0'; 
    agilor_ucs_pt_drop(&p);
    PointNumber();
**/

/*
printf("*******insertPoint测试********\n");
ucs_pt_t p = {};

p.v_type=9;//Type:float:9;String:0;int32:5;bool:11
strncpy(p.tag, "testPoint_5", sizeof(p.tag) - 1);
   p.tag[sizeof(p.tag) - 1] = '\0';
    p.rval =99;  
    time_t insrt_now;
    time(&insrt_now);
    p.ts =(int64_t)insrt_now*1000;
    agilor_ucs_pt_insert(&p);
**/

/*
printf("*******agilor_ucs_pt_query测试********\n");

    ucs_pt_t pt_list[50];

    int64_t start_time =1730876117000;
    time_t now;
    time(&now);
    int64_t end_time =(int64_t)now*1000;; 
    int64_t step = 0;            
    char tag[64] = "testPoint_5";
    int count = agilor_ucs_pt_query(tag,start_time,end_time,step,pt_list);
    for(int i = 0;i<count;i++){
          switch(pt_list[i].v_type){
	case 9:
	printf("%d:tag:%s,time:%ld,value:%f\n",i+1,pt_list[i].tag,pt_list[i].ts,pt_list[i].rval);
	break;
	case 0:
	printf("%d:tag:%s,time:%ld,value:%s\n",i+1,pt_list[i].tag,pt_list[i].ts,pt_list[i].sval);
	break;
	case 11:
	printf("%d:tag:%s,time:%ld,value:%d\n",i+1,pt_list[i].tag,pt_list[i].ts,pt_list[i].bval);
	break;
	case 5:
	printf("%d:tag:%s,time:%ld,value:%d\n",i+1,pt_list[i].tag,pt_list[i].ts,pt_list[i].lval);
	break;
	}
    }
**/

/*
printf("********agilor_ucs_pt_query_now测试***********\n");
     char tag[64] = "testPoint_5";
     ucs_pt_t pt;
     agilor_ucs_pt_query_now(tag, &pt);
          switch(pt.v_type){
	case 9:
	printf("tag:%s,time:%ld,value:%f\n",pt.tag,pt.ts,pt.rval);
	break;
	case 0:
	printf("tag:%s,time:%ld,value:%s\n",pt.tag,pt.ts,pt.sval);
	break;
	case 11:
	printf("tag:%s,time:%ld,value:%d\n",pt.tag,pt.ts,pt.bval);
	break;
	case 5:
	printf("tag:%s,time:%ld,value:%d\n",pt.tag,pt.ts,pt.lval);
	break;
	}
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




