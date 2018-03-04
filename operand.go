package amd64

import "fmt"

type Operand interface {
	fmt.Stringer
	// isOperand is unexported prevents external packages from
	// implementing Operand.
	isOperand()

	Rex(asm *Assembler, reg Register)
	ModRM(asm *Assembler, reg Register)
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
	panic("Imm.ModRM")
}

func (i Imm) String() string {
	return fmt.Sprintf("%v", i.Val)
}

type Register struct {
	Desc string
	Val  byte
	Bits byte
}

func (r Register) isOperand() {}
func (r Register) Rex(asm *Assembler, reg Register) {
	asm.rexBits(r.Bits, reg.Bits, reg.Val > 7, false, r.Val > 7)
}

func (r Register) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_REG, reg.Val&7, r.Val&7)
}

func (r Register) String() string {
	return r.Desc
}


type Indirect struct {
	base   Register
	offset int32
	bits   byte
}

func (i Indirect) short() bool {
	return int32(int8(i.offset)) == i.offset
}

func (i Indirect) isOperand() {}
func (i Indirect) Rex(asm *Assembler, reg Register) {
	asm.rexBits(reg.Bits, i.bits, reg.Val > 7, false, i.base.Val > 7)
}

func (i Indirect) ModRM(asm *Assembler, reg Register) {
	if i.base.Val == REG_SIB {
		SIB{i.offset, ESP, ESP, Scale1}.ModRM(asm, reg)
		return
	}
	if i.offset == 0 {
		asm.modrm(MOD_INDIR, reg.Val&7, i.base.Val&7)
	} else if i.short() {
		asm.modrm(MOD_INDIR_DISP8, reg.Val&7, i.base.Val&7)
		asm.byte(byte(i.offset))
	} else {
		asm.modrm(MOD_INDIR_DISP32, reg.Val&7, i.base.Val&7)
		asm.int32(uint32(i.offset))
	}
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

type PCRel struct {
	Addr uintptr
}

func (i PCRel) isOperand() {}
func (i PCRel) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, false, false)
}
func (i PCRel) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_INDIR, reg.Val&7, REG_DISP32)
	asm.rel32(i.Addr)
}

type Scale struct {
	scale byte
}

var (
	Scale1 = Scale{SCALE_1}
	Scale2 = Scale{SCALE_2}
	Scale4 = Scale{SCALE_4}
	Scale8 = Scale{SCALE_8}
)

type SIB struct {
	Offset      int32
	Base, Index Register
	Scale       Scale
}

func (s SIB) isOperand() {}
func (s SIB) Rex(asm *Assembler, reg Register) {
	asm.rex(reg.Bits == 64, reg.Val > 7, s.Index.Val > 7, s.Base.Val > 7)
}

func (s SIB) short() bool {
	return int32(int8(s.Offset)) == s.Offset
}

func (s SIB) ModRM(asm *Assembler, reg Register) {
	if s.Offset != 0 {
		if s.short() {
			asm.modrm(MOD_INDIR_DISP8, reg.Val&7, REG_SIB)
			asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
			asm.byte(uint8(s.Offset))
		} else {
			asm.modrm(MOD_INDIR_DISP32, reg.Val&7, REG_SIB)
			asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
			asm.int32(uint32(s.Offset))
		}
	} else {
		asm.modrm(MOD_INDIR, reg.Val&7, REG_SIB)
		asm.sib(s.Scale.scale, s.Index.Val&7, s.Base.Val&7)
	}
}