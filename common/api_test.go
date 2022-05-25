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
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Pinpoint_enable_utest()
	m.Run()
}

func TestConfigAgent(t *testing.T) {
	Pinpoint_enable_debug_report(true)
	Pinpoint_set_trace_limit(100)
	Pinpoint_set_collect_agent_host("tcp:10.10.10.1:9999")
	t.Log("pass")
	Pinpoint_enable_debug_report(false)
}

func TestAgentApi(t *testing.T) {

	Pinpoint_enable_debug_report(true)

	traceIdRoot := Pinpoint_start_trace(ROOT_TRACE)
	t.Logf("trace:%d", traceIdRoot)
	Pinpoint_add_clue("x", "xx", traceIdRoot, CurrentTraceLoc)
	Pinpoint_set_context("contextRoot", "xx", traceIdRoot)

	traceId1 := Pinpoint_start_trace(traceIdRoot)
	t.Logf("trace:%d", traceId1)

	Pinpoint_add_clue("x", "xx", traceId1, CurrentTraceLoc)
	Pinpoint_set_context("contextRoot", "xx", traceId1)

	Pinpoint_set_context("x", "xx", traceIdRoot)
	traceId2 := Pinpoint_start_trace(traceId1)
	t.Logf("trace:%d", traceId2)
	Pinpoint_mark_error("xx", "xxx", 0, traceId2)
	Pinpoint_add_clue("x", "xx", traceId2, CurrentTraceLoc)
	Pinpoint_set_context("x", "xx", traceIdRoot)
	traceId3 := Pinpoint_start_trace(traceId2)
	t.Logf("trace:%d", traceId3)
	Pinpoint_add_clue("x", "xx", traceId3, RootTraceLoc)
	Pinpoint_set_context("x", "xx", traceId3)

	Pinpoint_set_context("x", "xx", traceIdRoot)
	if Pinpoint_get_context("x", traceId3) != "xx" {
		t.Log("Pinpoint_get_context traceId3 failed")
		t.Log(Pinpoint_get_context("x", traceId3))
		t.Fail()
	}

	Pinpoint_set_int_context("intxx", 1025, traceIdRoot)
	if v, err := Pinpoint_get_int_context("intxx", traceId3); err != nil {
		t.Error(err)
	} else if v != 1025 {
		t.Error("Pinpoint_get_int_context 1025 failed")
	}

	if _, err := Pinpoint_get_int_context("intxx", 0); err == nil {
		t.Error("err should not be nil")
	}

	Pinpoint_end_trace(traceId3)
	childId := Pinpoint_start_trace_opt(traceId2, "TraceMinTimeMs:23", "TraceOnlyException")
	time.Sleep(time.Millisecond * 100)
	Pinpoint_add_exception("test exception", childId)
	Pinpoint_end_trace(childId)

	if Pinpoint_get_context("x", traceId2) != "xx" {
		t.Log("Pinpoint_get_context traceId2 failed")
		t.Fail()
	}
	Pinpoint_mark_error("xx", "xxx", 0, traceId3)
	Pinpoint_end_trace(traceId2)

	Pinpoint_wake_trace(traceId2)
	Pinpoint_end_trace(traceId2)

	if Pinpoint_get_context("x", traceId1) != "xx" {
		t.Log("Pinpoint_get_context traceId1 failed")
		t.Fail()
	}

	Pinpoint_end_trace(traceId1)

	if Pinpoint_get_context("contextRoot", traceIdRoot) != "xx" {
		t.Log("Pinpoint_get_context traceIdRoot failed")
		t.Fail()
	}

	Pinpoint_end_trace(traceIdRoot)
	ShowAgentStatus()
}

func TestUniqId(t *testing.T) {
	id := Pinpoint_unique_id()
	id1 := Pinpoint_unique_id()
	id2 := Pinpoint_unique_id()
	if id1*2 != id+id2 {
		t.Fail()
	}

	t.Logf("%d %d %d ", id, id1, id2)
}

func TestStartTime(t *testing.T) {
	if Pinpoint_start_time() > time.Now().Unix() {
		t.Errorf("start_time %d ", Pinpoint_start_time())
		t.Fail()
	}

	if Pinpoint_start_time() < 1621845863 {
		t.Errorf("start_time %d ", Pinpoint_start_time())
		t.Fail()
	}

}

func TestDropTrace(t *testing.T) {
	Pinpoint_enable_debug_report(true)
	traceIdRoot := Pinpoint_start_trace(ROOT_TRACE)
	Pinpoint_drop_trace(traceIdRoot)
	Pinpoint_mark_error("xx", "xxx", 0, traceIdRoot)
	Pinpoint_end_trace(traceIdRoot)
}

func TestTraceLimit(t *testing.T) {
	Pinpoint_enable_debug_report(true)
	Pinpoint_set_trace_limit(0)
	if Pinpoint_tracelimit() == false {
		t.Log("limited = 0")
		t.Fail()
	}

	Pinpoint_set_trace_limit(1)
	traceIdRoot := Pinpoint_start_trace(ROOT_TRACE)
	Pinpoint_add_clues("oknk", "xxx", traceIdRoot, CurrentTraceLoc)
	if Pinpoint_trace_is_root(traceIdRoot) == false {
		t.Logf("%d should be root", traceIdRoot)
		t.Fail()
	}

	Pinpoint_end_trace(traceIdRoot)
	if Pinpoint_trace_is_root(traceIdRoot) {
		t.Logf("%d should not root", traceIdRoot)
		t.Fail()
	}
	Pinpoint_tracelimit()
	Pinpoint_tracelimit()

	if Pinpoint_tracelimit() == false {
		t.Log("limited = 1")
		t.Fail()
	}

	time.Sleep(1 * time.Second)

	if Pinpoint_tracelimit() == true {
		t.Log("limited = no limited")
		t.Fail()
	}
}

func TestMissingCase(t *testing.T) {
	Pinpoint_gen_sid()
	Pinpoint_gen_tid()
	AddIgnoreUrls("/a", "/b", "/C")
	if !IsIgnore("/a") {
		t.Log("IsIgnore failed")
		t.Fail()
	}
}
