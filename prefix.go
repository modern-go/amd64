package amd64

// If a 66H override is used with REX and REX.W = 0, the operand size is 16 bits.
const Prefix16Bit = 0x66

// In 64-bit mode, the instruction’s default address size is 64 bits, 32 bit address size is supported using the prefix 67H.
const Prefix32Bit = 0x67

func REX(w, r, x, b bool) byte {
	bits := byte(0x40)
	if w {
		bits |= 0x08
	}
	if r {
		bits |= 0x04
	}
	if x {
		bits |= 0x02
	}
	if b {
		bits |= 0x01
	}
	return bits
}

func VEX2(r, vvvv, l, pp byte) byte {
	r = (^(r << 7)) & 0x80       // 10000000
	vvvv = (^(vvvv << 3)) & 0x78 // 01111000
	return r | vvvv | l<<2 | pp
}

func VEX31(r, x, b, m byte) byte {
	r = (^(r << 7)) & 0x80       // 10000000
	x = (^(x << 6)) & 0x40       // 01000000
	b = (^(x << 5)) & 0x20       // 00100000
	return r | x | b | m
}

func VEX32(w, vvvv, l, pp byte) byte {
	vvvv = (^(vvvv << 3)) & 0x78 // 01111000
	return w | vvvv | l<<2 | pp
}
