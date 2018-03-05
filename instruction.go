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
	opcode2  opcode
	prefix0F bool
	// OpcodeReg is encoded as reg in modrm
	opcodeReg opcode
	encoding  interface{}
	variants  map[[2]VariantKey]*instruction
}

func (insn *instruction) Opcode() byte {
	return byte(insn.opcode)
}

func (insn *instruction) Prefix0F() byte {
	if insn.prefix0F {
		return 0x0f
	}
	return 0x00
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
		if insn.prefix0F {
			variant.prefix0F = true
		}
	}
}

func (insn *instruction) findVariant(asm *Assembler, dst []VariantKey, src []VariantKey) (*instruction, [2]VariantKey) {
	if src == nil {
		for _, c := range dst {
			key := [2]VariantKey{c}
			variant := insn.variants[key]
			if variant != nil {
				return variant, key
			}
		}
		asm.ReportError(errors.New("no variant defined for this operand combination"))
		return nil, [2]VariantKey{}
	}
	for _, d := range dst {
		for _, s := range src {
			key := [2]VariantKey{d, s}
			variant := insn.variants[key]
			if variant != nil {
				return variant, key
			}
		}
	}
	asm.ReportError(fmt.Errorf(
		"no variant defined for this operand combination, dst: %v, src: %v",
		dst,
		src))
	return nil, [2]VariantKey{}
}
