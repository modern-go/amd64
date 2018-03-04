package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			DEC, EAX,
		},
		comment: "dec set opcode reg to 1",
		output: []uint8{
			aka(0xff, DEC.Opcode()),
			aka(0xc8, MODRM(ModeReg, DEC.OpcodeReg(), EAX.Value())),
		},
	}, {
		input: input{
			DEC, AL,
		},
		output: []uint8{
			aka(0xfe, DEC.Variant([2]VariantKey{{RM: 8}}).Opcode()),
			aka(0xc8, MODRM(ModeReg, DEC.OpcodeReg(), AL.Value())),
		},
	}}...)
}
