package amd64

type instruction struct {
	Mnemonic string
	opcode byte // default op code if not override by specific addressing modes
	assemble interface{}
}

func one_operand(a *Assembler, insn *instruction, operand1 Operand) {
	a.byte(insn.opcode)
	operand1.ModRM(a, Register{})
}