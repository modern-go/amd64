package tests

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	. "github.com/modern-go/amd64"
	"github.com/modern-go/test/must"
)

func Test_simd(t *testing.T) {
	t.Run("end to end", test.Case(func(ctx context.Context) {
		asm := &Assembler{}
		asm.Assemble(MOVD, XMM0, EDI)
		asm.Assemble(VPBROADCASTD, XMM0, XMM0)
		must.Nil(asm.Error)
		must.Equal([]byte{
			0xc5, 0xf9, 0x6e, 0xc7, // vmovd xmm0,edi
			0xc4, 0xe2, 0x79, 0x58, 0xc0, // vpbroadcastd xmm0,xmm0
		},
			asm.Buffer)
	}))
}
