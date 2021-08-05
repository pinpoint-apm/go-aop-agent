
#include "Inst.h"
#include <stdio.h>

void inst_str(Inst* inst,char* buf,int buf_size)
{
	int i = 0;
	if(buf == NULL || inst == NULL || buf_size <=0){
		return ;
	}
	//parse prefix
	char *pbuf = buf;
	int done = 0;
	for (i = 0; i < PREFIX_SIZE; i++) {
		Prefix p =inst->Prefix[i]; 
		if(p == 0){
			break;
		}

		if ( (p&PrefixImplicit) !=0){
			continue;
		}
		done = snprintf(pbuf,buf_size,"%hu",p);
		if(done<buf_size){
			buf_size -= done;
			pbuf += done;
		}else{
			return;
		}
	}
	//parse op
	done = snprintf(pbuf,buf_size,"%s ",op_name(inst->Op));
    if(done>= buf_size){
        return;
    }

    pbuf+=done;
    buf_size+=done;
	//parse args
	for (i = 0; i < ARGS_SIZE ; i++) {
		if( IS_NIL(inst->Args[i] )){
			break;
		}
        char argsBuf[128]={0};
		args_str(&inst->Args[i],argsBuf,sizeof(argsBuf));
        done = snprintf(pbuf,buf_size,"%s,",argsBuf);
		
        if(done>= buf_size || buf_size <=0 ){
            return;
        }
        pbuf+=done;
        buf_size+=done;
	}
}