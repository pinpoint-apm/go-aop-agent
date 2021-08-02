#pragma once

#include <stdint.h>
#include <stdio.h>

#ifdef __cplusplus
extern "C"
{
#endif

#define F_INVALID       0x01
#define F_PREFIX        0x02
#define F_REX           0x04
#define F_MODRM         0x08
#define F_SIB           0x10
#define F_DISP          0x20
#define F_IMM           0x40
#define F_RELATIVE      0x80

// the same as x86asm Inst
#define PREFIX_SIZE 14

typedef struct 
{
    uint16_t Prefix[PREFIX_SIZE];
    uint32_t Op;
    uint32_t Opcode;
    Imm    Args[4];
    int Mode;
    int AddrSize;
    int DataSize;
    int MemBytes;
    int Len;
    int PCRel;
    int PCRelOff;
}Inst;

unsigned int   m_to_asm( void *code, Inst *ld);
unsigned long  SizeOfProc( void *Proc );
void*          ResolveJmp( void *Proc );

#ifdef __cplusplus
}
#endif