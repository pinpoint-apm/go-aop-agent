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
// Copy from go asm
// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package x86asm
package asm

import (
	"encoding/hex"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/arch/x86/x86asm"
)

func TestDecode(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/decode.txt")
	if err != nil {
		t.Fatal(err)
	}
	all := string(data)
	for strings.Contains(all, "\t\t") {
		all = strings.Replace(all, "\t\t", "\t", -1)
	}
	for _, line := range strings.Split(all, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.SplitN(line, "\t", 4)
		i := strings.Index(f[0], "|")
		if i < 0 {
			t.Errorf("parsing %q: missing | separator", f[0])
			continue
		}
		if i%2 != 0 {
			t.Errorf("parsing %q: misaligned | separator", f[0])
		}
		size := i / 2
		code, err := hex.DecodeString(f[0][:i] + f[0][i+1:])
		if err != nil {
			t.Errorf("parsing %q: %v", f[0], err)
			continue
		}
		mode, err := strconv.Atoi(f[1])
		if err != nil {
			t.Errorf("invalid mode %q in: %s", f[1], line)
			continue
		}
		syntax, asm := f[2], f[3]
		inst, err := Decode(code, mode)
		var out string
		if err != nil {
			out = "error: " + err.Error()
		} else {
			switch syntax {
			case "gnu":
				out = x86asm.GNUSyntax(inst, 0, nil)
			case "intel":
				out = x86asm.IntelSyntax(inst, 0, nil)
			case "plan9": // [sic]
				out = x86asm.GoSyntax(inst, 0, nil)
			default:
				t.Errorf("unknown syntax %q", syntax)
				continue
			}
		}
		if out != asm || inst.Len != size {
			t.Errorf("Decode(%s) [%s] = %s, %d, want %s, %d | %s", f[0], syntax, out, inst.Len, asm, size, inst)
		}
	}
}

func TestDecodeDoesNotCrash(t *testing.T) {
	cases := [...][]byte{
		[]byte{},
		[]byte{0xc5},
		[]byte{0xc4},
	}
	for _, test := range cases {
		_, err := Decode([]byte(test), 64) // the only goal is that this line does not panic
		if err == nil {
			t.Errorf("expected error on invalid instruction %x", test)
		}
	}
}
