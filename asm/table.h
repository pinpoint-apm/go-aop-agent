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
#include <stdint.h>

#ifndef uint16
#define uint16 uint16_t
#endif

typedef uint16_t Prefix;

typedef int64_t Imm;

typedef uint8_t Reg;
typedef int32_t Rel;

#ifndef bool
typedef uint8_t bool;
#define true  ((uint8_t)1)
#define false ((uint8_t)0)
#endif

#define decodeOp uint16

#define PrefixImplicit   0x8000  // prefix is implied by instruction text
#define PrefixIgnored    0x4000  // prefix is ignored: either irrelevant or overridden by a later prefix
#define PrefixInvalid    0x2000  // prefix makes entire instruction invalid (bad LOCK)

// Memory segment overrides.
#define PrefixES   0x26  // ES segment override
#define PrefixCS   0x2E  // CS segment override
#define PrefixSS   0x36  // SS segment override
#define PrefixDS   0x3E  // DS segment override
#define PrefixFS   0x64  // FS segment override
#define PrefixGS   0x65  // GS segment override

// Branch prediction.
#define PrefixPN   0x12E  // predict not taken (conditional branch only)
#define PrefixPT   0x13E  // predict taken (conditional branch only)

// Size attributes.
#define PrefixDataSize   0x66  // operand size override
#define PrefixData16     0x166 
#define PrefixData32     0x266 
#define PrefixAddrSize   0x67  // address size override
#define PrefixAddr16     0x167 
#define PrefixAddr32     0x267 

// One of a kind.
#define PrefixLOCK       0xF0  // lock
#define PrefixREPN       0xF2  // repeat not zero
#define PrefixXACQUIRE   0x1F2 
#define PrefixBND        0x2F2 
#define PrefixREP        0xF3  // repeat
#define PrefixXRELEASE   0x1F3 

// The REX prefixes must be in the range [PrefixREX, PrefixREX+0x10).
// the other bits are set or not according to the intended use.
#define PrefixREX         0x40  // REX 64-bit extension prefix
#define PrefixREXW        0x08  // extension bit W (64-bit instruction width)
#define PrefixREXR        0x04  // extension bit R (r field in modrm)
#define PrefixREXX        0x02  // extension bit X (index field in sib)
#define PrefixREXB        0x01  // extension bit B (r/m field in modrm or base field in sib)
#define PrefixVEX2Bytes   0xC5  // Short form of vex prefix
#define PrefixVEX3Bytes   0xC4  // Long form of vex prefixs


enum E_DECODEOP{
	xFail  = 0 ,    // invalid instruction (return)
	xMatch                 ,// completed match
	xJump                  ,// jump to pc

	xCondByte     ,// switch on instruction byte value
	xCondSlashR   ,// read and switch on instruction /r value
	xCondPrefix   ,// switch on presence of instruction prefix
	xCondIs64     ,// switch on 64-bit processor mode
	xCondDataSize ,// switch on operand size
	xCondAddrSize ,// switch on address size
	xCondIsMem    ,// switch on memory vs register argument

	xSetOp ,// set instruction opcode

	xReadSlashR ,// read /r
	xReadIb     ,// read ib
	xReadIw     ,// read iw
	xReadId     ,// read id
	xReadIo     ,// read io
	xReadCb     ,// read cb
	xReadCw     ,// read cw
	xReadCd     ,// read cd
	xReadCp     ,// read cp
	xReadCm     ,// read cm

