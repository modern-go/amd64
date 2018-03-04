package amd64

type opcode byte

type variants map[condition]*instruction
type overrides map[string]interface{}

type instruction struct {
	Mnemonic string
	// if not match variant, this opcode will be used by default
	opcode    opcode
	assemble  interface{}
	variants  map[condition]*instruction
	overrides map[string]interface{}
}

type condition struct {
	rm byte // register or memory size
}

func oneOperand(a *Assembler, insn *instruction, operand1 Operand) {
	variant := insn.variants[condition{rm: operand1.Bits()}]
	if variant != nil {
		insn = variant
	}
	switch operand1.Bits() {
	case 64:
		operand1.Rex(a, Register{})
	case 16:
		a.byte(0x66)
	}
	a.byte(byte(insn.opcode))
	modrmRm, _ := insn.overrides["ModR/M r/m"].(byte)
	operand1.ModRM(a, Register{val: modrmRm})
}
