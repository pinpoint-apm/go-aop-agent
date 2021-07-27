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
package app

import (
	"context"

	"github.com/pinpoint-apm/go-aop-agent/common"
)

func init() {
	hook_common_func(TestUserFunc, hook_testuserfunc, hook_testuserfunc_trampoline)
	hook_common_func(TestComFunc, hook_testcomfunc, hook_testcomfunc_trampoline)
	hook_common_func(TestHighFunc, hook_testhighfunc, hook_testhighfunc_trampoline)
	hook_common_func(Add, hook_testaddfunc, hook_testaddfunc_trampoline)
	hook_common_func(Mul, hook_testmulfunc, hook_testmulfunc_trampoline)
	hook_common_func(Person.TestInheritFunc, hook_testinheritfunc, hook_testinheritfunc_trampoline)
	hook_common_func(TestLambdaFunc, hook_testlambdafunc, hook_testlambdafunc_trampoline)
	hook_common_func(Dinner.TestDecoratorFunc, hook_testdecoratorfunc, hook_testdecoratorfunc_trampoline)
	hook_common_func(Rice.WashRice, hook_testwashricefunc, hook_testwashricefunc_trampoline)
	hook_common_func(Water.AddWater, hook_testaddwaterfunc, hook_testaddwaterfunc_trampoline)
	hook_common_func(Book.TestAbstractFunc, hook_testabstractfunc, hook_testabstractfunc_trampoline)
	hook_common_func(TestGeneratorFunc, hook_testgeneratorfunc, hook_testgeneratorfunc_trampoline)
	hook_common_func(TestReturn, hook_testreturn, hook_testreturn_trampoline)
	hook_common_func(TestException, hook_exception, hook_exception_trampoline)
	hook_common_func(TestExpInRecursion, hook_expinrecursion, hook_expinrecursion_trampoline)
	hook_common_func(TestRecursion, hook_recursion, hook_recursion_trampoline)
}

////////////////////////app.TestUserFunc//////////////////////
//go:noinline
func hook_testuserfunc_trampoline(ctx context.Context, num1 int, num2 int) int{
	return 0
}

//go:noinline
func hook_testuserfunc(ctx context.Context, num1 int, num2 int) int{
	funcName :=get_func_name(TestUserFunc)
	if parentId, err :=common.GetParentId(ctx); err !=nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_testuserfunc_trampoline(ctx, num1, num2)
	}else {
		subTraceId := common.Pinpoint_start_trace(parentId) 
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName, num1, num2)
		response := hook_testuserfunc_trampoline(newCtx, num1, num2)
		onEnd(subTraceId, response)
		return response
		}
}

////////////////////////app.TestComFunc//////////////////////
//go:noinline
func hook_testcomfunc_trampoline(ctx context.Context, arg interface{}) interface{} {
	return nil
}

//go:noinline
func hook_testcomfunc(ctx context.Context, arg interface{}) interface{} {
	funcName := get_func_name(TestComFunc)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_testcomfunc_trampoline(ctx, arg)
	} else {
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName, arg)
		response := hook_testcomfunc_trampoline(newCtx, arg)
		onEnd(subTraceId, response)
		return response
	}
}

////////////////////////app.TestReturn//////////////////////
//go:noinline
func hook_testreturn_trampoline(ctx context.Context) func(int) {
	return nil
}

//go:noinline
func hook_testreturn(ctx context.Context) func(int) {
	funcName := get_func_name(TestReturn)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_testreturn_trampoline(ctx)
	} else {
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName)
		response := hook_testreturn_trampoline(newCtx)
		onEnd(subTraceId, response)
		return response
	}
}

////////////////////////app.TestHighFunc//////////////////////
//go:noinline
func hook_testhighfunc_trampoline(ctx context.Context, a, b int, T func(context.Context, int, int) int) int {
	return 0
}

//go:noinline
func hook_testhighfunc(ctx context.Context, a, b int, T func(context.Context, int, int) int) int {
	funcName := get_func_name(TestHighFunc)
	if parentId, err := common.GetParentId(ctx); err != nil {
			common.Logf("parentId is not traceId type. Dropped")
			return hook_testhighfunc_trampoline(ctx, a, b, T)
	} else {
			subTraceId := common.Pinpoint_start_trace(parentId)
			defer common.Pinpoint_end_trace(subTraceId)

			newCtx := onBefore(ctx, subTraceId, funcName, a, b, T)
			response := hook_testhighfunc_trampoline(newCtx, a, b, T)
			onEnd(subTraceId, response)
			return response
	}
}

////////////////////////app.TestAddFunc//////////////////////
//go:noinline
func hook_testaddfunc_trampoline(ctx context.Context, a, b int ) int {
	return 0
}

