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
	"net/http"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
)

func init() {
	if err := aop.AddHookP_CALL((*http.Transport).RoundTrip, hook_transport, hook_transport_trampoline); err != nil {
		common.Logf("Hook (*http.Transport).RoundTrip failed:%s", err)
		return
	}
	common.Logf("(*http.Transport).RoundTrip is hooked")
}
