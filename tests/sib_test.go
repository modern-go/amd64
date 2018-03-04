package tests

import (
	. "github.com/modern-go/amd64"
)

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			INC, DWORD_SIB(2, RAX, RBX, 0),
		},
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0x04, MODRM(ModeIndir, 0, RSP.Value())),
			aka(0x43, SIB(Scale2, RAX.Value(), RBX.Value())),
		},
	}, {
		input: input{
			INC, DWORD_SIB(0, RSP, RSP, 0),
		},
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0x04, MODRM(ModeIndir, 0, RSP.Value())),
			aka(0x24, SIB(Scale1, RSP.Value(), RSP.Value())),
		},
	}, {
		input: input{
			INC, DWORD(RSP, 0),
		},
		comment: "[rsp] is transformed to sib form",
		output: []uint8{
			aka(0xff, INC.Opcode()),
			aka(0x04, MODRM(ModeIndir, 0, RSP.Value())),
			aka(0x24, SIB(Scale1, RSP.Value(), RSP.Value())),
		},
	}}...)
}