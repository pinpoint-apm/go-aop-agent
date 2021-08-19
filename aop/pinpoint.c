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
#include "goX86asm.h"


static const BYTE NOP_1= (0xfa);
static const BYTE NOP_2[]= {0xd9,0xd0};
static const BYTE NOP_3[]= {0x0f,0x1f,0x00};

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
    // LOG_TRACE("start get neighbort page:%p",where);
    void* target  = get_page_boundary(where);
    // LOG_TRACE("where:%p target:%p",where,target);
    void* lo = get_2gb_low(target);
    void* hi = get_2gb_above(target);
    // LOG_TRACE("lo:%p hi:%p",lo,hi);
    void* try = NULL;
    // Try looking 1GB below or lower.
    if( try == NULL && target > (void*)0x40000000){
        // LOG_TRACE("try in > -1GB");
        try = try_get_page_from_addr_hi(lo,target - 0x40000000);
    }

    if( try == NULL && target < (void*)0xffffffff40000000){
        // LOG_TRACE("try in > +1GB");
        try = try_get_page_from_addr_lo(target +0x40000000,hi);
    }

    if( try == NULL && target > (void*)0x40000000){
        //  LOG_TRACE("try in <-1GB");
        try = try_get_page_from_addr_hi( target - 0x40000000 ,target);
    }

    if( try == NULL && target < (void*)0xffffffff40000000){
        // LOG_TRACE("try in <+1GB");
        try = try_get_page_from_addr_hi( target ,target + 0x40000000);
    }

    if(try == NULL){
        // LOG_TRACE("try in [-2GB,target)");
        try = try_get_page_from_addr_hi(lo ,target);
    }
    
    if(try == NULL){
        // LOG_TRACE("try in [target,+2GB)");
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
    LOG_TRACE("restore the trampoline_func inst to %p",trampoline->trampolineFunc.pTrampFunc);
    set_mm_area_opt(trampoline->trampolineFunc.pTrampFunc,trampoline->trampolineFunc.bakInstArLen,PROT_READ|PROT_WRITE|PROT_EXEC);
    memcpy(trampoline->trampolineFunc.pTrampFunc,trampoline->trampolineFunc.bakInstAr,trampoline->trampolineFunc.bakInstArLen);
    set_mm_area_opt(trampoline->target,trampoline->trampolineFunc.bakInstArLen,PROT_READ|PROT_EXEC);

    // munmap forward*
    // int page = getpagesize();
    // munmap(trampoline->forward,page);
    // munmap(trampoline->trampoline02,page);
    trampoline->forward = NULL;
    trampoline->back = NULL;
    // free trampoline
    free(ptr);
}

void place_safe_nop_inst(BYTE* p,int size)
{
    if(size <= 0){
        return ;
    }else if(size == 1){
        *p = NOP_1;
    }else if(size == 3){
        *(p) = NOP_3[0];
        *(p+1) = NOP_3[1];
        *(p+2) = NOP_3[2];
    }else if(size%2 == 0){
        for(int i =0;i<size;)
        {
            *(p+i) = NOP_2[0];
            *(p+i+1) = NOP_2[1];
            i+=2;
        }
    }else{
        // size -3 : must goto `size%2 == 0`
        place_safe_nop_inst(p,size-3);
        p += (size-3);
        place_safe_nop_inst(p,3);
    }
}

int32_t place_direct_jmp_inst(void*src,void* target,uint32_t placedSize)
{
    assert( placedSize >= JMP_INST_SIZE );
    set_mm_area_opt(src,JMP_INST_SIZE,PROT_READ|PROT_WRITE|PROT_EXEC);
    #if DTRACE
    { 
        BYTE* raw = src;
        LOG_TRACE("before:%p %X %X %X %X %X",src,raw[0],raw[1],raw[2],raw[3],raw[4]);
    }
    #endif
    BYTE* inst= src;
    *inst++ =  0xE9;
    *(int32_t*)inst++ = (int32_t)(target - src  - 5);
    place_safe_nop_inst(inst+JMP_INST_SIZE,placedSize - JMP_INST_SIZE);
    set_mm_area_opt(src,JMP_INST_SIZE,PROT_READ|PROT_EXEC);
    #if DTRACE
    {
        BYTE* raw = src;
        LOG_TRACE("[directly jmp] after:%p target:%p %X %X %X %X %X",src,target,raw[0],raw[1],raw[2],raw[3],raw[4]);
    }
    #endif 
    return JMP_INST_SIZE;
}

