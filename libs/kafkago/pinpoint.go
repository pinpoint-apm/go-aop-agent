package kafkago

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
	"github.com/segmentio/kafka-go"
)

func hook_common_func(f interface{}, hook_f interface{}, hook_f_trampoline interface{}) {
	funcName := get_func_name(f)
	common.Logf("try to hook " + funcName)
	if err := aop.AddHook(f, hook_f, hook_f_trampoline); err != nil {
		common.Logf("Hook "+funcName+" failed:%s", err)
		return
	}
	common.Logf(funcName + " is hooked")
}
func get_func_name(i interface{}) string {
	name := []byte(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
	for i := 0; i < len(name); i++ {
		if name[i] == '(' || name[i] == ')' {
			name = append(name[:i], name[i+1:]...)
		}
	}
	return string(name)
}

func commitMessages_onBefore(id common.TraceIdType, funcName string, reader *kafka.Reader, ctx context.Context, msgs ...kafka.Message) context.Context {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueSFunc := func(key, value string) {
		common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	}
	addClueFunc(common.PP_INTERCEPTOR_NAME, funcName)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_METHOD_CALL)
	if len(msgs) > 0 {
		addClueSFunc(common.PP_ARGS, fmt.Sprintf("commitMessages:%d ...", msgs[0].Offset))
	}

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return newCtx
}

func commitMessages_onEnd(id common.TraceIdType, err error) {
	// addClueSFunc := func(key, value string) {
	// 	common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	// }

	// addClueSFunc(common.PP_RETURN, fmt.Sprint(res))
}

func onException(id common.TraceIdType, err *error) {
	common.Logf("call onException")
	common.Pinpoint_add_exception(fmt.Sprint(*err), id)
}
