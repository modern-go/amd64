package tests

import (
	"context"
	"github.com/modern-go/amd64"
	. "github.com/modern-go/amd64"
	"github.com/modern-go/test"
	"github.com/modern-go/test/must"
	"testing"
)

func TestAssembler_MakeFunc(t *testing.T) {
	t.Run("by bytes", test.Case(func(ctx context.Context) {
		var ident func(i int) int
		assembler := amd64.Assembler{
			Buffer: []uint8{
				0x48, 0x8B, 0x44, 0x24, 0x08, // mov rax, qword ptr [rsp + 8]
				0x48, 0x89, 0x44, 0x24, 0x10, // mov qword ptr [rsp + 0x10], rax
				0xc3,                         // ret
			},
		}
		assembler.MakeFunc(&ident)
		must.Nil(assembler.Error)
		must.Equal(100, ident(100))
	}))
	t.Run("by asm", test.Case(func(ctx context.Context) {
		assembler := &amd64.Assembler{}
		assembler.Assemble(
			MOV, RAX, QWORD(RSP, 0x08),
			MOV, QWORD(RSP, 0x10), RAX,
			RET,
		)
		must.Nil(assembler.Error)
		must.Equal([]byte{
			0x48, 0x8B, 0x44, 0x24, 0x08, // mov rax, qword ptr [rsp + 8]
			0x48, 0x89, 0x44, 0x24, 0x10, // mov qword ptr [rsp + 0x10], rax
			0xc3,                         // ret
		}, assembler.Buffer)
		var ident func(i int) int
		assembler.MakeFunc(&ident)
		must.Nil(assembler.Error)
		must.Equal(100, ident(100))
	}))
}

func TestDump(t *testing.T) {
	t.Run("miss two", test.Case(func(ctx context.Context) {
		must.Equal("mov %miss%, %miss%", amd64.Dump(MOV))
	}))
	t.Run("miss one", test.Case(func(ctx context.Context) {
		must.Equal("mov rax, %miss%", amd64.Dump(MOV, RAX))
	}))
	t.Run("mov", test.Case(func(ctx context.Context) {
		must.Equal("mov rax, qword ptr [rsp+8]", amd64.Dump(MOV, RAX, QWORD(RSP, 8)))
	}))
	t.Run("two insn", test.Case(func(ctx context.Context) {
		must.Equal("mov rax, qword ptr [rsp+8]\nret", amd64.Dump(
			MOV, RAX, QWORD(RSP, 8),
			RET))
	}))
}

func TestSIB(t *testing.T) {
	t.Run("1xrax＋rbx＋0", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rax+rbx]", amd64.Dump(
			INC, DWORD_SIB(1, RAX, RBX, 0)))
	}))
	t.Run("2xrax＋rbx＋0", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rax*2+rbx]", amd64.Dump(
			INC, DWORD_SIB(2, RAX, RBX, 0)))
	}))
	t.Run("4xrax＋rbx＋0", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rax*4+rbx]", amd64.Dump(
			INC, DWORD_SIB(4, RAX, RBX, 0)))
	}))
	t.Run("8xrax＋rbx＋0", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rax*8+rbx]", amd64.Dump(
			INC, DWORD_SIB(8, RAX, RBX, 0)))
	}))
	t.Run("0xrsp＋rsp＋0", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rsp]", amd64.Dump(
			INC, DWORD_SIB(0, RSP, RSP, 0)))
	}))
	t.Run("0xrsp＋rsp＋16", test.Case(func(ctx context.Context) {
		must.Equal("inc dword ptr [rsp+16]", amd64.Dump(
			INC, DWORD_SIB(0, RSP, RSP, 16)))
	}))
}

type input []interface{}
type testCase struct {
	input    []interface{}
	output   []byte
	comment  string
	selected bool
}

var testCases []testCase

func TestAssembler_Assemble(t *testing.T) {
	for _, tc := range testCases {
		if tc.selected {
			testCases = []testCase{tc}
			break
		}
	}
	for _, testCase := range testCases {
		t.Run(Dump(testCase.input...), test.Case(func(ctx context.Context) {
			assembler := amd64.Assembler{}
			assembler.Assemble(testCase.input...)
			must.Nil(assembler.Error)
			must.Equal(testCase.output, assembler.Buffer)
		}))
	}
}
