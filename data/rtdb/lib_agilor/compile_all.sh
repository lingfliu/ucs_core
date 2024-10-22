#编译封装库
gcc -I./include -fpic -shared ./src/agilor_wrap.c -L./lib -lagilor -o ./lib/libagilor_wrap.so

#编译测试程序
gcc -I./include -L./lib ./src/test_con.c -lagilor -lagilor_wrap -Wl,-rpath,./lib -o Ctest
