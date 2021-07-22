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
	hook_common_func(TestComFunc, hook_testcomfunc, hook_testcomfunc_trampoline)
	hook_common_func(TestReturn, hook_testreturn, hook_testreturn_trampoline)
	hook_common_func(TestException, hook_exception, hook_exception_trampoline)
	hook_common_func(TestExpInRecursion, hook_expinrecursion, hook_expinrecursion_trampoline)
	hook_common_func(TestRecursion, hook_recursion, hook_recursion_trampoline)
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
