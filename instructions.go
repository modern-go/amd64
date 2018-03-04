package amd64

var INC = &instruction{
	Mnemonic: "inc",
	opcode:   0xff,
	assemble: one_operand,
}

var MOV = ""
var ADD = ""
var RET = ""
