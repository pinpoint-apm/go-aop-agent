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
#include "pinpoint.h"
#include "LDasm.h"


typedef struct page_mem_chunk_s{
    uint16_t freeOffset;
    struct page_mem_chunk_s * next;
}PageMemChunk;

static PageMemChunk* mem_chunk_header_g;
static void* get_neighbor_page(void* where);
static int set_mm_area_opt(void *ptr,int size, int prot);

static inline int value_to_cpu(int v)
{
    int to = v & ~(sizeof(long)- 1) ;

    return (to < v) ? (to+ (int)sizeof(long)) : (v);
}


static inline void append_page_chunk(void* page)
{
    PageMemChunk* newChunk = (PageMemChunk*) page;
    newChunk->next = NULL;
    newChunk->freeOffset = (uint16_t)sizeof(PageMemChunk);

    if(mem_chunk_header_g == NULL)
    {
        mem_chunk_header_g = newChunk;
        return;
    }

    PageMemChunk* end = mem_chunk_header_g;
    while(end->next != NULL){
        end = end->next;
    }

    end->next = newChunk;
}

static void* get_jmp32_mem_from_chunk(const void* dst,int size)
{
    size = value_to_cpu(size);

    PageMemChunk* pChunk = mem_chunk_header_g;
    while(pChunk){
        int maxSize = getpagesize() -  pChunk->freeOffset;
        void* pfree =  (char*)pChunk + pChunk->freeOffset;

        if(maxSize >= size && labs((long int )dst - (long int)pfree) < (long int)(0x7ff80000)){
            // update lastUsedOffset
            set_mm_area_opt(pChunk,sizeof(PageMemChunk),PROT_READ|PROT_WRITE|PROT_EXEC);
            pChunk->freeOffset += size;

            return pfree;
        }
        pChunk  = pChunk->next;
    }
    return NULL;
}


void* get_page_boundary(void* ptr)
{
    return (void*)((long)ptr & ~(getpagesize()-1));
}

static inline void* get_2gb_low(void* ptr)
{
    return (void*)(((long)ptr > 0x7ff80000) ? (long)ptr - 0x7ff80000 : 0x80000);
}

static inline void* get_2gb_above(void* ptr)
{
    return (void*)((unsigned long)ptr < 0xffffffff80000000 ? (unsigned long)ptr + 0x7ff80000 : 0xfffffffffff80000);
}

int set_mm_area_opt(void *ptr,int size, int prot)
{
    int page_size = getpagesize();
    void* p = get_page_boundary(ptr);
    
    for(;p <= ptr + size; p+=page_size){
        int ret = mprotect(p,page_size,prot);
        if (ret != 0){
            LOG_TRACE("mprotect on [start:%p,size:%d] opt:%d error:%s ",p,page_size,prot,strerror(errno));
            return -1;
        }
    }
    return 0;
}


static  void* try_get_page_from_addr_hi(void* low,void*hi)
{
    int page =getpagesize();
    void* mem = NULL;

    for(;low < hi; hi -= page){
        // force map memory from start  
        mem = mmap(hi,page,PROT_READ | PROT_EXEC |PROT_WRITE,MAP_ANONYMOUS| MAP_PRIVATE,-1,0);
        if(mem == NULL){
            LOG_TRACE("mmap:%p->%d failed:%s",hi,page,strerror(errno));
        }else if(mem!=hi){
            munmap(mem,page);
        }
        else{
            return mem;
        }
    }
    return mem;
}


static  void* try_get_page_from_addr_lo(void* low,void*hi)
{
    int page =getpagesize();
    void* mem = NULL;

    for(;low < hi;low += page){
        // force map memory from start  
        mem = mmap(low,page,PROT_READ | PROT_EXEC |PROT_WRITE,MAP_ANONYMOUS | MAP_PRIVATE,-1,0);
        if(mem == NULL){
            LOG_TRACE("mmap:%p->%d failed:%s",low,page,strerror(errno));
        }else if(mem!=low){
            munmap(mem,page);
        }
        else{
            return mem;
        }
    }
    return mem;
}


void* get_neighbor_mem(void* where, int size){
    void* mem = get_jmp32_mem_from_chunk(where,size);
    if(mem == NULL){
        void* page = get_neighbor_page(where);
        if(page == NULL){
            return NULL;
        }

        append_page_chunk(page);

        return get_neighbor_mem(where,size);
    }

    return mem;
} 

