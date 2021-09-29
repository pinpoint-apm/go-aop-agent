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

// #cgo pkg-config: pinpoint_common
// #include <pinpoint_common/common.h>
// #include <string.h>
//PPAgentT global_agent_info;
import "C"
import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"unsafe"
)

type TraceIdType int32

/////////////////////////////////////////
type LocationType int32

const (
	CurrentTraceLoc LocationType = 0x0
	RootTraceLoc    LocationType = 0x1
)

/////////////////////////////////////////

/////////////////////////////////////////
// copy from pinpoint-python
const (
	PP_HTTP_PINPOINT_PSPANID  = "HTTP_PINPOINT_PSPANID"
	PP_HTTP_PINPOINT_SPANID   = "HTTP_PINPOINT_SPANID"
	PP_HTTP_PINPOINT_TRACEID  = "HTTP_PINPOINT_TRACEID"
	PP_HTTP_PINPOINT_PAPPNAME = "HTTP_PINPOINT_PAPPNAME"
	PP_HTTP_PINPOINT_PAPPTYPE = "HTTP_PINPOINT_PAPPTYPE"
	PP_HTTP_PINPOINT_HOST     = "HTTP_PINPOINT_HOST"

	PP_HEADER_PINPOINT_PSPANID  = "Pinpoint-Pspanid"
	PP_HEADER_PINPOINT_SPANID   = "Pinpoint-Spanid"
	PP_HEADER_PINPOINT_TRACEID  = "Pinpoint-Traceid"
	PP_HEADER_PINPOINT_PAPPNAME = "Pinpoint-Pappname"
	PP_HEADER_PINPOINT_PAPPTYPE = "Pinpoint-Papptype"
	PP_HEADER_PINPOINT_HOST     = "Pinpoint-Host"
	PP_HEADER_PINPOINT_CLIENT   = "Pinpoint-Client"

	PP_HEADER_NGINX_PROXY  = "Pinpoint-ProxyNginx"
	PP_HTTP_NGINX_PROXY    = "HTTP_Pinpoint-ProxyNginx"
	PP_HEADER_APACHE_PROXY = "PINPOINT-PROXYAPACHE"
	PP_HTTP_APACHE_PROXY   = "HTTP_PINPOINT_PROXYAPACHE"

	PP_HEADER_PINPOINT_SAMPLED = "Pinpoint-Sampled"
	PP_HTTP_PINPOINT_SAMPLED   = "HTTP_PINPOINT_SAMPLED"

	PP_DESTINATION      = "dst"
	PP_INTERCEPTOR_NAME = "name"
	PP_APP_NAME         = "appname"
	PP_APP_ID           = "appid"
	PP_REQ_URI          = "uri"
	PP_REQ_CLIENT       = "client"
	PP_REQ_SERVER       = "server"
	PP_SERVER_TYPE      = "stp"
	PP_AGENT_TYPE       = "FT"

	PP_PARENT_SPAN_ID = "psid"
	PP_PARENT_NAME    = "pname"
	PP_PARENT_TYPE    = "ptype"
	PP_PARENT_HOST    = "Ah"

	PP_NGINX_PROXY    = "NP"
	PP_APACHE_PROXY   = "AP"
	PP_TRANSCATION_ID = "tid"
	PP_SPAN_ID        = "sid"
	PP_NOT_SAMPLED    = "s0"
	PP_SAMPLED        = "s1"
	PP_NEXT_SPAN_ID   = "nsid"
	PP_ADD_EXCEPTION  = "EXP"

	PP_SQL_FORMAT  = "SQL"
	PP_ARGS        = "-1"
	PP_RETURN      = "14"
	GOLANG         = "1800"
	PP_METHOD_CALL = "1801"
	PP_CELERY      = "1702"

	PP_REMOTE_METHOD = "9401"

	PP_HTTP_URL              = "40"
	PP_HTTP_PARAM            = "41"
	PP_HTTP_PARAM_ENTITY     = "42"
	PP_HTTP_COOKIE           = "45"
	PP_HTTP_STATUS_CODE      = "46"
	PP_HTTP_METHOD           = "206"
	PP_HTTP_INTERNAL_DISPLAY = 48
	PP_HTTP_IO               = 49
	PP_MESSAGE_QUEUE_URI     = 100

	PP_MYSQL                   = "2101"
	PP_REDIS                   = "8200"
	PP_REDIS_REDISSON          = "8203"
	PP_REDIS_REDISSON_INTERNAL = "8204"
	PP_POSTGRESQL              = "2501"
	PP_MEMCACHED               = "8050"
	PP_MEMCACHED_FUTURE_GET    = "8051"
	PP_MONGDB_EXE_QUERY        = "2651"
	PP_KAFKA                   = "8660"
	PP_KAFKA_TOPIC             = "140"
	PP_RABBITMQ_CLIENT         = "8300"
	PP_RABBITMQ_EXCHANGEKEY    = "130"
	PP_RABBITMQ_ROUTINGKEY     = "131"
)

