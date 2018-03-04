package amd64

import (
	"syscall"
	"unsafe"
	"github.com/modern-go/reflect2"
	"errors"
	"fmt"
)

const PageSize = 4096

type Assembler struct {
	Buffer []byte
	Error  error
}

func (a *Assembler) ReportError(err error) {
	if a.Error == nil {
		a.Error = err
	}
}

func (a *Assembler) Assemble(instructions ...interface{}) {
	for len(instructions) > 0 {
		insn, _ := instructions[0].(*instruction)
		if insn == nil {
			a.ReportError(fmt.Errorf("not instruction: %v", instructions))
			return
		}
		switch assemble := insn.assemble.(type) {
		case func(a *Assembler):
			assemble(a)
			instructions = instructions[1:]
		case func(a *Assembler, insn *instruction):
			assemble(a, insn)
			instructions = instructions[1:]
		case func(a *Assembler, insn *instruction, operand1 Operand):
			operand1, _ := instructions[1].(Operand)
			if operand1 == nil {
				a.ReportError(fmt.Errorf("not operand: %v", operand1))
				return
			}
			assemble(a, insn, operand1)
			instructions = instructions[2:]
		case func(a *Assembler, insn *instruction, operand1 Operand, operand2 Operand):
			operand1, _ := instructions[1].(Operand)
			if operand1 == nil {
				a.ReportError(fmt.Errorf("not operand: %v", operand1))
				return
			}
			operand2, _ := instructions[2].(Operand)
			if operand2 == nil {
				a.ReportError(fmt.Errorf("not operand: %v", operand2))
				return
			}
			assemble(a, insn, operand1, operand2)
			instructions = instructions[3:]
		default:
			a.ReportError(fmt.Errorf("unsupported: %v", insn))
			return
		}
	}
}

func (a *Assembler) byte(b byte) {
	a.Buffer = append(a.Buffer, b)
}

func (a *Assembler) int16(i uint16) {
	a.Buffer = append(a.Buffer, byte(i&0xFF), byte(i>>8))
}

func (a *Assembler) int32(i uint32) {
	a.Buffer = append(a.Buffer,
		byte(i&0xFF),
		byte(i>>8),
		byte(i>>16),
		byte(i>>24))
}

func (a *Assembler) int64(i uint64) {
	a.Buffer = append(a.Buffer,
		byte(i&0xFF),
		byte(i>>8),
		byte(i>>16),
		byte(i>>24),
		byte(i>>32),
		byte(i>>40),
		byte(i>>48),
		byte(i>>56))
}

func (a *Assembler) rel32(addr uintptr) {
	off := uintptr(addr) - uintptr(unsafe.Pointer(&a.Buffer[len(a.Buffer)-1])) - 4
	if uintptr(int32(off)) != off {
		a.ReportError(errors.New("call rel: target out of range"))
		return
	}
	a.int32(uint32(off))
}

func (a *Assembler) rex(w, r, x, b bool) {
	var bits byte
	if w {
		bits |= REXW
	}
	if r {
		bits |= REXR
	}
	if x {
		bits |= REXX
	}
	if b {
		bits |= REXB
	}
	if bits != 0 {
		a.byte(PFX_REX | bits)
	}
}

func (a *Assembler) rexBits(lsize, rsize byte, r, x, b bool) {
	if lsize != 0 && rsize != 0 && lsize != rsize {
		panic("mismatched instruction sizes")
	}
	lsize = lsize | rsize
	if lsize == 0 {
		lsize = 64
	}
	a.rex(lsize == 64, r, x, b)
}

func (a *Assembler) modrm(mod, reg, rm byte) {
	a.byte((mod << 6) | (reg << 3) | rm)
}

func (a *Assembler) sib(s, i, b byte) {
	a.byte((s << 6) | (i << 3) | b)
}

func (a *Assembler) MakeFunc(f interface{}) {
	pagesCount := (len(a.Buffer) / PageSize) + 1
	executableMem, err := syscall.Mmap(
		-1,
		0,
		pagesCount*PageSize,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		a.ReportError(err)
		return
	}
	copy(executableMem, a.Buffer)
	typ := reflect2.TypeOf(f)
	ptr := unsafe.Pointer(&executableMem)
	typ.UnsafeSet(reflect2.PtrOf(f), unsafe.Pointer(&ptr))
}