int32_t place_indirect_jmp_inst(void*src,void* target,uint32_t placedSize)
{
    assert(placedSize>=LONG_JMP_INST_SIZE);
    set_mm_area_opt(src,LONG_JMP_INST_SIZE,PROT_READ|PROT_WRITE|PROT_EXEC);
     #if DTRACE
    { 
        BYTE* raw = src;
        LOG_TRACE("before:%p %X %X %X %X %X %X",src,raw[0],raw[1],raw[2],raw[3],raw[4],raw[5]);
    }
    #endif
    BYTE* inst= src;
    *inst++ =  0xFF;
    *inst++ =  0x25;
    *(int32_t*)inst++= (int32_t)(target - (src + 6));
    place_safe_nop_inst(inst+LONG_JMP_INST_SIZE,placedSize - LONG_JMP_INST_SIZE);
    set_mm_area_opt(src,LONG_JMP_INST_SIZE,PROT_READ|PROT_EXEC);
    #if DTRACE
    {
        BYTE* raw = src;
        LOG_TRACE("[indirectly jmp] src:%p target:%p %X %X %X %X %X %X",src,target,raw[0],raw[1],raw[2],raw[3],raw[4],raw[5]);
    }
    #endif
    return LONG_JMP_INST_SIZE;
}


inline int32_t calc_inst_size(Reg* from,int len)
{
    Inst inst = {0};
    decode(from,len,&inst,64,false);
    return inst.Len;
}

