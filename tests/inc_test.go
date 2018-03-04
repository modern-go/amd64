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
			INC, ECX,
		},
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0xc1, MODRM(ModeReg, 0, ECX.Value())),
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
			aka(0xfe, INC.Variant([2]VariantKey{{RM:8}}).Opcode()),
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
			INC, QWORD(RAX, 0),
		},
		comment: "rax is 64 bit, need rex prefix",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, DWORD(RAX, 0),
		},
		comment: "32bit do not need rex prefix",
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, DWORD(EAX, 0),
		},
		comment: "32bit register need 0x67 prefi",
		output: []uint8{
			aka(0x67, Prefix32Bit),
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, EAX.Value())),
		},
	}, {
		input: input{
			INC, WORD(RAX, 0),
		},
		comment: "16bit need 0x66 prefix",
		output: []uint8{
			aka(0x66, Prefix16Bit),
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, WORD(EAX, 0),
		},
		comment: "16bit need 0x66 prefix, eax need 0x67 prefix",
		output: []uint8{
			aka(0x67, Prefix32Bit),
			aka(0x66, Prefix16Bit),
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, EAX.Value())),
		},
	}, {
		input: input{
			INC, BYTE(RAX, 0),
		},
		output: []uint8{
			aka(0xfe, INC.Variant([2]VariantKey{{RM:8}}).Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, BYTE(EAX, 0),
		},
		output: []uint8{
			aka(0x67, Prefix32Bit),
			aka(0xfe, INC.Variant([2]VariantKey{{RM:8}}).Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, QWORD(RAX, 0x10),
		},
		comment: "0x10 in the displacement",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x40, MODRM(ModeIndirDisp8, 0, RAX.Value())),
			0x10,
		},
	}, {
		input: input{
			INC, QWORD(RAX, 0x7f),
		},
		comment: "0x7f is still 8bit displacement",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x40, MODRM(ModeIndirDisp8, 0, RAX.Value())),
			0x7f,
		},
	}, {
		input: input{
			INC, QWORD(RAX, -0x7f),
		},
		comment: "-0x7f is still 8bit displacement",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x40, MODRM(ModeIndirDisp8, 0, RAX.Value())),
			0x81, // -0x7f
		},
	}, {
		input: input{
			INC, QWORD(RAX, 0x80),
		},
		comment: "0x80 need 32bit, so mode changed",
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x80, MODRM(ModeIndirDisp32, 0, RAX.Value())),
			0x80, 0x00, 0x00, 0x00,
		},
	}, {
		input: input{
			INC, QWORD(EAX, 0),
		},
		comment: "eax is 32bit, need 0x67 prefix",
		output: []uint8{
			aka(0x67, Prefix32Bit),
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0x00, MODRM(ModeIndir, 0, RAX.Value())),
		},
	}, {
		input: input{
			INC, ESP,
		},
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0xc4, MODRM(ModeReg, 0, ESP.Value())),
		},
	}, {
		input: input{
			INC, RSP,
		},
		output: []uint8{
			aka(0x48, REX(true, false, false, false)),
			aka(0xff, INC.Opcode()),
			aka(0xc4, MODRM(ModeReg, 0, ESP.Value())),
		},
	}}...)
}