	xArg1            ,// arg 1
	xArg3            ,// arg 3
	xArgAL           ,// arg AL
	xArgAX           ,// arg AX
	xArgCL           ,// arg CL
	xArgCR0dashCR7   ,// arg CR0-CR7
	xArgCS           ,// arg CS
	xArgDR0dashDR7   ,// arg DR0-DR7
	xArgDS           ,// arg DS
	xArgDX           ,// arg DX
	xArgEAX          ,// arg EAX
	xArgEDX          ,// arg EDX
	xArgES           ,// arg ES
	xArgFS           ,// arg FS
	xArgGS           ,// arg GS
	xArgImm16        ,// arg imm16
	xArgImm32        ,// arg imm32
	xArgImm64        ,// arg imm64
	xArgImm8         ,// arg imm8
	xArgImm8u        ,// arg imm8 but record as unsigned
	xArgImm16u       ,// arg imm8 but record as unsigned
	xArgM            ,// arg m
	xArgM128         ,// arg m128
	xArgM256         ,// arg m256
	xArgM1428byte    ,// arg m14/28byte
	xArgM16          ,// arg m16
	xArgM16and16     ,// arg m16&16
	xArgM16and32     ,// arg m16&32
	xArgM16and64     ,// arg m16&64
	xArgM16colon16   ,// arg m16=16
	xArgM16colon32   ,// arg m16=32
	xArgM16colon64   ,// arg m16=64
	xArgM16int       ,// arg m16int
	xArgM2byte       ,// arg m2byte
	xArgM32          ,// arg m32
	xArgM32and32     ,// arg m32&32
	xArgM32fp        ,// arg m32fp
	xArgM32int       ,// arg m32int
	xArgM512byte     ,// arg m512byte
	xArgM64          ,// arg m64
	xArgM64fp        ,// arg m64fp
	xArgM64int       ,// arg m64int
	xArgM8           ,// arg m8
	xArgM80bcd       ,// arg m80bcd
	xArgM80dec       ,// arg m80dec
	xArgM80fp        ,// arg m80fp
	xArgM94108byte   ,// arg m94/108byte
	xArgMm           ,// arg mm
	xArgMm1          ,// arg mm1
	xArgMm2          ,// arg mm2
	xArgMm2M64       ,// arg mm2/m64
	xArgMmM32        ,// arg mm/m32
	xArgMmM64        ,// arg mm/m64
	xArgMem          ,// arg mem
	xArgMoffs16      ,// arg moffs16
	xArgMoffs32      ,// arg moffs32
	xArgMoffs64      ,// arg moffs64
	xArgMoffs8       ,// arg moffs8
	xArgPtr16colon16 ,// arg ptr16=16
	xArgPtr16colon32 ,// arg ptr16=32
	xArgR16          ,// arg r16
	xArgR16op        ,// arg r16 with +rw in opcode
	xArgR32          ,// arg r32
	xArgR32M16       ,// arg r32/m16
	xArgR32M8        ,// arg r32/m8
	xArgR32op        ,// arg r32 with +rd in opcode
	xArgR64          ,// arg r64
	xArgR64M16       ,// arg r64/m16
	xArgR64op        ,// arg r64 with +rd in opcode
	xArgR8           ,// arg r8
	xArgR8op         ,// arg r8 with +rb in opcode
	xArgRAX          ,// arg RAX
	xArgRDX          ,// arg RDX
	xArgRM           ,// arg r/m
	xArgRM16         ,// arg r/m16
	xArgRM32         ,// arg r/m32
	xArgRM64         ,// arg r/m64
	xArgRM8          ,// arg r/m8
	xArgReg          ,// arg reg
	xArgRegM16       ,// arg reg/m16
	xArgRegM32       ,// arg reg/m32
	xArgRegM8        ,// arg reg/m8
	xArgRel16        ,// arg rel16
	xArgRel32        ,// arg rel32
	xArgRel8         ,// arg rel8
	xArgSS           ,// arg SS
	xArgST           ,// arg ST, aka ST(0)
	xArgSTi          ,// arg ST(i) with +i in opcode
	xArgSreg         ,// arg Sreg
	xArgTR0dashTR7   ,// arg TR0-TR7
	xArgXmm          ,// arg xmm
	xArgXMM0         ,// arg <XMM0>
	xArgXmm1         ,// arg xmm1
	xArgXmm2         ,// arg xmm2
	xArgXmm2M128     ,// arg xmm2/m128
	xArgYmm2M256     ,// arg ymm2/m256
	xArgXmm2M16      ,// arg xmm2/m16
	xArgXmm2M32      ,// arg xmm2/m32
	xArgXmm2M64      ,// arg xmm2/m64
	xArgXmmM128      ,// arg xmm/m128
	xArgXmmM32       ,// arg xmm/m32
	xArgXmmM64       ,// arg xmm/m64
	xArgYmm1         ,// arg ymm1
	xArgRmf16        ,// arg r/m16 but force mod=3
	xArgRmf32        ,// arg r/m32 but force mod=3
	xArgRmf64        ,// arg r/m64 but force mod=3
};

