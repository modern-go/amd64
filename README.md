# AMD64 Instruction Assembler

[![Sourcegraph](https://sourcegraph.com/github.com/modern-go/amd64/-/badge.svg)](https://sourcegraph.com/github.com/modern-go/amd64?badge)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/modern-go/amd64)
[![Build Status](https://travis-ci.org/modern-go/amd64.svg?branch=master)](https://travis-ci.org/modern-go/amd64)
[![codecov](https://codecov.io/gh/modern-go/amd64/branch/master/graph/badge.svg)](https://codecov.io/gh/modern-go/amd64)
[![rcard](https://goreportcard.com/badge/github.com/modern-go/amd64)](https://goreportcard.com/report/github.com/modern-go/amd64)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/modern-go/amd64/master/LICENSE)

* generate code during runtime: assembler, but run in your process
* Go assembly does not support all SIMD instruction

This does not support all instructions yet. 
But it has laid a ground work on instruction encoding abstraction.
New instruction will be added on demand basis.

# Tutorial

```go
import . "github.com/modern-go/amd64"

asm := &Assembler{}
asm.Assemble( 
    // RAX = i
    MOV, RAX, QWORD(RSP, 0x08),
    // j = RAX
    MOV, QWORD(RSP, 0x10), RAX,
    // return j
    RET,
)
// ident func does nothing
// return identical value out
var ident func(i int) (j int)
asm.MakeFunc(&ident)
fmt.Println(ident(100)) // will print 100
```

# SIMD

```go
asm := &Assembler{}
asm.Assemble(
    MOV, RDI, QWORD(RSP, 8),
    MOV, RSI, QWORD(RSP, 16),
    MOVD, XMM0, EDI,
    VPBROADCASTD, XMM0, XMM0,
    VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0),
    VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0x10),
    VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0x20),
    VPCMPEQD, XMM4, XMM0, XMMWORD(RSI, 0x30),
    VPACKSSDW, XMM1, XMM1, XMM2,
    VPACKSSDW, XMM2, XMM3, XMM4,
    VPACKSSWB, XMM1, XMM1, XMM2,
    VPMOVMSKB, ECX, XMM1,
    VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0x40),
    VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0x50),
    VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0x60),
    VPCMPEQD, XMM4, XMM0, XMMWORD(RSI, 0x70),
    VPACKSSDW, XMM1, XMM1, XMM2,
    VPACKSSDW, XMM2, XMM3, XMM4,
    VPACKSSWB, XMM1, XMM1, XMM2,
    VPMOVMSKB, EAX, XMM1,
    VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0x80),
    VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0x90),
    VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0xa0),
    VPCMPEQD, XMM4, XMM0, XMMWORD(RSI, 0xb0),
    VPACKSSDW, XMM1, XMM1, XMM2,
    VPACKSSDW, XMM2, XMM3, XMM4,
    VPACKSSWB, XMM1, XMM1, XMM2,
    VPMOVMSKB, EDX, XMM1,
    VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0xc0),
    VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0xd0),
    VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0xe0),
    VPCMPEQD, XMM0, XMM0, XMMWORD(RSI, 0xf0),
    VPACKSSDW, XMM1, XMM1, XMM2,
    VPACKSSDW, XMM0, XMM3, XMM0,
    VPACKSSWB, XMM0, XMM1, XMM0,
    VPMOVMSKB, ESI, XMM0,
    SHL, RSI, IMM(0x30),
    SHL, RDX, IMM(0x20),
    SHL, RAX, IMM(0x10),
    OR, RAX, RCX,
    OR, RAX, RDX,
    OR, RAX, RSI,
    MOV, QWORD(RSP, 0x18), RAX,
    RET,
)
var compareEqual func(key uint32, elements *[64]uint32) (ret uint64)
asm.MakeFunc(&compareEqual)
v1 := [64]uint32{
    3, 0, 0, 3, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 3, 3}
ret := compareEqual(3, &v1)
fmt.Println(strconv.FormatUint(uint64(ret), 2))
```

the output will be

```
1100000000000000000000000000000000000000000000000000000000001001
```

it searches the integer in the array faster by utilizing the SIMD instruction

# Acknowledgement 

* Initial implementation copied from https://github.com/nelhage/gojit
* https://medium.com/kokster/writing-a-jit-compiler-in-golang-964b61295f

