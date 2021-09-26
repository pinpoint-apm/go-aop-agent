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
#include "Args.h"
#include <stdio.h>
#include <assert.h>
ArgsT mk_imm_arg(Imm x){
	ArgsT args = {.type =E_IMM};
	args.value.imm=x;
	// assert(x!=1);
	return args;
}

ArgsT mk_rel_arg(Rel x){
	ArgsT args = {.type =E_REL};
	args.value.rel=x;
	return args;
}

ArgsT mk_reg_arg(Reg x){
	ArgsT args = {.type =E_REG};
	args.value.reg=x;
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

Reg* get_reg_arg(ArgsT* args)
{
	if(args->type == E_REG){
		return &args->value.reg;
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

Rel* get_rel_arg(ArgsT* args)
{
	if(args->type == E_REL){
		return &args->value.rel;
	}
	return NULL;
}

void args_str(ArgsT* args,char* buf ,int size)
{
	switch(args->type){
		case E_IMM:
			snprintf(buf,size,"%lx",args->value.imm);
		break;
		case E_MEM:
		{
			Mem* m = get_mem_arg(args);
			if(m == NULL) return ;
			char scaleBuf[32]={0};
			char dispBuf[32]={0};
			const char*baseStr = "" ,*plus="" ,*index = "",*scale = "";
			if (m->Base !=0){
				baseStr = reg_name(m->Base);
			}

			if(m->Scale !=0 ){
				if (m->Base !=0){
					plus ="+";	
				}
				if(m->Scale >1){
					snprintf(scaleBuf,32,"%u",m->Scale);
					scale = scaleBuf;
				}
				index = reg_name(m->Index);

			}

			if( m->Disp !=0 || (m->Base == 0 &&  m->Scale == 0)) {
				snprintf(dispBuf,32,"+0x%lx",m->Disp);
			}
			snprintf(buf,size,"[%s%s%s%s%s]",baseStr,plus,scale,index,dispBuf);
			break;
		}
		case E_REG:
		{
			snprintf(buf,size,"%s",reg_name(args->value.reg));
			break;
		}
		case E_REL:
		{
			snprintf(buf,size,"+%d",args->value.rel);
			break;
		}
		case E_NIL:
			snprintf(buf,size,"nil");
		break;
	}
}


ArgsT fixedArg[113] ={
	[xArg1]=    {.type=E_IMM, .value.imm=1},
	[xArg3]=    {.type=E_IMM, .value.imm=3},
	[xArgAL]=   {.type=E_REG, .value.reg = AL},
	[xArgAX]=   {.type=E_REG, .value.reg =AX},
	[xArgDX]=   {.type=E_REG, .value.reg =DX},
	[xArgEAX]=  {.type=E_REG, .value.reg =EAX},
	[xArgEDX]=  {.type=E_REG, .value.reg =EDX},
	[xArgRAX]=  {.type=E_REG, .value.reg =RAX},
	[xArgRDX]=  {.type=E_REG, .value.reg =RDX},
	[xArgCL]=   {.type=E_REG, .value.reg =CL},
	[xArgCS]=   {.type=E_REG, .value.reg =CS},
	[xArgDS]=   {.type=E_REG, .value.reg =DS},
	[xArgES]=   {.type=E_REG, .value.reg =ES},
	[xArgFS]=   {.type=E_REG, .value.reg =FS},
	[xArgGS]=   {.type=E_REG, .value.reg =GS},
	[xArgSS]=   {.type=E_REG, .value.reg =SS},
	[xArgST]=   {.type=E_REG, .value.reg =F0},
	[xArgXMM0]= {.type=E_REG, .value.reg =X0},
};
