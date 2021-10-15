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

package transport

import (
	"context"
	"net/http"

	"github.com/pinpoint-apm/go-aop-agent/common"
)

func onBefore(parentId common.TraceIdType, req *http.Request) (common.TraceIdType, *http.Request) {
	common.Logf("call onBefore")
	id := common.Pinpoint_start_trace(parentId)

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, "*http.Transport.RoundTrip")
	addClueFunc(common.PP_SERVER_TYPE, common.PP_REMOTE_METHOD)
	addClueSFunc(common.PP_HTTP_URL, req.URL.String())
	addClueFunc(common.PP_DESTINATION, req.URL.Host)

	ctx := context.WithValue(req.Context(), common.TRACE_ID, id)
	return id, (*req).WithContext(ctx)
}

func onEnd(id common.TraceIdType, response *http.Response, err *error) {
	common.Logf("call onEnd")
	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if response != nil {
		addClueSFunc(common.PP_HTTP_STATUS_CODE, response.Status)
	} else {
		addClueSFunc(common.PP_HTTP_STATUS_CODE, "500")
		addClueFunc(common.PP_ADD_EXCEPTION, "response is nil")
	}
	common.Pinpoint_end_trace(id)
}

//go:noinline
func hook_transport(t *http.Transport, req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if id, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type. (*http.Transport).RoundTrip dropped")
		return hook_transport_trampoline(t, req)
	} else {
		// trace limited
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, id) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			req.Header.Set(common.PP_HEADER_PINPOINT_SAMPLED, common.PP_NOT_SAMPLED)
			return hook_transport_trampoline(t, req)
		}

		subId, pinpointReq := onBefore(id, req)
		response, err := hook_transport_trampoline(t, pinpointReq)
		onEnd(subId, response, &err)
		return response, err
	}

}

//go:noinline
func hook_transport_trampoline(t *http.Transport, req *http.Request) (*http.Response, error) {
	return nil, nil
}
