package tests

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/amd64"
	. "github.com/modern-go/amd64"
	"github.com/modern-go/test/must"
)

func TestAssembler_Assemble(t *testing.T) {
	t.Run("simple case", test.Case(func(ctx context.Context) {
		var ident func(i int) int
		assembler := amd64.Assembler{
			Buffer: []uint8{
				0x48, 0x8B, 0x44, 0x24, 0x08, // mov
				0x48, 0x89, 0x44, 0x24, 0x10, // mov
				0xc3,                         // ret
			},
		}
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

type input []interface{}
type testCase struct {
	input    []interface{}
	output   []byte
	selected bool
}

var testCases []testCase

func TestAssembler_Write(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(Dump(testCase.input...), test.Case(func(ctx context.Context) {
			assembler := amd64.Assembler{}
			assembler.Assemble(testCase.input...)
			must.Nil(assembler.Error)
			must.Equal(testCase.output, assembler.Buffer)
		}))
	}
}
