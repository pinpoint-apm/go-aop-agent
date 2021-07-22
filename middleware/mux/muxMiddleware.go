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

package mux

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pinpoint-apm/go-aop-agent/common"
)

type pinpointResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (pp *pinpointResponseWriter) WriteHeader(code int) {
	pp.statusCode = code
	pp.ResponseWriter.WriteHeader(code)
}

func wrapperResponseWriter(w http.ResponseWriter) *pinpointResponseWriter {
	// default is 200
	return &pinpointResponseWriter{w, 200}
}

func PinpointMuxMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		// log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		pp := wrapperResponseWriter(w)
		// start trace
		id := common.Pinpoint_start_trace(common.ROOT_TRACE)
		addClueFunc := func(key, value string) {
			common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
		}
		// end trace
		defer func() {
			if pp != nil {
				common.Pinpoint_add_clues(common.PP_HTTP_STATUS_CODE, strconv.Itoa(pp.statusCode), id, common.CurrentTraceLoc)
			}
			common.Pinpoint_end_trace(id)
			// common.Logf("end trace:%d", id)
		}()

		nCtx := context.WithValue(r.Context(), common.TRACE_ID, id)
		pinpointRequest := r.WithContext(nCtx)

		// common.Logf("start trace:%d", id)
		addClueFunc(common.PP_APP_NAME, common.Appname)
		addClueFunc(common.PP_APP_ID, common.Appid)
		addClueFunc(common.PP_INTERCEPTOR_NAME, "mux middleware request")

		addClueFunc(common.PP_REQ_URI, pinpointRequest.RequestURI)
		addClueFunc(common.PP_REQ_SERVER, pinpointRequest.Host)
		addClueFunc(common.PP_REQ_CLIENT, pinpointRequest.RemoteAddr)
		addClueFunc(common.PP_SERVER_TYPE, common.GOLANG)
		common.Pinpoint_set_context(common.PP_SERVER_TYPE, common.GOLANG, id)
		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_PSPANID); value != "" {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}
		var sid string
		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_SPANID); value != "" {
			sid = value
		} else if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_SPANID); value != "" {
			sid = value
		} else {
			sid = common.Pinpoint_gen_sid()
		}
		addClueFunc(common.PP_SPAN_ID, sid)
		common.Pinpoint_set_context(common.PP_SPAN_ID, sid, id)

		var tid string
		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_TRACEID); value != "" {
			tid = value
		} else if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_TRACEID); value != "" {
			tid = value
		} else {
			tid = common.Pinpoint_gen_tid()
		}
		addClueFunc(common.PP_TRANSCATION_ID, tid)
		common.Pinpoint_set_context(common.PP_TRANSCATION_ID, tid, id)

		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_PAPPNAME); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, id)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_PAPPTYPE); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, id)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_HOST); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, id)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_PSPANID); value != "" {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_PAPPNAME); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, id)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_PAPPTYPE); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, id)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_HOST); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, id)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_NGINX_PROXY); value != "" {
			addClueFunc(common.PP_NGINX_PROXY, value)
		}

		if value := pinpointRequest.Header.Get(common.PP_HEADER_APACHE_PROXY); value != "" {
			addClueFunc(common.PP_APACHE_PROXY, value)
		}
		common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s1", id)

		if value1 := (pinpointRequest.Header.Get(common.PP_HTTP_PINPOINT_SAMPLED)); value1 == common.PP_NOT_SAMPLED || common.Pinpoint_tracelimit() {
			common.Pinpoint_drop_trace(id)
			common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s0", id)
		} else if value2 := pinpointRequest.Header.Get(common.PP_HEADER_PINPOINT_SAMPLED); value2 == common.PP_NOT_SAMPLED || common.Pinpoint_tracelimit() {
			common.Pinpoint_drop_trace(id)
			common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s0", id)
		}
		addClueFunc(common.PP_HTTP_METHOD, pinpointRequest.Method)
		next.ServeHTTP(pp, pinpointRequest)

	})
}
