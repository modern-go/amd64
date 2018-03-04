package amd64

type Instruction struct {
	Mnemonic string
	assemble interface{}
	imm_r    maybeByte
	imm_rm   ImmRm
	r_rm     maybeByte
	rm_r     maybeByte
	bits     byte
}

type ImmRm struct {
	op  maybeByte
	sub byte
}

type maybeByte interface {
	ok() bool
	value() byte
}

type j struct {
	val byte
}

func (j j) ok() bool    { return true }
func (j j) value() byte { return j.val }

type no struct{}

func (n no) ok() bool    { return false }
func (n no) value() byte { panic("no{}.value()!") }
