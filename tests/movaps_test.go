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
	}}...)
}
