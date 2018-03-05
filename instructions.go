package amd64

var INC = &instruction{
	mnemonic: "inc",
	opcode:   0xff,
	assemble: oneOperand,
	variants: variants{
		{{RM: 8}}:  {opcode: 0xfe},
		{{RM: 16}}: {},
		{{RM: 32}}: {},
		{{RM: 64}}: {},
	},
}
var DEC = &instruction{
	mnemonic:  "dec",
	opcode:    0xff,
	opcodeReg: 1,
	assemble:  oneOperand,
	variants: variants{
		{{RM: 8}}:  {opcode: 0xfe},
		{{RM: 16}}: {},
		{{RM: 32}}: {},
		{{RM: 64}}: {},
	},
}
var ADD = &instruction{
	mnemonic: "add",
	opcode:   0x00,
	assemble: twoOperands,
	variants: variants{
		{{RM: 8}, {R: 8}}: {},
		{{RM: 16}, {R: 16}}: {opcode: 0x01},
	},
}

var MOV = ""
var RET = ""

var allInstructions = []*instruction{
	INC,
	DEC,
}

func init() {
	for _, insn := range allInstructions {
		insn.initVariants()
	}
}
