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

#include "Inst.h"
#include "table.h"
#include "Args.h"

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