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

package redisv8

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/go-redis/redis/v8"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

var redisHost string

func init() {
	hook_common_func(redis.NewClient, hook_newclient, hook_newclient_trampoline)
}

type redisHook struct {
	redis.Hook
}

func (redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	common.Logf("call onBefore")
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. Dropped")
		return ctx, nil
	} else {
		id := common.Pinpoint_start_trace(parentId)

		addClueFunc := func(key, value string) {
			common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
		}

		addClueFunc(common.PP_INTERCEPTOR_NAME, fmt.Sprint(cmd.Args()[0]))
		addClueFunc(common.PP_SERVER_TYPE, common.PP_REDIS)
		addClueFunc(common.PP_DESTINATION, redisHost)

		newCtx := context.WithValue(ctx, common.TRACE_ID, id)
		return newCtx, nil
	}
}

func (redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	id, err := common.GetParentId(ctx)
	common.Pinpoint_end_trace(id)
	return err
}

// /////////////////////redis.NewClient.set///////////////////////////
//go:noinline
func hook_newclient_trampoline(opt *redis.Options) *redis.Client {
	return nil
}

//go:noinline
func hook_newclient(opt *redis.Options) *redis.Client {
	c := hook_newclient_trampoline(opt)
	c.AddHook(redisHook{})
	redisHost = c.String()
	return c
}

func hook_common_func(f interface{}, hook_f interface{}, hook_f_trampoline interface{}) {
	funcName := get_func_name(f)
	common.Logf("try to hook " + funcName)
	if err := aop.AddHook(f, hook_f, hook_f_trampoline); err != nil {
		common.Logf("Hook "+funcName+" failed:%s", err)
		return
	}
	common.Logf(funcName + " is hooked")
}

func get_func_name(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
