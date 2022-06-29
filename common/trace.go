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
package common

import (
	"context"
	"fmt"
	"net/url"
)

type PinTransactionHeader struct {
	Url        string
	Host       string
	RemoteAddr string
	ParentType string
	ParentName string
	ParentHost string
	ParentTid  string
	Err        error
}

type DeferFunc func(*error, ...interface{})

var emptyPinFunc = func(err *error, ret ...interface{}) {}

/**
 * PinFuncSum profile cumulative pefermance function
 * @param ctx context.Context
 */
func PinFuncSum(ctx context.Context, name string, args ...interface{}) (context.Context, DeferFunc) {

	if AgentIsDisabled() {
		return ctx, emptyPinFunc
	}

	if parentId, err := GetParentId(ctx); err != nil ||
		Pinpoint_get_context(PP_HEADER_PINPOINT_SAMPLED, parentId) == PP_NOT_SAMPLED {
		return ctx, emptyPinFunc
	} else {
		var id TraceIdType
		var nctx context.Context
		key := "[sum]" + name
		if v, err := Pinpoint_get_int_context(key, parentId); err == nil {
			// found v
			id = TraceIdType(v)
			Pinpoint_wake_trace(id)
			nctx = ctx
		} else {
			id = Pinpoint_start_trace(parentId)
			if id == TraceIdType(-1) {
				return ctx, emptyPinFunc
			}
			Pinpoint_set_int_context(key, int64(id), id)
			nctx = context.WithValue(ctx, TRACE_ID, id)
			Pinpoint_add_clue(PP_SERVER_TYPE, PP_METHOD_CALL, id, CurrentTraceLoc)
			Pinpoint_add_clue(PP_INTERCEPTOR_NAME, key, id, CurrentTraceLoc)
		}

		deferfunc := func(err *error, ret ...interface{}) {
			if err != nil && *err != nil {
				Pinpoint_add_clue(PP_ADD_EXCEPTION, (*err).Error(), id, CurrentTraceLoc)
			}

			Pinpoint_end_trace(id)
		}
		return nctx, deferfunc
	}
}

/**
 * PinFuncOnce profile  function once
 */
func PinFuncOnce(ctx context.Context, name string, args ...interface{}) (context.Context, DeferFunc) {

	if AgentIsDisabled() {
		return ctx, emptyPinFunc
	}

	if parentId, err := GetParentId(ctx); err != nil ||
		Pinpoint_get_context(PP_HEADER_PINPOINT_SAMPLED, parentId) == PP_NOT_SAMPLED {
		return ctx, emptyPinFunc
	} else {

		id := Pinpoint_start_trace(parentId)
		if id == TraceIdType(-1) {
			return ctx, emptyPinFunc
		}
		nctx := context.WithValue(ctx, TRACE_ID, id)
		addClueFunc := func(key, value string) {
			Pinpoint_add_clue(key, value, id, CurrentTraceLoc)
		}
		addClueSFunc := func(key, value string) {
			Pinpoint_add_clues(key, value, id, CurrentTraceLoc)
		}
		addClueFunc(PP_SERVER_TYPE, PP_METHOD_CALL)
		addClueFunc(PP_INTERCEPTOR_NAME, name)
		addClueSFunc(PP_ARGS, fmt.Sprint(args...))

		deferfunc := func(err *error, ret ...interface{}) {
			if err != nil && *err != nil {
				addClueFunc(PP_ADD_EXCEPTION, (*err).Error())
			}

			if len(ret) > 0 {
				addClueSFunc(PP_RETURN, fmt.Sprint(ret...))
			}

			Pinpoint_end_trace(id)
		}
		return nctx, deferfunc
	}
}

type FuncPile func(context.Context)

func GenerateTid() string {
	return Pinpoint_gen_tid()
}

