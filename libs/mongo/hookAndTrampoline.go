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

package mongo

import (
	"context"

	"github.com/pinpoint-apm/go-aop-agent/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:noinline
func hook_trampoline_update(coll *mongo.Collection, ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return nil
}

//go:noinline
func hook_update(coll *mongo.Collection, ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	const stub = "*mongo.Collection.FindOneAndUpdate"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_update(coll, ctx, filter, update, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_update(coll, ctx, filter, update, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		res := hook_trampoline_update(coll, *newCtx, filter, update, opts...)
		common.Logf("call onEnd %s ", stub)
		onEnd(subTraceId, err)
		return res
	}
}

//go:noinline
func hook_trampoline_delete(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return nil
}

//go:noinline
func hook_delete(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	const stub = "*mongo.Collection.FindOneAndDelete"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_delete(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_delete(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		res := hook_trampoline_delete(coll, *newCtx, filter, opts...)
		common.Logf("call onEnd %s ", stub)
		onEnd(subTraceId, err)
		return res
	}
}

//go:noinline
func hook_trampoline_replace(coll *mongo.Collection, ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return nil
}

//go:noinline
func hook_replace(coll *mongo.Collection, ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {

	const stub = "*mongo.Collection.FindOneAndReplace"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_replace(coll, ctx, filter, replacement, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_replace(coll, ctx, filter, replacement, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		res := hook_trampoline_replace(coll, *newCtx, filter, replacement, opts...)
		common.Logf("call onEnd %s ", stub)
		onEnd(subTraceId, err)
		return res
	}

}

//go:noinline
func hook_trampoline_insertone(coll *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}

//go:noinline
func hook_insertone(coll *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	var stub = "*mongo.Collection.InsertOne"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_insertone(coll, ctx, document, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_insertone(coll, ctx, document, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call onBefore %s ", stub)
		newCtx := insertOneOnBefore(stub, subTraceId, coll, ctx, document, opts...)
		res, err := hook_trampoline_insertone(coll, *newCtx, document, opts...)
		common.Logf("call onEnd %s ", stub)
		insertOneOnEnd(subTraceId, res, err)
		return res, err
	}
}

//go:noinline
func hook_trampoline_find(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}

//go:noinline
func hook_find(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	var stub = "*mongo.Collection.Find"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_find(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_find(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)
		// debug.PrintStack()
		common.Logf("call OnBefore %s ", stub)
		newCtx := findOnBefore(stub, subTraceId, coll, ctx, filter)
		cursor, result := hook_trampoline_find(coll, *newCtx, filter, opts...)
		common.Logf("call OnEnd %s ", stub)
		findOnEnd(subTraceId, cursor, result)
		return cursor, result
	}
}

//go:noinline
func hook_trampoline_findone(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}

//go:noinline
func hook_findone(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	var stub = "*mongo.Collection.FindOne"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_findone(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_findone(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := findOnBefore(stub, subTraceId, coll, ctx, filter)
		result := hook_trampoline_findone(coll, *newCtx, filter, opts...)
		common.Logf("call OnEnd %s ", stub)
		onEndResult(subTraceId, result)
		return result
	}
}

//go:noinline
func hook_insertmany(coll *mongo.Collection, ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	var stub = "*mongo.Collection.InsertMany"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_insertmany(coll, ctx, documents, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_insertmany(coll, ctx, documents, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := hook_trampoline_insertmany(coll, *newCtx, documents, opts...)
		common.Logf("call OnEnd %s ", stub)
		insertManyOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_trampoline_insertmany(coll *mongo.Collection, ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return nil, nil
}

//go:noinline
func hook_updateone(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var stub = "*mongo.Collection.UpdateOne"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_updateone(coll, ctx, filter, update, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_updateone(coll, ctx, filter, update, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := hook_trampoline_updateone(coll, *newCtx, filter, update, opts...)
		common.Logf("call OnEnd %s ", stub)
		updateOneOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_trampoline_updateone(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

//go:noinline
func hook_trampoline_updatebyid(coll *mongo.Collection, ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

//go:noinline
func hook_updatebyid(coll *mongo.Collection, ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var stub = "*mongo.Collection.UpdateByID"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_updateone(coll, ctx, id, update, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_updateone(coll, ctx, id, update, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := hook_trampoline_updateone(coll, *newCtx, id, update, opts...)
		common.Logf("call OnEnd %s ", stub)
		updateOneOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_updatebmany(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var stub = "*mongo.Collection.UpdateMany"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_updatebmany(coll, ctx, filter, update, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_updatebmany(coll, ctx, filter, update, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := hook_trampoline_updatebmany(coll, *newCtx, filter, update, opts...)
		common.Logf("call OnEnd %s ", stub)
		updateOneOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_trampoline_updatebmany(coll *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

//go:noinline
func hook_replaceone(coll *mongo.Collection, ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	var stub = "*mongo.Collection.ReplaceOne"
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_trampoline_replaceone(coll, ctx, filter, replacement, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_trampoline_replaceone(coll, ctx, filter, replacement, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := hook_trampoline_replaceone(coll, *newCtx, filter, replacement, opts...)
		common.Logf("call OnEnd %s ", stub)
		updateOneOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_trampoline_replaceone(coll *mongo.Collection, ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return nil, nil
}

//go:noinline
func hook_trampoline_deleteone(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var stub = "*mongo.Collection.DeleteOne"
	trampline := hook_deleteone
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return trampline(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return trampline(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := trampline(coll, *newCtx, filter, opts...)
		common.Logf("call OnEnd %s ", stub)
		deleteOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_deleteone(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}

//go:noinline
func hook_deletemany(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var stub = "*mongo.Collection.DeleteMany"
	trampline := hook_trampoline_deletemany
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return trampline(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return trampline(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		result, err := trampline(coll, *newCtx, filter, opts...)
		common.Logf("call OnEnd %s ", stub)
		deleteOnEnd(subTraceId, result, err)
		return result, err
	}
}

//go:noinline
func hook_trampoline_deletemany(coll *mongo.Collection, ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return nil, nil
}

//go:noinline
func hook_drop(coll *mongo.Collection, ctx context.Context) error {
	var stub = "*mongo.Collection.Drop"
	trampline := hook_trampoline_drop
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return trampline(coll, ctx)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return trampline(coll, ctx)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		err := trampline(coll, *newCtx)
		common.Logf("call OnEnd %s ", stub)
		onEnd(subTraceId, err)
		return err
	}
}

//go:noinline
func hook_trampoline_drop(coll *mongo.Collection, ctx context.Context) error {
	return nil
}

//go:noinline
func hook_countdocments(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var stub = "*mongo.Collection.CountDocuments"
	trampline := hook_trampoline_countdocments
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return trampline(coll, ctx, filter, opts...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return trampline(coll, ctx, filter, opts...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		common.Logf("call OnBefore %s ", stub)
		newCtx := onBefore(stub, subTraceId, coll, ctx)
		val, err := trampline(coll, *newCtx, filter, opts...)
		common.Logf("call OnEnd %s ", stub)
		countOnEnd(subTraceId, val, err)
		return val, err
	}
}

//go:noinline
func hook_trampoline_countdocments(coll *mongo.Collection, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return 0, nil
}
