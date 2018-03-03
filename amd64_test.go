package amd64_test

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
				0x48, 0x8B, 0x44, 0x24, 0x08,
				0x48, 0x89, 0x44, 0x24, 0x10,
				0xc3, // ret
			},
		}
		assembler.Assemble(&ident)
		must.Nil(assembler.Error)
		must.Equal(100, ident(100))
	}))
}

func TestAssembler_Write(t *testing.T) {
	t.Run("simple case", test.Case(func(ctx context.Context) {
		var ident func(i int) int
		assembler := amd64.Assembler{}
		assembler.Write(
			MOV, RAX, QWORD(RSP, 8),
			MOV, QWORD(RSP, 16), RAX,
			RET,
		)
		assembler.Assemble(&ident)
		must.Nil(assembler.Error)
		must.Equal(100, ident(100))
	}))
}
