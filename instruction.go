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
	r byte // register, size
	m byte // memory, size
	rm byte // register or memory, size
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

func (insn *instruction) findVariant(dst []condition, src []condition) *instruction {
	if src == nil {
		for _, c := range dst {
			variant := insn.variants[[2]condition{c}]
			if variant != nil {
				return variant
			}
		}
		return nil
	}
	for _, s := range src {
		for _, d := range dst {
			variant := insn.variants[[2]condition{d, s}]
			if variant != nil {
				return variant
			}
		}
	}
	return nil
}

func oneOperand(asm *Assembler, insn *instruction, dst Operand) {
	variant := insn.findVariant(dst.Conditions(), nil)
	if variant == nil {
		asm.ReportError(errors.New("no variant defined for this operand combination"))
		return
	}
	insn = variant
	dst.Prefix(asm, Register{})
	asm.byte(byte(insn.opcode))
	modrmRm, _ := insn.overrides["ModR/M r/m"].(byte)
	dst.ModRM(asm, Register{val: modrmRm})
}

func twoOperands(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	variant := insn.findVariant(dst.Conditions(), src.Conditions())
	if variant == nil {
		asm.ReportError(errors.New("no variant defined for this operand combination"))
		return
	}
	insn = variant
	src.Prefix(asm, Register{})
	asm.byte(byte(insn.opcode))
	dstRegister, _ := dst.(Register)
	src.ModRM(asm, dstRegister)
}