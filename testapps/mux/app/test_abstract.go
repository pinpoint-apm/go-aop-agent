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

package app

import (
	"fmt"
	"context"
)

type Info interface{
	TestAbstractFunc()
}

type Book struct{
	Color string
}

func Function(In Info){
	In.TestAbstractFunc()
}

func (B Book) TestAbstractFunc(ctx context.Context) string{
	fmt.Println("the book color is", B.Color)
	return "the book color is RED"
}

