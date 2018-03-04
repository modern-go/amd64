package amd64

var INC = &instruction{
	Mnemonic: "inc",
	opcode:   0xff,
	assemble: oneOperand,
	rm8:      0xfe,
}

var MOV = ""
var ADD = ""
var RET = ""