/**
 * @brief Get the page neighbor mem object
 * size is 1 pagesize(4096)
 * @param where 
 * @return void* 
 */
void* get_neighbor_page(void* where)
{
    LOG_TRACE("start get neighbort page:%p",where);
    void* target  = get_page_boundary(where);
    LOG_TRACE("where:%p target:%p",where,target);
    void* lo = get_2gb_low(target);
    void* hi = get_2gb_above(target);
    LOG_TRACE("lo:%p hi:%p",lo,hi);
    void* try = NULL;
    // Try looking 1GB below or lower.
    if( try == NULL && target > (void*)0x40000000){
        LOG_TRACE("try in > -1GB");
        try = try_get_page_from_addr_hi(lo,target - 0x40000000);
    }

    if( try == NULL && target < (void*)0xffffffff40000000){
        LOG_TRACE("try in > +1GB");
        try = try_get_page_from_addr_lo(target +0x40000000,hi);
    }

    if( try == NULL && target > (void*)0x40000000){
         LOG_TRACE("try in <-1GB");
        try = try_get_page_from_addr_hi( target - 0x40000000 ,target);
    }

    if( try == NULL && target < (void*)0xffffffff40000000){
        LOG_TRACE("try in <+1GB");
        try = try_get_page_from_addr_hi( target ,target + 0x40000000);
    }

    if(try == NULL){
        LOG_TRACE("try in [-2GB,target)");
        try = try_get_page_from_addr_hi(lo ,target);
    }
    
    if(try == NULL){
        LOG_TRACE("try in [target,+2GB)");
        try = try_get_page_from_addr_lo(target ,hi);
    }

    LOG_TRACE("where:%p neighbor:%p range:%ld",where,try,where-try);
    return try;
}

void unhook(void* ptr)
{
    if(ptr == NULL){
        return;
    }

    // restore target inst
    Trampoline* trampoline = (Trampoline*)ptr;
    LOG_TRACE("restore the backup inst to %p",trampoline->target);
    set_mm_area_opt(trampoline->target,trampoline->fromInstBackUp.instBackupSize,PROT_READ|PROT_WRITE|PROT_EXEC);
    memcpy(trampoline->target,trampoline->fromInstBackUp.instBackUp,trampoline->fromInstBackUp.instBackupSize);
    set_mm_area_opt(trampoline->target,trampoline->fromInstBackUp.instBackupSize,PROT_READ|PROT_EXEC);

    // restore trampoline_func
    LOG_TRACE("restore the trampoline_func inst to %p",trampoline->trampoline_func);
    set_mm_area_opt(trampoline->trampoline_func,trampoline->trampolineBackUpSize,PROT_READ|PROT_WRITE|PROT_EXEC);
    memcpy(trampoline->trampoline_func,trampoline->trampolineFuncInstBackUp,trampoline->trampolineBackUpSize);
    set_mm_area_opt(trampoline->target,trampoline->trampolineBackUpSize,PROT_READ|PROT_EXEC);

    // munmap trampoline01*
    // int page = getpagesize();
    // munmap(trampoline->trampoline01,page);
    // munmap(trampoline->trampoline02,page);
    trampoline->trampoline01 = NULL;
    trampoline->trampoline02 = NULL;
    // free trampoline
    free(ptr);
}

void* insert_direct_jmp_inst(void*src,void* target)
{
    set_mm_area_opt(src,JMP_INST_SIZE,PROT_READ|PROT_WRITE|PROT_EXEC);
    INST* inst= src;
    *inst++ =  0xE9;
    *(int32_t*)inst++ = (int32_t)(target - src  - 5);
    set_mm_area_opt(src,JMP_INST_SIZE,PROT_READ|PROT_EXEC);
    #ifdef DTRACE
    {
        INST* raw = src;
        LOG_TRACE("src:%p -> target:%p %X -> %X %X %X %X",src,target,raw[0],raw[1],raw[2],raw[3],raw[4]);
    }
    #endif 
    return inst;
}

void* insert_indirect_jmp_inst(void*src,void* target)
{
    set_mm_area_opt(src,6,PROT_READ|PROT_WRITE|PROT_EXEC);
    INST* inst= src;
    *inst++ =  0xFF;
    *inst++ =  0x25;
    *(int32_t*)inst++= (int32_t)(target - (src + 6));  
    set_mm_area_opt(src,6,PROT_READ|PROT_EXEC);
    #ifdef DTRACE
    {
        INST* raw = src;
        LOG_TRACE("src:%p -> target:%p %X %X -> %X %X %X %X",src,target,raw[0],raw[1],raw[2],raw[3],raw[4],raw[5]);
    }
    #endif
    return inst;
}