func PinTranscation(header *PinTransactionHeader, pile FuncPile, parentCtx context.Context) {

	catchPanic := true
	if AgentIsDisabled() {
		pile(parentCtx)
	} else {

		id := Pinpoint_start_trace(ROOT_TRACE)

		newCtx, cancel := context.WithCancel(parentCtx)
		defer cancel()
		//note: update context
		pinctx := context.WithValue(newCtx, TRACE_ID, id)

		addClueFunc := func(key, value string) {
			Pinpoint_add_clue(key, value, id, CurrentTraceLoc)
		}

		sid := Pinpoint_gen_sid()
		addClueFunc(PP_SPAN_ID, sid)
		Pinpoint_set_context(PP_SPAN_ID, sid, id)

		addClueFunc(PP_APP_NAME, Appname)
		addClueFunc(PP_APP_ID, Appid)
		addClueFunc(PP_INTERCEPTOR_NAME, "pinpoint middleware")

		addClueFunc(PP_REQ_URI, header.Url)
		addClueFunc(PP_REQ_SERVER, header.Host)
		addClueFunc(PP_REQ_CLIENT, header.RemoteAddr)
		addClueFunc(PP_SERVER_TYPE, GOLANG)
		Pinpoint_set_context(PP_HEADER_PINPOINT_SAMPLED, "s1", id)

		if header.ParentType != "" {
			addClueFunc(PP_PARENT_TYPE, header.ParentType)
		}

		if header.ParentName != "" {
			addClueFunc(PP_PARENT_NAME, header.ParentName)
		}

		if header.ParentHost != "" {
			addClueFunc(PP_PARENT_HOST, header.ParentHost)
		}

		var tid string

		if header.ParentTid != "" {
			tid = header.ParentTid
			addClueFunc(PP_PARENT_SPAN_ID, tid)
			addClueFunc(PP_NEXT_SPAN_ID, sid)
		} else {
			tid = GenerateTid()
		}

		addClueFunc(PP_TRANSCATION_ID, tid)
		Pinpoint_set_context(PP_TRANSCATION_ID, tid, id)
		// end transcation

		defer func() {
			if catchPanic {
				Pinpoint_mark_error("PinpointMiddleWare found a panic! o_o ....", "", 0, id)
			}
			Pinpoint_end_trace(id)
		}()
		// note: must be an error
		// just return
		if header.Err != nil {
			catchPanic = false
			Pinpoint_mark_error(header.Err.Error(), "trace.go", 0, id)
			return
		}
		pile(pinctx)
		catchPanic = false
	}

}

func GetAppName() string {
	return Appname
}

func PinHttpClientFunc(ctx context.Context, name, remoteUrl string, option []string, args ...interface{}) (context.Context, DeferFunc) {
	if AgentIsDisabled() {
		return ctx, emptyPinFunc
	}

	if parentId, err := GetParentId(ctx); err != nil ||
		Pinpoint_get_context(PP_HEADER_PINPOINT_SAMPLED, parentId) == PP_NOT_SAMPLED {
		return ctx, emptyPinFunc
	} else {
		var id TraceIdType
		if option == nil {
			id = Pinpoint_start_trace(parentId)
		} else {
			id = Pinpoint_start_trace_opt(parentId, option...)
		}

		if id == TraceIdType(-1) {
			return ctx, emptyPinFunc
		}

		nctx := context.WithValue(ctx, TRACE_ID, id)
		addClueFunc := func(key, value string) {
			Pinpoint_add_clue(key, value, id, CurrentTraceLoc)
		}
		addClueSFunc := func(key, value string) {
			Pinpoint_add_clues(key, value, id, CurrentTraceLoc)
		}
		addClueFunc(PP_SERVER_TYPE, PP_REMOTE_METHOD)
		addClueFunc(PP_INTERCEPTOR_NAME, name)
		addClueSFunc(PP_HTTP_URL, remoteUrl)
		u, err := url.Parse(remoteUrl)
		if err == nil {
			addClueFunc(PP_DESTINATION, u.Host)
		}

		deferfunc := func(err *error, ret ...interface{}) {
			if err != nil && *err != nil {
				Pinpoint_add_exception((*err).Error(), id)
			}

			if len(ret) > 0 {
				addClueSFunc(PP_RETURN, fmt.Sprint(ret...))
			}

			Pinpoint_end_trace(id)
		}
		return nctx, deferfunc
	}
}
