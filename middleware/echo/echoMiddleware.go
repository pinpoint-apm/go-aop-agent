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

package echo

import (
	"context"
	"strconv"

	"github.com/labstack/echo"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func PinpointMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	if common.AgentIsDisabled() {
		return func(c echo.Context) (err error) {
			return next(c)
		}
	}

	return func(c echo.Context) (err error) {
		// check while list
		url := c.Request().RequestURI
		if common.IsIgnore(url) {
			common.Logf("%s is ignore by setting", url)
			return next(c)
		}

		id := common.Pinpoint_start_trace(common.ROOT_TRACE)
		addClueFunc := func(key, value string) {
			common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
		}

		addCluesFunc := func(key, value string) {
			common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
		}
		catchPanic := true
		defer func() {
			if catchPanic {
				common.Pinpoint_mark_error("PinpointMiddleWare found a panic! o_o ....", "", 0, id)
			}

			if c.Response() != nil {
				common.Pinpoint_add_clues(common.PP_HTTP_STATUS_CODE, strconv.Itoa(c.Response().Status), id, common.CurrentTraceLoc)
			}
			common.Pinpoint_end_trace(id)
			// common.Logf("end trace:%d", id)
		}()

		nCtx := context.WithValue(c.Request().Context(), common.TRACE_ID, id)
		nReq := c.Request().WithContext(nCtx)

		// common.Logf("start trace:%d", id)
		addClueFunc(common.PP_APP_NAME, common.Appname)
		addClueFunc(common.PP_APP_ID, common.Appid)
		addClueFunc(common.PP_INTERCEPTOR_NAME, "echo middleware request")

		addClueFunc(common.PP_REQ_URI, url)
		addClueFunc(common.PP_REQ_SERVER, c.Request().Host)
		addClueFunc(common.PP_REQ_CLIENT, c.Request().RemoteAddr)
		addClueFunc(common.PP_SERVER_TYPE, common.GOLANG)
		common.Pinpoint_set_context(common.PP_SERVER_TYPE, common.GOLANG, id)

		var sid string
		if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_SPANID); value != "" {
			sid = value
		} else {
			sid = common.Pinpoint_gen_sid()
		}

		addClueFunc(common.PP_SPAN_ID, sid)
		common.Pinpoint_set_context(common.PP_SPAN_ID, sid, id)

		var tid string
		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_TRACEID); value != "" {
			tid = value
		} else if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_TRACEID); value != "" {
			tid = value
		} else {
			tid = common.Pinpoint_gen_tid()
		}

		addClueFunc(common.PP_TRANSCATION_ID, tid)
		common.Pinpoint_set_context(common.PP_TRANSCATION_ID, tid, id)

		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_PAPPTYPE); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, id)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_PAPPTYPE); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, id)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_HOST); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, id)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_HOST); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, id)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_PSPANID); value != "" {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_PSPANID); value != "" {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_PINPOINT_PAPPNAME); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, id)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_PAPPNAME); value != "" {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, id)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_NGINX_PROXY); value != "" {
			addClueFunc(common.PP_NGINX_PROXY, value)
		}

		if value := c.Request().Header.Get(common.PP_HEADER_APACHE_PROXY); value != "" {
			addClueFunc(common.PP_APACHE_PROXY, value)
		}
		common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s1", id)

		if value := c.Request().Header.Get(common.PP_HTTP_PINPOINT_SAMPLED); value == common.PP_NOT_SAMPLED || common.Pinpoint_tracelimit() {
			common.Pinpoint_drop_trace(id)
			common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s0", id)
		}
		addCluesFunc(common.PP_HTTP_METHOD, c.Request().Method)
		c.SetRequest(nReq)
		err = next(c)
		catchPanic = false
		return err
	}
}

func PinpointErrorHandler(originHandler func(e error, c echo.Context)) func(e error, c echo.Context) {

	return func(e error, c echo.Context) {
		originHandler(e, c)
		if id, err := common.GetParentId(c.Request().Context()); err == nil {
			common.Pinpoint_mark_error(e.Error(), "", 0, id)
		}
	}
}
