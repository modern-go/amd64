package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			ADD, AL, BL,
		},
		output: []uint8{
			aka(0x00, ADD.Opcode()),
			aka(0xd8, MODRM(ModeReg, BL.Value(), AL.Value())),
		},
		selected: true,
	}, {
		input: input{
			ADD, BYTE(RBX, 0), BL,
		},
		output: []uint8{
			0x00, 0x03,
		},
	}}...)
}
