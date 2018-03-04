package tests

func aka(b1 byte, b2 byte) byte {
	if b1 != b2 {
		panic("aka not equal")
	}
	return b1
}