//go:noinline
func hook_testaddfunc(ctx context.Context, a, b int) int {
	funcName := get_func_name(Add)
	if parentId, err := common.GetParentId(ctx); err != nil {
			common.Logf("parentId is not traceId type. Dropped")
			return hook_testaddfunc_trampoline(ctx, a, b)
	} else {
			subTraceId := common.Pinpoint_start_trace(parentId)
			defer common.Pinpoint_end_trace(subTraceId)

			newCtx := onBefore(ctx, subTraceId, funcName, a, b)
			response := hook_testaddfunc_trampoline(newCtx, a, b)
			onEnd(subTraceId, response)
			return response
	}
}

////////////////////////app.TestMulFunc//////////////////////
//go:noinline
func hook_testmulfunc_trampoline(ctx context.Context, a, b int ) int {
	return 0
}

//go:noinline
func hook_testmulfunc(ctx context.Context, a, b int) int {
	funcName := get_func_name(Mul)
	if parentId, err := common.GetParentId(ctx); err != nil {
			common.Logf("parentId is not traceId type. Dropped")
			return hook_testmulfunc_trampoline(ctx, a, b)
	} else {
			subTraceId := common.Pinpoint_start_trace(parentId)
			defer common.Pinpoint_end_trace(subTraceId)

			newCtx := onBefore(ctx, subTraceId, funcName, a, b)
			response := hook_testmulfunc_trampoline(newCtx, a, b)
			onEnd(subTraceId, response)
			return response
	}
}

////////////////////////app.TestInheritFunc//////////////////////
//go:noinline
func hook_testinheritfunc_trampoline(p Person, ctx context.Context) (string, string, int){
    return "Hello", "World", 0
}

//go:noinline
func hook_testinheritfunc(p Person, ctx context.Context) (string, string, int){
    funcName :=get_func_name(Person.TestInheritFunc)
    if parentId, err :=common.GetParentId(ctx); err !=nil {
	         common.Logf("parentId is not traceId type. Dropped")
	         return hook_testinheritfunc_trampoline(p, ctx)
    } else {
	        subTraceId := common.Pinpoint_start_trace(parentId) 
	        defer common.Pinpoint_end_trace(subTraceId)

	        newCtx := onBefore(ctx, subTraceId, funcName, p)
	        response, response0, response1:= hook_testinheritfunc_trampoline(p, newCtx)
	        onEnd(subTraceId, response, response0, response1)
	        return response, response0, response1
	}
}

////////////////////////app.TestLambdaFunc//////////////////////
//go:noinline
func hook_testlambdafunc_trampoline(ctx context.Context) int {
    return 0
}

//go:noinline
func hook_testlambdafunc(ctx context.Context) int{
    funcName := get_func_name(TestLambdaFunc)
    if parentId, err := common.GetParentId(ctx); err != nil {
		     common.Logf("parentId is not traceId type. Dropped")
		     return hook_testlambdafunc_trampoline(ctx)
    } else {
		     subTraceId := common.Pinpoint_start_trace(parentId)
		     defer common.Pinpoint_end_trace(subTraceId)

		     newCtx := onBefore(ctx, subTraceId, funcName)
		     response := hook_testlambdafunc_trampoline(newCtx)
		     onEnd(subTraceId, response)
		     return response
    }
}

////////////////////////app.TestDecoratorFunc//////////////////////
//go:noinline
func hook_testdecoratorfunc_trampoline(D Dinner, ctx context.Context) (string, string, string, string, string) {
    return "A","B","C","D","E"
}

//go:noinline
func hook_testdecoratorfunc(D Dinner, ctx context.Context) (string, string, string, string, string) {
    funcName := get_func_name(Dinner.TestDecoratorFunc)
    if parentId, err := common.GetParentId(ctx); err != nil {
		     common.Logf("parentId is not traceId type. Dropped")
		     return hook_testdecoratorfunc_trampoline(D, ctx)
    } else {
		     subTraceId := common.Pinpoint_start_trace(parentId)
		     defer common.Pinpoint_end_trace(subTraceId)

		     newCtx := onBefore(ctx, subTraceId, funcName)
		     response, response0, response1,response2,response3 := hook_testdecoratorfunc_trampoline(D, newCtx)
		     onEnd(subTraceId, response, response1,response2,response3)
		     return response, response0, response1,response2,response3
    }
}

////////////////////////app.TestWashFunc//////////////////////
//go:noinline
func hook_testwashricefunc_trampoline(R Rice, ctx context.Context) string{
	return ""
}

//go:noinline
func hook_testwashricefunc(R Rice, ctx context.Context) string {
	funcName := get_func_name(Rice.WashRice)
	if parentId, err := common.GetParentId(ctx); err != nil {
			common.Logf("parentId is not traceId type. Dropped")
			return hook_testwashricefunc_trampoline(R, ctx)
	} else {
			subTraceId := common.Pinpoint_start_trace(parentId)
			defer common.Pinpoint_end_trace(subTraceId)

			newCtx := onBefore(ctx, subTraceId, funcName, R)
			response := hook_testwashricefunc_trampoline(R, newCtx)
			onEnd(subTraceId, response)
			return response
	}
}