/////////////////////////////////////////

const (
	ROOT_TRACE = 0
	TRACE_ID   = "trace_id"
)

var (
	Appname string
	Appid   string
)

var logEnable = strings.ToLower(os.Getenv("PINPOINT_LOG_ENABLE"))

var logCallBack = log.Printf

var ignoreUrls = map[string]bool{}

func init() {
	C.global_agent_info.agent_type = C.int(1800)
	C.global_agent_info.trace_limit = C.long(-1)
	Appname = "notset"
	Appid = "notset"
}

func Logf(format string, v ...interface{}) {
	if logCallBack != nil && logEnable == "true" {
		logCallBack(format, v...)
	}
}

/**
 * @description: logger callback. if you want to specify your own logger callback
 * 		common.SetLogCallBack(callback) // nil => drop everything
 * @param {string} format
 * @param {...interface{}} v
 * @return {*}
 */
func SetLogCallBack(callback func(format string, v ...interface{})) {
	logCallBack = callback
}

/**
 * @description: For Debug, trace the pinpoint
 * @param {bool} enable
 * @return {*}
 */
func Pinpoint_enable_debug_report(enable bool) {
	if enable {
		Logf("enable debug report")
		C.global_agent_info.inter_flag |= C.uchar(1)
	} else {
		C.global_agent_info.inter_flag &= C.uchar(0xFE)
	}
}

/**
 * @description: unittest only
 * @param {*}
 * @return {*}
 */
func Pinpoint_enable_utest() {
	C.global_agent_info.inter_flag |= C.uchar(0x4)
}

/**
 * @description: set trace_limit.
 * @param {int32} limitPerSec times per second.(-1 means no limit)
 * @return {*}
 */
func Pinpoint_set_trace_limit(limitPerSec int32) {
	C.global_agent_info.trace_limit = C.long(limitPerSec)
}

/**
 * @description:  Set collector-agent host
 * @param {string} host: tcp:dev.collector:9999
 * @return {*}
 */
func Pinpoint_set_collect_agent_host(host string) {
	cstr := C.CString(host)
	defer C.free(unsafe.Pointer(cstr))

	C.strncpy((*C.char)(&C.global_agent_info.co_host[0]), (*C.char)(cstr), C.ulong(256))
}

/**
 * @description: Create an new trace tree(id=-1) or add a new trace into current trace tree (id>0)
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_start_trace(id TraceIdType) TraceIdType {
	return TraceIdType(C.pinpoint_start_trace(C.int(id)))
}

/**
 * @description: End trace node(id) or trace tree(If current id the root node)
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_end_trace(id TraceIdType) TraceIdType {
	return TraceIdType(C.pinpoint_end_trace(C.int(id)))
}

/**
* @description: Attach some information on current trace node
* @param {TraceIdType} id: trace node identifier
* @param {string} key
* @param {string} value
* @return {*}
 */
func Pinpoint_add_clue(key, value string, id TraceIdType, loc LocationType) {
	ckey := C.CString(key)
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cvalue))
	C.pinpoint_add_clue(C.int(id), ckey, cvalue, C.E_NODE_LOC(C.E_CURRENT_LOC))
}

/**
 * @description: check current trace node is root node or not
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_trace_is_root(id TraceIdType) bool {
	if C.pinpoint_trace_is_root(C.int(id)) == 1 {
		return true
	} else {
		return false
	}
}

/**
 * @description: Store some information on current trace tree.
 * context will be free when trace tree end.
 * @param {*} key
 * @param {string} value
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_set_context(key, value string, id TraceIdType) {
	ckey := C.CString(key)
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cvalue))
	C.pinpoint_set_context_key(C.int(id), ckey, cvalue)

}

/**
 * @description: Get current trace tree context by key
 * @param {string} key
 * @param {TraceIdType} id
 * @return {*} "" if not exist! So DO NOT set "" into context.
 */
