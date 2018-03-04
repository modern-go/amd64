package amd64

type instruction struct {
	Mnemonic string
	opcode byte // default op code if not override by specific addressing modes
	assemble interface{}
}

func oneOperand(a *Assembler, insn *instruction, operand1 Operand) {
	if operand1.Bits() == 64 {
		operand1.Rex(a, Register{})
	}
	a.byte(insn.opcode)
	operand1.ModRM(a, Register{})
}