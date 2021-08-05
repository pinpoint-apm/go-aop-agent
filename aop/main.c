#include "goX86.h"
#include <stdio.h>
#include <assert.h>

int main(int argc,const char* argv[]){
    (void)argc;
    (void)argv;

    Inst inst = {0};
    Reg reg[]={0x48, 0x8b,0x05,0x02,0x74,0x21,0x00};
    int ret = decode(reg,sizeof(reg),&inst,64,false);
    assert(ret == 0);
    char buf[128]={0};
    inst_str(&inst,buf,sizeof(buf));
    printf("%s\n",buf);
    return 0;
}