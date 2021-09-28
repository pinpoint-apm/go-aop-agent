package micro

import (
	"context"

	"github.com/micro/go-micro/v2/server"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func PinpointHandle(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		traceId := common.Pinpoint_start_trace(common.ROOT_TRACE)
		addClueFunc := func(key, value string) {
			common.Pinpoint_add_clue(key, value, traceId, common.CurrentTraceLoc)
		}
		addCluesFunc := func(key, value string) {
			common.Pinpoint_add_clues(key, value, traceId, common.CurrentTraceLoc)
		}

		catchPanic := true
		defer func() {
			if catchPanic {
				common.Pinpoint_mark_error("PinpointHandle found a panic! o_o ....", "", 0, traceId)
			}

			// if req.Response() != nil {
			// 	common.Pinpoint_add_clues(common.PP_HTTP_STATUS_CODE, strconv.Itoa(c.Response().Status), id, common.CurrentTraceLoc)
			// }

			common.Pinpoint_end_trace(traceId)
			// common.Logf("end trace:%d", id)
		}()

		// update context
		nCtx := context.WithValue(ctx, common.TRACE_ID, traceId)
		addClueFunc(common.PP_APP_NAME, common.Appname)
		addClueFunc(common.PP_APP_ID, common.Appid)
		addClueFunc(common.PP_INTERCEPTOR_NAME, "echo middleware request")

		addClueFunc(common.PP_REQ_URI, req.Service())
		addClueFunc(common.PP_REQ_SERVER, req.Service())
		addClueFunc(common.PP_REQ_CLIENT, req.Endpoint())
		addClueFunc(common.PP_SERVER_TYPE, common.GOLANG)
		common.Pinpoint_set_context(common.PP_SERVER_TYPE, common.GOLANG, traceId)

		header := req.Header()
		var sid string
		if value, ok := header[common.PP_HEADER_PINPOINT_SPANID]; ok {
			sid = value
		} else {
			sid = common.Pinpoint_gen_sid()
		}

		addClueFunc(common.PP_SPAN_ID, sid)
		common.Pinpoint_set_context(common.PP_SPAN_ID, sid, traceId)

		var tid string
		if value, OK := header[common.PP_HTTP_PINPOINT_TRACEID]; OK {
			tid = value
		} else if value, OK := header[common.PP_HEADER_PINPOINT_TRACEID]; OK {
			tid = value
		} else {
			tid = common.Pinpoint_gen_tid()
		}
		addClueFunc(common.PP_TRANSCATION_ID, tid)
		common.Pinpoint_set_context(common.PP_TRANSCATION_ID, tid, traceId)

		if value, OK := header[common.PP_HTTP_PINPOINT_PAPPTYPE]; OK {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, traceId)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value, OK := header[common.PP_HEADER_PINPOINT_PAPPTYPE]; OK {
			common.Pinpoint_set_context(common.PP_PARENT_TYPE, value, traceId)
			addClueFunc(common.PP_PARENT_TYPE, value)
		}

		if value, OK := header[common.PP_HTTP_PINPOINT_HOST]; OK {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, traceId)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value, ok := header[common.PP_HEADER_PINPOINT_HOST]; ok {
			common.Pinpoint_set_context(common.PP_PARENT_HOST, value, traceId)
			addClueFunc(common.PP_PARENT_HOST, value)
		}

		if value, ok := header[common.PP_HTTP_PINPOINT_PSPANID]; ok {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}

		if value, ok := header[common.PP_HEADER_PINPOINT_PSPANID]; ok {
			addClueFunc(common.PP_PARENT_SPAN_ID, value)
		}

		if value, ok := header[common.PP_HEADER_PINPOINT_PAPPNAME]; ok {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, traceId)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value, ok := header[common.PP_HTTP_PINPOINT_PAPPNAME]; ok {
			common.Pinpoint_set_context(common.PP_PARENT_NAME, value, traceId)
			addClueFunc(common.PP_PARENT_NAME, value)
		}

		if value, ok := header[common.PP_HEADER_NGINX_PROXY]; ok {
			addClueFunc(common.PP_NGINX_PROXY, value)
		}

		if value, ok := header[common.PP_HEADER_APACHE_PROXY]; ok {
			addClueFunc(common.PP_APACHE_PROXY, value)
		}
		common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s1", traceId)

		if value, ok := header[common.PP_HTTP_PINPOINT_SAMPLED]; ok && (value == common.PP_NOT_SAMPLED || common.Pinpoint_tracelimit()) {
			common.Pinpoint_drop_trace(traceId)
			common.Pinpoint_set_context(common.PP_HEADER_PINPOINT_SAMPLED, "s0", traceId)
		}
		addCluesFunc(common.PP_HTTP_METHOD, "9162")

		err := fn(nCtx, req, rsp)
		catchPanic = false
		return err
	}
}
