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
#include "Args.h"
#include "table.h"

// the same as x86asm Inst
#define PREFIX_SIZE 14
#define ARGS_SIZE 4
typedef uint32_t OpcodeType;
#define OPCODE_1(Opcode) ((Opcode>>24))
#define OPCODE_2(Opcode) (((Opcode>>16)&0xFF))
typedef struct 
{
    uint16_t Prefix[PREFIX_SIZE];
    uint32_t Op;
    uint32_t Opcode;
    ArgsT Args[ARGS_SIZE];
    int Mode;
    int AddrSize;
    int DataSize;
    int MemBytes;
    int Len;
    int PCRel;
    int PCRelOff;
}Inst;

void inst_str(Inst* inst,char* buf,int buf_size);

