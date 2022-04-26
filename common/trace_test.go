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
	"log"
	"testing"
	"time"
)

func init() {
	SetLogCallBack(log.Printf)
	init_pinpoint()
}

func init_pinpoint() {
	Pinpoint_enable_debug_report(true)
	Pinpoint_set_collect_agent_host("tcp:127.0.0.1:9999")
	Pinpoint_set_trace_limit(-1)
	Appname = "test-go-echo"
	Appid = "test-go-echo-id"
}

func callSth(ctx context.Context, a int, b string) int {
	// do something
	time.Sleep(1 * time.Second)
	return 10
}

func TestCallsth(t *testing.T) {
	ctx := context.Background()
	trans := PinTransactionHeader{
		Url:        "/xxx",
		Host:       "/xxx",
		RemoteAddr: "127.0.0.1",
		ParentType: "8660",
		ParentName: "kafka",
	}

	pile := func(ctx context.Context) {
		callSth(ctx, 19, "xxx")
	}

	PinTranscation(&trans, pile, ctx)

}
