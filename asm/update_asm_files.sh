#!/bin/bash
rm -f Args.c Args.h Inst.h Inst.c table.h table.c goX86asm.h goX86asm.c

ln -s ../aop/Args.h Args.h
ln -s ../aop/Args.c Args.c
ln -s ../aop/Inst.c Inst.c
ln -s ../aop/Inst.h Inst.h
ln -s ../aop/table.h table.h
ln -s ../aop/table.c table.c
ln -s ../aop/goX86asm.h goX86asm.h
ln -s ../aop/goX86asm.c goX86asm.c

