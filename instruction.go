package amd64

import (
	"errors"
	"fmt"
)

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
	encoding  interface{}
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
	R   byte   // register, size
	M   byte   // memory, size
	RM  byte   // register or memory, size
	IMM byte   // immediate, size
	REG string // special register
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

func (insn *instruction) findVariant(asm *Assembler, dst []VariantKey, src []VariantKey) *instruction {
	if src == nil {
		for _, c := range dst {
			variant := insn.variants[[2]VariantKey{c}]
			if variant != nil {
				return variant
			}
		}
		asm.ReportError(errors.New("no variant defined for this operand combination"))
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
	asm.ReportError(fmt.Errorf(
		"no variant defined for this operand combination, dst: %v, src: %v",
		dst,
		src))
	return nil
}

func zeroOperand(asm *Assembler, insn *instruction) {
	asm.byte(byte(insn.opcode))
}

func oneOperand(asm *Assembler, insn *instruction, dst Operand) {
	variant := insn.findVariant(asm, dst.variantKeys(), nil)
	if variant == nil {
		return
	}
	insn = variant
	dst.prefix(asm, nil)
	asm.byte(byte(insn.opcode))
	dst.operands(asm, nil, encodingParams{
		opcodeReg: insn.opcodeReg,
	})
}

func twoOperands(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	variant := insn.findVariant(asm, dst.variantKeys(), src.variantKeys())
	if variant == nil {
		return
	}
	insn = variant
	if insn.encoding != nil {
		encode := insn.encoding.(func(*Assembler, *instruction, Operand, Operand))
		encode(asm, insn, dst, src)
		return
	}
	if src.isMemory() {
		// memory can only be encoded as dst
		src, dst = dst, src
	}
	dst.prefix(asm, src)
	asm.byte(byte(insn.opcode))
	dst.operands(asm, src, encodingParams{
		opcodeReg: insn.opcodeReg,
	})
}

// without MODRM
func encodingI(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	dst.prefix(asm, src)
	asm.byte(byte(insn.opcode))
	dst.operands(asm, src, encodingParams{
		opcodeReg:    insn.opcodeReg,
		withoutMODRM: true,
	})
}