inline int32_t calc_inst_size(void* from)
{
    ldasm_data _data;
    return ldasm(from,&_data);
}

int32_t calcRelativeOffset(INST* inst,int32_t size,INST* instBaseAddr){

    INST* cur = inst;
    ldasm_data ld = { 0 };

    // copy from https://github.com/DarthTon/Blackbone/pull/420
    const int64_t diffMinVals[] = {0ll, -128ll, -32768ll, -8388608ll, -2147483648ll, -549755813888ll, -140737488355328ll, -36028797018963968ll, -9223372036854775807ll};
    const int64_t diffMaxVals[] = {0ll, 127ll, 32767ll, 8388607ll, 2147483647ll, 549755813887ll, 140737488355327ll, 36028797018963967ll, 9223372036854775807ll};

    do{
        uint32_t len = ldasm( cur, &ld );
        if (ld.flags & F_INVALID
             || (len == 1 && (cur[ld.opcd_offset] == 0xCC || cur[ld.opcd_offset] == 0xC3))
             || (len == 3 && cur[ld.opcd_offset] == 0xC2)) // 0xCC -> INT 0xC3/0xC2-> RETN
        {
            break;
        }

        if (ld.flags & F_RELATIVE)
        {
            // NOTE: how about big-endian
	        int32_t diff = 0;
            const uintptr_t ofst = (ld.disp_offset != 0 ? ld.disp_offset : ld.imm_offset);
            const uintptr_t sz = ld.disp_size != 0 ? ld.disp_size : ld.imm_size;

            memcpy( &diff, cur + ofst, sz );

            int64_t newDiff = ((int64_t) diff) + (instBaseAddr -inst);

            if (newDiff < diffMinVals[sz]  ||  newDiff > diffMaxVals[sz]) {
                LOG_TRACE("invalid offset. newDiff:%ld inst offset:%ld diff size:%ld",newDiff,ofst,sz);
            	return -1;
            }
            memcpy(cur + ofst, &newDiff, sz);
        }
        cur += len;
    }while(cur < inst + size);

    return 0;
}


/**
 * @brief generate `back trampoline` for call `origin function`
 * 1. `from` must be writable/executeable
 * 
 * @param from function address of `origin function template`
 * @param to origin function/method entry inst
 * @param p_inst saved inst from `Foo`
 * @param inst_size size of p_inst
 * @return void* 
 */

/**
 * @brief generate `back trampoline` for call `origin function`
 * 1. `from` must be writable/executeable
 * 
 * @param from t
 * @param from_inst 
 * @return void* 
 */
void* insert_back_trampoline(void* trampoline_func, FromInstBackUp* from_inst)
{
    void* origin_func = from_inst->instBaseAddr+from_inst->instBackupSize;
    LOG_TRACE("trampoline_func:%p origin_func:%p",trampoline_func,origin_func);

    int32_t saved_inst_size =from_inst->instBackupSize;
    void* func_saved_inst = from_inst->instBackUp;

    #ifdef DTRACE
    {
        int i = 0;
        LOG_TRACE("saved instructions: ");
        for(;i<saved_inst_size;i++){
            fprintf(stderr," %X ",func_saved_inst[i]);
        }
        fprintf(stderr,"\n");
    }
    #endif

    TrampolineBack* back = (TrampolineBack* )get_neighbor_mem(trampoline_func,sizeof(TrampolineBack)+ saved_inst_size);
    back->restoreInstSize = saved_inst_size;
    back->toAddress = (long)origin_func;
    // insert jmp: trampoline_func to trampoline memory inst address
    insert_direct_jmp_inst(trampoline_func,back->inst);

    // restore from_inst
    memcpy(back->inst,func_saved_inst,saved_inst_size);

    //note: after reloactedOffset, some inst could be changed
    if(calcRelativeOffset(back->inst,saved_inst_size,from_inst->instBaseAddr) != 0){
        return NULL;
    }

    // jmp trampoline to origin function
    INST* next_i = back->inst + saved_inst_size;

   // check range size
    if( labs((long ) next_i - (long) origin_func) >> 31 == 0 ){
        // use  directly jmp
        insert_direct_jmp_inst(next_i,origin_func);
    }else{
        // insert jmp: trampoline to func
        insert_indirect_jmp_inst(next_i,&back->toAddress);
    }

    return back;
}

