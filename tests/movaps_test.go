package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			MOVAPS, XMM1, XMM2,
		},
		output: []uint8{
			aka(0x0f, MOVAPS.Prefix0F()),
			aka(0x28, MOVAPS.Opcode()),
			aka(0xca, MODRM(ModeReg, XMM1.Value(), XMM2.Value())),
		},
	}, {
		input: input{
			MOV, XMM1, XMM2,
		},
		output: []uint8{
			aka(0x0f, MOVAPS.Prefix0F()),
			aka(0x28, MOVAPS.Opcode()),
			aka(0xca, MODRM(ModeReg, XMM1.Value(), XMM2.Value())),
		},
	}, {
		input: input{
			VMOVAPS, XMM1, XMM2,
		},
		output: []uint8{
			aka(0xc5, 0xc5),
			aka(0xf8, 0xf8),
			aka(0x28, MOVAPS.Opcode()),
			aka(0xca, MODRM(ModeReg, XMM1.Value(), XMM2.Value())),
		},
	}, {
		input: input{
			VMOVAPS, XMM1, XMM2,
		},
		output: []uint8{
			aka(0xc5, VMOVAPS.PrefixC5()),
			aka(0xf8, VEX2(0, 0, 0, 0)),
			aka(0x28, VMOVAPS.Opcode()),
			aka(0xca, MODRM(ModeReg, XMM1.Value(), XMM2.Value())),
		},
	}, {
		input: input{
			MOVAPS, XMM1, XMMWORD(RBX, 0),
		},
		output: []uint8{
			aka(0x0f, MOVAPS.Prefix0F()),
			aka(0x28, MOVAPS.Opcode()),
			aka(0x0b, MODRM(ModeIndir, XMM1.Value(), RBX.Value())),
		},
	}, {
		input: input{
			MOVAPS, XMMWORD(RBX, 0), XMM1,
		},
		output: []uint8{
			aka(0x0f, MOVAPS.Prefix0F()),
			aka(0x29, MOVAPS.Variant([2]VariantKey{{REG: "xmm", M: 128}, {REG: "xmm"}}).Opcode()),
			aka(0x0b, MODRM(ModeIndir, XMM1.Value(), RBX.Value())),
		},
	}, {
		input: input{
			VMOVAPS, XMMWORD(RBX, 0), XMM1,
		},
		output: []uint8{
			aka(0xc5, VMOVAPS.PrefixC5()),
			aka(0xf8, VEX2(0, 0, 0, 0)),
			aka(0x29, MOVAPS.Variant([2]VariantKey{{REG: "xmm", M: 128}, {REG: "xmm"}}).Opcode()),
			aka(0x0b, MODRM(ModeIndir, XMM1.Value(), RBX.Value())),
		},
	}}...)
}