int32_t calcRelativeOffset(BYTE* reg,int32_t size,BYTE* instBaseAddr){

    Reg* cur = reg;
    

    // copy from https://github.com/DarthTon/Blackbone/pull/420
    const int64_t diffMinVals[] = {0ll, -128ll, -32768ll, -8388608ll, -2147483648ll, -549755813888ll, -140737488355328ll, -36028797018963968ll, -9223372036854775807ll};
    const int64_t diffMaxVals[] = {0ll, 127ll, 32767ll, 8388607ll, 2147483647ll, 549755813887ll, 140737488355327ll, 36028797018963967ll, 9223372036854775807ll};

    do{
        Inst inst = {0};
        E_RET_TYPE ret = decode(cur,size,&inst,64,false);
        if (ret != E_OK 
             || (inst.Len == 1 && (inst.Op == INT || inst.Op == RET))
             || (inst.Len == 3 && inst.Op == 0xC2)) // 0xCC -> INT 0xC3/0xC2-> RETN
        {
            break;
        }

        if (inst.PCRel != 0)
        {
            assert(inst.PCRelOff != 0);
            // NOTE: how about big-endian
	        int32_t relative = 0;

            const uintptr_t ofst = inst.PCRelOff;
            const uintptr_t sz = inst.PCRel;

            memcpy( &relative, cur + ofst, sz );

            int64_t newRel = instBaseAddr + (int64_t)relative - reg ; //relative + (instBaseAddr - reg);

            if (newRel < diffMinVals[sz]  ||  newRel > diffMaxVals[sz]) {
                LOG_ETRACE("invalid offset. newDiff:%ld inst offset:%ld diff size:%ld",newRel,ofst,sz);
            	return -1;
            }

            LOG_TRACE("origin relative:%x updated relative:%lx ",relative,newRel);
            memcpy(cur + ofst, &newRel, sz);
        }
        cur += inst.Len;
        size -= inst.Len;
    }while( size > 0);

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
 * @return TrampolineBack* 
 */
TrampolineBack* insert_back_trampoline(TrampolineFuncT* trampolineFunc, FromInstBackUp* bakInst)
{
    void* origin_func = bakInst->instBaseAddr+bakInst->instBackupSize;
    LOG_TRACE("trampoline_func:%p origin_func:%p",trampolineFunc->pTrampFunc,origin_func);

    BYTE* inst  = bakInst->instBackUp;
    int32_t len = bakInst->instBackupSize;

    TrampolineBack* back = (TrampolineBack* )get_neighbor_mem(trampolineFunc->pTrampFunc,sizeof(TrampolineBack)+ len);
    back->restoreInstSize = len;
    back->toAddress = (long)origin_func;
    // insert jmp: trampoline_func to trampoline memory inst address
    place_direct_jmp_inst(trampolineFunc->pTrampFunc,back->inst,trampolineFunc->bakInstArLen);

    // restore from_inst
    memcpy(back->inst,inst,len);

    //note: after reloactedOffset, some inst could be changed
    if(calcRelativeOffset(back->inst,len,bakInst->instBaseAddr) != 0){
        return NULL;
    }

    // jmp trampoline to origin function
    BYTE* jmpInst = back->inst + len;

   // check range size
    if( labs((long ) jmpInst - (long) origin_func) >> 31 == 0 ){
        // use  directly jmp
        place_direct_jmp_inst(jmpInst,origin_func,JMP_INST_SIZE);
    }else{
        // insert jmp: trampoline to func
        place_indirect_jmp_inst(jmpInst,&back->toAddress,LONG_JMP_INST_SIZE);
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
 * @return TrampolineForward* address of trampoline. NULL: no need to trampoline, directly jmp is enough 
 */
TrampolineForward* insert_forward_trampoline(void* from, void* to)
{   
    // check range size
    if( labs((long ) from - (long) to) >> 31  == 0 ){
        // use  directly jmp
        place_direct_jmp_inst(from,to,JMP_INST_SIZE);
        return NULL;
    }

    // use directly jmp and indirect jmp

    TrampolineForward* forward = (TrampolineForward*)get_neighbor_mem(from,sizeof(TrampolineForward));
    forward->toAddress = (long)to;
    // x64 
    place_indirect_jmp_inst(forward->inst,&forward->toAddress,LONG_JMP_INST_SIZE);
    place_direct_jmp_inst(from,forward->inst,JMP_INST_SIZE);

    return forward;
}

// exclude  jmp +imm32
int32_t detour_does_code_end_function_e9(OpcodeType  pbCode)
{
    if (OPCODE_1(pbCode) == 0xeb ||    // jmp +imm8
        // OPCODE_1(pbCode) == 0xe9 ||    // jmp +imm32
        OPCODE_1(pbCode) == 0xe0 ||    // jmp eax
        OPCODE_1(pbCode) == 0xc2 ||    // ret +imm8
        OPCODE_1(pbCode) == 0xc3 ||    // ret
        OPCODE_1(pbCode) == 0xcc) {    // brk
        return 1;
    }
    else if (OPCODE_1(pbCode) == 0xf3 && OPCODE_2(pbCode) == 0xc3) {  // rep ret
        return 1;
    }
    else if (OPCODE_1(pbCode) == 0xff && OPCODE_2(pbCode) == 0x25) {  // jmp [+imm32]
        return 1;
    }
   
    return 0;
}

// copy from detour https://github.com/microsoft/Detours
int32_t detour_does_code_end_function(OpcodeType  pbCode)
{
    if ( OPCODE_1(pbCode) == 0xeb ||    // jmp +imm8
        OPCODE_1(pbCode) == 0xe9 ||    // jmp +imm32
        OPCODE_1(pbCode) == 0xe0 ||    // jmp eax
        OPCODE_1(pbCode) == 0xc2 ||    // ret +imm8
        OPCODE_1(pbCode) == 0xc3 ||    // ret
        OPCODE_1(pbCode) == 0xcc) {    // brk
        return 1;
    }
    else if ( OPCODE_1(pbCode) == 0xf3 &&  OPCODE_1(pbCode) == 0xc3) {  // rep ret
        return 1;
    }
    else if (OPCODE_1(pbCode) == 0xff && OPCODE_2(pbCode) == 0x25) {  // jmp [+imm32]
        return 1;
    }
   
    return 0;
}


int32_t make_space_for_jmp_boundary(void* address,const int32_t minSpace,BYTE* bakInstBytes,int32_t maxBakSize)
{
    assert(minSpace <= maxBakSize);

    BYTE* pByte = address;
    do{
        Inst inst = {0};
        E_RET_TYPE ret= decode(pByte,BACKUP_INST_SIZE,&inst,64,false);

        if( ret != E_OK ){
            LOG_ETRACE("unknown instruction");
            return -1;
        }

        if(detour_does_code_end_function(inst.Opcode)){
            char buf[128]={0};
		    inst_str(&inst,buf,sizeof(buf));
            LOG_ETRACE("met end function when scan space: %s",buf);
            return -1;
        }

        #if DTRACE
        { 
            char buf[128]={0};
		    inst_str(&inst,buf,sizeof(buf));
            LOG_TRACE("pInst:%p inst:{%s},len:%d",pByte,buf,inst.Len);
        }
        #endif

        pByte += inst.Len;
    }while(pByte - (BYTE*)address < minSpace );

    // backup the inst
    int32_t usageLen = pByte - (BYTE*)address;
    if(usageLen > maxBakSize){
        LOG_ETRACE("[🐛]backup space is too small. address:%p usageLen:%d maxBakSize:%d",address,usageLen,maxBakSize);
        return -1;
    }
    memcpy(bakInstBytes,address,usageLen);
    return usageLen;
}

void* located_nearest_call_target(void*start)
{
    #define UNKNOWN_SIZE 32
    BYTE* pInst = start;
    Inst _inst = {0};
    do{
        E_RET_TYPE ret = decode(pInst,UNKNOWN_SIZE,&_inst,64,false);
        if( ret != E_OK ){
            break;
        }
        // get pbCode and address
        BYTE pbCode = _inst.Op;
        // uint8_t pbCodeSize = _data.opcd_size;
        if(_inst.Len == CALL_INST_SIZE && pbCode == CALL){
            Rel* rel = get_rel_arg(&_inst.Args[0]);
            if(rel == NULL){
                LOG_TRACE("no such rel in `callq` inst");
                break;
            }
            BYTE* pfunc = pInst+ CALL_INST_SIZE +*rel;
            LOG_TRACE("origin start: %p real func: %p",start,pfunc);
            return pfunc;
        }
        pInst+=_inst.Len;
    }while(detour_does_code_end_function(_inst.Opcode) != 1);

    LOG_TRACE("start:%p no such callq",start);
    return NULL;
}


void* located_nearest_jmp_target(void*start)
{
    #define UNKNOWN_SIZE 32
    BYTE* pInst = start;
    Inst _inst = {0};
    do{
        E_RET_TYPE ret = decode(pInst,UNKNOWN_SIZE,&_inst,64,false);
        if( ret != E_OK ){
            break;
        }

        if(_inst.Len == JMP_INST_SIZE && _inst.Op == JMP && OPCODE_1(_inst.Opcode) == 0xE9){
            Rel* rel = get_rel_arg(&_inst.Args[0]);
            if(rel == NULL){
                LOG_TRACE("no such rel in `callq` inst");
                break;
            }
            BYTE* pfunc = pInst+ JMP_INST_SIZE + *rel;
            return pfunc;
        }
        pInst+=_inst.Len;
    }while(detour_does_code_end_function_e9(_inst.Opcode) != 1);
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
void* hook(void* from,void* to,void* callFrom)
{
    if(from == to ||callFrom == from || callFrom == to ){
        LOG_ETRACE("input is invalid");
        return NULL;
    }

    Trampoline* trampoline = (Trampoline*)malloc(sizeof(Trampoline));
    if(trampoline == NULL) {
        LOG_ETRACE("malloc %ld failed",sizeof(Trampoline));
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
        int32_t size = make_space_for_jmp_boundary(callFrom,JMP_INST_SIZE,trampoline->trampolineFunc.bakInstAr,BACKUP_INST_SIZE);
        if(size ==-1){
            return NULL;
        }
        trampoline->trampolineFunc.bakInstArLen = size;
        trampoline->trampolineFunc.pTrampFunc = callFrom;
    }

    //1. insert `forward` trampoline
    trampoline->forward = insert_forward_trampoline(from,to);
    LOG_TRACE("forward trampoline:%p from:%p to:%p backup:{base:%p size:%d} ",trampoline->forward, \
        from,to,trampoline->fromInstBackUp.instBaseAddr,trampoline->fromInstBackUp.instBackupSize);

    //2. insert `to` trampoline
    TrampolineBack* back = insert_back_trampoline(&trampoline->trampolineFunc,&trampoline->fromInstBackUp);
    if( back == NULL){
        LOG_ETRACE("hook：from:%p to:%p trampoline:%p failed",from,to,callFrom);
        return NULL;
    }
    trampoline->back = back;
    LOG_TRACE("trampoline_func:%p -> origin landing:%lx ",callFrom,back->toAddress);
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
    {
        BYTE inst[32] = {0};
        assert(make_space_for_jmp_boundary(test_hook,5,inst,32)>0 );
        assert(make_space_for_jmp_boundary(empty,5,inst,32)>0);
    }
    // BYTE inst[9]={0x64, 0x48, 0x8b, 0x0c, 0x25 ,0xf8 ,0xff ,0xff ,0xff};
    {
        BYTE inst[9] ={0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF};
        BYTE movInst[9]={0x64, 0x48, 0x8b, 0x0c, 0x25 ,0xf8 ,0xff ,0xff ,0xff};
        BYTE buf[32]={0};
        assert(make_space_for_jmp_boundary(inst,1,buf,32) == -1);
        assert(make_space_for_jmp_boundary(inst,3,buf,32) == -1);
        assert(make_space_for_jmp_boundary(inst,4,buf,32) == -1);
        assert(make_space_for_jmp_boundary(inst,5,buf,32) == -1);
        assert(make_space_for_jmp_boundary(movInst,sizeof(movInst),buf,32) == sizeof(inst));
    }

    LOG_TRACE("passed");
}

void testAsmCall()
{
    // BYTE inst[]= {0xe8 ,   0x7b  ,  0xfe  ,  0xff  ,  0xff  ,  0x48 };
    // ldasm_data _data;
    // do{
    //     int len = ldasm(inst,&_data);
    //     assert(len == 5);
    //     LOG_TRACE("%d \n",len);
    // }while(0);
    // LOG_TRACE("passed");
}

void testLea(){

    BYTE inst[]= {0x48,0x8D,0x05,0xA9,0xA2,0x020,0x40};
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

void test_invalid_hook(){
    hook(retStr,NULL,NULL);
    hook(retStr,retStr,NULL);
    hook(retStr,HookRetStr,HookRetStr);
    hook(retStr,HookRetStrTrampoline,HookRetStrTrampoline);
}

void test_located_nearest_jmp_target()
{
    BYTE inst[] ={0xe8,0xdb,0xe4,0xf5,0xff,0xe9,0x5d,0xff,0xff,0xff};
    void* land= located_nearest_jmp_target(inst);
    assert( land !=inst || land !=NULL);
    LOG_TRACE("%p",land);

}

void test_try_get_page_from_addr_hi()
{
    void* hi = &test_located_nearest_jmp_target;
    assert(try_get_page_from_addr_hi(hi-getpagesize()*128 ,hi));
}

int main()
{
    printf("-------test_make_space---------------------------- \n");
    test_make_space();
    printf("-------test_hook---------------------------- \n");
    test_hook();
    printf("-------test_hook---------------------------- \n");
    test_hook();
    printf("-------testAsmCall---------------------------- \n");
    testAsmCall();
    printf("-------test_call_nearest_func----------------------------\n");
    test_call_nearest_func();
    printf("-------test_call_get_jmp32_mem_from_chunk----------------------------\n");
    test_call_get_jmp32_mem_from_chunk();
    printf("-------test_32_range----------------------------\n");
    test_32_range();
    printf("-------test_retStr----------------------------\n");
    test_retStr();
    printf("-------test_invalid_hook----------------------------\n");
    test_invalid_hook();
    printf("-------testLea----------------------------\n");
    testLea();
    printf("-------test_located_nearest_jmp_target----------------------------\n");
    test_located_nearest_jmp_target();
    printf("-------test_try_get_page_from_addr_hi----------------------------\n");
    test_try_get_page_from_addr_hi();
    return 0;
}

#endif