////////////////////////app.TestWaterFunc//////////////////////
//go:noinline
func hook_testaddwaterfunc_trampoline(W Water, ctx context.Context) string{
	return ""
}

//go:noinline
func hook_testaddwaterfunc(W Water, ctx context.Context) string {
	funcName := get_func_name(Water.AddWater)
	if parentId, err := common.GetParentId(ctx); err != nil {
			common.Logf("parentId is not traceId type. Dropped")
			return hook_testaddwaterfunc_trampoline(W, ctx)
	} else {
			subTraceId := common.Pinpoint_start_trace(parentId)
			defer common.Pinpoint_end_trace(subTraceId)

			newCtx := onBefore(ctx, subTraceId, funcName, W)
			response := hook_testaddwaterfunc_trampoline(W, newCtx)
			onEnd(subTraceId, response)
			return response
	}
}

////////////////////////app.TestInheritFunc//////////////////////
//go:noinline
func hook_testabstractfunc_trampoline(b Book, ctx context.Context) string{
    return "the book color is RED"
}

//go:noinline
func hook_testabstractfunc(b Book, ctx context.Context) string{
    funcName :=get_func_name(Book.TestAbstractFunc)
    if parentId, err :=common.GetParentId(ctx); err !=nil {
	         common.Logf("parentId is not traceId type. Dropped")
	         return hook_testabstractfunc_trampoline(b, ctx)
    } else {
	         subTraceId := common.Pinpoint_start_trace(parentId) 
	         defer common.Pinpoint_end_trace(subTraceId)

	         newCtx := onBefore(ctx, subTraceId, funcName, b)
	         response := hook_testabstractfunc_trampoline(b, newCtx)
	         onEnd(subTraceId, response)
	         return response
	}
}

////////////////////////app.TestGeneratorFunc//////////////////////
//go:noinline
func hook_testgeneratorfunc_trampoline(ctx context.Context, c chan int) (int, int, int, int ){
    return 1,2,3,4
}

//go:noinline
func hook_testgeneratorfunc(ctx context.Context, c chan int) (int, int, int, int){
    funcName := get_func_name(TestGeneratorFunc)
    if parentId, err := common.GetParentId(ctx); err != nil {
		     common.Logf("parentId is not traceId type. Dropped")
		     return hook_testgeneratorfunc_trampoline(ctx, c)
    } else {
		     subTraceId := common.Pinpoint_start_trace(parentId)
		     defer common.Pinpoint_end_trace(subTraceId)

		     newCtx := onBefore(ctx, subTraceId, funcName, c)
		     response, response0, response1, response2:= hook_testgeneratorfunc_trampoline(newCtx, c)
		     onEnd(subTraceId, response, response0, response1, response2)
		    return response, response0, response1, response2
    }
}

////////////////////////app.TestException//////////////////////
//go:noinline
func hook_exception_trampoline(ctx context.Context, varDividee int, varDivider int) (result int, err error) {
	return 0, nil
}

//go:noinline
func hook_exception(ctx context.Context, varDividee int, varDivider int) (result int, err error) {
	funcName := get_func_name(TestException)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_exception_trampoline(ctx, varDividee, varDivider)
	} else {
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName, varDividee, varDivider)
		response, err := hook_exception_trampoline(newCtx, varDividee, varDivider)
		if err != nil {
			onException(subTraceId, &err)
		}
		onEnd(subTraceId, response)
		return response, err
	}
}

////////////////////////app.TestExpInRecursion//////////////////////
//go:noinline
func hook_expinrecursion_trampoline(ctx context.Context, i int) (int, error) {
	return 0, nil
}

//go:noinline
func hook_expinrecursion(ctx context.Context, i int) (int, error) {
	funcName := get_func_name(TestExpInRecursion)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_expinrecursion_trampoline(ctx, i)
	} else {
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName, i)
		response, err := hook_expinrecursion_trampoline(newCtx, i)
		if err != nil {
			onException(subTraceId, &err)
		}
		onEnd(subTraceId, response)
		return response, err
	}
}

////////////////////////app.TestRecursion//////////////////////
//go:noinline
func hook_recursion_trampoline(ctx context.Context, i int) int {
	return 0
}

//go:noinline
func hook_recursion(ctx context.Context, i int) int {
	funcName := get_func_name(TestRecursion)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return hook_recursion_trampoline(ctx, i)
	} else {
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(ctx, subTraceId, funcName, i)
		response := hook_recursion_trampoline(newCtx, i)
		onEnd(subTraceId, response)
		return response
	}
}
