#include "goX86.h"
#include "table.h"

typedef uint16_t Prefix;
#ifndef bool
typedef uint8_t bool;
#define true  ((uint8_t)1)
#define false ((uint8_t)0)
#endif
// typedef [14]uint16_t Prefixes;


Prefix PrefixImplicit  = 0x8000; // prefix is implied by instruction text
Prefix PrefixIgnored   = 0x4000; // prefix is ignored: either irrelevant or overridden by a later prefix
Prefix PrefixInvalid   = 0x2000; // prefix makes entire instruction invalid (bad LOCK)

// Memory segment overrides.
Prefix PrefixES  = 0x26; // ES segment override
Prefix PrefixCS  = 0x2E; // CS segment override
Prefix PrefixSS  = 0x36; // SS segment override
Prefix PrefixDS  = 0x3E; // DS segment override
Prefix PrefixFS  = 0x64; // FS segment override
Prefix PrefixGS  = 0x65; // GS segment override

// Branch prediction.
Prefix PrefixPN  = 0x12E; // predict not taken (conditional branch only)
Prefix PrefixPT  = 0x13E; // predict taken (conditional branch only)

// Size attributes.
const Prefix PrefixDataSize  = 0x66; // operand size override
Prefix PrefixData16    = 0x166;
Prefix PrefixData32    = 0x266;
const Prefix PrefixAddrSize  = 0x67; // address size override
Prefix PrefixAddr16    = 0x167;
Prefix PrefixAddr32    = 0x267;

// One of a kind.
Prefix PrefixLOCK      = 0xF0; // lock
Prefix PrefixREPN      = 0xF2; // repeat not zero
Prefix PrefixXACQUIRE  = 0x1F2;
Prefix PrefixBND       = 0x2F2;
Prefix PrefixREP       = 0xF3; // repeat
Prefix PrefixXRELEASE  = 0x1F3;

// The REX prefixes must be in the range [PrefixREX, PrefixREX+0x10).
// the other bits are set or not according to the intended use.
Prefix PrefixREX        = 0x40; // REX 64-bit extension prefix
Prefix PrefixREXW       = 0x08; // extension bit W (64-bit instruction width)
Prefix PrefixREXR       = 0x04; // extension bit R (r field in modrm)
Prefix PrefixREXX       = 0x02; // extension bit X (index field in sib)
Prefix PrefixREXB       = 0x01; // extension bit B (r/m field in modrm or base field in sib)
Prefix PrefixVEX2Bytes  = 0xC5; // Short form of vex prefix
Prefix PrefixVEX3Bytes  = 0xC4; // Long form of vex prefixs


Inst instPrefix(Reg b,int mode );
Inst truncated(const Reg* b,int len,int mode );

Reg baseRegForBits(int bits);

uint16_t toLittleEndianU16(const uint8_t* x) {
	return (x[1] << 8) | x[0];
}

uint32_t toLittleEndianU32(const uint8_t* x) {
	return x[3]<<24 | x[2]<<16 | (x[1] << 8) | x[0];
}

uint64_t toLittleEndianU64(const uint8_t* x) {
	return x[7]<< 56 |x[6]<<48| x[5]<<40 | x[4]<<32 |x[3]<<24 | x[2]<<16 | (x[1] << 8) | x[0];
}

bool IsREX( Prefix prefix ){
	return (prefix & 0xF0 ) == PrefixREX ? true : false ;
}

Reg baseRegForBits(int  bits ) {
	switch (bits) {
	case 8:
		return AL;
	case 16:
		return AX;
	case 32:
		return EAX;
	case 64:
		return RAX;
	}
	return 0;
}

Inst truncated(const Reg* b,int len,int mode )
{
	if(len == 0){
		Inst inst = {0};
		return inst;
	}
	return instPrefix(b[0],mode);
}

Reg prefixToSegment(Prefix  p ) {
	switch (p & ~PrefixImplicit ){
	case PrefixCS:
		return CS;
	case PrefixDS:
		return DS;
	case PrefixES:
		return ES;
	case PrefixFS:
		return FS;
	case PrefixGS:
		return GS;
	case PrefixSS:
		return SS;
	}
	return 0;
}


Inst instPrefix(Reg b,int mode )
{
	// When tracing it is useful to see what called instPrefix to report an error.
	// if trace { todo add trace
	// 	_, file, line, _ := runtime.Caller(1)
	// 	fmt.Printf("%s:%d\n", file, line)
	// }
	Prefix p = b;
	switch (p) {
	case PrefixDataSize:
		if (mode == 16) {
			p = PrefixData32;
		} else {
			p = PrefixData16;
		}
		break;
	case PrefixAddrSize:
		if (mode == 32) {
			p = PrefixAddr16;
		} else {
			p = PrefixAddr32;
		}
		break;
	}
	// Note: using composite literal with Prefix key confuses 'bundle' tool.
	Inst inst = {0};
	inst.Len = 1;
	inst.Prefix[0] = p;
	return inst;
}