/**
 * @brief 
 *  * check the jmp range size
 *  * create a trampoline struct recording src
 *  * generate src jmp, trampoline jmp
 *  * return trampoline pointer
 *  note: go-layer must free trampoline pointer
 * @param src 
 * @param dst 
 * @return void* address of trampoline. NULL: no need to trampoline, directly jmp is enough 
 */
void* insert_forward_trampoline(void* from, void* to)
{   
    LOG_TRACE("from:%p to:%p",from,to);

    // check range size
    if( labs((long ) from - (long) to) >> 31  == 0 ){
        // use  directly jmp
        insert_direct_jmp_inst(from,to);
        return NULL;
    }

    // use directly jmp and indirect jmp

    TrampolineForward* forward = (TrampolineForward*)get_neighbor_mem(from,sizeof(TrampolineForward));
    forward->toAddress = (long)to;
    // x64 
    insert_indirect_jmp_inst(forward->inst,&forward->toAddress);
    insert_direct_jmp_inst(from,forward->inst);

    return forward;
}

// exclude  jmp +imm32
int32_t detour_does_code_end_function_e9(INST* pbCode)
{
    if (pbCode[0] == 0xeb ||    // jmp +imm8
        // pbCode[0] == 0xe9 ||    // jmp +imm32
        pbCode[0] == 0xe0 ||    // jmp eax
        pbCode[0] == 0xc2 ||    // ret +imm8
        pbCode[0] == 0xc3 ||    // ret
        pbCode[0] == 0xcc) {    // brk
        return 1;
    }
    else if (pbCode[0] == 0xf3 && pbCode[1] == 0xc3) {  // rep ret
        return 1;
    }
    else if (pbCode[0] == 0xff && pbCode[1] == 0x25) {  // jmp [+imm32]
        return 1;
    }
    else if ((pbCode[0] == 0x26 ||      // jmp es:
              pbCode[0] == 0x2e ||      // jmp cs:
              pbCode[0] == 0x36 ||      // jmp ss:
              pbCode[0] == 0x3e ||      // jmp ds:
              pbCode[0] == 0x64 ||      // jmp fs:
              pbCode[0] == 0x65) &&     // jmp gs:
             pbCode[1] == 0xff &&       // jmp [+imm32]
             pbCode[2] == 0x25) {
        return 1;
    }
    return 0;
}


// copy from detour https://github.com/microsoft/Detours
int32_t detour_does_code_end_function(INST* pbCode)
{
    if (pbCode[0] == 0xeb ||    // jmp +imm8
        pbCode[0] == 0xe9 ||    // jmp +imm32
        pbCode[0] == 0xe0 ||    // jmp eax
        pbCode[0] == 0xc2 ||    // ret +imm8
        pbCode[0] == 0xc3 ||    // ret
        pbCode[0] == 0xcc) {    // brk
        return 1;
    }
    else if (pbCode[0] == 0xf3 && pbCode[1] == 0xc3) {  // rep ret
        return 1;
    }
    else if (pbCode[0] == 0xff && pbCode[1] == 0x25) {  // jmp [+imm32]
        return 1;
    }
    else if ((pbCode[0] == 0x26 ||      // jmp es:
              pbCode[0] == 0x2e ||      // jmp cs:
              pbCode[0] == 0x36 ||      // jmp ss:
              pbCode[0] == 0x3e ||      // jmp ds:
              pbCode[0] == 0x64 ||      // jmp fs:
              pbCode[0] == 0x65) &&     // jmp gs:
             pbCode[1] == 0xff &&       // jmp [+imm32]
             pbCode[2] == 0x25) {
        return 1;
    }
    return 0;
}



int32_t make_space_for_jmp_boundary(void* address,int32_t min_size,INST* out_backup_inst,int32_t out_size)
{
    assert(min_size <= out_size);

    ldasm_data _data;
    INST* pInst = address;
    do{
        uint32_t inst_len= ldasm(pInst,&_data);
        if( (_data.flags & F_INVALID) !=0){
            break;
        }
        if(detour_does_code_end_function(pInst)){
            LOG_TRACE("met end function when scan space: %X %X",pInst[0],pInst[1]);
            return -1;
        }
        pInst+=inst_len;
    }while((long)pInst - (long)address <min_size);

    // backup the inst
    out_size = (long)pInst - (long)address;
    memcpy(out_backup_inst,address,out_size);
    return out_size;
}

