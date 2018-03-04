package amd64

const (
	ModeIndir       = 0x0
	ModeIndirDisp8  = 0x1
	ModeIndirDisp32 = 0x2
	ModeReg         = 0x3
)

func MODRM(mod byte, reg byte, rm byte) byte {
	return mod << 6 | reg << 3 | rm
}

func SIB(scale byte, index byte, base byte) byte {
	return scale << 6 | index << 3 | base
}
