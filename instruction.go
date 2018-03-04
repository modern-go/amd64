package amd64

type opcode byte
type instruction struct {
	Mnemonic string
	opcode   opcode // default op code if not override by specific addressing modes
	assemble interface{}
	rm8      opcode
}

func oneOperand(a *Assembler, insn *instruction, operand1 Operand) {
	switch operand1.Bits() {
	case 64:
		operand1.Rex(a, Register{})
		a.byte(byte(insn.opcode))
	case 32:
		a.byte(byte(insn.opcode))
	case 16:
		panic("not implemented")
	case 8:
		a.byte(byte(insn.rm8))
	}
	operand1.ModRM(a, Register{})
}
