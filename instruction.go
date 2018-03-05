package amd64

import "errors"

type opcode byte

type variants map[[2]VariantKey]*instruction
type overrides map[string]interface{}

type instruction struct {
	mnemonic string
	// if not match variant, this opcode will be used by default
	opcode opcode
	// secondary opcode
	opcode2 opcode
	// OpcodeReg is encoded as reg in modrm
	opcodeReg opcode
	assemble  interface{}
	variants  map[[2]VariantKey]*instruction
}

func (insn *instruction) Opcode() byte {
	return byte(insn.opcode)
}

func (insn *instruction) Variant(key [2]VariantKey) *instruction {
	return insn.variants[key]
}

func (insn *instruction) OpcodeReg() byte {
	return byte(insn.opcodeReg)
}

type VariantKey struct {
	R  byte // register, size
	M  byte // memory, size
	RM byte // register or memory, size
	XMM bool // xmm register
}

func (insn *instruction) initVariants() {
	for _, variant := range insn.variants {
		if variant.opcode == 0 {
			variant.opcode = insn.opcode
		}
		if variant.opcode2 == 0 {
			variant.opcode2 = insn.opcode2
		}
		if variant.opcodeReg == 0 {
			variant.opcodeReg = insn.opcodeReg
		}
	}
}

func (insn *instruction) findVariant(dst []VariantKey, src []VariantKey) *instruction {
	if src == nil {
		for _, c := range dst {
			variant := insn.variants[[2]VariantKey{c}]
			if variant != nil {
				return variant
			}
		}
		return nil
	}
	for _, s := range src {
		for _, d := range dst {
			variant := insn.variants[[2]VariantKey{d, s}]
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
	dst.Prefix(asm, nil)
	asm.byte(byte(insn.opcode))
	dst.Operands(asm, nil, insn.opcodeReg)
}

func twoOperands(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	variant := insn.findVariant(dst.Conditions(), src.Conditions())
	if variant == nil {
		asm.ReportError(errors.New("no variant defined for this operand combination"))
		return
	}
	insn = variant
	dst.Prefix(asm, Register{})
	asm.byte(byte(insn.opcode))
	srcRegister, _ := src.(Register)
	dst.Operands(asm, srcRegister, insn.opcodeReg)
}
