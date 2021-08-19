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

package httpClient

import (
	"context"
	"net/http"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func generatePinpointHeader(id common.TraceIdType, req *http.Request) {
	common.Logf("generatePinpointHeader")
	req.Header.Set(common.PP_HEADER_PINPOINT_PAPPTYPE, common.GOLANG)
	req.Header.Set(common.PP_HEADER_PINPOINT_PAPPNAME, common.Appname)
	req.Header.Set("Pinpoint-Flags", "0")
	req.Header.Set(common.PP_HEADER_PINPOINT_HOST, req.URL.Host)
	if tid := common.Pinpoint_get_context(common.PP_TRANSCATION_ID, id); tid != "" {
		req.Header.Set(common.PP_HEADER_PINPOINT_TRACEID, tid)
	}

	if sid := common.Pinpoint_get_context(common.PP_SPAN_ID, id); sid != "" {
		req.Header.Set(common.PP_HEADER_PINPOINT_PSPANID, sid)
	}

	nextSid := common.Pinpoint_gen_sid()

	common.Pinpoint_set_context(common.PP_NEXT_SPAN_ID, nextSid, id)
	req.Header.Set(common.PP_HEADER_PINPOINT_SPANID, nextSid)
}

func onBefore_Do(parentId common.TraceIdType, c *http.Client, req *http.Request) (common.TraceIdType, *http.Request) {
	common.Logf("call onBefore_Do")
	id := common.Pinpoint_start_trace(parentId)

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, "*http.Client.Do")
	addClueFunc(common.PP_SERVER_TYPE, common.PP_REMOTE_METHOD)
	addClueSFunc(common.PP_HTTP_URL, req.URL.String())
	addClueFunc(common.PP_DESTINATION, req.URL.Host)
	// add pinpoint header
	generatePinpointHeader(id, req)

	// return a wrapped request
	ctx := context.WithValue(req.Context(), common.TRACE_ID, id)
	return id, (*req).WithContext(ctx)
}

func onEnd_Do(id common.TraceIdType, response *http.Response, err *error) {
	common.Logf("call onEnd_Do")
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	if value := common.Pinpoint_get_context(common.PP_NEXT_SPAN_ID, id); value != "" {
		addClueFunc(common.PP_NEXT_SPAN_ID, value)
	}

	addClueSFunc(common.PP_HTTP_STATUS_CODE, response.Status)
	common.Pinpoint_end_trace(id)
}

//go:noinline
func hook_Do_trampoline(c *http.Client, req *http.Request) (*http.Response, error) {
	return nil, nil
}

// obsolated
// func getParentId(req *http.Request) (common.TraceIdType, error) {
// 	parentId := req.Context().Value(common.TRACE_ID)
// 	if parentId == nil {
// 		common.Logf("no parentId")
// 		return common.TraceIdType(-1), errors.New("no parentId")
// 	} else {
// 		if id, OK := parentId.(common.TraceIdType); !OK {
// 			common.Logf("parentId is not traceId type")
// 			return common.TraceIdType(-1), errors.New("parentId is not traceId type")
// 		} else {
// 			return id, nil
// 		}
// 	}
// }

//go:noinline
func hook_Do(c *http.Client, req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if id, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. client.Do dropped")
		return hook_Do_trampoline(c, req)
	} else {
		// trace limited
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, id) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			req.Header.Set(common.PP_HEADER_PINPOINT_SAMPLED, common.PP_NOT_SAMPLED)
			return hook_Do_trampoline(c, req)
		}

		subId, pinpointReq := onBefore_Do(id, c, req)
		response, err := hook_Do_trampoline(c, pinpointReq)
		onEnd_Do(subId, response, &err)
		return response, err
	}

}

func hook_client_Do() {

	if err := aop.AddHookP_CALL((*http.Client).Do, hook_Do, hook_Do_trampoline); err != nil {
		common.Logf("Hook (*http.Client).Do failed:%s", err)
		return
	}
	common.Logf("client.Do() is hooked")
}

func init() {
	common.Logf("try to hook client.Do()")
	hook_client_Do()
}
