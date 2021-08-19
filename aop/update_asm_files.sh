#!/bin/bash
rm -f Args.c Args.h Inst.h Inst.c table.h table.c goX86asm.h goX86asm.c

ln -s ../asm/Args.h Args.h
ln -s ../asm/Args.c Args.c
ln -s ../asm/Inst.c Inst.c
ln -s ../asm/Inst.h Inst.h
ln -s ../asm/table.h table.h
ln -s ../asm/table.c table.c
ln -s ../asm/goX86asm.h goX86asm.h
ln -s ../asm/goX86asm.c goX86asm.c

