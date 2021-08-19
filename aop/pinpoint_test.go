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
	"runtime"
	"testing"
)

//go:noinline
func foo(a, b int32) int32 {
	return a + b
}

//go:noinline
func foo_tramp(a, b int32) int32 {
	return a + b
}

//go:noinline
func hook_foo(a, b int32) int32 {
	before := func(a, b *int32) {
		// fmt.Print("call before \n")
		*a += *a
		*b += *b
	}

	before(&a, &b)

	ret := foo_tramp(a, b)

	end := func(ret *int32) {
		*ret *= 2
		// fmt.Printf("update ret to %d \n", *ret)
	}
	end(&ret)
	return ret
}

// as same function as foo
//go:noinline
func foo_test(a, b int32) int32 {
	return a + b
}

func BenchmarkFoo(t *testing.B) {
	for i := 0; i < t.N; i++ {
		if foo_test(1, 2) != 3 {
			t.Fail()
		}
	}
}

func BenchmarkHookFoo(t *testing.B) {
	AddHook(foo, hook_foo, foo_tramp)
	for i := 0; i < t.N; i++ {
		ret := foo(1, 2)
		if ret != 12 {
			t.Fail()
		}
	}

}

//go:noinline
func hook_foo_test(a, b int32) int32 {

	before := func(a, b *int32) {
		// fmt.Print("call before \n")
		*a += *a
		*b += *b
	}

	before(&a, &b)

	ret := foo_test(a, b)

	end := func(ret *int32) {
		*ret *= 2
		// fmt.Printf("update ret to %d \n", *ret)
	}
	end(&ret)
	return ret
}

func BenchmarkCallFooWithoutAop(t *testing.B) {

	for i := 0; i < t.N; i++ {
		ret := hook_foo_test(1, 2)
		if ret != 12 {
			t.Fail()
		}
	}

}

func TestHook(t *testing.T) {
	AddHook(foo, hook_foo, foo_tramp)
	if foo_test(1, 2) != 3 {
		t.Log("oh oh")
		t.Fail()
	}

	if foo(1, 2) != 12 {
		t.Log("oh no")
		t.Fail()
	}

}

/////////////////////////////////////////////////////////////////////////
// test struct

type Foo struct {
	vi32 int32
	// vi64    int64
	// vint    int
	vString string
	// vbool   bool
}

//go:noinline
func (f *Foo) output(v int32) int32 {
	f.vi32 = v
	return f.vi32
}

//go:noinline
func (f *Foo) outputVString(v string) string {
	f.vString = v
	return f.vString
}

//go:noinline
func Foo_output_tramp(f *Foo, v int32) int32 {
	return 0
}

//go:noinline
func hook_Foo_output(f *Foo, v int32) int32 {
	before := func(a *int32) {
		// fmt.Print("call before \n")
		*a += *a
	}

	before(&v)

	ret := Foo_output_tramp(f, v)

	end := func(ret *int32) {
		*ret *= 2
	}
	end(&ret)
	return ret
}

func TestFooOutputVi32(t *testing.T) {
	f := Foo{vi32: 10}
	if f.output(11) != 11 {
		t.Fail()
	}

	err := AddHook((*Foo).output, hook_Foo_output, Foo_output_tramp)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if f.output(11) != 44 {
		t.Log(f.output(11))
		t.Fail()
	}
	UnHook((*Foo).output)
}

//go:noinline
func Foo_outputVString_tramp(f *Foo, v string) string {
	return ""
}

//go:noinline
func hook_Foo_output_VString(f *Foo, v string) string {
	before := func(a *string) {
		// fmt.Print("call before \n")
		*a = "hello pinpoint"
	}

	before(&v)

	ret := Foo_outputVString_tramp(f, v)

	end := func(ret *string) {
		*ret += "!!!"
	}
	end(&ret)
	return ret
}

func TestFooOutputVString(t *testing.T) {
	f := Foo{vString: "hello world"}
	if f.outputVString("it's a sun day") != "it's a sun day" {
		t.Fail()
	}

	err := AddHook((*Foo).outputVString, hook_Foo_output_VString, Foo_outputVString_tramp)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if f.outputVString("it's a sun day") != "hello pinpoint!!!" {
		t.Log(f.outputVString("it's a sun day"))
		t.Fail()
	}

	UnHook((*Foo).output)
}

