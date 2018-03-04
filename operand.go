package amd64

type Operand interface {
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

type Register struct {
	Val  byte
	Bits byte
}

func (r Register) isOperand() {}
func (i Register) Rex(asm *Assembler, reg Register) {
	asm.rexBits(i.Bits, reg.Bits, reg.Val > 7, false, i.Val > 7)
}

func (r Register) ModRM(asm *Assembler, reg Register) {
	asm.modrm(MOD_REG, reg.Val&7, r.Val&7)
}


type Indirect struct {
	Base   Register
	Offset int32
	Bits   byte
}

func (i Indirect) short() bool {
	return int32(int8(i.Offset)) == i.Offset
}

func (i Indirect) isOperand() {}
func (i Indirect) Rex(asm *Assembler, reg Register) {
	asm.rexBits(reg.Bits, i.Bits, reg.Val > 7, false, i.Base.Val > 7)
}

func (i Indirect) ModRM(asm *Assembler, reg Register) {
	if i.Base.Val == REG_SIB {
		SIB{i.Offset, ESP, ESP, Scale1}.ModRM(asm, reg)
		return
	}
	if i.Offset == 0 {
		asm.modrm(MOD_INDIR, reg.Val&7, i.Base.Val&7)
	} else if i.short() {
		asm.modrm(MOD_INDIR_DISP8, reg.Val&7, i.Base.Val&7)
		asm.byte(byte(i.Offset))
	} else {
		asm.modrm(MOD_INDIR_DISP32, reg.Val&7, i.Base.Val&7)
		asm.int32(uint32(i.Offset))
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