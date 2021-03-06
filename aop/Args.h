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
#pragma once
#include "table.h"

typedef enum {
	E_NIL=0,E_MEM,E_IMM,E_REG,E_REL
}E_ARGS_TYPE;
typedef struct 
{
	E_ARGS_TYPE type;
	union 
	{
		Mem mem;
		Imm imm;
		Reg reg;
		Rel rel;
	} value;
}ArgsT;

ArgsT mk_imm_arg(Imm x);
ArgsT mk_mem_arg(Mem x);
ArgsT mk_reg_arg(Reg x);
ArgsT mk_nil_arg(void);
ArgsT mk_rel_arg(Rel x);


Imm* get_imm_arg(ArgsT* args);
Mem* get_mem_arg(ArgsT* args);
Reg* get_reg_arg(ArgsT* args);
Rel* get_rel_arg(ArgsT* args);

#define IS_IMM(x) (x.type == E_IMM)
#define IS_MEM(x) (x.type == E_MEM)
#define IS_NIL(x) (x.type == E_NIL)
#define IS_REG(x) (x.type == E_REG)
#define IS_REL(x) (x.type == E_REL)

void args_str(ArgsT*,char* buf ,int size);
extern ArgsT fixedArg[113];