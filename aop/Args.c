#include "Args.h"
#include <stdio.h>
ArgsT mk_imm_arg(Imm x){
	ArgsT args = {.type =E_IMM};
	args.value.imm=x;
	return args;
}

ArgsT mk_nil_arg(void){
	ArgsT args = {.type =E_NIL};
	return args;
}

ArgsT mk_mem_arg(Mem x){
	ArgsT args = {.type =E_MEM};
	args.value.mem=x;
	return args;
}

Imm* get_imm_arg(ArgsT* args)
{
	if(args->type == E_IMM){
		return &args->value.imm;
	}
	return NULL;
}

Mem* get_mem_arg(ArgsT* args)
{
	if(args->type == E_MEM){
		return &args->value.mem;
	}
	return NULL;
}

void args_str(ArgsT* args,char* buf ,int size)
{
	switch(args->type){
		case E_IMM:
			snprintf(buf,size,"%x",args->value.imm);
		break;
		case E_MEM:
			Mem* m = get_mem_arg(args);

			if(m == NULL) return ;
			char scaleBuf[32]={0};
			char dispBuf[32]={0};
			const char*baseStr =NULL ,*plus=NULL ,*index = NULL;
			if (m->Base !=0){
				baseStr = reg_name(m->Base);
			}

			if(m->Scale !=0 ){
				if (m->Base !=0){
					plus ="+";	
				}
				if(m->Scale >1){
					snprintf(scaleBuf,32,"%u",m->Scale);
				}
				index = reg_name(m->Index);
			}

			if(m->Disp !=0 || m->Base == 0 &&  m->Scale == 0) {
				snprintf(dispBuf,32,"%lld",m->Disp);
			}
			snprintf(buf,size,"[%s%s%s%s%s]",baseStr,plus,scaleBuf,index,dispBuf);
			break;
		break;
		case E_NIL:
			snprintf(buf,size,"nil");
		break;
	}
}