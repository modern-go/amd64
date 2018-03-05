package amd64

var INC = &instruction{
	mnemonic: "inc",
	opcode:   0xff,
	encoding: oneOperand,
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
	encoding:  oneOperand,
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
	encoding: twoOperands,
	variants: variants{
		{{RM: 8}, {R: 8}}:      {},
		{{RM: 16}, {R: 16}}:    {opcode: 0x01},
		{{RM: 32}, {R: 32}}:    {opcode: 0x01},
		{{RM: 64}, {R: 64}}:    {opcode: 0x01},
		{{R: 8}, {RM: 8}}:      {opcode: 0x02},
		{{R: 16}, {RM: 16}}:    {opcode: 0x03},
		{{R: 32}, {RM: 32}}:    {opcode: 0x03},
		{{R: 64}, {RM: 64}}:    {opcode: 0x03},
		{{REG: "al"}, {IMM: 8}}: {opcode: 0x04, encoding: encodingI}, // 0x04 is shorter than 0x80 form
		{{RM: 8}, {IMM: 8}}:    {opcode: 0x80},
	},
}

var MOV = &instruction{
	mnemonic: "mov",
	opcode: 0x88,
	encoding: twoOperands,
	variants: variants{
		{{RM:8},{R:8}}:{},
		{{RM:16},{R:16}}:{opcode: 0x89},
		{{RM:32},{R:32}}:{opcode: 0x89},
		{{RM:64},{R:64}}:{opcode: 0x89},
		{{R:8},{RM:8}}:{opcode: 0x8a},
		{{R:16},{RM:16}}:{opcode: 0x8b},
		{{R:32},{RM:32}}:{opcode: 0x8b},
		{{R:64},{RM:64}}:{opcode: 0x8b},
	},
}

var RET = &instruction{
	mnemonic: "ret",
	opcode: 0xc3,
	encoding: zeroOperand,
}

var allInstructions = []*instruction{
	INC,
	DEC,
}

func init() {
	for _, insn := range allInstructions {
		insn.initVariants()
	}
}
