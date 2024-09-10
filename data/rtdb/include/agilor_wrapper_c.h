#ifndef _AGILOR_WRAPPER_C_H_
#define _AGILOR_WRAPPER_C_H_

extern "C" {
    void* agilor_create(const char* name);
    void agilor_destroy(void* p);
    void agilor_set(void* p, const char* name, const char* value);
    const char* agilor_get(void* p, const char* name);
    void agilor_save(void* p);
    void agilor_load(void* p);
}
#endif // _AGILOR_WRAPPER_C_H_