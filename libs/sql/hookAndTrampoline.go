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

package sql

import (
	"context"
	"database/sql"

	"github.com/pinpoint-apm/go-aop-agent/common"
)

func init() {
	hook_common_func(sql.Open, hook_open, hook_open_trampoline)
	hook_common_func((*sql.DB).QueryContext, hook_query, hook_query_trampoline)
	hook_common_func((*sql.DB).ExecContext, hook_exec, hook_exec_trampoline)
	hook_common_func((*sql.DB).PingContext, hook_ping, hook_ping_trampoline)
}

/////////////////////sql.Open///////////////////////////
var DBMap map[*sql.DB]DSN

type DSN struct {
	driverName     string
	dataSourceName string
}

//go:noinline
func hook_open_trampoline(driverName, dataSourceName string) (*sql.DB, error) {
	return nil, nil
}

//go:noinline
func hook_open(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := hook_open_trampoline(driverName, dataSourceName)
	dsn := DSN{driverName: driverName, dataSourceName: dataSourceName}
	DBMap = make(map[*sql.DB]DSN)
	DBMap[db] = dsn
	return db, err
}

/////////////////////sql.DB.QueryContext///////////////////////////
//go:noinline
func hook_query_trampoline(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

//go:noinline
func hook_query(db *sql.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	funcName := get_func_name((*sql.DB).QueryContext)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_query_trampoline(db, ctx, query, args...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_query_trampoline(db, ctx, query, args...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(subTraceId, funcName, db, ctx, query, args...)
		res, err := hook_query_trampoline(db, *newCtx, query, args...)
		if err != nil {
			onException(subTraceId, &err)
		}
		queryonEnd(subTraceId, res)
		return res, err
	}
}

/////////////////////sql.DB.ExecContext///////////////////////////
//go:noinline
func hook_exec_trampoline(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

//go:noinline
func hook_exec(db *sql.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	funcName := get_func_name((*sql.DB).ExecContext)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_exec_trampoline(db, ctx, query, args...)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_exec_trampoline(db, ctx, query, args...)
		}

		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := onBefore(subTraceId, funcName, db, ctx, query, args...)
		res, err := hook_exec_trampoline(db, *newCtx, query, args...)
		if err != nil {
			onException(subTraceId, &err)
		}
		execonEnd(subTraceId, res)
		return res, err
	}
}

/////////////////////sql.DB.PingContext///////////////////////////
//go:noinline
func hook_ping_trampoline(db *sql.DB, ctx context.Context) error {
	return nil
}

//go:noinline
func hook_ping(db *sql.DB, ctx context.Context) error {
	funcName := get_func_name((*sql.DB).PingContext)
	if parentId, err := common.GetParentId(ctx); err != nil {
		common.Logf("parentId is not traceId type")
		return hook_ping_trampoline(db, ctx)
	} else {
		if common.Pinpoint_get_context(common.PP_HEADER_PINPOINT_SAMPLED, parentId) == common.PP_NOT_SAMPLED {
			common.Logf("trace dropped")
			return hook_ping_trampoline(db, ctx)
		}
		subTraceId := common.Pinpoint_start_trace(parentId)
		defer common.Pinpoint_end_trace(subTraceId)

		newCtx := pingonBefore(subTraceId, funcName, db, ctx)
		err := hook_ping_trampoline(db, *newCtx)
		if err != nil {
			onException(subTraceId, &err)
		}
		pingonEnd(subTraceId, &err)
		return err
	}
}
