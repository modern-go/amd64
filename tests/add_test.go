package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			ADD, AL, BL,
		},
		comment: "r8,r8",
		output: []uint8{
			aka(0x02, ADD.Variant(VariantKey{{R: 8}, {RM: 8}}).Opcode()),
			aka(0xc3, MODRM(ModeReg, AL.Value(), BL.Value())),
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
			aka(0x03, ADD.Variant(VariantKey{{R: 16}, {RM: 16}}).Opcode()),
			aka(0xc3, MODRM(ModeReg, AX.Value(), BX.Value())),
		},
	}, {
		input: input{
			ADD, WORD(RBX, 0), AX,
		},
		comment: "m16,r16",
		output: []uint8{
			aka(0x66, Prefix16Bit),
			aka(0x01, ADD.Variant(VariantKey{{RM: 16}, {R: 16}}).Opcode()),
			aka(0x03, MODRM(ModeIndir, AX.Value(), RBX.Value())),
		},
	}, {
		input: input{
			ADD, BL, IMM(2),
		},
		comment: "r8,imm8",
		output: []uint8{
			aka(0x80, ADD.Variant(VariantKey{{RM: 8}, {IMM: 8}}).Opcode()),
			aka(0xC3, MODRM(ModeReg, 0, BL.Value())),
			0x02,
		},
	}, {
		input: input{
			ADD, AL, IMM(2),
		},
		comment: "al,imm8",
		output: []uint8{
			aka(0x04, ADD.Variant(VariantKey{{REG: "al"}, {IMM: 8}}).Opcode()),
			0x02,
		},
	}}...)
}
