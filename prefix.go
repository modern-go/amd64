package amd64

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