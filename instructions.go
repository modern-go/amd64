package amd64

var INC = &instruction{
	Mnemonic: "inc",
	opcode:   0xff,
	assemble: oneOperand,
}

var MOV = ""
var ADD = ""
var RET = ""
