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
	}}...)
}
