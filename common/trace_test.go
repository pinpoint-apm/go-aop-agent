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
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func init_test() {
	SetLogCallBack(log.Printf)
	init_pinpoint()
}

func init_pinpoint() {
	Pinpoint_enable_debug_report(true)
	Pinpoint_set_collect_agent_host("tcp:127.0.0.1:9999")
	Pinpoint_set_trace_limit(-1)
	Appname = "test-go-echo"
	Appid = "test-go-echo-id"
	fmt.Println("----------------------------")
}

func httpclient(ctx context.Context) {
	_, deferfun := PinHttpClientFunc(ctx, "userHttpclient", "http://www.naver.com1/index.html", nil)
	ret := make([]string, 1)
	defer deferfun(nil, ret)
	ret[0] = "success"
}

func httpclientV1(ctx context.Context) {
	_, deferfun := PinHttpClientFunc(ctx, "userHttpclientV1", "http://www.naver.com1/index.html", []string{
		"TraceMinTimeMs:23",
		"TraceOnlyException",
	})
	ret := make([]string, 1)
	var err error
	defer deferfun(&err, ret)
	time.Sleep(1 * time.Second)
	err = errors.New("test exception")
	ret[0] = "success"
}

func callSum(ctx context.Context) {
	_, deferfun := PinFuncSum(ctx, "callSum")
	defer deferfun(nil)

}

func callOnce(ctx context.Context) {
	_, deferfun := PinFuncOnce(ctx, "callOnce")
	defer deferfun(nil)

}

func callSth(ctx context.Context, a int, b string) int {
	// do something
	time.Sleep(1 * time.Second)
	// call userhttclient
	httpclient(ctx)
	httpclientV1(ctx)
	for i := 0; i < 1000; i++ {
		callSum(ctx)
	}
	callOnce(ctx)
	return 10
}

func TestCallsth(t *testing.T) {
	init_test()
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
