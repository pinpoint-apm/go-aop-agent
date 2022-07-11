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

package kafkago

import (
	"context"
	"errors"

	"github.com/pinpoint-apm/go-aop-agent/common"
	"github.com/segmentio/kafka-go"
)

func init() {
	hook_common_func((*kafka.Reader).CommitMessages, hook_commitMessages, hook_commitMessages_trampoline)
	hook_common_func((*kafka.Writer).WriteMessages, hook_writeMessages, hook_writeMessages_trampoline)
}

//go:noinline
func hook_writeMessages(writer *kafka.Writer, ctx context.Context, msgs ...kafka.Message) error {
	funcName := "kafka.Writer.WriteMessages"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_writeMessages_trampoline(writer, ctx, msgs...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_writeMessages_trampoline(writer, ctx, msgs...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := writeMessages_onBefore(subTraceId, funcName, writer, ctx, msgs...)
		err := hook_writeMessages_trampoline(writer, newCtx, msgs...)
		if err != nil {
			onException(subTraceId, &err)
		}
		commitMessages_onEnd(subTraceId, err)
		return err
	}
}

//go:noinline
func hook_writeMessages_trampoline(writer *kafka.Writer, ctx context.Context, msgs ...kafka.Message) error {
	return errors.New("")
}

//go:noinline
func hook_commitMessages(reader *kafka.Reader, ctx context.Context, msgs ...kafka.Message) error {
	funcName := "kafka.Reader.CommitMessages"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_commitMessages_trampoline(reader, ctx, msgs...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_commitMessages_trampoline(reader, ctx, msgs...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := commitMessages_onBefore(subTraceId, funcName, reader, ctx, msgs...)
		err := hook_commitMessages_trampoline(reader, newCtx, msgs...)
		if err != nil {
			onException(subTraceId, &err)
		}
		commitMessages_onEnd(subTraceId, err)
		return err
	}

}

//go:noinline
func hook_commitMessages_trampoline(reader *kafka.Reader, ctx context.Context, msgs ...kafka.Message) error {
	return errors.New("")
}
