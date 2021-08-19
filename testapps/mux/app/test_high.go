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

import "context"

type TestHighType func(context.Context, int) int

//go:noinline
func Add(ctx context.Context, a, b int) int {
	return (a + b)
}

//go:noinline
func Mul(ctx context.Context, a, b int) int {
	return (a * b)
}

//go:noinline
func TestHighFunc(ctx context.Context, a, b int, T func(context.Context, int, int) int) int {
	return T(ctx, a, b)
}
