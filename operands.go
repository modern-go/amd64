package amd64

func QWORD(base Register, offset int) Indirect {
	return Indirect{
		base:   base,
		offset: int32(offset),
		bits:   64,
	}
}

var (
	EAX = Register{"eax", 0, 32}
	RAX = Register{"rax", 0, 64}
	AL = Register{"al", 0, 8}
	ECX = Register{"ecx", 1, 32}
	RCX = Register{"rcx", 1, 64}
	EDX = Register{"edx", 2, 32}
	RDX = Register{"rdx", 2, 64}
	EBX = Register{"ebx", 3, 32}
	RBX = Register{"rbx", 3, 64}
	ESP = Register{"esp", 4, 32}
	RSP = Register{"rsp", 4, 64}
	EBP = Register{"ebp", 5, 32}
	RBP = Register{"rbp", 5, 64}
	ESI = Register{"esi", 6, 32}
	RSI = Register{"rsi", 6, 64}
	EDI = Register{"edi", 7, 32}
	RDI = Register{"rdi", 7, 64}

	R8D  = Register{"r8d", 8, 32}
	R8   = Register{"r8", 8, 64}
	R9D  = Register{"r9d", 9, 32}
	R9   = Register{"r9", 9, 64}
	R10D = Register{"r10d", 10, 32}
	R10  = Register{"r10", 10, 64}
	R11D = Register{"r11d", 11, 32}
	R11  = Register{"r11", 11, 64}
	R12D = Register{"r12d", 12, 32}
	R12  = Register{"r12", 12, 64}
	R13D = Register{"r13d", 13, 32}
	R13  = Register{"r13", 13, 64}
	R14D = Register{"r14d", 14, 32}
	R14  = Register{"r14", 14, 64}
	R15D = Register{"r15d", 15, 32}
	R15  = Register{"r15", 15, 64}
)
