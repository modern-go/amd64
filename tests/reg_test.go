package tests

import (
	. "github.com/modern-go/amd64"
)

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			INC, EAX,
		},
		comment: "rex prefix not required because eax is 32 bit",
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0xc0, MODRM(ModeReg, 0, EAX.Value())),
		},
	}, {
		input: input{
			INC, RAX,
		},
		comment: "rax is 64 bit, requires rex prefix",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0xc0, MODRM(ModeReg, 0, EAX.Value())),
		},
	}, {
		input: input{
			INC, AL,
		},
		comment: "al is 8 bit, has a different opcode",
		output: []uint8{
			aka(0xfe, INC.Variant([2]VariantKey{{RM: 8}}).Opcode()),
			aka(0xc0, MODRM(ModeReg, 0, AL.Value())),
		},
	}, {
		input: input{
			INC, AX,
		},
		comment: "ax is 16 bit, need 0x66 prefix",
		output: []uint8{
			aka(0x66, Prefix16Bit),
			aka(0xff, INC.Opcode()),
			aka(0xc0, MODRM(ModeReg, 0, AX.Value())),
		},
	}, {
		input: input{
			INC, ECX,
		},
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0xc1, MODRM(ModeReg, 0, ECX.Value())),
		},
	}, {
		input: input{
			INC, R11,
		},
		output: []uint8{
			aka(0x49, REX(true, false, false, true)),
			aka(0xff, INC.Opcode()),
			aka(0xc3, MODRM(ModeReg, 0, R11.Value() - 8)),
		},
	}}...)
}