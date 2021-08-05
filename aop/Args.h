#pragma once
#include "table.h"

typedef enum {
	E_MEM,E_IMM,E_NIL
}E_ARGS_TYPE;
typedef struct 
{
	E_ARGS_TYPE type;
	union 
	{
		Mem mem;
		Imm imm;
	} value;
}ArgsT;

ArgsT mk_imm_arg(Imm x);
ArgsT mk_mem_arg(Mem x);
ArgsT mk_nil_arg(void);
void args_str(ArgsT*,char* buf ,int size);
Imm* get_imm_arg(ArgsT* args);
Mem* get_mem_arg(ArgsT* args);
#define IS_IMM(x) (x.type == E_IMM)
#define IS_MEM(x) (x.type == E_MEM)
#define IS_NIL(x) (x.type == E_NIL)