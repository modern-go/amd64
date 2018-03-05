package amd64

func zeroOperand(asm *Assembler, insn *instruction) {
	asm.byte(byte(insn.opcode))
}

func oneOperand(asm *Assembler, insn *instruction, dst Operand) {
	variant, _ := insn.findVariant(asm, dst.variantKeys(), nil)
	if variant == nil {
		return
	}
	insn = variant
	encoding, _ := insn.encoding.(func(*Assembler, *instruction, Operand))
	if encoding == nil {
		encoding = encodingM
	}
	encoding(asm, insn, dst)
}

func twoOperands(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	variant, key := insn.findVariant(asm, dst.variantKeys(), src.variantKeys())
	if variant == nil {
		return
	}
	insn = variant
	encoding, _ := insn.encoding.(func(*Assembler, *instruction, Operand, Operand))
	if encoding == nil {
		if key[0].RM > 0 {
			encoding = encodingMR
		} else {
			encoding = encodingRM
		}
		if key[1].IMM > 0 {
			encoding = encodingMI
		}
	}
	encoding(asm, insn, dst, src)
}

// dst: rm
func encodingM(asm *Assembler, insn *instruction, dst Operand) {
	dst.rex(asm, nil)
	asm.byte(byte(insn.opcode))
	dst.modrm(asm, byte(insn.opcodeReg))
}

// dst: rm
// src: imm
func encodingMI(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	dst.rex(asm, src)
	asm.byte(byte(insn.opcode))
	dst.modrm(asm, 0)
	asm.imm(src.(Immediate))
}

// dst: rm
// src: reg
func encodingMR(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	dst.rex(asm, src)
	asm.byte(byte(insn.opcode))
	dst.modrm(asm, src.(Register).Value())
}

// dst: reg
// src: rm
func encodingRM(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	src.rex(asm, dst)
	asm.byte(byte(insn.opcode))
	src.modrm(asm, dst.(Register).Value())
}

// dst: al/ax/eax/rax
// src: imm8/16/32
func encodingI(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	dst.rex(asm, src)
	asm.byte(byte(insn.opcode))
	asm.imm(src.(Immediate))
}

// dst: reg
// src: rm
func encodingA(asm *Assembler, insn *instruction, dst Operand, src Operand) {
	src.rex(asm, dst)
	if insn.prefix0F {
		asm.byte(0x0f)
	}
	asm.byte(byte(insn.opcode))
	src.modrm(asm, dst.(Register).Value())
}