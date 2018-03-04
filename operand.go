package amd64

import (
	"fmt"
	"errors"
)

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

type Imm struct {
	Val int32
}

func U32(u uint32) int32 {
	return int32(u)
}

func (i Imm) isOperand() {}
func (i Imm) Rex(asm *Assembler, reg Register) {
	panic("Imm.Rex")
}
func (i Imm) ModRM(asm *Assembler, reg Register) {
	panic("Imm.AssembleOperands")
}

func (i Imm) String() string {
	return fmt.Sprintf("%v", i.Val)
}

type Register struct {
	desc string
	val  byte
	bits byte
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
	base   Register
	offset int32
	bits   byte
	conditions []VariantKey
}

func (i Indirect) short() bool {
	return int32(int8(i.offset)) == i.offset
}

func (i Indirect) isOperand() {}

func (i Indirect) Prefix(asm *Assembler, src Operand) {
	switch i.bits {
	case 64:
		switch i.base.bits {
		case 64:
		case 32:
			asm.byte(Prefix32Bit)
		default:
			asm.ReportError(errors.New("unsupported register"))
			return
		}
		asm.byte(REX(i.bits == 64, false, false, i.base.val > 7))
	case 32:
		panic("not implemented")
	case 16:
		panic("not implemented")
	case 8:
		panic("not implemented")
	}
}

func (i Indirect) Operands(asm *Assembler, src Operand, regOpcode opcode) {
	//if i.base.val == REG_SIB {
	//	SIB{i.offset, ESP, ESP, Scale1}.AssembleOperands(asm, reg)
	//	return
	//}
	if i.offset == 0 {
		asm.byte(MODRM(ModeIndir, byte(regOpcode), i.base.val&7))
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

//type PCRel struct {
//	Addr uintptr
//}
//
//func (i PCRel) isOperand() {}
//func (i PCRel) Rex(asm *Assembler, reg Register) {
//	asm.rex(reg.bits == 64, reg.val > 7, false, false)
//}
//func (i PCRel) ModRM(asm *Assembler, reg Register) {
//	asm.modrm(ModeIndir, reg.val&7, REG_DISP32)
//	asm.rel32(i.Addr)
//}

//type Scale struct {
//	scale byte
//}
//
//var (
//	Scale1 = Scale{SCALE_1}
//	Scale2 = Scale{SCALE_2}
//	Scale4 = Scale{SCALE_4}
//	Scale8 = Scale{SCALE_8}
//)

//type SIB struct {
//	Offset      int32
//	Base, Index Register
//	Scale       Scale
//}
//
//func (s SIB) isOperand() {}
//func (s SIB) Rex(asm *Assembler, reg Register) {
//	asm.rex(reg.bits == 64, reg.val > 7, s.Index.val > 7, s.Base.val > 7)
//}
//
//func (s SIB) short() bool {
//	return int32(int8(s.Offset)) == s.Offset
//}
//
//func (s SIB) AssembleOperands(asm *Assembler, reg Register) {
//	if s.Offset != 0 {
//		if s.short() {
//			asm.modrm(ModeIndirDisp8, reg.val&7, REG_SIB)
//			asm.sib(s.Scale.scale, s.Index.val&7, s.Base.val&7)
//			asm.byte(uint8(s.Offset))
//		} else {
//			asm.modrm(ModeIndirDisp32, reg.val&7, REG_SIB)
//			asm.sib(s.Scale.scale, s.Index.val&7, s.Base.val&7)
//			asm.int32(uint32(s.Offset))
//		}
//	} else {
//		asm.modrm(ModeIndir, reg.val&7, REG_SIB)
//		asm.sib(s.Scale.scale, s.Index.val&7, s.Base.val&7)
//	}
//}