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
	"fmt"
	"reflect"
	"runtime"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func onBefore(ctx context.Context, id common.TraceIdType, funcName string, arg ...interface{}) context.Context {
	// common.Logf("call onBefore")

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, funcName)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_METHOD_CALL)
	addClueSFunc(common.PP_ARGS, fmt.Sprint(arg))

	ctx = context.WithValue(ctx, common.TRACE_ID, id)
	return ctx
}

func onEnd(id common.TraceIdType, res ...interface{}) {
	// common.Logf("call onEnd")

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc(common.PP_RETURN, fmt.Sprint(res))

	// common.Pinpoint_end_trace(id)
}

func onException(id common.TraceIdType, err *error) {
	// common.Logf("call onException")
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	addClueFunc(common.PP_ADD_EXCEPTION, fmt.Sprint(*err))
	common.Pinpoint_mark_error(fmt.Sprint(*err), "", 0, id)
}

func hook_common_func(f interface{}, hook_f interface{}, hook_f_trampoline interface{}) {
	// funcName := get_func_name(f)
	// common.Logf("try to hook " + funcName)
	if err := aop.AddHook(f, hook_f, hook_f_trampoline); err != nil {
		// common.Logf("Hook "+funcName+" failed:%s", err)
		return
	}
	// common.Logf(funcName + " is hooked")
}

func get_func_name(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
