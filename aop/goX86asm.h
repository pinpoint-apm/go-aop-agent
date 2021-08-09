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
typedef enum {
    E_UNRECOGNIZED = -1024,
    E_TRUNCATED,
    E_INVALID_MODE,
    E_INTERNAL,
    E_PARA_INVALID,
    E_OK=0,
}E_RET_TYPE;

E_RET_TYPE decode(Reg *src,int len, Inst *ld, int mode,bool gnuCompat);
#ifdef __cplusplus
}
#endif