func Pinpoint_get_context(key string, id TraceIdType) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.pinpoint_get_context_key(C.int(id), ckey)
	if cvalue == nil {
		return ""
	} else {
		return C.GoString(cvalue)
	}
}

/**
 * @description:  The same as `Pinpoint_add_clue`. API for add annotation.
 * @param {*} key
 * @param {string} value
 * @param {TraceIdType} id
 * @param {LocationType} loc
 * @return {*}
 */
func Pinpoint_add_clues(key, value string, id TraceIdType, loc LocationType) {
	ckey := C.CString(key)
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cvalue))
	C.pinpoint_add_clues(C.int(id), ckey, cvalue, C.E_NODE_LOC(C.E_CURRENT_LOC))
}

/**
 * @description: An unique id per host
 * @param {*}
 * @return {*}
 */
func Pinpoint_unique_id() int64 {
	return int64(C.generate_unique_id())
}

/**
 * @description: A random number from [0,2147483647)
 * @param {*}
 * @return {*}
 */
func Pinpoint_gen_sid() string {
	return fmt.Sprintf("%d", rand.Int31n(2147483647))
}

/**
 * @description: geneate a transaction id for a span.
 * format: APPID^Start time^A random number
 * @param {*}
 * @return {*}
 */
func Pinpoint_gen_tid() string {
	return fmt.Sprintf("%s^%d^%d", Appid, Pinpoint_start_time(), Pinpoint_unique_id())
}

/**
 * @description: Pinpoint-web doesn't know which span is error until you tell him.
 * @param {*} emsg
 * @param {string} error_filename
 * @param {uint32} error_lineno
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_mark_error(emsg, error_filename string, error_lineno uint32, id TraceIdType) {
	msg := C.CString(emsg)
	file_name := C.CString(error_filename)
	lineno := C.uint(error_lineno)
	C.catch_error(C.int(id), msg, file_name, lineno)
}

/**
 * @description: Agent first run time
 * @param {*}
 * @return {*}
 */
func Pinpoint_start_time() int64 {
	return int64(C.pinpoint_start_time())
}

/**
 * @description: Drop current trace tree(trace id).
 *  A dropped trace tree will not send to pinpoint-collector
 * @param {TraceIdType} id
 * @return {*}
 */
func Pinpoint_drop_trace(id TraceIdType) {
	C.mark_current_trace_status(C.int(id), C.int(C.E_TRACE_BLOCK))
}

/**
 * @description: Check sample speed is reached the limit or not.
 * @param {*}
 * @return {*} true: current trace should be dropped. false: not limited
 */
func Pinpoint_tracelimit() bool {
	if C.check_tracelimit(-1) == 1 {
		return true
	} else {
		return false
	}
}

/**
 * @description: get parent id ctx from context.Context
 * @param {context.Context} ctx
 * @return {*}
 */
func GetParentId(ctx context.Context) (TraceIdType, error) {
	if ctx == nil {
		return TraceIdType(-1), errors.New("no ctx")
	} else if parentId := ctx.Value(TRACE_ID); parentId == nil {
		// debug.PrintStack()
		Logf("no parentId")
		return TraceIdType(-1), errors.New("no parentId")
	} else {
		if id, OK := parentId.(TraceIdType); !OK {
			Logf("parentId is not traceId type")
			return TraceIdType(-1), errors.New("parentId is not traceId type")
		} else {
			return id, nil
		}
	}
}

/**
 * @description: Middleware use this to exclude some urls
 * @param {...string} urls
 * @return {*}
 */
func AddIgnoreUrls(urls ...string) {
	for _, url := range urls {
		// add or replace
		ignoreUrls[url] = true
	}
}

/**
 * @description: Check url is ignore by `AddIgnoreUrls`
 * @param {string} url
 * @return {*}
 */
func IsIgnore(url string) bool {
	_, OK := ignoreUrls[url]
	return OK
}

/**
 * @description:
 *	if FORCE_DISABLE_PINPOINT_AGENT ==true, aop,middleware could not working after restart
 *	user can evn FORCE_DISABLE_PINPOINT_AGENT=true to disable pinpoint agent without recompiling binary program
 */
func AgentIsDisabled() bool {
	if flag := os.Getenv("FORCE_DISABLE_PINPOINT_AGENT"); strings.ToLower(flag) == "true" {
		return true
	} else {
		return false
	}
}
