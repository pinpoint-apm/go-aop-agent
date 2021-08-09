package aop

//#cgo CFLAGS: -DNTRACE
//#include "goX86asm.h"
import "C"
import (
	"errors"
	"unsafe"

	"golang.org/x/arch/x86/x86asm"
)

func Decode(code []byte, mode int) (goInst x86asm.Inst, err error) {
	inst := C.Inst{}
	codeLen := int32(len(code))
	ret := C.decode((*C.uint8_t)(unsafe.Pointer(&code[0])), C.int(codeLen), &inst, C.int(mode), C.uchar(0))
	goInst = x86asm.Inst{
		Op:       x86asm.Op(inst.Op),
		Opcode:   uint32(inst.Opcode),
		Mode:     int(inst.Mode),
		AddrSize: int(inst.AddrSize),
		DataSize: int(inst.DataSize),
		MemBytes: int(inst.MemBytes),
		Len:      int(inst.Len),
		PCRel:    int(inst.PCRel),
		PCRelOff: int(inst.PCRelOff),
	}
	// transfer c_inst into goInst
	for i := 0; i < 14; i++ {
		goInst.Prefix[i] = x86asm.Prefix(inst.Prefix[i])
	}

	for i := 0; i < 4; i++ {
		switch inst.Args[i]._type {
		case C.E_MEM:
			var c_mem *C.Mem
			c_mem = ((*C.Mem)(unsafe.Pointer(&inst.Args[i].value)))
			// fmt.Println(c_mem)
			mem := x86asm.Mem{}
			mem.Base = x86asm.Reg(c_mem.Base)
			mem.Segment = x86asm.Reg(c_mem.Segment)
			mem.Index = x86asm.Reg(c_mem.Index)
			mem.Scale = uint8(c_mem.Scale)
			mem.Disp = int64(c_mem.Disp)
			goInst.Args[i] = mem
		case C.E_NIL:
			goInst.Args[i] = nil
		case C.E_REG:
			c_reg := ((*C.Reg)(unsafe.Pointer(&inst.Args[i].value)))
			var reg x86asm.Reg = x86asm.Reg(*c_reg)
			goInst.Args[i] = reg
		case C.E_IMM:
			c_imm := ((*C.Imm)(unsafe.Pointer(&inst.Args[i].value)))
			var imm x86asm.Imm = x86asm.Imm(*c_imm)
			goInst.Args[i] = imm
		case C.E_REL:
			c_rel := ((*C.Rel)(unsafe.Pointer(&inst.Args[i].value)))
			var rel x86asm.Rel = x86asm.Rel(*c_rel)
			goInst.Args[i] = rel
		}
	}
	//  transfer ret into err
	if ret != C.E_OK {
		switch ret {
		case C.E_UNRECOGNIZED:
			err = errors.New("unrecognized instruction")
		case C.E_TRUNCATED:
			err = errors.New("truncated instruction")
		case C.E_INVALID_MODE:
			err = errors.New("invalid x86 mode in Decode")
		case C.E_INTERNAL:
			err = errors.New("internal error")
		}
		return
	}
	return
}
