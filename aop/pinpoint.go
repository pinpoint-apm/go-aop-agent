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

package aop

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/pinpoint-apm/go-aop-agent/common"
)

// #cgo CFLAGS: -DNTEST -DTRACE
// #include "pinpoint.h"
import "C"

var trampolineMap map[uintptr]unsafe.Pointer

func init() {
	trampolineMap = make(map[uintptr]unsafe.Pointer)
}

func AddHookP_CALL(iSrc, iTarget, iTrampoline_func interface{}) error {

	if common.AgentIsDisabled() {
		return errors.New("agent disabled")
	}

	src := reflect.ValueOf(iSrc)
	target := reflect.ValueOf(iTarget)
	trampoline_func := reflect.ValueOf(iTrampoline_func)

	// bug: check the signature of all functions
	if src.Kind() != reflect.Func {
		return errors.New("src is not function")
	} else if target.Kind() != reflect.Func {
		return errors.New("target is not function")
	} else if trampoline_func.Kind() != reflect.Func {
		return errors.New("trampoline_func is not function")
	}

	if src.Type() != target.Type() || target.Type() != trampoline_func.Type() {
		return errors.New("src, target,trampoline_func signature must be the same")
	}

	newSrcPointer := C.located_nearest_call_target(unsafe.Pointer(src.Pointer()))
	if newSrcPointer == nil {
		return errors.New("located nearest target failed")
	}

	if _, ok := trampolineMap[uintptr(newSrcPointer)]; ok {
		return errors.New("src exist")
	}

	trampoline := C.hook(newSrcPointer, unsafe.Pointer(target.Pointer()), unsafe.Pointer(trampoline_func.Pointer()))
	if trampoline == nil {
		return errors.New("hook failed, check the output under \"-DDEBUG\" ")
	} else {
		// store into trampoline map
		trampolineMap[uintptr(newSrcPointer)] = trampoline
		return nil
	}
}

func AddHookP_JMP(iSrc, iTarget, iTrampoline_func interface{}) error {
	if common.AgentIsDisabled() {
		return errors.New("agent disabled")
	}

	src := reflect.ValueOf(iSrc)
	target := reflect.ValueOf(iTarget)
	trampoline_func := reflect.ValueOf(iTrampoline_func)
	// skip Kind checking

	newSrcPointer := C.located_nearest_jmp_target(unsafe.Pointer(src.Pointer()))
	if newSrcPointer == nil {
		return errors.New("located nearest target failed")
	}

	if _, ok := trampolineMap[uintptr(newSrcPointer)]; ok {
		return errors.New("src exist")
	}

	trampoline := C.hook(newSrcPointer, unsafe.Pointer(target.Pointer()), unsafe.Pointer(trampoline_func.Pointer()))
	if trampoline == nil {
		return errors.New("hook failed, check the output under \"-DDEBUG\" ")
	} else {
		// store into trampoline map
		trampolineMap[uintptr(newSrcPointer)] = trampoline
		return nil
	}
}

/**
 * @description: Bind `iTarget` and `iTrampoline_func` on `iSrc`
 * 1. iSrc,iTarget,and iTrampoline_func  must be addressable.
 * 2. The signature of them must be the same.
 * @param {*} iSrc
 * @param {*} iTarget
 * @param {interface{}} iTrampoline_func
 * @return {*}
 */
func AddHook(iSrc, iTarget, iTrampoline_func interface{}) error {

	if common.AgentIsDisabled() {
		return errors.New("agent disabled")
	}

	src := reflect.ValueOf(iSrc)
	target := reflect.ValueOf(iTarget)
	trampoline_func := reflect.ValueOf(iTrampoline_func)

	if src.Kind() != reflect.Func {
		return errors.New("src is not function")
	} else if target.Kind() != reflect.Func {
		return errors.New("target is not function")
	} else if trampoline_func.Kind() != reflect.Func {
		return errors.New("trampoline_func is not function")
	}

	if _, ok := trampolineMap[src.Pointer()]; ok {
		return errors.New("src exist")
	}

	if src.Type() == target.Type() && target.Type() == trampoline_func.Type() {

		trampoline := C.hook(unsafe.Pointer(src.Pointer()), unsafe.Pointer(target.Pointer()), unsafe.Pointer(trampoline_func.Pointer()))
		if trampoline == nil {
			return errors.New("hook failed, check the output under \"-DDEBUG\" ")
		} else {
			// store into trampoline map
			trampolineMap[src.Pointer()] = trampoline
			return nil
		}
	} else {
		return errors.New("src, target,trampoline_func signature must be the same")
	}
}

/**
 * @description: remove the hook on `iSrc`
 * 1. DO NOT CALL this function when function is called by other goroutines
 * 2. The memory will not free until the program exit
 * @param {interface{}} iSrc
 * @return {*}
 */
func UnHook(iSrc interface{}) {
	src := reflect.ValueOf(iSrc)
	if trampoline, ok := trampolineMap[src.Pointer()]; ok {
		C.unhook(trampoline)
		delete(trampolineMap, src.Pointer())
	}
}
