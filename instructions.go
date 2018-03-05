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
		{{RM: 8}, {R: 8}}:         {},
		{{RM: 16}, {R: 16}}:       {opcode: 0x01},
		{{RM: 32}, {R: 32}}:       {opcode: 0x01},
		{{RM: 64}, {R: 64}}:       {opcode: 0x01},
		{{R: 8}, {RM: 8}}:         {opcode: 0x02},
		{{R: 16}, {RM: 16}}:       {opcode: 0x03},
		{{R: 32}, {RM: 32}}:       {opcode: 0x03},
		{{R: 64}, {RM: 64}}:       {opcode: 0x03},
		{{REG: "al"}, {IMM: 8}}:   {opcode: 0x04, encoding: encodingI}, // 0x04 is shorter than 0x80 form
		{{REG: "ax"}, {IMM: 16}}:  {opcode: 0x05, encoding: encodingI},
		{{REG: "eax"}, {IMM: 32}}: {opcode: 0x05, encoding: encodingI},
		{{RM: 8}, {IMM: 8}}:       {opcode: 0x80},
		{{RM: 16}, {IMM: 16}}:     {opcode: 0x81},
		{{RM: 32}, {IMM: 32}}:     {opcode: 0x81},
		{{RM: 64}, {IMM: 64}}:     {opcode: 0x81},
		{{RM: 16}, {IMM: 8}}:      {opcode: 0x83},
		{{RM: 32}, {IMM: 8}}:      {opcode: 0x83},
		{{RM: 64}, {IMM: 8}}:      {opcode: 0x83},
	},
}

var MOV = &instruction{
	mnemonic: "mov",
	opcode:   0x88,
	encoding: twoOperands,
	variants: variants{
		{{RM: 8}, {R: 8}}:   {},
		{{RM: 16}, {R: 16}}: {opcode: 0x89},
		{{RM: 32}, {R: 32}}: {opcode: 0x89},
		{{RM: 64}, {R: 64}}: {opcode: 0x89},
		{{R: 8}, {RM: 8}}:   {opcode: 0x8a},
		{{R: 16}, {RM: 16}}: {opcode: 0x8b},
		{{R: 32}, {RM: 32}}: {opcode: 0x8b},
		{{R: 64}, {RM: 64}}: {opcode: 0x8b},
	},
}

var MOVAPS = &instruction{
	mnemonic: "movaps",
	prefix0F: true,
	opcode:   0x28,
	encoding: twoOperands,
	variants: variants{
		{{REG:"xmm"},{REG:"xmm",M:128}}:{encoding:encodingA},
		{{REG:"xmm",M:128},{REG:"xmm"}}:{opcode: 0x29, encoding:encodingB},
	},
}

var RET = &instruction{
	mnemonic: "ret",
	opcode:   0xc3,
	encoding: zeroOperand,
}

var allInstructions = []*instruction{
	INC,
	DEC,
	ADD,
	MOV,
	RET,
	MOVAPS,
}

func init() {
	for _, insn := range allInstructions {
		insn.initVariants()
	}
}
