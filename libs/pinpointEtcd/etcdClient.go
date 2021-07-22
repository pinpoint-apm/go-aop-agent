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

package pinpointEtcd

import (
	"context"
	"runtime"
	_ "unsafe"

	"github.com/coreos/etcd/clientv3"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

const stub = "github.com/coreos/etcd/clientv3.newClient"

type packKV struct {
	kv  clientv3.KV
	dst string
}

type wrapperKV struct {
	packKv packKV
}

func combineStr(arr []string) string {
	var outStr string
	for i := 0; i < len(arr); i++ {
		outStr += arr[i] + " "
	}
	return outStr
}

func onBefore(id common.TraceIdType, funcName string, kv *wrapperKV, ctx *context.Context, key *string) *context.Context {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, funcName)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_REDIS)
	if key != nil {
		addClueSFunc(common.PP_ARGS, "key:"+*key)
	}
	addClueFunc(common.PP_DESTINATION, kv.packKv.dst)
	newCtx := context.WithValue(*ctx, common.TRACE_ID, id)

	return &newCtx
}

func onEnd(id common.TraceIdType, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
	// not implement right now
	// addClueFunc(common.PP_RETURN,  res ....)
}

func (kv *wrapperKV) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if parentId, err := common.GetParentId(ctx); err != nil {
		return kv.packKv.kv.Put(ctx, key, val, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return kv.packKv.kv.Put(ctx, key, val, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore Put ")
		newCtx := onBefore(subTraceId, "clientv3.Put", kv, &ctx, &key)
		res, err := kv.packKv.kv.Put(*newCtx, key, val, opts...)
		common.Logf("call onEnd Put")
		onEnd(subTraceId, err)
		return res, err
	}
}

func (kv *wrapperKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if parentId, err := common.GetParentId(ctx); err != nil {
		return kv.packKv.kv.Get(ctx, key, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return kv.packKv.kv.Get(ctx, key, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore Get ")
		newCtx := onBefore(subTraceId, "clientv3.Get", kv, &ctx, &key)
		res, err := kv.packKv.kv.Get(*newCtx, key, opts...)
		common.Logf("call onEnd Get")
		onEnd(subTraceId, err)
		return res, err
	}
}

func (kv *wrapperKV) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	if parentId, err := common.GetParentId(ctx); err != nil {
		return kv.packKv.kv.Delete(ctx, key, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return kv.packKv.kv.Delete(ctx, key, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore Delete ")
		newCtx := onBefore(subTraceId, "clientv3.Delete", kv, &ctx, &key)
		res, err := kv.packKv.kv.Delete(*newCtx, key, opts...)
		common.Logf("call onEnd Delete ")
		onEnd(subTraceId, err)
		return res, err
	}
}

func (kv *wrapperKV) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	if parentId, err := common.GetParentId(ctx); err != nil {
		return kv.packKv.kv.Do(ctx, op)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return kv.packKv.kv.Do(ctx, op)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore Do ")
		newCtx := onBefore(subTraceId, "clientv3.Do", kv, &ctx, nil)
		res, err := kv.packKv.kv.Do(*newCtx, op)
		common.Logf("call onEnd Do ")
		onEnd(subTraceId, err)
		return res, err
	}
}

func (kv *wrapperKV) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	if parentId, err := common.GetParentId(ctx); err != nil {
		return kv.packKv.kv.Compact(ctx, rev, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return kv.packKv.kv.Compact(ctx, rev, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore Compact ")
		newCtx := onBefore(subTraceId, "clientv3.Compact", kv, &ctx, nil)
		res, err := kv.packKv.kv.Compact(*newCtx, rev, opts...)
		common.Logf("call onEnd Compact ")
		onEnd(subTraceId, err)
		return res, err
	}
}

func (kv *wrapperKV) Txn(ctx context.Context) clientv3.Txn {
	return kv.packKv.kv.Txn(ctx)
}

//go:noinline
func hook_trampoline(cfg *clientv3.Config) (*clientv3.Client, error) {
	return nil, nil
}

//go:noinline
func hook(cfg *clientv3.Config) (*clientv3.Client, error) {
	client, err := hook_trampoline(cfg)
	if err == nil {
		// originKv := client.KV
		// nKv := pinpointKV{originKv: oriKv{kv: originKv}, cfg: cfg, dst: combineStr(cfg.Endpoints)}
		// kvs.PushBack(nKv)
		// client.KV = &nKv
		originKv := client.KV
		// nKv := new(pinpointKV)
		nKv := wrapperKV{}
		nKv.packKv = packKV{kv: originKv, dst: combineStr(cfg.Endpoints)}
		client.KV = &nKv
		// nKv.cfg = cfg
		// nKv.dst = combineStr(cfg.Endpoints)
		// kvs.PushBack(client.KV)
		// TODO, then test is done, remove gc
		common.Logf("replace client KV with pinpointKV,AOP working  %p %p ...", originKv, client.KV)
		runtime.GC()
		common.Logf("gc once")
		runtime.GC()
		common.Logf("gc once")
		runtime.GC()
		common.Logf("gc once")
	}
	return client, err
}

//go:linkname newClient github.com/coreos/etcd/clientv3.newClient
func newClient(*clientv3.Config) (*clientv3.Client, error)

func hook_etcdClient() {
	if err := aop.AddHook(newClient, hook, hook_trampoline); err != nil {
		common.Logf("Hook clientv3.newClient failed:%s", err)
		return
	}
	common.Logf("%s is hooked", stub)
}

func init() {
	common.Logf("try to hook %s", stub)
	hook_etcdClient()
}
