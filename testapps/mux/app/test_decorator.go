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

type Decoratorer interface{
	TestDecoratorFunc(context.Context)
}

type Cooker struct {}

type Rice struct {}

type Water struct {}

type Dinner struct {
	C Cooker
	R Rice
	W Water
}

func (C Cooker) TestDecoratorFunc(ctx context.Context) {
	fmt.Println("wash the pot") 
}

func (R Rice) TestDecoratorFunc(ctx context.Context) {
	fmt.Println("do sth before washing rice")
	R.WashRice(ctx)
	//return "do sth before washing rice" 
}

func (R Rice) WashRice(ctx context.Context) {
	fmt.Println("wash rice")
}

func (W Water) TestDecoratorFunc(ctx context.Context) {
	fmt.Println("do sth before adding water")
	W.AddWater(ctx)
	//return "do sth before adding water"
}

func (W Water) AddWater(ctx context.Context) {
	fmt.Println("add water")
}

func (D Dinner) TestDecoratorFunc(ctx context.Context) (string, string, string, string, string){
	D.C.TestDecoratorFunc(ctx)
	D.R.TestDecoratorFunc(ctx)
	D.W.TestDecoratorFunc(ctx)
	return "wash the pot\n", "do sth before washing rice\n", "wash rice\n", "do sth before adding water\n", "add water"

}
