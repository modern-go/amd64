package amd64

import (
	"errors"
	"fmt"
)

type opcode byte

type VariantKey [3]Qualifier
type variants map[VariantKey]*instruction

type vexForm byte

var form0F vexForm = 1
var formVEX2 vexForm = 2
var formVEX3 vexForm = 3

type instruction struct {
	mnemonic string
	// if not match variant, this opcode will be used by default
	opcode opcode
	// secondary opcode
	opcode2 opcode
	vexForm vexForm
	vexPP   byte
	// OpcodeReg is encoded as reg in modrm
	opcodeReg opcode
	encoding  interface{}
	variants  map[VariantKey]*instruction
}

func (insn *instruction) Opcode() byte {
	return byte(insn.opcode)
}

func (insn *instruction) Prefix0F() byte {
	if insn.vexForm == form0F {
		return 0x0f
	}
	return 0x00
}

func (insn *instruction) PrefixC5() byte {
	if insn.vexForm == formVEX2 {
		return 0xc5
	}
	return 0x00
}

func (insn *instruction) Variant(key VariantKey) *instruction {
	return insn.variants[key]
}

func (insn *instruction) OpcodeReg() byte {
	return byte(insn.opcodeReg)
}

type Qualifier struct {
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
		if variant.vexForm == 0 {
			variant.vexForm = insn.vexForm
		}
		if variant.vexPP == 0 {
			variant.vexPP = insn.vexPP
		}
	}
}

func (insn *instruction) findVariant(
	asm *Assembler, dst []Qualifier, src1 []Qualifier, src2 []Qualifier) (*instruction, VariantKey) {
	if src1 == nil && src2 == nil {
		for _, c := range dst {
			key := VariantKey{c}
			variant := insn.variants[key]
			if variant != nil {
				return variant, key
			}
		}
		asm.ReportError(errors.New("no variant defined for this operand combination"))
		return nil, VariantKey{}
	}
	if src2 == nil {
		for _, d := range dst {
			for _, s := range src1 {
				key := VariantKey{d, s}
				variant := insn.variants[key]
				if variant != nil {
					return variant, key
				}
			}
		}
		asm.ReportError(fmt.Errorf(
			"no variant defined for this operand combination, dst: %v, src: %v",
			dst,
			src1))
		return nil, VariantKey{}
	}
	for _, d := range dst {
		for _, s1 := range src1 {
			for _, s2 := range src2 {
				key := VariantKey{d, s1, s2}
				variant := insn.variants[key]
				if variant != nil {
					return variant, key
				}
			}
		}
	}
	asm.ReportError(fmt.Errorf(
		"no variant defined for this operand combination, dst: %v, src1: %v, src2: %v",
		dst, src1, src2))
	return nil, VariantKey{}
}