typedef int32_t OFFSETTYPE;
#pragma pack (push,1)
typedef struct call_s
{
    INST opcode;
    OFFSETTYPE address;
}CallInst;
#pragma pack(pop)

void* located_nearest_call_target(void*start)
{
    ldasm_data _data;
    INST* pInst = start;
    do{
        uint32_t inst_len= ldasm(pInst,&_data);
        if( (_data.flags & F_INVALID) !=0){
            break;
        }
        // get opcode and address
        INST* opcode = (uint8_t*)pInst + _data.opcd_offset;
        uint8_t opCodeSize = _data.opcd_size;
        if(inst_len == CALL_INST_SIZE && opCodeSize ==1 && opcode[0] == 0xE8){
            CallInst* pCall = ( CallInst*) pInst;
            INST* pfunc = pInst+ CALL_INST_SIZE + pCall->address;
            LOG_TRACE("origin start: %p real func: %p",start,pfunc);
            return pfunc;
        }
        pInst+=inst_len;
    }while(detour_does_code_end_function(pInst) != 1);
    LOG_TRACE("start:%p no such callq",start);
    return NULL;
}


void* located_nearest_jmp_target(void*start)
{
    ldasm_data _data;
    INST* pInst = start;
    do{
        uint32_t inst_len= ldasm(pInst,&_data);
        if( (_data.flags & F_INVALID) !=0){
            break;
        }
        // get opcode and address
        INST* opcode = (uint8_t*)pInst + _data.opcd_offset;
        uint8_t opCodeSize = _data.opcd_size;
        if(inst_len == JMP_INST_SIZE && opCodeSize ==1 && opcode[0] == 0xE9){
            CallInst* pCall = ( CallInst*) pInst;
            INST* pfunc = pInst+ JMP_INST_SIZE + pCall->address;
            LOG_TRACE("origin start: %p  real func: %p",start,pfunc);
            return pfunc;
        }
        pInst+=inst_len;
    }while(detour_does_code_end_function_e9(pInst) != 1);
    LOG_TRACE("start:%p no such callq",start);
    return NULL;
}

// /**
//  * @brief patch for hook
//  *   parse the machine code, update the from address from callq
//  * exp: callq platform is complex, needs more test
//  * @param from 
//  * @param to 
//  * @param trampoline_func 
//  * @return void* 
//  */
// void* hookP0_1(void* from,void* to,void* trampoline_func)
// {
//     void*n_from = located_nearest_call_target(from);
//     if(n_from){
//         return hook(n_from,to,trampoline_func);
//     }else{
//         return NULL;
//     }
// }



/**
 * @brief 
 * 
 * @param from 
 * @param to 
 * @param trampoline 
 * @return void* 
 */
void* hook(void* from,void* to,void* trampoline_func)
{
    if(from == to ||trampoline_func == from || trampoline_func == to ){
        LOG_TRACE("input is invalid");
        return NULL;
    }

    Trampoline* trampoline = (Trampoline*)malloc(sizeof(Trampoline));
    if(trampoline == NULL) {
        LOG_TRACE("malloc %ld failed",sizeof(Trampoline));
        return NULL;
    }

    // locate the safe inst boundary for jmp-from inst
    {
        int32_t size = make_space_for_jmp_boundary(from,JMP_INST_SIZE,trampoline->fromInstBackUp.instBackUp,BACKUP_INST_SIZE);
        if(size ==-1){
            return NULL;
        }
        trampoline->fromInstBackUp.instBackupSize = size;
        trampoline->fromInstBackUp.instBaseAddr = from;
        trampoline->target = from;
    }

    // locate the safe inst boundary for jmp-trampoline_func inst
    {
        int32_t size = make_space_for_jmp_boundary(trampoline_func,JMP_INST_SIZE,trampoline->trampolineFuncInstBackUp,BACKUP_INST_SIZE);
        if(size ==-1){
            return NULL;
        }
        trampoline->trampolineBackUpSize = size;
        trampoline->trampoline_func = trampoline_func;
    }

    //1. insert `forward` trampoline
    trampoline->trampoline01 = insert_forward_trampoline(from,to);
    //2. insert `to` trampoline
    void* trampoline02 = insert_back_trampoline(trampoline_func,&trampoline->fromInstBackUp);
    if( trampoline02 == NULL){
        LOG_TRACE("hook：from:%p to:%p trampoline:%p failed",from,to,trampoline_func);
        return NULL;
    }

    trampoline->trampoline02 = trampoline02;
    //3. return trampoline for unhook
    return trampoline;
}