unsigned int decode(Reg *src,int len, Inst *ld, int mode,bool gnuCompat)
{
	if(src == NULL || ld == NULL){
		return -1;
	}

	switch(mode)
	{
	case 16:
    case 32:
    case 64:
        break;
		// ok
		// TODO(rsc): 64-bit mode not tested, probably not working.
	default:
		return -1;
	}

	// Maximum instruction size is 15 bytes.
	// If we need to read more, return 'truncated instruction.
	if (len > 15) {
		len = 15;
	}

		// prefix decoding information
	int	pos           = 0;// position reading src
	int	nprefix       = 0;// number of prefixes
	int	lockIndex     = -1;// index of LOCK prefix in src and inst.Prefix
	int	repIndex      = -1;// index of REP/REPN prefix in src and inst.Prefix
	int	segIndex      = -1;// index of Group 2 prefix in src and inst.Prefix
	int	dataSizeIndex = -1; // index of Group 3 prefix in src and inst.Prefix
	int	addrSizeIndex = -1; // index of Group 4 prefix in src and inst.Prefix
	Prefix	rex  =0; // rex byte if present (or 0)
	Prefix	rexUsed  =0;      // bits used in rex byte
	int	rexIndex      = -1 ;  // index of rex byte
	Prefix	vex = 0      ;     // use vex encoding
	int	vexIndex      = -1;   // index of vex prefix

	int	addrMode = mode ;// address mode (width in bits)
	int	dataMode = mode ;// operand mode (width in bits)

		// decoded ModR/M fields
	bool haveModrm = false;
	int	modrm    =0;
	int	mod  =0;
	int	regop   =0 ;
	int	rm       =0;

		// if ModR/M is memory reference, Mem form
	Mem	mem     ={0};
	bool haveMem = false ;

		// decoded SIB fields
	bool	haveSIB =false;
	int sib     =0;
	int	scale   =0;
	int	index   =0;
	int	base    =0;
	int	displen =0;
	int	dispoff =0;

		// decoded immediate values
	int64_t	imm   =0; 
	int8_t	imm8  =0; 
	int64_t	immc  =0; 
	int	immcpos=0;

		// output
	int	opshift = 0;
	Inst	inst  ={0};  
	int 	narg   ={0}; // number of arguments written to inst
	

	if (mode == 64) {
		dataMode = 32;
	}

	// Prefixes are certainly the most complex and underspecified part of
	// decoding x86 instructions. Although the manuals say things like
	// up to four prefixes, one from each group, nearly everyone seems to
	// agree that in practice as many prefixes as possible, including multiple
	// from a particular group or repetitions of a given prefix, can be used on
	// an instruction, provided the total instruction length including prefixes
	// does not exceed the agreed-upon maximum of 15 bytes.
	// Everyone also agrees that if one of these prefixes is the LOCK prefix
	// and the instruction is not one of the instructions that can be used with
	// the LOCK prefix or if the destination is not a memory operand,
	// then the instruction is invalid and produces the #UD exception.
	// However, that is the end of any semblance of agreement.
	//
	// What happens if prefixes are given that conflict with other prefixes?
	// For example, the memory segment overrides CS, DS, ES, FS, GS, SS
	// conflict with each other: only one segment can be in effect.
	// Disassemblers seem to agree that later prefixes take priority over
	// earlier ones. I have not taken the time to write assembly programs
	// to check to see if the hardware agrees.
	//
	// What happens if prefixes are given that have no meaning for the
	// specific instruction to which they are attached? It depends.
	// If they really have no meaning, they are ignored. However, a future
	// processor may assign a different meaning. As a disassembler, we
	// don't really know whether we're seeing a meaningless prefix or one
	// whose meaning we simply haven't been told yet.
	//
	// Combining the two questions, what happens when conflicting
	// extension prefixes are given? No one seems to know for sure.
	// For example, MOVQ is 66 0F D6 /r, MOVDQ2Q is F2 0F D6 /r,
	// and MOVQ2DQ is F3 0F D6 /r. What is '66 F2 F3 0F D6 /r'?
	// Which prefix wins? See the xCondPrefix prefix for more.
	//
	// Writing assembly test cases to divine which interpretation the
	// CPU uses might clarify the situation, but more likely it would
	// make the situation even less clear.

	// Read non-REX prefixes.
	int pos =0;
	Prefix p =0;
	for (; pos < len; pos++) {
		p = (Prefix)src[pos];
		switch (p) {

		// Group 1 - lock and repeat prefixes
		// According to Intel, there should only be one from this set,
		// but according to AMD both can be present.
		case 0xF0:
			if (lockIndex >= 0) {
				inst.Prefix[lockIndex] |= PrefixIgnored;
			}
			lockIndex = pos;
			break;
		case 0xF2:
		case 0xF3:
			if (repIndex >= 0) {
				inst.Prefix[repIndex] |= PrefixIgnored;
			}
			repIndex = pos;
			break;
		// Group 2 - segment override / branch hints
		case 0x26:
		case 0x2E: 
		case 0x36:
		case 0x3E:
			if (mode == 64) {
				p |= PrefixIgnored;
				break;
			}
			// fallthrough
		case 0x64:
		case 0x65:
			if (segIndex >= 0) {
				inst.Prefix[segIndex] |= PrefixIgnored;
			}
			segIndex = pos;
			break;
		// Group 3 - operand size override
		case 0x66:
			if (mode == 16) {
				dataMode = 32;
				p = PrefixData32;
			} else {
				dataMode = 16;
				p = PrefixData16;
			}
			if (dataSizeIndex >= 0) {
				inst.Prefix[dataSizeIndex] |= PrefixIgnored;
			}
			dataSizeIndex = pos;
			break;
		// Group 4 - address size override
		case 0x67:
			if (mode == 32) {
				addrMode = 16;
				p = PrefixAddr16;
			} else {
				addrMode = 32;
				p = PrefixAddr32;
			}
			if (addrSizeIndex >= 0 ){
				inst.Prefix[addrSizeIndex] |= PrefixIgnored;
			}
			addrSizeIndex = pos;
			break;
		//Group 5 - Vex encoding
		case 0xC5:
			if( pos == 0 && pos+1 < len && (mode == 64 || (mode == 32 && src[pos+1]&0xc0 == 0xc0)) ){
				vex = p;
				vexIndex = pos;
				inst.Prefix[pos] = p;
				inst.Prefix[pos+1] = (Prefix)src[pos+1];
				pos += 1;
				continue;
			} else {
				nprefix = pos;
				goto ReadPrefixes;
			}
		case 0xC4:
			if (pos == 0 && pos+2 < len && (mode == 64 || (mode == 32 && src[pos+1]&0xc0 == 0xc0)) ){
				vex = p;
				vexIndex = pos;
				inst.Prefix[pos] = p;
				inst.Prefix[pos+1] = (Prefix)src[pos+1];
				inst.Prefix[pos+2] = (Prefix)src[pos+2];
				pos += 2;
				continue;
			} else {
				nprefix = pos;
				goto ReadPrefixes;
			}
		default:
			nprefix = pos;
			break;
		}

		if (pos >= PREFIX_SIZE ) {
			*ld = instPrefix(src[0],mode);
			return 0; // instPrefix(src[0]); // too long
		}
		inst.Prefix[pos] = p;
	}

// goto here	
ReadPrefixes:

	// Read REX prefix.
	if (pos < len && mode == 64 && IsREX(src[pos]) && vex == 0 ){
		rex = src[pos];
		rexIndex = pos;
		if (pos >= PREFIX_SIZE ) {
			*ld = instPrefix(src[0],mode);
			return 0;
		}
		inst.Prefix[pos] = rex;
		pos++;
		if (rex & PrefixREXW != 0 ){
			dataMode = 64;
			if (dataSizeIndex >= 0 ){
				inst.Prefix[dataSizeIndex] |= PrefixIgnored;
			}
		}
	}

	// Decode instruction stream, interpreting decoding instructions.
	// opshift gives the shift to use when saving the next
	// opcode byte into inst.Opcode.
	opshift = 24;

	// Decode loop, executing decoder program.
	int oldPC = 0, prevPC = 0;
	int pc = 1;

// Decode:

	for (; ; ){ // TODO uint
		oldPC = prevPC;
		prevPC = pc;
		// if trace {
		// 	println("run", pc)
		// }

		uint16_t x = decoder[pc];
		// ignore cover report
		// if decoderCover != nil {
		// 	decoderCover[pc] = true
		// }
		pc++;

		// Read and decode ModR/M if needed by opcode.
		switch (x) {
		case xCondSlashR:
		case xReadSlashR:
			if (haveModrm) {
				ld->Len = pos;
				//todo errors.New("internal error")print
				return -1;
				// return Inst{Len: pos}, errInternal
			}

			haveModrm = true;
			if (pos >= len) {
				*ld = truncated(src, len,mode);
				return 0;
			}

			modrm = src[pos];
			pos++;
			if (opshift >= 0 ){
				inst.Opcode |= uint32(modrm) << uint(opshift);
				opshift -= 8;
			}
			mod = modrm >> 6;
			regop = (modrm >> 3) & 07;
			rm = modrm & 07;
			if (rex&PrefixREXR != 0 ){
				rexUsed |= PrefixREXR;
				regop |= 8;
			}
			if (addrMode == 16) {
				// 16-bit modrm form
				if( mod != 3 ){
					haveMem = true;
					mem = addr16[rm];
					if( rm == 6 && mod == 0 ){
						mem.Base = 0;
					}

					// Consume disp16 if present.
					if (mod == 0 && rm == 6 || mod == 2 ){
						if(pos+2 > len ){
							*ld = truncated(src, len,mode);
							return 0;
						}
						mem.Disp = toLittleEndianU16(&src[pos]);
						pos += 2;
					}

					// Consume disp8 if present.
					if (mod == 1 ){
						if (pos >= len ){
							*ld = truncated(src, len,mode);
							return 0;
						}
						mem.Disp = src[pos];
						pos++;
					}
				}
			} else {
				haveMem = (mod != 3) ? true : false;

				// 32-bit or 64-bit form
				// Consume SIB encoding if present.
				if (rm == 4 && mod != 3) {
					haveSIB = true;
					if(  pos >= len) {
						*ld = truncated(src, len,mode);
						return 0;
					}
					sib = src[pos];
					pos++;
					if (opshift >= 0 ){
						inst.Opcode |= (uint32_t)sib << (uint32_t)opshift;
						opshift -= 8;
					}
					scale = sib >> 6;
					index = (sib >> 3) & 07;
					base = sib & 07;
					if( rex & PrefixREXB != 0 || vex == 0xC4 && inst.Prefix[vexIndex+1]&0x20 == 0 ){
						rexUsed |= PrefixREXB;
						base |= 8;
					}
					if( rex&PrefixREXX != 0 || vex == 0xC4 && inst.Prefix[vexIndex+1]&0x40 == 0 ){
						rexUsed |= PrefixREXX;
						index |= 8;
					}

					mem.Scale = 1 << (scale);
					if (index == 4) {
						// no mem.Index
					} else {
						mem.Index = baseRegForBits(addrMode) + (Reg)index;
					}

					if (base&7 == 5 && mod == 0 ){
						// no mem.Base
					} else {
						mem.Base = baseRegForBits(addrMode) + (Reg)(base);
					}
				} else {
					if ( rex&PrefixREXB != 0 ){
						rexUsed |= PrefixREXB;
						rm |= 8;
					}
					if (mod == 0 && rm&7 == 5 || rm&7 == 4 ){
						// base omitted
					} else if (mod != 3) {
						mem.Base = baseRegForBits(addrMode) + (Reg)(rm);
					}
				}

				// Consume disp32 if present.
				if ( mod == 0 && (rm&7 == 5 || haveSIB && base&7 == 5) || mod == 2 ){
					if (pos+4 > len ){
						*ld = truncated(src, len,mode);
						return 0;
					}
					dispoff = pos;
					displen = 4;
					mem.Disp = toLittleEndianU32(&src[pos]);
					pos += 4;
				}

				// Consume disp8 if present.
				if( mod == 1 ){
					if (pos >= len ){
						*ld = truncated(src, len,mode);
						return 0;
					}
					dispoff = pos;
					displen = 1;
					mem.Disp = src[pos];
					pos++;
				}

				// In 64-bit, mod=0 rm=5 is PC-relative instead of just disp.
				// See Vol 2A. Table 2-7.
				if( mode == 64 && mod == 0 && rm&7 == 5 ){
					if (addrMode == 32) {
						mem.Base = EIP;
					} else {
						mem.Base = RIP;
					}
				}
			}

			if (segIndex >= 0 ){
				mem.Segment = prefixToSegment(inst.Prefix[segIndex]);
			}
		}

		// Execute single opcode.
		switch ( (decodeOp)(x) ){
		default:
			// todo: error
			// println("bad op", x, "at", pc-1, "from", oldPC)
			//  Inst{Len: pos}, errInternal
			return -1;

		case xFail:
			inst.Op = 0;
			goto BREAK_DECODE;

		case xMatch:
			goto BREAK_DECODE;

		case xJump:
			pc = decoder[pc];

		// Conditional branches.
			break;
		case xCondByte:
			if (pos >= len ){
				*ld = truncated(src,len, mode);
				return 0;
			}
			Reg b = src[pos];
			int n = decoder[pc],i = 0;
			pc++;
			for (; i < n; i++ ){
				// int16_t xb, xpc := decoder[pc], int(decoder[pc+1])
				Reg xb =(Reg)decoder[pc];
				int xpc = decoder[pc+1];
				pc += 2;
				if (b == xb) {
					pc = xpc;
					pos++;
					if (opshift >= 0 ){
						inst.Opcode |= (uint32_t)(b) << (opshift);
						opshift -= 8;
					}
					goto CONTINUE_DECODE;
				}
			}
			// xCondByte is the only conditional with a fall through,
			// so that it can be used to pick off special cases before
			// an xCondSlash. If the fallthrough instruction is xFail,
			// advance the position so that the decoded instruction
			// size includes the byte we just compared against.
			if ((decodeOp)(decoder[pc]) == xJump ){
				pc = decoder[pc+1];
			}
			if ( (decodeOp)(decoder[pc]) == xFail ){
				pos++;
			}
			break;
		case xCondIs64:
			if (mode == 64) {
				pc = (decoder[pc+1]);
			} else {
				pc = (decoder[pc]);
			}
			break;
		case xCondIsMem:
			bool mem = haveMem;
			if (!haveModrm) {
				if (pos >= len) {
					*ld = instPrefix(src[0], mode); // too long
					return 0;
				}
				mem = src[pos]>>6 != 3;
			}
			if (mem) {
				pc = (decoder[pc+1]);
			} else {
				pc = (decoder[pc]);
			}
			break;
		case xCondDataSize:
			switch (dataMode) {
			case 16:
				if (dataSizeIndex >= 0) {
					inst.Prefix[dataSizeIndex] |= PrefixImplicit;
				}
				pc = (decoder[pc]);
			case 32:
				if ( dataSizeIndex >= 0 ){
					inst.Prefix[dataSizeIndex] |= PrefixImplicit;
				}
				pc = (decoder[pc+1]);
			case 64:
				rexUsed |= PrefixREXW;
				pc = (decoder[pc+2]);
			}

		case xCondAddrSize:
			switch (addrMode) {
			case 16:
				if (addrSizeIndex >= 0) {
					inst.Prefix[addrSizeIndex] |= PrefixImplicit;
				}
				pc = (decoder[pc]);
				break;
			case 32:
				if (addrSizeIndex >= 0) {
					inst.Prefix[addrSizeIndex] |= PrefixImplicit;
				}
				pc = (decoder[pc+1]);
				break;
			case 64:
				pc = (decoder[pc+2]);
				break;
			}
			break;
		case xCondPrefix:
			// Conditional branch based on presence or absence of prefixes.
			// The conflict cases here are completely undocumented and
			// differ significantly between GNU libopcodes and Intel xed.
			// I have not written assembly code to divine what various CPUs
			// do, but it wouldn't surprise me if they are not consistent either.
			//
			// The basic idea is to switch on the presence of a prefix, so that
			// for example:
			//
			//	xCondPrefix, 4
			//	0xF3, 123,
			//	0xF2, 234,
			//	0x66, 345,
			//	0, 456
			//
			// branch to 123 if the F3 prefix is present, 234 if the F2 prefix
			// is present, 66 if the 345 prefix is present, and 456 otherwise.
			// The prefixes are given in descending order so that the 0 will be last.
			//
			// It is unclear what should happen if multiple conditions are
			// satisfied: what if F2 and F3 are both present, or if 66 and F2
			// are present, or if all three are present? The one chosen becomes
			// part of the opcode and the others do not. Perhaps the answer
			// depends on the specific opcodes in question.
			//
			// The only clear example is that CRC32 is F2 0F 38 F1 /r, and
			// it comes in 16-bit and 32-bit forms based on the 66 prefix,
			// so 66 F2 0F 38 F1 /r should be treated as F2 taking priority,
			// with the 66 being only an operand size override, and probably
			// F2 66 0F 38 F1 /r should be treated the same.
			// Perhaps that rule is specific to the case of CRC32, since no
			// 66 0F 38 F1 instruction is defined (today) (that we know of).
			// However, both libopcodes and xed seem to generalize this
			// example and choose F2/F3 in preference to 66, and we
			// do the same.
			//
			// Next, what if both F2 and F3 are present? Which wins?
			// The Intel xed rule, and ours, is that the one that occurs last wins.
			// The GNU libopcodes rule, which we implement only in gnuCompat mode,
			// is that F3 beats F2 unless F3 has no special meaning, in which
			// case F3 can be a modified on an F2 special meaning.
			//
			// Concretely,
			//	66 0F D6 /r is MOVQ
			//	F2 0F D6 /r is MOVDQ2Q
			//	F3 0F D6 /r is MOVQ2DQ.
			//
			//	F2 66 0F D6 /r is 66 + MOVDQ2Q always.
			//	66 F2 0F D6 /r is 66 + MOVDQ2Q always.
			//	F3 66 0F D6 /r is 66 + MOVQ2DQ always.
			//	66 F3 0F D6 /r is 66 + MOVQ2DQ always.
			//	F2 F3 0F D6 /r is F2 + MOVQ2DQ always.
			//	F3 F2 0F D6 /r is F3 + MOVQ2DQ in Intel xed, but F2 + MOVQ2DQ in GNU libopcodes.
			//	Adding 66 anywhere in the prefix section of the
			//	last two cases does not change the outcome.
			//
			// Finally, what if there is a variant in which 66 is a mandatory
			// prefix rather than an operand size override, but we know of
			// no corresponding F2/F3 form, and we see both F2/F3 and 66.
			// Does F2/F3 still take priority, so that the result is an unknown
			// instruction, or does the 66 take priority, so that the extended
			// 66 instruction should be interpreted as having a REP/REPN prefix?
			// Intel xed does the former and GNU libopcodes does the latter.
			// We side with Intel xed, unless we are trying to match libopcodes
			// more closely during the comparison-based test suite.
			//
			// In 64-bit mode REX.W is another valid prefix to test for, but
			// there is less ambiguity about that. When present, REX.W is
			// always the first entry in the table.
			uint16 n = decoder[pc];
			pc++;
			bool sawF3 = false;
			int j = 0;
			for (; j < n; j++ ){
				Prefix prefix = (decoder[pc+2*j]);

				if (IsREX(prefix)) {
					rexUsed |= prefix;
					if (rex&prefix == prefix) {
						pc = decoder[pc+2*j+1];
						// continue Decode
						goto CONTINUE_DECODE;
					}
					continue;
				}

				bool ok = false;
				if (prefix == 0) {
					ok = true;
				} else if (IsREX(prefix)) {
					rexUsed |= prefix;
					if (rex&prefix == prefix) {
						ok = true;
					}
				} else if (prefix == 0xC5 || prefix == 0xC4 ){
					if (vex == prefix ){
						ok = true;
					}
				} else if (vex != 0 && (prefix == 0x0F || prefix == 0x0F38 || prefix == 0x0F3A ||
					prefix == 0x66 || prefix == 0xF2 || prefix == 0xF3) ){
					Prefix vexM, vexP ;
					if (vex == 0xC5 ){
						vexM = 1 ;// 2 byte vex always implies 0F
						vexP = inst.Prefix[vexIndex+1];
					} else {
						vexM = inst.Prefix[vexIndex+1];
						vexP = inst.Prefix[vexIndex+2];
					}
					switch (prefix) {
					case 0x66:
						ok = vexP&3 == 1;
						break;
					case 0xF3:
						ok = vexP&3 == 2;
						break;
					case 0xF2:
						ok = vexP&3 == 3;
						break;
					case 0x0F:
						ok = vexM&3 == 1;
						break;
					case 0x0F38:
						ok = vexM&3 == 2;
						break;
					case 0x0F3A:
						ok = vexM&3 == 3;
						break;
					}
				} else {
					if (prefix == 0xF3) {
						sawF3 = true;
					}
					switch (prefix) {
					case PrefixLOCK:
						if (lockIndex >= 0) {
							inst.Prefix[lockIndex] |= PrefixImplicit;
							ok = true;
						}
						break;
					case PrefixREP:
					case PrefixREPN:
						if (repIndex >= 0 && inst.Prefix[repIndex]&0xFF == prefix ){
							inst.Prefix[repIndex] |= PrefixImplicit;
							ok = true;
						}
						if ( gnuCompat && !ok && prefix == 0xF3 && repIndex >= 0 && (j+1 >= n || decoder[pc+2*(j+1)] != 0xF2) ){
							// Check to see if earlier prefix F3 is present.
							int i = 0;
							for( i = repIndex - 1; i >= 0; i-- ){
								if (inst.Prefix[i]&0xFF == prefix ){
									inst.Prefix[i] |= PrefixImplicit;
									ok = true;
								}
							}
						}
						if (gnuCompat && !ok && prefix == 0xF2 && repIndex >= 0 && !sawF3 && inst.Prefix[repIndex]&0xFF == 0xF3 ){
							// Check to see if earlier prefix F2 is present.
							int i = 0;
							for( i = repIndex - 1; i >= 0; i-- ){
								if( inst.Prefix[i]&0xFF == prefix ){
									inst.Prefix[i] |= PrefixImplicit;
									ok = true;
								}
							}
						}
						break;
					case PrefixCS:
					case PrefixDS:
					case PrefixES:
					case PrefixFS:
					case PrefixGS:
					case PrefixSS:
						if (segIndex >= 0 && inst.Prefix[segIndex]&0xFF == prefix ){
							inst.Prefix[segIndex] |= PrefixImplicit;
							ok = true;
						}
						break;
					case PrefixDataSize:
						// Looking for 66 mandatory prefix.
						// The F2/F3 mandatory prefixes take priority when both are present.
						// If we got this far in the xCondPrefix table and an F2/F3 is present,
						// it means the table didn't have any entry for that prefix. But if 66 has
						// special meaning, perhaps F2/F3 have special meaning that we don't know.
						// Intel xed works this way, treating the F2/F3 as inhibiting the 66.
						// GNU libopcodes allows the 66 to match. We do what Intel xed does
						// except in gnuCompat mode.
						if (repIndex >= 0 && !gnuCompat ){
							inst.Op = 0;
							goto BREAK_DECODE;
						}
						if (dataSizeIndex >= 0 ){
							inst.Prefix[dataSizeIndex] |= PrefixImplicit;
							ok = true;
						}
					case PrefixAddrSize:
						if( addrSizeIndex >= 0 ){
							inst.Prefix[addrSizeIndex] |= PrefixImplicit;
							ok = true;
						}
					}
				}
				if (ok) {
					pc = (decoder[pc+2*j+1]);
					goto CONTINUE_DECODE;
				}
			}
			inst.Op = 0;
			goto BREAK_DECODE;

		case xCondSlashR:
			pc = (decoder[pc+regop&7]);
			break;
		// Input.

		case xReadSlashR:
			// done above
			break;
		case xReadIb:
			if (pos >= len) {
				*ld = truncated(src, len,mode);
				return 0;
			}
			imm8 = (src[pos]);
			pos++;
			break;
		case xReadIw:
			if (pos+2 > len) {
					*ld = truncated(src, len,mode);
				return 0;
			}
			imm = toLittleEndianU16(&src[pos]);
			pos += 2 ;
			break;
		case xReadId:
			if (pos+4 > len) {
					*ld = truncated(src, len,mode);
				return 0;
			}
			imm = toLittleEndianU32(&src[pos]);
			pos += 4;
			break;
		case xReadIo:
			if (pos+8 > len ){
					*ld = truncated(src, len,mode);
				return 0;
			}
			imm = toLittleEndianU64(&src[pos]);
			pos += 8;
			break;
		case xReadCb:
			if (pos >= len) {
				*ld = truncated(src, len,mode);
				return 0;
			}
			immcpos = pos;
			immc = src[pos];
			pos++;
			break;
		case xReadCw:
			if (pos+2 > len){
				*ld = truncated(src, len,mode);
				return 0;
			}
			immcpos = pos;
			immc = toLittleEndianU16(&src[pos]);
			pos += 2;
			break;
		case xReadCm:
			immcpos = pos;
			if (addrMode == 16) {
				if (pos+2 > len ){
					*ld = truncated(src, len,mode);
					return 0;
				}
				immc = toLittleEndianU16(&src[pos]);
				pos += 2;
			} else if (addrMode == 32) {
				if (pos+4 > len ){
					*ld = truncated(src, len,mode);
					return 0;
				}
				immc = toLittleEndianU32(src[pos]);
				pos += 4;
			} else {
				if (pos+8 > len ){
					*ld = truncated(src, len,mode);
					return 0;
				}
				immc = toLittleEndianU64(&src[pos]);
				pos += 8;
			}
			break;
		case xReadCd:
			immcpos = pos;
			if(pos+4 > len){
				*ld = truncated(src, len,mode);
				return 0;
			}
			immc = toLittleEndianU32(&src[pos]);
			pos += 4;
			break;
		case xReadCp:
			immcpos = pos;
			if (pos+6 > len) {
				*ld = truncated(src, len,mode);
				return 0;
			}
			int64_t w = toLittleEndianU32(&src[pos]);
			int64_t w2 =toLittleEndianU16(&src[pos+4]);
			immc = (w2)<<32 | (w);
			pos += 6;
			break;
		// Output.

		case xSetOp:
			inst.Op = (decoder[pc]);
			pc++;
			break;
		case 	xArg1:
		case	xArg3:
		case 	xArgAL:
		case	xArgAX:
		case	xArgCL:
		case 	xArgCS:
		case	xArgDS:
		case 	xArgDX:
		case	xArgEAX:
		case	xArgEDX:
		case	xArgES:
		case	xArgFS:
		case	xArgGS:
		case	xArgRAX:
		case	xArgRDX:
		case	xArgSS:
		case	xArgST:
		case	xArgXMM0:
			inst.Args[narg] = fixedArg[x];
			narg++;
			break;
		case xArgImm8:
			inst.Args[narg] = imm8;
			narg++;
			break;
		case xArgImm8u:
			inst.Args[narg] = imm8;
			narg++;
			break;
		case xArgImm16:
			inst.Args[narg] = imm;
			narg++;
			break;
		case xArgImm16u:
			inst.Args[narg] = imm;
			narg++;
			break;
		case xArgImm32:
			inst.Args[narg] = imm;
			narg++;
			break;
		case xArgImm64:
			inst.Args[narg] =imm;
			narg++;
			break;
		case xArgM:
		case xArgM128:
		case xArgM256:
		case xArgM1428byte:
		case xArgM16:
		case	xArgM16and16:
		case	xArgM16and32:
		case	xArgM16and64:
		case	xArgM16colon16:
		case	xArgM16colon32:
		case	xArgM16colon64:
		case	xArgM16int:
		case	xArgM2byte:
		case	xArgM32:
		case	xArgM32and32:
		case	xArgM32fp:
		case	xArgM32int:
		case	xArgM512byte:
		case	xArgM64:
		case	xArgM64fp:
		case	xArgM64int:
		case	xArgM8:
		case	xArgM80bcd:
		case	xArgM80dec:
		case	xArgM80fp:
		case	xArgM94108byte:
		case	xArgMem:
			if (!haveMem) {
				inst.Op = 0;
				goto BREAK_DECODE;
			}
			inst.Args[narg] = mem;
			inst.MemBytes = memBytes[x];
			if (mem.Base == RIP) {
				inst.PCRel = displen;
				inst.PCRelOff = dispoff;
			}
			narg++;
			break;
		case xArgPtr16colon16:
			inst.Args[narg] = (immc >> 16);
			inst.Args[narg+1] = (immc & (1<<16 - 1));
			narg += 2;
			break;
		case xArgPtr16colon32:
			inst.Args[narg] = (immc >> 32);
			inst.Args[narg+1] = (immc & (1<<32 - 1));
			narg += 2;
			break;
		case xArgMoffs8:
		case xArgMoffs16:
		case xArgMoffs32:
		case xArgMoffs64:
			{	
				// TODO(rsc): Can address be 64 bits?
				Mem mem = {.Disp= immc};
				if (segIndex >= 0 ){
					mem.Segment = prefixToSegment(inst.Prefix[segIndex]);
					inst.Prefix[segIndex] |= PrefixImplicit;
				}
				inst.Args[narg] = mem;
				inst.MemBytes = memBytes[x];
				if (mem.Base == RIP) {
					inst.PCRel = displen;
					inst.PCRelOff = dispoff;
				}
				narg++;
				break;
			}
		case xArgYmm1:
			Reg base = baseReg[x];
			Reg index = (regop);
			if (inst.Prefix[vexIndex+1]&0x80 == 0 ){
				index += 8;
			}
			inst.Args[narg] = base + index;
			narg++;
			break;
		case xArgR8:
		case xArgR16:
		case xArgR32:
		case xArgR64:
		case xArgXmm:
		case xArgXmm1:
		case xArgDR0dashDR7:
			Reg base = baseReg[x];
			Reg index = regop;
			if (rex != 0 && base == AL && index >= 4 ){
				rexUsed |= PrefixREX;
				index -= 4;
				base = SPB;
			}
			inst.Args[narg] = base + index;
			narg++;
			break;
		case xArgMm:
		case xArgMm1:
		case xArgTR0dashTR7:
			inst.Args[narg] = baseReg[x] + (regop&7);
			narg++;
			break;
		case xArgCR0dashCR7:
			// AMD documents an extension that the LOCK prefix
			// can be used in place of a REX prefix in order to access
			// CR8 from 32-bit mode. The LOCK prefix is allowed in
			// all modes, provided the corresponding CPUID bit is set.
			if (lockIndex >= 0 ){
				inst.Prefix[lockIndex] |= PrefixImplicit;
				regop += 8;
			}
			inst.Args[narg] = CR0 + (regop);
			narg++;
			break;
		case xArgSreg:
			regop &= 7;
			if (regop >= 6) {
				inst.Op = 0;
				goto BREAK_DECODE;
			}
			inst.Args[narg] = ES + (regop);
			narg++;
			break;

		case xArgRmf16:
		case xArgRmf32:
		case xArgRmf64:
			Reg base = baseReg[x];
			Reg index = modrm & 07;
			if (rex&PrefixREXB != 0 ){
				rexUsed |= PrefixREXB;
				index += 8;
			}
			inst.Args[narg] = base + index;
			narg++;
			break;
		case xArgR8op:
		case xArgR16op:
		case xArgR32op:
		case xArgR64op:
		case xArgSTi:
			Reg n = inst.Opcode >> uint(opshift+8) & 07;
			Reg base = baseReg[x];
			Reg index = n;
			if (rex&PrefixREXB != 0 && x != xArgSTi ){
				rexUsed |= PrefixREXB;
				index += 8;
			}
			if (rex != 0 && base == AL && index >= 4 ){
				rexUsed |= PrefixREX;
				index -= 4;
				base = SPB;
			}
			inst.Args[narg] = base + index;
			narg++;
			break;
		case xArgRM8:
		case xArgRM16:
		case xArgRM32:
		case xArgRM64:
		case xArgR32M16:
		case xArgR32M8:
		case xArgR64M16:
		case xArgMmM32:
		case xArgMmM64:
		case xArgMm2M64:
		case xArgXmm2M16:
		case xArgXmm2M32:
		case xArgXmm2M64:
		case xArgXmmM64:
		case xArgXmmM128:
		case xArgXmmM32: 
		case  xArgXmm2M128:
	    case xArgYmm2M256:
			if (haveMem) {
				inst.Args[narg] = mem;
				inst.MemBytes = memBytes[(x)];
				if (mem.Base == RIP) {
					inst.PCRel = displen;
					inst.PCRelOff = dispoff;
				}
			} else {
				Reg base = baseReg[x];
				Reg index = rm;
				switch (x) {
				case xArgMmM32:
				case xArgMmM64:
				case xArgMm2M64:
					// There are only 8 MMX registers, so these ignore the REX.X bit.
					index &= 7;
					break;
				case xArgRM8:
					if (rex != 0 && index >= 4 ){
						rexUsed |= PrefixREX;
						index -= 4;
						base = SPB;
					}
					break;
				case xArgYmm2M256:
					if (vex == 0xC4 && inst.Prefix[vexIndex+1]&0x40 == 0x40 ){
						index += 8;
					}
					break;
				}
				inst.Args[narg] = base + index;
			}
			narg++;
			break;
		case xArgMm2: // register only; TODO(rsc): Handle with tag modrm_regonly tag
			if (haveMem) {
				inst.Op = 0;
				goto BREAK_DECODE;
			}
			inst.Args[narg] = baseReg[x] + (rm&7);
			narg++;
			break;
		case xArgXmm2: // register only; TODO(rsc): Handle with tag modrm_regonly tag
			if (haveMem) {
				inst.Op = 0;
				goto BREAK_DECODE;
			}
			inst.Args[narg] = baseReg[x] + (rm);
			narg++;
			break;
		case xArgRel8:
			inst.PCRelOff = immcpos;
			inst.PCRel = 1;
			inst.Args[narg] = (Rel)(immc);
			narg++;
			break;
		case xArgRel16:
			inst.PCRelOff = immcpos;
			inst.PCRel = 2;
			inst.Args[narg] = (Rel)(immc);
			narg++;
			break;
		case xArgRel32:
			inst.PCRelOff = immcpos;
			inst.PCRel = 4;
			inst.Args[narg] = (Rel)(immc);
			narg++;
			break;
		}

CONTINUE_DECODE:
		// continue;
	}

BREAK_DECODE:

	if (inst.Op == 0) {
		// Invalid instruction.
		if (nprefix > 0) {
			//todo invalid instruction
			*ld = instPrefix(src[0], mode); // invalid instruction
			return -1;
		}
		Inst inst = {.Len= pos};
		*ld = inst;
		//todo ErrUnrecognized
		return -1; //ErrUnrecognized
	}

	// Matched! Hooray!

	// 90 decodes as XCHG EAX, EAX but is NOP.
	// 66 90 decodes as XCHG AX, AX and is NOP too.
	// 48 90 decodes as XCHG RAX, RAX and is NOP too.
	// 43 90 decodes as XCHG R8D, EAX and is *not* NOP.
	// F3 90 decodes as REP XCHG EAX, EAX but is PAUSE.
	// It's all too special to handle in the decoding tables, at least for now.
	if (inst.Op == XCHG && inst.Opcode>>24 == 0x90 ){
		if ( inst.Args[0] == RAX || inst.Args[0] == EAX || inst.Args[0] == AX ){
			inst.Op = NOP;
			if ( dataSizeIndex >= 0) {
				// bit clear
				// inst.Prefix[dataSizeIndex] &^= PrefixImplicit;
				inst.Prefix[dataSizeIndex] &= ~(inst.Prefix[dataSizeIndex] & PrefixImplicit);
			}
			inst.Args[0] = -1;
			inst.Args[1] = -1;
		}
		if (repIndex >= 0 && inst.Prefix[repIndex] == 0xF3 ){
			inst.Prefix[repIndex] |= PrefixImplicit;
			inst.Op = PAUSE;
			inst.Args[0] = -1;
			inst.Args[1] = -1;
		} else if (gnuCompat) {
			int i;
			for (i = nprefix - 1; i >= 0; i-- ){
				if (inst.Prefix[i]&0xFF == 0xF3 ){
					inst.Prefix[i] |= PrefixImplicit;
					inst.Op = PAUSE;
					inst.Args[0] = -1;
					inst.Args[1] = -1;
					break;
				}
			}
		}
	}

	// defaultSeg returns the default segment for an implicit
	// memory reference: the final override if present, or else DS.
	defaultSeg := func() Reg {
		if segIndex >= 0 {
			inst.Prefix[segIndex] |= PrefixImplicit
			return prefixToSegment(inst.Prefix[segIndex])
		}
		return DS
	}

	// Add implicit arguments not present in the tables.
	// Normally we shy away from making implicit arguments explicit,
	// following the Intel manuals, but adding the arguments seems
	// the best way to express the effect of the segment override prefixes.
	// TODO(rsc): Perhaps add these to the tables and
	// create bytecode instructions for them.
	usedAddrSize := false
	switch inst.Op {
	case INSB, INSW, INSD:
		inst.Args[0] = Mem{Segment: ES, Base: baseRegForBits(addrMode) + DI - AX}
		inst.Args[1] = DX
		usedAddrSize = true

	case OUTSB, OUTSW, OUTSD:
		inst.Args[0] = DX
		inst.Args[1] = Mem{Segment: defaultSeg(), Base: baseRegForBits(addrMode) + SI - AX}
		usedAddrSize = true

	case MOVSB, MOVSW, MOVSD, MOVSQ:
		inst.Args[0] = Mem{Segment: ES, Base: baseRegForBits(addrMode) + DI - AX}
		inst.Args[1] = Mem{Segment: defaultSeg(), Base: baseRegForBits(addrMode) + SI - AX}
		usedAddrSize = true

	case CMPSB, CMPSW, CMPSD, CMPSQ:
		inst.Args[0] = Mem{Segment: defaultSeg(), Base: baseRegForBits(addrMode) + SI - AX}
		inst.Args[1] = Mem{Segment: ES, Base: baseRegForBits(addrMode) + DI - AX}
		usedAddrSize = true

	case LODSB, LODSW, LODSD, LODSQ:
		switch inst.Op {
		case LODSB:
			inst.Args[0] = AL
		case LODSW:
			inst.Args[0] = AX
		case LODSD:
			inst.Args[0] = EAX
		case LODSQ:
			inst.Args[0] = RAX
		}
		inst.Args[1] = Mem{Segment: defaultSeg(), Base: baseRegForBits(addrMode) + SI - AX}
		usedAddrSize = true

	case STOSB, STOSW, STOSD, STOSQ:
		inst.Args[0] = Mem{Segment: ES, Base: baseRegForBits(addrMode) + DI - AX}
		switch inst.Op {
		case STOSB:
			inst.Args[1] = AL
		case STOSW:
			inst.Args[1] = AX
		case STOSD:
			inst.Args[1] = EAX
		case STOSQ:
			inst.Args[1] = RAX
		}
		usedAddrSize = true

	case SCASB, SCASW, SCASD, SCASQ:
		inst.Args[1] = Mem{Segment: ES, Base: baseRegForBits(addrMode) + DI - AX}
		switch inst.Op {
		case SCASB:
			inst.Args[0] = AL
		case SCASW:
			inst.Args[0] = AX
		case SCASD:
			inst.Args[0] = EAX
		case SCASQ:
			inst.Args[0] = RAX
		}
		usedAddrSize = true

	case XLATB:
		inst.Args[0] = Mem{Segment: defaultSeg(), Base: baseRegForBits(addrMode) + BX - AX}
		usedAddrSize = true
	}

	// If we used the address size annotation to construct the
	// argument list, mark that prefix as implicit: it doesn't need
	// to be shown when printing the instruction.
	if haveMem || usedAddrSize {
		if addrSizeIndex >= 0 {
			inst.Prefix[addrSizeIndex] |= PrefixImplicit
		}
	}

	// Similarly, if there's some memory operand, the segment
	// will be shown there and doesn't need to be shown as an
	// explicit prefix.
	if haveMem {
		if segIndex >= 0 {
			inst.Prefix[segIndex] |= PrefixImplicit
		}
	}

	// Branch predict prefixes are overloaded segment prefixes,
	// since segment prefixes don't make sense on conditional jumps.
	// Rewrite final instance to prediction prefix.
	// The set of instructions to which the prefixes apply (other then the
	// Jcc conditional jumps) is not 100% clear from the manuals, but
	// the disassemblers seem to agree about the LOOP and JCXZ instructions,
	// so we'll follow along.
	// TODO(rsc): Perhaps this instruction class should be derived from the CSV.
	if isCondJmp[inst.Op] || isLoop[inst.Op] || inst.Op == JCXZ || inst.Op == JECXZ || inst.Op == JRCXZ {
	PredictLoop:
		for i := nprefix - 1; i >= 0; i-- {
			p := inst.Prefix[i]
			switch p & 0xFF {
			case PrefixCS:
				inst.Prefix[i] = PrefixPN
				break PredictLoop
			case PrefixDS:
				inst.Prefix[i] = PrefixPT
				break PredictLoop
			}
		}
	}

	// The BND prefix is part of the Intel Memory Protection Extensions (MPX).
	// A REPN applied to certain control transfers is a BND prefix to bound
	// the range of possible destinations. There's surprisingly little documentation
	// about this, so we just do what libopcodes and xed agree on.
	// In particular, it's unclear why a REPN applied to LOOP or JCXZ instructions
	// does not turn into a BND.
	// TODO(rsc): Perhaps this instruction class should be derived from the CSV.
	if isCondJmp[inst.Op] || inst.Op == JMP || inst.Op == CALL || inst.Op == RET {
		for i := nprefix - 1; i >= 0; i-- {
			p := inst.Prefix[i]
			if p&^PrefixIgnored == PrefixREPN {
				inst.Prefix[i] = PrefixBND
				break
			}
		}
	}

	// The LOCK prefix only applies to certain instructions, and then only
	// to instances of the instruction with a memory destination.
	// Other uses of LOCK are invalid and cause a processor exception,
	// in contrast to the "just ignore it" spirit applied to all other prefixes.
	// Mark invalid lock prefixes.
	hasLock := false
	if lockIndex >= 0 && inst.Prefix[lockIndex]&PrefixImplicit == 0 {
		switch inst.Op {
		// TODO(rsc): Perhaps this instruction class should be derived from the CSV.
		case ADD, ADC, AND, BTC, BTR, BTS, CMPXCHG, CMPXCHG8B, CMPXCHG16B, DEC, INC, NEG, NOT, OR, SBB, SUB, XOR, XADD, XCHG:
			if isMem(inst.Args[0]) {
				hasLock = true
				break
			}
			fallthrough
		default:
			inst.Prefix[lockIndex] |= PrefixInvalid
		}
	}

	// In certain cases, all of which require a memory destination,
	// the REPN and REP prefixes are interpreted as XACQUIRE and XRELEASE
	// from the Intel Transactional Synchroniation Extensions (TSX).
	//
	// The specific rules are:
	// (1) Any instruction with a valid LOCK prefix can have XACQUIRE or XRELEASE.
	// (2) Any XCHG, which always has an implicit LOCK, can have XACQUIRE or XRELEASE.
	// (3) Any 0x88-, 0x89-, 0xC6-, or 0xC7-opcode MOV can have XRELEASE.
	if isMem(inst.Args[0]) {
		if inst.Op == XCHG {
			hasLock = true
		}

		for i := len(inst.Prefix) - 1; i >= 0; i-- {
			p := inst.Prefix[i] &^ PrefixIgnored
			switch p {
			case PrefixREPN:
				if hasLock {
					inst.Prefix[i] = inst.Prefix[i]&PrefixIgnored | PrefixXACQUIRE
				}

			case PrefixREP:
				if hasLock {
					inst.Prefix[i] = inst.Prefix[i]&PrefixIgnored | PrefixXRELEASE
				}

				if inst.Op == MOV {
					op := (inst.Opcode >> 24) &^ 1
					if op == 0x88 || op == 0xC6 {
						inst.Prefix[i] = inst.Prefix[i]&PrefixIgnored | PrefixXRELEASE
					}
				}
			}
		}
	}

	// If REP is used on a non-REP-able instruction, mark the prefix as ignored.
	if repIndex >= 0 {
		switch inst.Prefix[repIndex] {
		case PrefixREP, PrefixREPN:
			switch inst.Op {
			// According to the manuals, the REP/REPE prefix applies to all of these,
			// while the REPN applies only to some of them. However, both libopcodes
			// and xed show both prefixes explicitly for all instructions, so we do the same.
			// TODO(rsc): Perhaps this instruction class should be derived from the CSV.
			case INSB, INSW, INSD,
				MOVSB, MOVSW, MOVSD, MOVSQ,
				OUTSB, OUTSW, OUTSD,
				LODSB, LODSW, LODSD, LODSQ,
				CMPSB, CMPSW, CMPSD, CMPSQ,
				SCASB, SCASW, SCASD, SCASQ,
				STOSB, STOSW, STOSD, STOSQ:
				// ok
			default:
				inst.Prefix[repIndex] |= PrefixIgnored
			}
		}
	}

	// If REX was present, mark implicit if all the 1 bits were consumed.
	if rexIndex >= 0 {
		if rexUsed != 0 {
			rexUsed |= PrefixREX
		}
		if rex&^rexUsed == 0 {
			inst.Prefix[rexIndex] |= PrefixImplicit
		}
	}

	inst.DataSize = dataMode
	inst.AddrSize = addrMode
	inst.Mode = mode
	inst.Len = pos
	return inst, nil
}

}