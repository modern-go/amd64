package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			ADD, AL, BL,
		},
		comment: "r8,r8",
		output: []uint8{
			aka(0x00, ADD.Opcode()),
			aka(0xd8, MODRM(ModeReg, BL.Value(), AL.Value())),
		},
	}, {
		input: input{
			ADD, BYTE(RBX, 0), AL,
		},
		comment: "m8,r8",
		output: []uint8{
			aka(0x00, ADD.Opcode()),
			aka(0x03, MODRM(ModeIndir, AL.Value(), RBX.Value())),
		},
	}, {
		input: input{
			ADD, AX, BX,
		},
		comment: "r16,r16",
		output: []uint8{
			aka(0x66, Prefix16Bit),
			aka(0x01, ADD.Variant([2]VariantKey{{RM:16},{R:16}}).Opcode()),
			aka(0xd8, MODRM(ModeReg, BX.Value(), AX.Value())),
		},
	}, {
		input: input{
			ADD, WORD(RBX, 0), AX,
		},
		comment: "m16,r16",
		output: []uint8{
			aka(0x66, Prefix16Bit),
			aka(0x01, ADD.Variant([2]VariantKey{{RM:16},{R:16}}).Opcode()),
			aka(0x03, MODRM(ModeIndir, AX.Value(), RBX.Value())),
		},
	}}...)
}
