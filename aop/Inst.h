#pragma once
#include "Args.h"
#include "table.h"

// the same as x86asm Inst
#define PREFIX_SIZE 14
#define ARGS_SIZE 4

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

