package tests

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	. "github.com/modern-go/amd64"
	"github.com/modern-go/test/must"
	"strconv"
)

func Test_simd(t *testing.T) {
	t.Run("end to end", test.Case(func(ctx context.Context) {
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
			// vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0x40]
			VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0x40),
			// vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0x50]
			VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0x50),
			// vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0x60]
			VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0x60),
			// vpcmpeqd xmm4,xmm0,XMMWORD PTR [rsi+0x70]
			VPCMPEQD, XMM4, XMM0, XMMWORD(RSI, 0x70),
			VPACKSSDW, XMM1, XMM1, XMM2,
			VPACKSSDW, XMM2, XMM3, XMM4,
			VPACKSSWB, XMM1, XMM1, XMM2,
			// vpmovmskb eax,xmm1
			VPMOVMSKB, EAX, XMM1,
			// vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0x80]
			VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0x80),
			// vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0x90]
			VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0x90),
			// vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0xa0]
			VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0xa0),
			// vpcmpeqd xmm4,xmm0,XMMWORD PTR [rsi+0xb0]
			VPCMPEQD, XMM4, XMM0, XMMWORD(RSI, 0xb0),
			VPACKSSDW, XMM1, XMM1, XMM2,
			VPACKSSDW, XMM2, XMM3, XMM4,
			VPACKSSWB, XMM1, XMM1, XMM2,
			// vpmovmskb edx,xmm1
			VPMOVMSKB, EDX, XMM1,
			// vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0xc0]
			VPCMPEQD, XMM1, XMM0, XMMWORD(RSI, 0xc0),
			// vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0xd0]
			VPCMPEQD, XMM2, XMM0, XMMWORD(RSI, 0xd0),
			// vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0xe0]
			VPCMPEQD, XMM3, XMM0, XMMWORD(RSI, 0xe0),
			// vpcmpeqd xmm0,xmm0,XMMWORD PTR [rsi+0xf0]
			VPCMPEQD, XMM0, XMM0, XMMWORD(RSI, 0xf0),
			// vpackssdw xmm1,xmm1,xmm2
			VPACKSSDW, XMM1, XMM1, XMM2,
			// vpackssdw xmm0,xmm3,xmm0
			VPACKSSDW, XMM0, XMM3, XMM0,
			// vpacksswb xmm0,xmm1,xmm0
			VPACKSSWB, XMM0, XMM1, XMM0,
			// vpmovmskb esi,xmm0
			VPMOVMSKB, ESI, XMM0,
			SHL, RSI, IMM(0x30),
			SHL, RDX, IMM(0x20),
			SHL, RAX, IMM(0x10),
			OR, RAX, RCX,
			OR, RAX, RDX,
			OR, RAX, RSI,
			// mov    QWORD PTR [rsp+0x18],rax
			MOV, QWORD(RSP, 0x18), RAX,
			RET,
		)
		must.Nil(asm.Error)
		must.Equal([]byte{
			0x48, 0x8B, 0x7c, 0x24, 0x08, // mov rdi,QWORD PTR [rsp+0x8]
			0x48, 0x8B, 0x74, 0x24, 0x10, // mov rsi,QWORD PTR [rsp+0x10]
			0xc5, 0xf9, 0x6e, 0xc7,       // vmovd xmm0,edi
			0xc4, 0xe2, 0x79, 0x58, 0xc0, // vpbroadcastd xmm0,xmm0
			0xc5, 0xf9, 0x76, 0x0e,       // vpcmpeqd xmm1, xmm0, xmmword ptr [rsi]
			0xc5, 0xf9, 0x76, 0x56, 0x10, // vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0x10]
			0xc5, 0xf9, 0x76, 0x5e, 0x20, // vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0x20]
			0xc5, 0xf9, 0x76, 0x66, 0x30, // vpcmpeqd xmm4,xmm0,XMMWORD PTR [rsi+0x30]
			0xc5, 0xf1, 0x6b, 0xca,       // vpackssdw xmm1, xmm1, xmm2
			0xc5, 0xe1, 0x6b, 0xd4,       // vpackssdw xmm2, xmm3, xmm4
			0xc5, 0xf1, 0x63, 0xca,       // vpacksswb xmm1, xmm1, xmm2
			0xc5, 0xf9, 0xd7, 0xc9,       // vpmovmskb ecx, xmm1
			0xc5, 0xf9, 0x76, 0x4e, 0x40, // vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0x40]
			0xc5, 0xf9, 0x76, 0x56, 0x50, // vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0x50]
			0xc5, 0xf9, 0x76, 0x5e, 0x60, // vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0x60]
			0xc5, 0xf9, 0x76, 0x66, 0x70, // vpcmpeqd xmm4,xmm0,XMMWORD PTR [rsi+0x70]
			0xc5, 0xf1, 0x6b, 0xca,       // vpackssdw xmm1, xmm1, xmm2
			0xc5, 0xe1, 0x6b, 0xd4,       // vpackssdw xmm2, xmm3, xmm4
			0xc5, 0xf1, 0x63, 0xca,       // vpacksswb xmm1, xmm1, xmm2
			0xc5, 0xf9, 0xd7, 0xc1,       // vpmovmskb eax,xmm1
			0xc5, 0xf9, 0x76, 0x8e,
			0x80, 0x00, 0x00, 0x00, // vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0x80]
			0xc5, 0xf9, 0x76, 0x96,
			0x90, 0x00, 0x00, 0x00, // vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0x90]
			0xc5, 0xf9, 0x76, 0x9e,
			0xa0, 0x00, 0x00, 0x00, // vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0xa0]
			0xc5, 0xf9, 0x76, 0xa6,
			0xb0, 0x00, 0x00, 0x00, // vpcmpeqd xmm4,xmm0,XMMWORD PTR [rsi+0xb0]
			0xc5, 0xf1, 0x6b, 0xca, // vpackssdw xmm1, xmm1, xmm2
			0xc5, 0xe1, 0x6b, 0xd4, // vpackssdw xmm2, xmm3, xmm4
			0xc5, 0xf1, 0x63, 0xca, // vpacksswb xmm1, xmm1, xmm2
			0xc5, 0xf9, 0xd7, 0xd1, // vpmovmskb edx,xmm1
			0xc5, 0xf9, 0x76, 0x8e,
			0xc0, 0x00, 0x00, 0x00, // vpcmpeqd xmm1,xmm0,XMMWORD PTR [rsi+0xc0]
			0xc5, 0xf9, 0x76, 0x96,
			0xd0, 0x00, 0x00, 0x00, // vpcmpeqd xmm2,xmm0,XMMWORD PTR [rsi+0xd0]
			0xc5, 0xf9, 0x76, 0x9e,
			0xe0, 0x00, 0x00, 0x00, // vpcmpeqd xmm3,xmm0,XMMWORD PTR [rsi+0xe0]
			0xc5, 0xf9, 0x76, 0x86,
			0xf0, 0x00, 0x00, 0x00,       // vpcmpeqd xmm0,xmm0,XMMWORD PTR [rsi+0xf0]
			0xc5, 0xf1, 0x6b, 0xca,       // vpackssdw xmm1,xmm1,xmm2
			0xc5, 0xe1, 0x6b, 0xc0,       // vpackssdw xmm0,xmm3,xmm0
			0xc5, 0xf1, 0x63, 0xc0,       // vpacksswb xmm0,xmm1,xmm0
			0xc5, 0xf9, 0xd7, 0xf0,       // vpmovmskb esi,xmm0
			0x48, 0xc1, 0xe6, 0x30,       // shl rsi, 0x30
			0x48, 0xc1, 0xe2, 0x20,       // shl rdx, 0x20
			0x48, 0xc1, 0xe0, 0x10,       // shl rax, 0x10
			0x48, 0x09, 0xc8,             // or rax, rcx
			0x48, 0x09, 0xd0,             // or rax, rdx,
			0x48, 0x09, 0xf0,             // or rax, rsi,
			0x48, 0x89, 0x44, 0x24, 0x18, // mov QWORD PTR [rsp+0x18],rax
			0xc3,                         //ret
		},
			asm.Buffer)
		var compareEqual func(key uint32, elements *[64]uint32) (ret uint64)
		asm.MakeFunc(&compareEqual)
		must.Nil(asm.Error)
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
		if "1100000000000000000000000000000000000000000000000000000000001001" != strconv.FormatUint(uint64(ret), 2) {
			t.Fail()
		}
	}))
}
