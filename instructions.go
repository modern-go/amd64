package amd64

var INC = &instruction{
	Mnemonic: "inc",
	opcode:   0xff,
	assemble: oneOperand,
	variants: variants{
		{{rm: 8}}: {opcode: 0xfe},
		{{rm: 16}}: {},
		{{rm: 32}}: {},
		{{rm: 64}}: {},
	},
}
var DEC = &instruction{
	Mnemonic: "dec",
	opcode:   0xff,
	assemble: oneOperand,
	overrides: overrides{
		"ModR/M r/m": byte(1),
	},
	variants: variants{
		{{rm: 8}}: {opcode: 0xfe},
		{{rm: 16}}: {},
		{{rm: 32}}: {},
		{{rm: 64}}: {},
	},
}
var ADD = &instruction{
	Mnemonic: "add",
	opcode: 0,
	assemble: twoOperands,
	variants: variants{
		{{rm:8},{r:8}}:{},
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
