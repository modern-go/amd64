package amd64

import (
	"errors"
	"fmt"
)

const RegESP = 4
const RegEBP = 5
const Scale1 = 0
const Scale2 = 1
const Scale4 = 2
const Scale8 = 3

type Operand interface {
	fmt.Stringer
	// isOperand is unexported prevents external packages from
	// implementing Operand.
	isOperand()

	Prefix(asm *Assembler, src Operand)
	Operands(asm *Assembler, src Operand, opcodeReg opcode)
	Conditions() []VariantKey
	Bits() byte
}

type Immediate struct {
	val uint32
	bits byte
}

func (i Immediate) Prefix(asm *Assembler, src Operand) {
	panic("can not use immediate as dst operand")
}

func (i Immediate) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	panic("can not use immediate as dst operand")
}

func (i Immediate) isOperand() {
	panic("Imm.AssembleOperands")
}

func (i Immediate) String() string {
	return fmt.Sprintf("%v", i.val)
}

func (i Immediate) Bits() byte {
	return i.bits
}

func (i Immediate) Conditions() []VariantKey {
	return []VariantKey{
		{IMM: i.bits},
	}
}

type Register struct {
	desc       string
	val        byte
	bits       byte
	conditions []VariantKey
}

func (r Register) isOperand() {}

func (r Register) Prefix(asm *Assembler, src Operand) {
	switch r.bits {
	case 64:
		srcReg, _ := src.(Register)
		asm.byte(REX(r.bits == 64, srcReg.val > 7, false, r.val > 7))
	case 32:
	case 16:
		asm.byte(Prefix16Bit)
	case 8:
	default:
		asm.ReportError(errors.New("register size is invalid"))
		return
	}
}

func (r Register) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	srcReg, isSrcReg := src.(Register)
	if isSrcReg {
		asm.byte(MODRM(ModeReg, byte(srcReg.val), r.val&7))
	} else {
		asm.byte(MODRM(ModeReg, byte(regOpcode), r.val&7))
	}
}

func (r Register) String() string {
	return r.desc
}

func (r Register) Conditions() []VariantKey {
	return r.conditions
}

func (r Register) Value() byte {
	return r.val
}

func (r Register) Bits() byte {
	return r.bits
}

type Indirect struct {
	base       Register
	offset     int32
	bits       byte
	conditions []VariantKey
}

func (i Indirect) short() bool {
	return int32(int8(i.offset)) == i.offset
}

func (i Indirect) isOperand() {}

func (i Indirect) Prefix(asm *Assembler, src Operand) {
	switch i.base.bits {
	case 64:
	case 32:
		asm.byte(Prefix32Bit)
	default:
		asm.ReportError(errors.New("unsupported register"))
		return
	}
	switch i.bits {
	case 64:
		asm.byte(REX(i.bits == 64, false, false, i.base.val > 7))
	case 32:
	case 16:
		asm.byte(Prefix16Bit)
	case 8:
	default:
		asm.ReportError(errors.New("invalid size"))
		return
	}
}

func (i Indirect) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	if i.offset == 0 {
		if i.base.val == RegEBP {
			asm.byte(MODRM(ModeIndirDisp8, byte(regOpcode), i.base.val&7))
			asm.byte(0)
		} else {
			asm.byte(MODRM(ModeIndir, byte(regOpcode), i.base.val&7))
		}
	} else if i.short() {
		asm.byte(MODRM(ModeIndirDisp8, byte(regOpcode), i.base.val&7))
		asm.byte(byte(i.offset))
	} else {
		asm.byte(MODRM(ModeIndirDisp32, byte(regOpcode), i.base.val&7))
		asm.int32(uint32(i.offset))
	}
}

func (i Indirect) Conditions() []VariantKey {
	return i.conditions
}

func (i Indirect) Bits() byte {
	return i.bits
}

func (i Indirect) String() string {
	sizeDirective := ""
	switch i.bits {
	case 64:
		sizeDirective = "qword ptr"
	case 32:
		sizeDirective = "dword ptr"
	case 16:
		sizeDirective = "word ptr"
	case 8:
		sizeDirective = "byte ptr"
	default:
		sizeDirective = "invalid"
	}
	if i.offset >= 0 {
		return fmt.Sprintf("%s [%v+%v]", sizeDirective, i.base, i.offset)
	} else {
		return fmt.Sprintf("%s [%v%v]", sizeDirective, i.base, i.offset)
	}
}

type RipIndirect struct {
	Indirect
}

func (i RipIndirect) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	asm.byte(MODRM(ModeIndir, byte(regOpcode), RegEBP))
	asm.int32(uint32(i.offset))
}

type AbsoluteIndirect struct {
	Indirect
}

func (i AbsoluteIndirect) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	asm.byte(MODRM(ModeIndir, byte(regOpcode), RegESP))
	asm.byte(SIB(Scale1, RegESP, RegEBP))
	asm.int32(uint32(i.offset))
}

type ScaledIndirect struct {
	scale byte
	index Register
	Indirect
}

func (i ScaledIndirect) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	if i.offset == 0 {
		asm.byte(MODRM(ModeIndir, byte(regOpcode), RegESP))
		asm.byte(SIB(i.scale, i.index.val&7, i.base.val&7))
	} else if i.short() {
		asm.byte(MODRM(ModeIndirDisp8, byte(regOpcode), RegESP))
		asm.byte(SIB(i.scale, i.index.val&7, i.base.val&7))
		asm.byte(byte(i.offset))
	} else {
		asm.byte(MODRM(ModeIndirDisp32, byte(regOpcode), RegESP))
		asm.byte(SIB(i.scale, i.index.val&7, i.base.val&7))
		asm.int32(uint32(i.offset))
	}
}

func (i ScaledIndirect) String() string {
	var desc []byte
	switch i.bits {
	case 64:
		desc = append(desc, "qword ptr ["...)
	case 32:
		desc = append(desc, "dword ptr ["...)
	case 16:
		desc = append(desc, "word ptr ["...)
	case 8:
		desc = append(desc, "byte ptr ["...)
	default:
		desc = append(desc, "invalid ["...)
	}
	scale := 0
	if i.index.val != RegESP {
		scale = 1 << i.scale
	}
	if scale == 0 {
		// skip
	} else if scale == 1 {
		desc = append(desc, fmt.Sprintf("%v+", i.index)...)
	} else {
		desc = append(desc, fmt.Sprintf("%v*%v+", i.index, scale)...)
	}
	if i.base.val != RegEBP {
		desc = append(desc, i.base.String()...)
	}
	if i.offset > 0 {
		desc = append(desc, fmt.Sprintf("+%v", i.offset)...)
	} else if i.offset < 0 {
		desc = append(desc, fmt.Sprintf("%v", i.offset)...)
	}
	desc = append(desc, ']')
	return string(desc)
}