// opecode code define
enum E_OP{
	// _ Op = iota
	AAA = 1,
	AAD,
	AAM,
	AAS,
	ADC,
	ADD,
	ADDPD,
	ADDPS,
	ADDSD,
	ADDSS,
	ADDSUBPD,
	ADDSUBPS,
	AESDEC,
	AESDECLAST,
	AESENC,
	AESENCLAST,
	AESIMC,
	AESKEYGENASSIST,
	AND,
	ANDNPD,
	ANDNPS,
	ANDPD,
	ANDPS,
	ARPL,
	BLENDPD,
	BLENDPS,
	BLENDVPD,
	BLENDVPS,
	BOUND,
	BSF,
	BSR,
	BSWAP,
	BT,
	BTC,
	BTR,
	BTS,
	CALL,
	CBW,
	CDQ,
	CDQE,
	CLC,
	CLD,
	CLFLUSH,
	CLI,
	CLTS,
	CMC,
	CMOVA,
	CMOVAE,
	CMOVB,
	CMOVBE,
	CMOVE,
	CMOVG,
	CMOVGE,
	CMOVL,
	CMOVLE,
	CMOVNE,
	CMOVNO,
	CMOVNP,
	CMOVNS,
	CMOVO,
	CMOVP,
	CMOVS,
	CMP,
	CMPPD,
	CMPPS,
	CMPSB,
	CMPSD,
	CMPSD_XMM,
	CMPSQ,
	CMPSS,
	CMPSW,
	CMPXCHG,
	CMPXCHG16B,
	CMPXCHG8B,
	COMISD,
	COMISS,
	CPUID,
	CQO,
	CRC32,
	CVTDQ2PD,
	CVTDQ2PS,
	CVTPD2DQ,
	CVTPD2PI,
	CVTPD2PS,
	CVTPI2PD,
	CVTPI2PS,
	CVTPS2DQ,
	CVTPS2PD,
	CVTPS2PI,
	CVTSD2SI,
	CVTSD2SS,
	CVTSI2SD,
	CVTSI2SS,
	CVTSS2SD,
	CVTSS2SI,
	CVTTPD2DQ,
	CVTTPD2PI,
	CVTTPS2DQ,
	CVTTPS2PI,
	CVTTSD2SI,
	CVTTSS2SI,
	CWD,
	CWDE,
	DAA,
	DAS,
	DEC,
	DIV,
	DIVPD,
	DIVPS,
	DIVSD,
	DIVSS,
	DPPD,
	DPPS,
	EMMS,
	ENTER,
	EXTRACTPS,
	F2XM1,
	FABS,
	FADD,
	FADDP,
	FBLD,
	FBSTP,
	FCHS,
	FCMOVB,
	FCMOVBE,
	FCMOVE,
	FCMOVNB,
	FCMOVNBE,
	FCMOVNE,
	FCMOVNU,
	FCMOVU,
	FCOM,
	FCOMI,
	FCOMIP,
	FCOMP,
	FCOMPP,
	FCOS,
	FDECSTP,
	FDIV,
	FDIVP,
	FDIVR,
	FDIVRP,
	FFREE,
	FFREEP,
	FIADD,
	FICOM,
	FICOMP,
	FIDIV,
	FIDIVR,
	FILD,
	FIMUL,
	FINCSTP,
	FIST,
	FISTP,
	FISTTP,
	FISUB,
	FISUBR,
	FLD,
	FLD1,
	FLDCW,
	FLDENV,
	FLDL2E,
	FLDL2T,
	FLDLG2,
	FLDLN2,
	FLDPI,
	FLDZ,
	FMUL,
	FMULP,
	FNCLEX,
	FNINIT,
	FNOP,
	FNSAVE,
	FNSTCW,
	FNSTENV,
	FNSTSW,
	FPATAN,
	FPREM,
	FPREM1,
	FPTAN,
	FRNDINT,
	FRSTOR,
	FSCALE,
	FSIN,
	FSINCOS,
	FSQRT,
	FST,
	FSTP,
	FSUB,
	FSUBP,
	FSUBR,
	FSUBRP,
	FTST,
	FUCOM,
	FUCOMI,
	FUCOMIP,
	FUCOMP,
	FUCOMPP,
	FWAIT,
	FXAM,
	FXCH,
	FXRSTOR,
	FXRSTOR64,
	FXSAVE,
	FXSAVE64,
	FXTRACT,
	FYL2X,
	FYL2XP1,
	HADDPD,
	HADDPS,
	HLT,
	HSUBPD,
	HSUBPS,
	ICEBP,
	IDIV,
	IMUL,
	IN,
	INC,
	INSB,
	INSD,
	INSERTPS,
	INSW,
	INT,
	INTO,
	INVD,
	INVLPG,
	INVPCID,
	IRET,
	IRETD,
	IRETQ,
	JA,
	JAE,
	JB,
	JBE,
	JCXZ,
	JE,
	JECXZ,
	JG,
	JGE,
	JL,
	JLE,
	JMP,
	JNE,
	JNO,
	JNP,
	JNS,
	JO,
	JP,
	JRCXZ,
	JS,
	LAHF,
	LAR,
	LCALL,
	LDDQU,
	LDMXCSR,
	LDS,
	LEA,
	LEAVE,
	LES,
	LFENCE,
	LFS,
	LGDT,
	LGS,
	LIDT,
	LJMP,
	LLDT,
	LMSW,
	LODSB,
	LODSD,
	LODSQ,
	LODSW,
	LOOP,
	LOOPE,
	LOOPNE,
	LRET,
	LSL,
	LSS,
	LTR,
	LZCNT,
	MASKMOVDQU,
	MASKMOVQ,
	MAXPD,
	MAXPS,
	MAXSD,
	MAXSS,
	MFENCE,
	MINPD,
	MINPS,
	MINSD,
	MINSS,
	MONITOR,
	MOV,
	MOVAPD,
	MOVAPS,
	MOVBE,
	MOVD,
	MOVDDUP,
	MOVDQ2Q,
	MOVDQA,
	MOVDQU,
	MOVHLPS,
	MOVHPD,
	MOVHPS,
	MOVLHPS,
	MOVLPD,
	MOVLPS,
	MOVMSKPD,
	MOVMSKPS,
	MOVNTDQ,
	MOVNTDQA,
	MOVNTI,
	MOVNTPD,
	MOVNTPS,
	MOVNTQ,
	MOVNTSD,
	MOVNTSS,
	MOVQ,
	MOVQ2DQ,
	MOVSB,
	MOVSD,
	MOVSD_XMM,
	MOVSHDUP,
	MOVSLDUP,
	MOVSQ,
	MOVSS,
	MOVSW,
	MOVSX,
	MOVSXD,
	MOVUPD,
	MOVUPS,
	MOVZX,
	MPSADBW,
	MUL,
	MULPD,
	MULPS,
	MULSD,
	MULSS,
	MWAIT,
	NEG,
	NOP,
	NOT,
	OR,
	ORPD,
	ORPS,
	OUT,
	OUTSB,
	OUTSD,
	OUTSW,
	PABSB,
	PABSD,
	PABSW,
	PACKSSDW,
	PACKSSWB,
	PACKUSDW,
	PACKUSWB,
	PADDB,
	PADDD,
	PADDQ,
	PADDSB,
	PADDSW,
	PADDUSB,
	PADDUSW,
	PADDW,
	PALIGNR,
	PAND,
	PANDN,
	PAUSE,
	PAVGB,
	PAVGW,
	PBLENDVB,
	PBLENDW,
	PCLMULQDQ,
	PCMPEQB,
	PCMPEQD,
	PCMPEQQ,
	PCMPEQW,
	PCMPESTRI,
	PCMPESTRM,
	PCMPGTB,
	PCMPGTD,
	PCMPGTQ,
	PCMPGTW,
	PCMPISTRI,
	PCMPISTRM,
	PEXTRB,
	PEXTRD,
	PEXTRQ,
	PEXTRW,
	PHADDD,
	PHADDSW,
	PHADDW,
	PHMINPOSUW,
	PHSUBD,
	PHSUBSW,
	PHSUBW,
	PINSRB,
	PINSRD,
	PINSRQ,
	PINSRW,
	PMADDUBSW,
	PMADDWD,
	PMAXSB,
	PMAXSD,
	PMAXSW,
	PMAXUB,
	PMAXUD,
	PMAXUW,
	PMINSB,
	PMINSD,
	PMINSW,
	PMINUB,
	PMINUD,
	PMINUW,
	PMOVMSKB,
	PMOVSXBD,
	PMOVSXBQ,
	PMOVSXBW,
	PMOVSXDQ,
	PMOVSXWD,
	PMOVSXWQ,
	PMOVZXBD,
	PMOVZXBQ,
	PMOVZXBW,
	PMOVZXDQ,
	PMOVZXWD,
	PMOVZXWQ,
	PMULDQ,
	PMULHRSW,
	PMULHUW,
	PMULHW,
	PMULLD,
	PMULLW,
	PMULUDQ,
	POP,
	POPA,
	POPAD,
	POPCNT,
	POPF,
	POPFD,
	POPFQ,
	POR,
	PREFETCHNTA,
	PREFETCHT0,
	PREFETCHT1,
	PREFETCHT2,
	PREFETCHW,
	PSADBW,
	PSHUFB,
	PSHUFD,
	PSHUFHW,
	PSHUFLW,
	PSHUFW,
	PSIGNB,
	PSIGND,
	PSIGNW,
	PSLLD,
	PSLLDQ,
	PSLLQ,
	PSLLW,
	PSRAD,
	PSRAW,
	PSRLD,
	PSRLDQ,
	PSRLQ,
	PSRLW,
	PSUBB,
	PSUBD,
	PSUBQ,
	PSUBSB,
	PSUBSW,
	PSUBUSB,
	PSUBUSW,
	PSUBW,
	PTEST,
	PUNPCKHBW,
	PUNPCKHDQ,
	PUNPCKHQDQ,
	PUNPCKHWD,
	PUNPCKLBW,
	PUNPCKLDQ,
	PUNPCKLQDQ,
	PUNPCKLWD,
	PUSH,
	PUSHA,
	PUSHAD,
	PUSHF,
	PUSHFD,
	PUSHFQ,
	PXOR,
	RCL,
	RCPPS,
	RCPSS,
	RCR,
	RDFSBASE,
	RDGSBASE,
	RDMSR,
	RDPMC,
	RDRAND,
	RDTSC,
	RDTSCP,
	RET,
	ROL,
	ROR,
	ROUNDPD,
	ROUNDPS,
	ROUNDSD,
	ROUNDSS,
	RSM,
	RSQRTPS,
	RSQRTSS,
	SAHF,
	SAR,
	SBB,
	SCASB,
	SCASD,
	SCASQ,
	SCASW,
	SETA,
	SETAE,
	SETB,
	SETBE,
	SETE,
	SETG,
	SETGE,
	SETL,
	SETLE,
	SETNE,
	SETNO,
	SETNP,
	SETNS,
	SETO,
	SETP,
	SETS,
	SFENCE,
	SGDT,
	SHL,
	SHLD,
	SHR,
	SHRD,
	SHUFPD,
	SHUFPS,
	SIDT,
	SLDT,
	SMSW,
	SQRTPD,
	SQRTPS,
	SQRTSD,
	SQRTSS,
	STC,
	STD,
	STI,
	STMXCSR,
	STOSB,
	STOSD,
	STOSQ,
	STOSW,
	STR,
	SUB,
	SUBPD,
	SUBPS,
	SUBSD,
	SUBSS,
	SWAPGS,
	SYSCALL,
	SYSENTER,
	SYSEXIT,
	SYSRET,
	TEST,
	TZCNT,
	UCOMISD,
	UCOMISS,
	UD0,
	UD1,
	UD2,
	UNPCKHPD,
	UNPCKHPS,
	UNPCKLPD,
	UNPCKLPS,
	VERR,
	VERW,
	VMOVDQA,
	VMOVDQU,
	VMOVNTDQ,
	VMOVNTDQA,
	VZEROUPPER,
	WBINVD,
	WRFSBASE,
	WRGSBASE,
	WRMSR,
	XABORT,
	XADD,
	XBEGIN,
	XCHG,
	XEND,
	XGETBV,
	XLATB,
	XOR,
	XORPD,
	XORPS,
	XRSTOR,
	XRSTOR64,
	XRSTORS,
	XRSTORS64,
	XSAVE,
	XSAVE64,
	XSAVEC,
	XSAVEC64,
	XSAVEOPT,
	XSAVEOPT64,
	XSAVES,
	XSAVES64,
	XSETBV,
	XTEST
};

