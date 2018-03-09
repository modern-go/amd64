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

	isMemory() bool
	rex(asm *Assembler, reg Operand)
	vex(asm *Assembler, insn *instruction, reg Operand)
	modrm(asm *Assembler, reg byte)
	Qualifiers() []Qualifier
	Bits() byte
}

type encodingParams struct {
	opcodeReg    opcode
	withoutMODRM bool
}

type Immediate struct {
	val        uint32
	bits       byte
	qualifiers []Qualifier
}

func (i Immediate) rex(asm *Assembler, reg Operand) {
	panic("can not use immediate as dst operand")
}

func (i Immediate) vex(asm *Assembler, insn *instruction, reg Operand) {
	panic("can not use immediate as dst operand")
}

func (i Immediate) modrm(asm *Assembler, reg byte) {
	panic("can not use immediate as dst operand")
}

func (i Immediate) isMemory() bool {
	return false
}

func (i Immediate) String() string {
	return fmt.Sprintf("%v", i.val)
}

func (i Immediate) Bits() byte {
	return i.bits
}

func (i Immediate) Qualifiers() []Qualifier {
	return i.qualifiers
}

type Register struct {
	desc       string
	val        byte
	bits       byte
	qualifiers []Qualifier
}

func (r Register) isMemory() bool { return false }

func (r Register) rex(asm *Assembler, reg Operand) {
	switch r.bits {
	case 128:
	case 64:
		register, _ := reg.(Register)
		asm.byte(REX(r.bits == 64, register.val > 7, false, r.val > 7))
	case 32:
	case 16:
		asm.byte(Prefix16Bit)
	case 8:
	default:
		asm.ReportError(errors.New("register size is invalid"))
		return
	}
}

func (r Register) vex(asm *Assembler, insn *instruction, reg Operand) {
	switch insn.vexForm {
	case 0:
	case form0F:
		asm.byte(0x0f)
	case formVEX2:
		asm.byte(0xc5)
		asm.byte(VEX2(0, 0, 0, insn.vexPP))
	case formVEX3:
		asm.byte(0xc4)
		asm.byte(VEX31(0, 0, 0, 2))
		asm.byte(VEX32(0, 0, 0, insn.vexPP))
	default:
		asm.ReportError(fmt.Errorf("unknown vex form: %v", insn.vexForm))
		return
	}
}

func (r Register) modrm(asm *Assembler, reg byte) {
	asm.byte(MODRM(ModeReg, reg, r.val&7))
}

func (r Register) Qualifiers() []Qualifier {
	return r.qualifiers
}

func (r Register) Bits() byte {
	return r.bits
}

func (r Register) Value() byte {
	return r.val
}

func (r Register) String() string {
	return r.desc
}

type Indirect struct {
	base       Register
	offset     int32
	bits       byte
	qualifiers []Qualifier
}

func (i Indirect) short() bool {
	return int32(int8(i.offset)) == i.offset
}

func (i Indirect) isMemory() bool {
	return true
}

func (i Indirect) rex(asm *Assembler, reg Operand) {
	switch i.base.bits {
	case 128:
	case 64:
	case 32:
		asm.byte(Prefix32Bit)
	default:
		asm.ReportError(errors.New("unsupported register"))
		return
	}
	switch i.bits {
	case 128:
	case 64:
		register, _ := reg.(Register)
		asm.byte(REX(i.bits == 64, register.val > 7, false, i.base.val > 7))
	case 32:
	case 16:
		asm.byte(Prefix16Bit)
	case 8:
	default:
		asm.ReportError(errors.New("invalid size"))
		return
	}
}

func (i Indirect) vex(asm *Assembler, insn *instruction, reg Operand) {
	switch insn.vexForm {
	case 0:
	case form0F:
		asm.byte(0x0f)
	case formVEX2:
		asm.byte(0xc5)
		asm.byte(VEX2(0, 0, 0, insn.vexPP))
	default:
		asm.ReportError(fmt.Errorf("unknown vex form: %v", insn.vexForm))
		return
	}
}

func (i Indirect) modrm(asm *Assembler, reg byte) {
	if i.offset == 0 {
		if i.base.val == RegEBP {
			asm.byte(MODRM(ModeIndirDisp8, reg, i.base.val&7))
			asm.byte(0)
		} else {
			asm.byte(MODRM(ModeIndir, reg, i.base.val&7))
		}
	} else if i.short() {
		asm.byte(MODRM(ModeIndirDisp8, reg, i.base.val&7))
		asm.byte(byte(i.offset))
	} else {
		asm.byte(MODRM(ModeIndirDisp32, reg, i.base.val&7))
		asm.int32(uint32(i.offset))
	}
}

func (i Indirect) Qualifiers() []Qualifier {
	return i.qualifiers
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

func (i RipIndirect) modrm(asm *Assembler, reg byte) {
	asm.byte(MODRM(ModeIndir, reg, RegEBP))
	asm.int32(uint32(i.offset))
}

type AbsoluteIndirect struct {
	Indirect
}

func (i AbsoluteIndirect) modrm(asm *Assembler, reg byte) {
	asm.byte(MODRM(ModeIndir, reg, RegESP))
	asm.byte(SIB(Scale1, RegESP, RegEBP))
	asm.int32(uint32(i.offset))
}

type ScaledIndirect struct {
	scale byte
	index Register
	Indirect
}

func (i ScaledIndirect) rex(asm *Assembler, reg Operand) {
	switch i.base.bits {
	case 128:
	case 64:
	case 32:
		asm.byte(Prefix32Bit)
	default:
		asm.ReportError(errors.New("unsupported register"))
		return
	}
	switch i.bits {
	case 128:
	case 64:
		register, _ := reg.(Register)
		asm.byte(REX(i.bits == 64, register.val > 7, i.index.val > 7, i.base.val > 7))
	case 32:
	case 16:
		asm.byte(Prefix16Bit)
	case 8:
	default:
		asm.ReportError(errors.New("invalid size"))
		return
	}
}

func (i ScaledIndirect) modrm(asm *Assembler, reg byte) {
	if i.offset == 0 {
		asm.byte(MODRM(ModeIndir, reg, RegESP))
		asm.byte(SIB(i.scale, i.index.val&7, i.base.val&7))
	} else if i.short() {
		asm.byte(MODRM(ModeIndirDisp8, reg, RegESP))
		asm.byte(SIB(i.scale, i.index.val&7, i.base.val&7))
		asm.byte(byte(i.offset))
	} else {
		asm.byte(MODRM(ModeIndirDisp32, reg, RegESP))
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
