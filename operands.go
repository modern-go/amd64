package amd64

func QWORD(base Register, offset int) interface{} {
	return newIndirect(64, base, offset)
}

func QWORD_SIB(scale byte, index Register, base Register, offset int) interface{} {
	return newSibIndirect(64, scale, index, base, offset)
}

func DWORD(base Register, offset int) interface{} {
	return newIndirect(32, base, offset)
}

func DWORD_SIB(scale byte, index Register, base Register, offset int) interface{} {
	return newSibIndirect(32, scale, index, base, offset)
}

func WORD(base Register, offset int) interface{} {
	return newIndirect(16, base, offset)
}

func WORD_SIB(scale byte, index Register, base Register, offset int) interface{} {
	return newSibIndirect(16, scale, index, base, offset)
}

func BYTE(base Register, offset int) interface{} {
	return newIndirect(8, base, offset)
}

func BYTE_SIB(scale byte, index Register, base Register, offset int) interface{} {
	return newSibIndirect(8, scale, index, base, offset)
}

func newIndirect(bits byte, base Register, offset int) interface{} {
	if base.val == RegESP {
		return DWORD_SIB(0, base, base, offset)
	}
	indirect := Indirect{
		base:   base,
		offset: int32(offset),
		bits:   bits,
		conditions: []VariantKey{{
			M: bits,
		}, {
			RM: bits,
		}},
	}
	switch base.desc {
	case RIP.desc:
		return RipIndirect{indirect}
	case ABSOLUTE.desc:
		return AbsoluteIndirect{indirect}
	}
	return indirect
}

func newSibIndirect(bits byte, scale byte, index Register, base Register, offset int) interface{} {
	switch scale {
	case 0:
		if index.val != RegESP {
			panic("scale 0 can only applied to esp")
		}
		if base.val != RegESP {
			return newIndirect(bits, base, offset)
		}
		scale = Scale1
	case 1:
		scale = Scale1
	case 2:
		scale = Scale2
	case 4:
		scale = Scale4
	case 8:
		scale = Scale8
	default:
		panic("invalid scale")
	}
	return ScaledIndirect{
		scale: scale,
		index: index,
		Indirect: Indirect{
			base:   base,
			offset: int32(offset),
			bits:   bits,
			conditions: []VariantKey{{
				M: bits,
			}, {
				RM: bits,
			}},
		},
	}
}

var RIP = Register{desc: "RIP", bits: 64}
var ABSOLUTE = Register{desc: "ABSOLUTE", bits: 64}

var (
	AL = Register{"al", 0, 8, []VariantKey{
		{R: 8},
		{RM: 8},
	}}
	AX = Register{"ax", 0, 16, []VariantKey{
		{R: 16},
		{RM: 16},
	}}
	EAX = Register{"eax", 0, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RAX = Register{"rax", 0, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	ECX = Register{"ecx", 1, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RCX = Register{"rcx", 1, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	EDX = Register{"edx", 2, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RDX = Register{"rdx", 2, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	BL = Register{"bl", 3, 8, []VariantKey{
		{R: 8},
		{RM: 8},
	}}
	BX = Register{"bl", 3, 16, []VariantKey{
		{R: 16},
		{RM: 16},
	}}
	EBX = Register{"ebx", 3, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RBX = Register{"rbx", 3, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	ESP = Register{"esp", 4, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RSP = Register{"rsp", 4, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	EBP = Register{"ebp", 5, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RBP = Register{"rbp", 5, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	ESI = Register{"esi", 6, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RSI = Register{"rsi", 6, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	EDI = Register{"edi", 7, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	RDI = Register{"rdi", 7, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}

	R8D = Register{"r8d", 8, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R8 = Register{"r8", 8, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R9D = Register{"r9d", 9, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R9 = Register{"r9", 9, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R10D = Register{"r10d", 10, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R10 = Register{"r10", 10, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R11D = Register{"r11d", 11, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R11 = Register{"r11", 11, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R12D = Register{"r12d", 12, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R12 = Register{"r12", 12, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R13D = Register{"r13d", 13, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R13 = Register{"r13", 13, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R14D = Register{"r14d", 14, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R14 = Register{"r14", 14, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
	R15D = Register{"r15d", 15, 32, []VariantKey{
		{R: 32},
		{RM: 32},
	}}
	R15 = Register{"r15", 15, 64, []VariantKey{
		{R: 64},
		{RM: 64},
	}}
)
