#ifndef _LDASM_
#define _LDASM_

#include <stdint.h>
#include <string.h>

#if defined (_M_AMD64) || defined (__x86_64__)
    #define USE64
#elif defined (_M_IX86) || defined (__i386__)
    #define USE32
#else
    #error "Unknown or unsupported platform"
#endif

#ifdef USE64
    #define is_x64 1
#else
    #define is_x64 0
#endif//USE64

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

typedef struct _ldasm_data
{
    uint8_t  flags;
    uint8_t  rex;
    uint8_t  modrm;
    uint8_t  sib;
    uint8_t  opcd_offset;
    uint8_t  opcd_size;
    uint8_t  disp_offset;
    uint8_t  disp_size;
    uint8_t  imm_offset;
    uint8_t  imm_size;
} ldasm_data;

unsigned int   ldasm( void *code, ldasm_data *ld);
unsigned long  SizeOfProc( void *Proc );
void*          ResolveJmp( void *Proc );

#ifdef __cplusplus
}
#endif

#endif//_LDASM_