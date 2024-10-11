#ifndef _AGILOR_WRAPPER_C_H_
#define _AGILOR_WRAPPER_C_H_

#ifdef __cplusplus
extern "C" {
#endif
    typedef struct {
        int64_t id;
        int64_t node_id;
        char[64] name;
        void* pt_value;
        int64_t ts;
    } ucs_pt_t;

    void agilor_ucs_pt_create(ucs_pt_t* p); //创建一个ucs点位
    void agilor_ucs_pt_drop(ucs_pt_t* p); //删除一个ucs点位
    void agilor_ucs_pt_insert(ucs_pt_t* p); //插入一个ucs点位值
    void agilor_ucs_pt_query(char* pt_id, int64_t start_time, int64_t end_time, int64_t step); //查询范围内点位值
    void agilor_ucs_pt_query_now(char* pt_id); //查询范围内点位值
    void agilor_ucs_pt_remove_before(char* pt_id, int64_t before_time) //查询范围内点位值
#ifdef __cplusplus
}
#endif

#endif // _AGILOR_WRAPPER_C_H_