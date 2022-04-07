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

type FuncPile func(context.Context)

func GenerateTid() string {
	return Pinpoint_gen_tid()
}

func PinAsTranscation(header *PinTransactionHeader, pile FuncPile, ctx context.Context) {

	catchPanic := true
	if AgentIsDisabled() {
		goto CALL_PILE
	} else {

		id := Pinpoint_start_trace(ROOT_TRACE)
		//note: update context
		ctx = context.WithValue(ctx, TRACE_ID, id)

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
		addClueFunc(PP_PARENT_TYPE, header.ParentType)
		addClueFunc(PP_PARENT_NAME, header.ParentName)
		if header.ParentHost != "" {
			addClueFunc(PP_PARENT_HOST, header.ParentHost)
		}
		var tid string

		if header.ParentTid != "" {
			tid = header.ParentTid
			addClueFunc(PP_PARENT_SPAN_ID, tid)
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
	}

CALL_PILE:
	pile(ctx)
	catchPanic = false
}

func GetAppName() string {
	return Appname
}
