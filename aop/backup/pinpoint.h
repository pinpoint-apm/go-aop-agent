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
#include<stdio.h>

//linux 
#include <sys/mman.h>
#include <unistd.h>
#include <string.h>
#include <stdint.h>
#include <errno.h>
#include <stdlib.h>
#include <assert.h>

#ifndef NTRACE
#define LOG_TRACE(fmt,args...)  fprintf(stderr,"[%s:%d] %s: " fmt "\n",__FILE__,__LINE__,__FUNCTION__,##args)
#else
#define LOG_TRACE(fmt,args...) 
#endif
typedef unsigned char INST;

typedef struct trampoline_forward_s{
    long toAddress;
    INST inst[6];   // indirectly jmp：0xFF 0x25
}TrampolineForward;

typedef struct trampoline_back_s{
    long toAddress;
    int32_t restoreInstSize;
    INST inst[0]; // include restore inst and jmp
}TrampolineBack;

#define BACKUP_INST_SIZE 32
#define JMP_INST_SIZE 5
#define LONG_JMP_INST_SIZE 6
#define CALL_INST_SIZE 5

typedef struct {
    INST instBackUp[BACKUP_INST_SIZE];
    uint8_t instBackupSize;
    /*
    * ┌──────────────┐
    * │ push xxx     │  ◄─────── instBaseAddr
    * ├──────────────┤
    * │ push xxx     │
    * ├──────────────┤
    * │ lea xx, [rip]│
    * ├──────────────┤
    * │              │
    * └──────────────┘
    */
    INST* instBaseAddr; // base address of backup
}FromInstBackUp;

typedef struct trampoline_s{
    /**
     * @brief 
     * why 32, 5 is the jmp inst size, 13 is the maximun inst size
     * NOTE!!! those space must be executeable
     */

    FromInstBackUp fromInstBackUp;

    INST trampolineFuncInstBackUp[BACKUP_INST_SIZE];
    uint8_t trampolineBackUpSize;
    // store the src address for restore
    void* target;
    void* trampoline_func;
    TrampolineForward* trampoline01;
    TrampolineBack* trampoline02;
}Trampoline;

void* hook(void* from,void* to,void* trampoline_func);
void* located_nearest_call_target(void*start);
void* located_nearest_jmp_target(void*start);
void  unhook(void* ptr);