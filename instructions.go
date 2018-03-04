package amd64

var MOV = &Instruction{"mov", arithmetic, j{0xB8}, ImmRm{j{0xc7}, 0}, j{0x89}, j{0x8b}, 64}
var ADD = &Instruction{"add", arithmetic,no{}, ImmRm{j{0x81}, 0}, j{0x01}, j{0x03}, 64}
var RET = &Instruction{Mnemonic: "ret", assemble: func(a *Assembler) {
	a.byte(0xc3)
}}