#ifndef NTEST
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                  test zone                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

int foo(int a,int b){
    return a+b;
}

int foo1(int a,int b){
    return a-b;
}

int before(int a,int b)
{
    (void)a;
    (void)b;
    printf("call before \n");
    return 0;
}

int end(int* a)
{
    printf("call end  %d\n",*a);
    *a+=*a;
    return 0;
}

int hook_foo(int a,int b)
{
    int ret = before(a,b);
    printf("start call origin \n");
    a = 10;
    b= 12;
    ret = foo1(a,b);
    printf("end call origin \n");
    end(&ret);
    return ret;
}

#define PRINT(fmt,func) LOG_TRACE(#func " = "fmt,func )

void test_hook()
{
    void* trampline = hook(foo,hook_foo,foo1);
    PRINT("%d",foo(1,3));
    LOG_TRACE(" \n unhook \n");
    unhook(trampline);
    PRINT("%d",foo(1,3));
    LOG_TRACE("passed");
}

void empty(){}
void test_make_space()
{
    INST inst[32];
    assert(make_space_for_jmp_boundary(test_hook,5,inst,32)>0 );
    assert(make_space_for_jmp_boundary(empty,5,inst,32)>0);
    LOG_TRACE("passed");
}

void testAsmCall()
{
    INST inst[]= {0xe8 ,   0x7b  ,  0xfe  ,  0xff  ,  0xff  ,  0x48 };
    ldasm_data _data;
    do{
        int len = ldasm(inst,&_data);
        assert(len == 5);
        LOG_TRACE("%d \n",len);
    }while(0);
    LOG_TRACE("passed");
}

//@obsolated
void testLea(){

    INST inst[]= {0x48,0x8D,0x05,0xA9,0xA2,0x020,0x40};
    calcRelativeOffset(inst,sizeof(inst),inst+ 1024);
    uint32_t i = 0;
    for(; i< sizeof(inst);i++){
        printf("0x%x ",inst[i]);
    }

}

void B(){
    LOG_TRACE("---------just some test-------------------");
}

void A(){
    B();
}

void test_call_nearest_func()
{
    void* print_func =located_nearest_call_target(A);
    if(print_func != B){
        LOG_TRACE("failed: located before failed");
        assert(print_func == B);
    }
    LOG_TRACE("passed");
}

void test_call_get_jmp32_mem_from_chunk()
{
    assert(value_to_cpu(7)  == 8);
    assert(value_to_cpu(8)  == 8);
    assert(value_to_cpu(9)  == 16);
    assert(value_to_cpu(0)  == 0);

   assert(get_neighbor_mem(printf,12)!=NULL);

}

void test_32_range()
{
    int32_t  zero_com = ~0x0;
    int32_t  negative_1 = -1;
    assert(zero_com  == negative_1);
}


static const char* Msg="I'm string";

static const char* retStr(){
    return Msg;
}

char N_Msg[128] =  {0};

const char* HookRetStrTrampoline();

const char* HookRetStr(){

    const char* ret = HookRetStrTrampoline();

    return strcat(strcat(N_Msg,ret),"[hooked]");
}
const char* NULLStr="NULL";
const char* HookRetStrTrampoline(){
    return NULLStr;
}

void test_retStr()
{
    LOG_TRACE("%s",retStr());
    void* trampline = hook(retStr,HookRetStr,HookRetStrTrampoline);
    if(trampline !=NULL){
        assert(strcmp(retStr(),"I'm string[hooked]") == 0);
        unhook(trampline);
        assert(strcmp(retStr(),Msg) == 0);
        LOG_TRACE("test_retStr passed");
        return ;
    }
    LOG_TRACE(" hook(HookRetStr,HookRetStr,HookRetStrTrampoline) failed ");
}

int main()
{
    test_make_space();
    test_hook();
    test_hook();
    testAsmCall();
    test_call_nearest_func();
    test_call_get_jmp32_mem_from_chunk();
    test_32_range();
    test_retStr();
    return 0;
}

// gcc -g common.c LDasm.c -Wall
#endif