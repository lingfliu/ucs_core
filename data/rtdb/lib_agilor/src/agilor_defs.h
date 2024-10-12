#ifndef AGILORDB_SRC_COMMON_AGILOR_DEFS_H_
#define AGILORDB_SRC_COMMON_AGILOR_DEFS_H_

#undef NULL
#define NULL 0

// for linux
#ifndef _WIN32
typedef int SOCKET;
#endif
#include<stdbool.h>
typedef char int8;
typedef short int16;
typedef int int32;
typedef long long int64;

typedef unsigned char uint8;
typedef unsigned short uint16;
typedef unsigned int uint32;
typedef unsigned long long uint64;

typedef unsigned char byte;

typedef float float32;
typedef double float64;

typedef bool agibool;
#define agitrue           ((agibool)1)
#define agifalse          ((agibool)0)

typedef long long agirecordset;  // 结果集

// 暂时没用到
// typedef unsigned long long TIMESTAMP;
// typedef int HRESULT;

/*
//****** 关键字段长度定义 ******
const uint32 C_PROJECT_NAME_LEN = 64;
const uint32 C_FILE_PATH_LEN = 1024;

const uint32 C_FULL_TAGNAME_LEN = 80;  // SERVERNAME.TAGNAME, sucha as LGCAG.ZL_AI1001

const uint32 C_USERNAME_LEN = 32;
const uint32 C_PASSWORD_LEN = 16;
const uint32 C_TAGNAME_LEN = 64;  // maybe some tags on different server have the same name
const uint32 C_TAGDESC_LEN = 32;
const uint32 C_TAGUNIT_LEN = 16;
const uint32 C_DEVICENAME_LEN = 32;
const uint32 C_GROUPNAME_LEN = 32;
const uint32 C_STRINGVALUE_LEN = 128;
const uint32 C_SOURCETAG_LEN = 128;  // the physical tag on devices
const uint32 C_ENUMDESC_LEN = 128;
*/
// 点状态
const int32 STAT_QUALITY_GOOD = 0x00000010L;
const int32 STAT_QUALITY_BAD = 0x00000020L;
const int32 STAT_QUALITY_UNKNOWN = 0x00000040L;

#if !defined(AGILOR_EXPORT)

// #if defined(AGILOR_SHARED_LIBRARY)
#if defined(_WIN32)  // for windows
#define AGILOR_EXPORT __declspec(dllexport)
#else  // defined(_WIN32)  // for linux
#define AGILOR_EXPORT __attribute__((visibility("default")))
#endif  // defined(_WIN32)

// #else  // defined(AGILOR_SHARED_LIBRARY)
// #define AGILOR_EXPORT
// #endif

#endif  // !defined(AGILOR_EXPORT)

#endif
