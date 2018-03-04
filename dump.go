package amd64

func Dump(instructions ...interface{}) string {
	var desc []byte
	isFirst := true
	for len(instructions) > 0 {
		if isFirst {
			isFirst = false
		} else {
			desc = append(desc, '\n')
		}
		insn, _ := instructions[0].(*Instruction)
		instructions = instructions[1:]
		if insn == nil {
			desc = append(desc, " %invalid%"...)
			continue
		}
		desc = append(desc, insn.Mnemonic...)
		operandsCount := 0
		switch insn.assemble.(type) {
		case func(a *Assembler):
		case func(a *Assembler, insn *Instruction):
		case func(a *Assembler, insn *Instruction, operand1 Operand):
			operandsCount = 1
		case func(a *Assembler, insn *Instruction, operand1 Operand, operand2 Operand):
			operandsCount = 2
		default:
			desc = append(desc, " %unknown%"...)
		}
		for i := 0; i < operandsCount; i++ {
			if i != 0 {
				desc = append(desc, ","...)
			}
			var operand Operand
			if len(instructions) > 0 {
				operand, _ = instructions[0].(Operand)
			}
			if operand == nil {
				desc = append(desc, " %miss%"...)
			} else {
				desc = append(desc, ' ')
				desc = append(desc, operand.String()...)
				instructions = instructions[1:]
			}
		}
	}
	return string(desc)
}
