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
			0xff, aka(0xc0, MODRM(ModeReg, 0, EAX.Value())),
		},
		selected: true,
	}, {
		input: input{
			INC, ECX,
		},
		output: []uint8{
			0xff, MODRM(ModeReg, 0, ECX.Value()),
		},
	}, {
		input: input{
			INC, RAX,
		},
		comment: "rax is 64 bit, requires rex prefix",
		output: []uint8{
			REX(true, false, false, false),
			0xff, MODRM(ModeReg, 0, EAX.Value()),
		},
	}, {
		input: input{
			INC, AL,
		},
		comment: "al is 8 bit, has a different opcode",
		output: []uint8{
			0xfe, MODRM(ModeReg, 0, AL.Value()),
		},
	}, {
		input: input{
			INC, AX,
		},
		comment: "ax is 16 bit, need 0x66 prefix",
		output: []uint8{
			0x66, 0xff, 0xc0,
		},
	}, {
		input: input{
			INC, QWORD(RAX, 0),
		},
		comment: "rax is 64 bit, need rex prefix",
		output: []uint8{
			0x48, 0xff, 0x00,
		},
	}, {
		input: input{
			INC, QWORD(RAX, 16),
		},
		comment: "16 in the displacement",
		output: []uint8{
			0x48, 0xff, 0x40, 0x10,
		},
	}, {
		input: input{
			INC, QWORD(EAX, 0),
		},
		comment: "eax is 32bit, need 0x67 prefix",
		output: []uint8{
			0x67, 0x48, 0xff, 0x00,
		},
	}}...)
}
