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
	"fmt"

	"github.com/pinpoint-apm/go-aop-agent/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertOneOnBefore(stub string, id common.TraceIdType, coll *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) *context.Context {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, stub)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_MONGDB_EXE_QUERY)
	addClueFunc(common.PP_DESTINATION, coll.Database().Name())

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return &newCtx
}

func insertOneOnEnd(id common.TraceIdType, res *mongo.InsertOneResult, err error) {

	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}

	insertId := fmt.Sprintf("%s", res.InsertedID)
	addClueFunc(common.PP_RETURN, insertId)
}

func findOnBefore(stub string, id common.TraceIdType, coll *mongo.Collection, ctx context.Context, filter interface{}) *context.Context {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	// addCluesFunc := func(key, value string) {
	// 	common.Pinpoint_add_clues(key, value, id, common.CurrentTraceLoc)
	// }

	addClueFunc(common.PP_INTERCEPTOR_NAME, stub)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_MONGDB_EXE_QUERY)
	addClueFunc(common.PP_DESTINATION, coll.Database().Name())

	// if bs, ok := filter.([]byte); ok {
	// 	// Slight optimization so we'll just use MarshalBSON and not go through the codec machinery.
	// 	args := fmt.Sprint("%s", bson.Raw(bs))
	// 	addCluesFunc(common.PP_ARGS, args)
	// }

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return &newCtx
}

func findOnEnd(id common.TraceIdType, cursor *mongo.Cursor, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}

}

func onBefore(stub string, id common.TraceIdType, coll *mongo.Collection, ctx context.Context) *context.Context {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	addClueFunc(common.PP_INTERCEPTOR_NAME, stub)
	addClueFunc(common.PP_SERVER_TYPE, common.PP_MONGDB_EXE_QUERY)
	addClueFunc(common.PP_DESTINATION, coll.Database().Name())

	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
	return &newCtx
}

func onEnd(id common.TraceIdType, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
}

func countOnEnd(id common.TraceIdType, val int64, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}

	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
	count := fmt.Sprintf("%d", val)
	addClueFunc(common.PP_RETURN, count)
}

func onEndResult(id common.TraceIdType, result *mongo.SingleResult) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	if err := result.Err(); err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
}

// func insertManyOnBefore(stub string, id common.TraceIdType, coll *mongo.Collection, ctx context.Context) *context.Context {
// 	addClueFunc := func(key, value string) {
// 		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
// 	}

// 	addClueFunc(common.PP_INTERCEPTOR_NAME, stub)
// 	addClueFunc(common.PP_SERVER_TYPE, common.PP_MONGDB_EXE_QUERY)
// 	addClueFunc(common.PP_DESTINATION, coll.Database().Name())

// 	newCtx := context.WithValue(ctx, common.TRACE_ID, id)
// 	return &newCtx
// }

func insertManyOnEnd(id common.TraceIdType, result *mongo.InsertManyResult, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
}

func updateOneOnEnd(id common.TraceIdType, result *mongo.UpdateResult, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
}

func deleteOnEnd(id common.TraceIdType, result *mongo.DeleteResult, err error) {
	addClueFunc := func(key, value string) {
		common.Pinpoint_add_clue(key, value, id, common.CurrentTraceLoc)
	}
	if err != nil {
		addClueFunc(common.PP_ADD_EXCEPTION, err.Error())
		return
	}
}