enum E_REG {
	// _ Reg = iota

	// 8-bit
	AL = 1,
	CL,
	DL,
	BL,
	AH,
	CH,
	DH,
	BH,
	SPB,
	BPB,
	SIB,
	DIB,
	R8B,
	R9B,
	R10B,
	R11B,
	R12B,
	R13B,
	R14B,
	R15B,

	// 16-bit
	AX,
	CX,
	DX,
	BX,
	SP,
	BP,
	SI,
	DI,
	R8W,
	R9W,
	R10W,
	R11W,
	R12W,
	R13W,
	R14W,
	R15W,

	// 32-bit
	EAX,
	ECX,
	EDX,
	EBX,
	ESP,
	EBP,
	ESI,
	EDI,
	R8L,
	R9L,
	R10L,
	R11L,
	R12L,
	R13L,
	R14L,
	R15L,

	// 64-bit
	RAX,
	RCX,
	RDX,
	RBX,
	RSP,
	RBP,
	RSI,
	RDI,
	R8,
	R9,
	R10,
	R11,
	R12,
	R13,
	R14,
	R15,

	// Instruction pointer.
	IP,  // 16-bit
	EIP, // 32-bit
	RIP, // 64-bit

	// 387 floating point registers.
	F0,
	F1,
	F2,
	F3,
	F4,
	F5,
	F6,
	F7,

