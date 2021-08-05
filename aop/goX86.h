#pragma once
#include "Inst.h"
#include "table.h"

#include <stdint.h>
#include <stdio.h>

#ifdef __cplusplus
extern "C"
{
#endif

//the same as x86asm Inst

// unsigned int   m_to_asm( void *code, Inst *ld);
// unsigned long  SizeOfProc( void *Proc );
// void*          ResolveJmp( void *Proc );
uint32_t decode(Reg *src,int len, Inst *ld, int mode,bool gnuCompat);
#ifdef __cplusplus
}
#endif