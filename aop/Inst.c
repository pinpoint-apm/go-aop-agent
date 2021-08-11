
/*
 * Copyright 2021 NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "Inst.h"
#include <stdio.h>
#include <string.h>

#define SUFFIX_STR ", "

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
    buf_size -=done;
	//parse args
	for (i = 0; i < ARGS_SIZE ; i++) {
		if( IS_NIL(inst->Args[i])){
			break;
		}
        char argsBuf[128]={0};
		args_str(&inst->Args[i],argsBuf,sizeof(argsBuf));
        done = snprintf(pbuf,buf_size,"%s"SUFFIX_STR,argsBuf);
		
        if(done>= buf_size || buf_size <=0 ){
            return;
        }
        pbuf+=done;
        buf_size -= done;
	}

	//rollback suffix_str
	*(pbuf - strlen(SUFFIX_STR))= '\0';
}