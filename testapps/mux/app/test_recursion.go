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
	"context"
	"fmt"
)

//go:noinline
func TestRecursion(ctx context.Context, i int) int {
	if i == 0 {
		return 1
	}
	return i * TestRecursion(ctx, i-1)
}

//go:noinline
func TestExpInRecursion(ctx context.Context, i int) (int, error) {
	if i == 0 {
		return 0, fmt.Errorf(`Exception in recursion when i is %b`, i)
	} else {
		p, _ := TestExpInRecursion(ctx, i-1)
		return i * p, nil
	}
}
