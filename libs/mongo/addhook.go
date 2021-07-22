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
	"github.com/pinpoint-apm/go-aop-agent/aop"
	"github.com/pinpoint-apm/go-aop-agent/common"
	"go.mongodb.org/mongo-driver/mongo"
)

func addhook() {
	// find*
	if err := aop.AddHook((*mongo.Collection).FindOneAndUpdate, hook_update, hook_trampoline_update); err != nil {
		common.Logf("Hook (*mongo.Collection).FindOneAndUpdate failed:%s", err)
		return
	}
	common.Logf(" (*mongo.Collection).FindOneAndUpdate is hooked")

	if err := aop.AddHook((*mongo.Collection).FindOneAndDelete, hook_delete, hook_trampoline_delete); err != nil {
		common.Logf("Hook (*mongo.Collection).FindOneAndDelete failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).FindOneAndDelete is hooked ")

	if err := aop.AddHook((*mongo.Collection).FindOneAndReplace, hook_replace, hook_trampoline_replace); err != nil {
		common.Logf("Hook (*mongo.Collection).FindOneAndReplace failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).FindOneAndReplace is hooked ")

	if err := aop.AddHook((*mongo.Collection).FindOne, hook_findone, hook_trampoline_findone); err != nil {
		common.Logf("Hook (*mongo.Collection).FindOne failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).FindOne is hooked ")

	if err := aop.AddHook((*mongo.Collection).Find, hook_find, hook_trampoline_find); err != nil {
		common.Logf("Hook (*mongo.Collection).Find failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).Find is hooked ")

	// insert
	if err := aop.AddHook((*mongo.Collection).InsertOne, hook_insertone, hook_trampoline_insertone); err != nil {
		common.Logf("Hook (*mongo.Collection).InsertOne failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).InsertOne is hooked ")

	if err := aop.AddHook((*mongo.Collection).InsertMany, hook_insertmany, hook_trampoline_insertmany); err != nil {
		common.Logf("Hook (*mongo.Collection).InsertMany failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).InsertMany is hooked ")
	// update*
	if err := aop.AddHook((*mongo.Collection).UpdateOne, hook_updateone, hook_trampoline_updateone); err != nil {
		common.Logf("Hook (*mongo.Collection).UpdateOne failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).UpdateOne is hooked ")

	if err := aop.AddHook((*mongo.Collection).UpdateByID, hook_updatebyid, hook_trampoline_updatebyid); err != nil {
		common.Logf("Hook (*mongo.Collection).UpdateByID failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).UpdateByID is hooked ")

	if err := aop.AddHook((*mongo.Collection).UpdateMany, hook_updatebmany, hook_trampoline_updatebmany); err != nil {
		common.Logf("Hook (*mongo.Collection).UpdateMany failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).UpdateMany is hooked ")
	// replace
	if err := aop.AddHook((*mongo.Collection).ReplaceOne, hook_replaceone, hook_trampoline_replaceone); err != nil {
		common.Logf("Hook (*mongo.Collection).ReplaceOne failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).ReplaceOne is hooked ")

	//delete
	if err := aop.AddHook((*mongo.Collection).DeleteOne, hook_deleteone, hook_trampoline_deleteone); err != nil {
		common.Logf("Hook (*mongo.Collection).DeleteOne failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).DeleteOne is hooked ")

	if err := aop.AddHook((*mongo.Collection).DeleteMany, hook_deletemany, hook_trampoline_deletemany); err != nil {
		common.Logf("Hook (*mongo.Collection).DeleteMany failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).DeleteMany is hooked ")

	if err := aop.AddHook((*mongo.Collection).Drop, hook_drop, hook_trampoline_drop); err != nil {
		common.Logf("Hook (*mongo.Collection).Drop failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).Drop is hooked ")

	if err := aop.AddHook((*mongo.Collection).CountDocuments, hook_countdocments, hook_trampoline_countdocments); err != nil {
		common.Logf("Hook (*mongo.Collection).CountDocuments failed:%s", err)
		return
	}
	common.Logf("(*mongo.Collection).CountDocuments is hooked ")

}

func init() {
	common.Logf("try to add hook on mongo api")
	addhook()
	common.Logf("hook on mongo api is done")
}
