package amd64

var INC = &instruction{
	Mnemonic: "inc",
	opcode:   0xff,
	assemble: oneOperand,
	variants: variants{
		{rm: 8}: {opcode: 0xfe},
	},
}
var DEC = &instruction{
	Mnemonic: "dec",
	opcode:   0xff,
	assemble: oneOperand,
	variants: variants{
		{rm: 8}: {
			opcode: 0xfe,
			overrides: overrides{
				"ModR/M r/m": byte(1),
			},
		},
	},
	overrides: overrides{
		"ModR/M r/m": byte(1),
	},
}

var MOV = ""
var ADD = ""
var RET = ""