	// MMX registers.
	M0,
	M1,
	M2,
	M3,
	M4,
	M5,
	M6,
	M7,

	// XMM registers.
	X0,
	X1,
	X2,
	X3,
	X4,
	X5,
	X6,
	X7,
	X8,
	X9,
	X10,
	X11,
	X12,
	X13,
	X14,
	X15,

	// Segment registers.
	ES,
	CS,
	SS,
	DS,
	FS,
	GS,

	// System registers.
	GDTR,
	IDTR,
	LDTR,
	MSW,
	TASK,

	// Control registers.
	CR0,
	CR1,
	CR2,
	CR3,
	CR4,
	CR5,
	CR6,
	CR7,
	CR8,
	CR9,
	CR10,
	CR11,
	CR12,
	CR13,
	CR14,
	CR15,

	// Debug registers.
	DR0,
	DR1,
	DR2,
	DR3,
	DR4,
	DR5,
	DR6,
	DR7,
	DR8,
	DR9,
	DR10,
	DR11,
	DR12,
	DR13,
	DR14,
	DR15,

	// Task registers.
	TR0,
	TR1,
	TR2,
	TR3,
	TR4,
	TR5,
	TR6,
	TR7,
};


typedef struct{
	Reg Segment;
	Reg Base;
	uint8_t Scale;
	Reg Index; 
	int64_t Disp;
}Mem;

extern Mem addr16[16];
extern int8_t memBytes[123];
extern Reg baseReg[127];
extern bool isCondJmp[XTEST+1];
extern bool isLoop[XTEST + 1];
extern const char* opNames[614];
extern uint16 decoder[13430] ;

const char* op_name(uint16 op);
const char* reg_name(Reg reg);
const char* op_str(uint16 op);