type foo_inter interface {
	Foo(int) int
}

type foo_s struct {
}

//
func (f *foo_s) Foo(int) int {
	return 0
}

type KV interface {
	Get(key string) string
}

type KVImplement struct {
}

//go:noinline
func (i *KVImplement) Get(key string) string {
	return "notfound"
}

type KVImplementWrapper struct {
	KVImplement
	a int
}

type Client struct {
	Kv KV
}

//go:noinline
func newKvImplement() KV {
	return &KVImplement{}
}

//go:noinline
func newKvImplementHook() KV {
	// kv := newKvImplementTrampoline()
	return &KVImplementWrapper{a: 10}
}

//go:noinline
func newKvImplementTrampoline() KV {
	return nil
}

func TestWrapperInterface(t *testing.T) {
	AddHook(newKvImplement, newKvImplementHook, newKvImplementTrampoline)

	client := Client{Kv: newKvImplement()}
	if client.Kv.Get("COVID-19") != "notfound" {
		t.Log(client.Kv.Get("COVID-19"))
		t.Fail()
	}
	runtime.GC()
}

type Base struct {
}

const baseString = "I'm Base"
const baseStringTheMore = baseString + " the more"
const baseStringFooTheMore = baseString + "Foo"

//go:noinline
func (base *Base) DoSomeThing() string {
	return baseString
}

type FooBase struct {
	Base
}

//go:noinline
func hook_DoSomeThing(base *Base) string {
	return hook_DoSomeThing_trampoline(base) + " the more"
}

//go:noinline
func hook_DoSomeThing_trampoline(base *Base) string {
	return ""
}

//go:noinline
func hook_DoSomeThing_FooBase(foo *FooBase) string {
	return hook_DoSomeThing_trampoline_FooBase(foo) + "Foo"
}

//go:noinline
func hook_DoSomeThing_trampoline_FooBase(foo *FooBase) string {
	return ""
}

func TestInherit(t *testing.T) {
	// // fmt.Println("AddHook on Base")
	baseError := AddHook((*Base).DoSomeThing, hook_DoSomeThing, hook_DoSomeThing_trampoline)
	if baseError != nil {
		t.Log("FooBase add hook failed")
		t.Fail()
	}
	foo := FooBase{}
	varl := foo.DoSomeThing()
	if varl != baseStringTheMore {
		t.Log(varl)
		t.Fail()
	}
	UnHook((*Base).DoSomeThing)

	fooError := AddHookP_JMP((*FooBase).DoSomeThing, hook_DoSomeThing_FooBase, hook_DoSomeThing_trampoline_FooBase)
	if fooError != nil {
		t.Fatalf("AddHookP_JMP failed:%s", fooError.Error())
	}

	foo = FooBase{}
	varl = foo.DoSomeThing()
	if varl != baseStringFooTheMore {
		t.Log(varl)
		t.Fail()
	}

}

func TestInvalid(t *testing.T) {
	baseError := AddHook((*Base).DoSomeThing, hook_DoSomeThing_FooBase, hook_DoSomeThing_FooBase)
	if baseError == nil {
		t.Log("AddHook checking failed")
		t.Fail()
	}
	t.Log(baseError.Error())

	baseError = AddHook((*Base).DoSomeThing, (*Base).DoSomeThing, hook_DoSomeThing_FooBase)
	if baseError == nil {
		t.Log("AddHook checking failed")
		t.Fail()
	}
	t.Log(baseError.Error())
	baseError = AddHook((*Base).DoSomeThing, (*Base).DoSomeThing, TestInherit)
	if baseError == nil {
		t.Log("AddHook checking failed")
		t.Fail()
	}
	t.Log(baseError.Error())
}

//go:noinline
func raw() int {
	return 1024
}

//go:noinline
func callRaw() int {
	return raw() * 3
}

//go:noinline
func hookRawTrampoline() int {
	return 1024
}

//go:noinline
func hookRaw() int {
	return hookRawTrampoline() + 2
}

func TestAddHookP_CALL(t *testing.T) {

	baseError := AddHookP_CALL(callRaw, hookRaw, hookRawTrampoline)
	if baseError != nil {
		t.Log("AddHook checking failed")
		t.Fail()
	}

	v := callRaw()
	if v != 3078 { // (1024+2) * 3
		t.Logf("callRaw not working :%d", v)
		t.Fail()
	}

}
