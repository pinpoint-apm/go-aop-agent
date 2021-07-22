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

package main

// only import libs you cared

import (
	"log"

	"github.com/pinpoint-apm/go-aop-agent/common"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/httpClient"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/mongo"
)

func init() {
	common.SetLogCallBack(log.Printf)
	init_pinpoint()
}

func init_pinpoint() {
	common.Pinpoint_enable_debug_report(true)
	common.Pinpoint_set_collect_agent_host("tcp:127.0.0.1:9999")
	common.Pinpoint_set_trace_limit(-1)
	common.Appname = "test-go-echo"
	common.Appid = "test-go-echo-id"
}
