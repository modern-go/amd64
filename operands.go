package amd64

func QWORD(base Register, offset int) Indirect {
	return Indirect{
		Base:   base,
		Offset: int32(offset),
	}
}

var (
	EAX = Register{0, 32}
	RAX = Register{0, 64}
	ECX = Register{1, 32}
	RCX = Register{1, 64}
	EDX = Register{2, 32}
	RDX = Register{2, 64}
	EBX = Register{3, 32}
	RBX = Register{3, 64}
	ESP = Register{4, 32}
	RSP = Register{4, 64}
	EBP = Register{5, 32}
	RBP = Register{5, 64}
	ESI = Register{6, 32}
	RSI = Register{6, 64}
	EDI = Register{7, 32}
	RDI = Register{7, 64}

	R8D  = Register{8, 32}
	R8   = Register{8, 64}
	R9D  = Register{9, 32}
	R9   = Register{9, 64}
	R10D = Register{10, 32}
	R10  = Register{10, 64}
	R11D = Register{11, 32}
	R11  = Register{11, 64}
	R12D = Register{12, 32}
	R12  = Register{12, 64}
	R13D = Register{13, 32}
	R13  = Register{13, 64}
	R14D = Register{14, 32}
	R14  = Register{14, 64}
	R15D = Register{15, 32}
	R15  = Register{15, 64}
)