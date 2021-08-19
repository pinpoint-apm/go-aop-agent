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
	"github.com/pinpoint-apm/go-aop-agent/common"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/httpClient"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/mongo"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/redisv8"
	_ "github.com/pinpoint-apm/go-aop-agent/libs/sql"
)

func init() {
	init_pinpoint()
}

func init_pinpoint() {
	common.Pinpoint_enable_debug_report(true)
	common.Pinpoint_set_collect_agent_host("tcp:127.0.0.1:9999")
	common.Pinpoint_set_trace_limit(10)
	common.Appname = "go-agent1"
	common.Appid = "Go-Agent1"
}
