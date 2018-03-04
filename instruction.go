package amd64

import "errors"

type opcode byte

type variants map[[2]condition]*instruction
type overrides map[string]interface{}

type instruction struct {
	Mnemonic string
	// if not match variant, this opcode will be used by default
	opcode    opcode
	assemble  interface{}
	variants  map[[2]condition]*instruction
	overrides map[string]interface{}
}

type condition struct {
	r byte // only register, size
	m byte // only memory, size
	rm byte // register or register based memory, size
}

func (insn *instruction) initVariants() {
	for _, variant := range insn.variants {
		if variant.opcode == 0 {
			variant.opcode = insn.opcode
		}
		if variant.overrides == nil {
			variant.overrides = insn.overrides
		}
	}
}

func oneOperand(a *Assembler, insn *instruction, operand1 Operand) {
	variant := insn.variants[[2]condition{{rm: operand1.Bits()}}]
	if variant == nil {
		a.ReportError(errors.New("no variant defined for this operand combination"))
		return
	}
	insn = variant
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
