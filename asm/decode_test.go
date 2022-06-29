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
	"log"
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
		inst, err := Decode([]byte(test), 64) // the only goal is that this line does not panic
		if err == nil {
			t.Errorf("expected error on invalid instruction %x", test)
		}
		log.Println(inst)
	}
}

func TestDecodeDebug(t *testing.T) {
	cases := [...][]byte{
		[]byte{0x64, 0x48, 0x8b, 0x0c, 0x25, 0xf8, 0xff, 0xff, 0xff},
	}
	for _, test := range cases {
		inst, err := Decode([]byte(test), 64) // the only goal is that this line does not panic
		if err != nil {
			t.Errorf("expected error on invalid instruction %x", test)
		}
		log.Println(inst)
	}
}

func TestGODecodeDebug(t *testing.T) {
	cases := [...][]byte{
		[]byte{0x64, 0x48, 0x8b, 0x0c, 0x25, 0xf8, 0xff, 0xff, 0xff},
	}
	for _, test := range cases {
		inst, err := x86asm.Decode([]byte(test), 64) // the only goal is that this line does not panic
		if err != nil {
			t.Errorf("expected error on invalid instruction %x", test)
		}
		log.Println(inst)
	}
}

func Benchmark_C_Decode(b *testing.B) {
	codes := [][]byte{
		{0x48, 0x8b, 0x05, 0x02, 0x74, 0x21, 0x00},
		{0x0f, 0x01, 0xf8},
		{0x66, 0x90},
		{0x0f, 0xba, 0x30, 0x11},
		{0x26, 0xa0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		{0x66, 0x90},
		{0xc5},
		{0xc4},
		{0xe8, 0x76, 0x7f, 0xff, 0xff},
		{0xf3, 0xc3},
		{0x66, 0xe9, 0x11, 0x22, 0x33, 0x44},
		{0x65, 0xff, 0x25, 0x11, 0x22, 0x33, 0x44},
		{0x64, 0x48, 0x8b, 0x0c, 0x25, 0xf8, 0xff, 0xff, 0xff}, //11
	}

	for i := 0; i < len(codes); i++ {
		Decode(codes[1], 64)
	}
}

func Benchmark_go_Decode(b *testing.B) {

	codes := [][]byte{
		{0x48, 0x8b, 0x05, 0x02, 0x74, 0x21, 0x00},
		{0x0f, 0x01, 0xf8},
		{0x66, 0x90},
		{0x0f, 0xba, 0x30, 0x11},
		{0x26, 0xa0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		{0x66, 0x90},
		{0xc5},
		{0xc4},
		{0xe8, 0x76, 0x7f, 0xff, 0xff},
		{0xf3, 0xc3},
		{0x66, 0xe9, 0x11, 0x22, 0x33, 0x44},
		{0x65, 0xff, 0x25, 0x11, 0x22, 0x33, 0x44},
		{0x64, 0x48, 0x8b, 0x0c, 0x25, 0xf8, 0xff, 0xff, 0xff}, //11
	}

	for i := 0; i < len(codes); i++ {
		x86asm.Decode(codes[1], 64)
	}
}
