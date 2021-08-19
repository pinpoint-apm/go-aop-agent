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

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"runtime"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func onBefore(id common.TraceIdType, funcName string, db *sql.DB, ctx context.Context, query string, args ...interface{}) *context.Context {
	common.Logf("call onBefore")

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, funcName)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_MYSQL)
	addClueFunc(common.PP_DESTINATION, gethost(DBMap[db].dataSourceName))
	addClueFunc(common.PP_SQL_FORMAT, query)

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return &newCtx
}

func pingonBefore(id common.TraceIdType, funcName string, db *sql.DB, ctx context.Context) *context.Context {
	common.Logf("call onBefore")

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, funcName)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_MYSQL)
	addClueFunc(common.PP_DESTINATION, gethost(DBMap[db].dataSourceName))

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return &newCtx
}

func queryonEnd(id common.TraceIdType, res *sql.Rows) {

}

func execonEnd(id common.TraceIdType, res sql.Result) {

}

func pingonEnd(id common.TraceIdType, err *error) {

}

func onException(id common.TraceIdType, err *error) {
	common.Logf("call onException")
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	addClueFunc(common.PP_ADD_EXCEPTION, fmt.Sprint(*err))
	common.Pinpoint_mark_error(fmt.Sprint(*err), "", 0, id)
}

func hook_common_func(f interface{}, hook_f interface{}, hook_f_trampoline interface{}) {
	// fmt.Println(reflect.ValueOf(f).Type())
	// fmt.Println(reflect.ValueOf(hook_f).Type())
	// fmt.Println(reflect.ValueOf(hook_f_trampoline).Type())
	funcName := get_func_name(f)
	common.Logf("try to hook " + funcName)
	if err := aop.AddHook(f, hook_f, hook_f_trampoline); err != nil {
		common.Logf("Hook "+funcName+" failed:%s", err)
		return
	}
	common.Logf(funcName + " is hooked")
}

func get_func_name(i interface{}) string {
	name := []byte(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
	for i := 0; i < len(name); i++ {
		if name[i] == '(' || name[i] == ')' {
			name = append(name[:i], name[i+1:]...)
		}
	}
	return string(name)
}

func gethost(dsn string) string {
	var host string
	var begin, end int
	if dsn != "" {
		for i := len(dsn) - 1; i >= 0; i-- {
			if dsn[i] == '(' {
				begin = i
				break
			} else if dsn[i] == ')' {
				end = i
			}
		}
		host = dsn[begin+1 : end]
		return host
	} else {
		return ""
	}
}
