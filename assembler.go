package amd64

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"syscall"
	"unsafe"
)

const PageSize = 4096

type Assembler struct {
	Buffer []byte
	Error  error
}

func (asm *Assembler) ReportError(err error) {
	if asm.Error == nil {
		asm.Error = err
	}
}

func (asm *Assembler) Assemble(instructions ...interface{}) {
	for len(instructions) > 0 {
		insn, _ := instructions[0].(*instruction)
		if insn == nil {
			asm.ReportError(fmt.Errorf("not instruction: %v", instructions))
			return
		}
		switch assemble := insn.encoding.(type) {
		case func(a *Assembler):
			assemble(asm)
			instructions = instructions[1:]
		case func(a *Assembler, insn *instruction):
			assemble(asm, insn)
			instructions = instructions[1:]
		case func(a *Assembler, insn *instruction, operand1 Operand):
			operand1, _ := instructions[1].(Operand)
			if operand1 == nil {
				asm.ReportError(fmt.Errorf("not operand: %v", operand1))
				return
			}
			assemble(asm, insn, operand1)
			instructions = instructions[2:]
		case func(a *Assembler, insn *instruction, operand1 Operand, operand2 Operand):
			operand1, _ := instructions[1].(Operand)
			if operand1 == nil {
				asm.ReportError(fmt.Errorf("not operand: %v", operand1))
				return
			}
			operand2, _ := instructions[2].(Operand)
			if operand2 == nil {
				asm.ReportError(fmt.Errorf("not operand: %v", operand2))
				return
			}
			assemble(asm, insn, operand1, operand2)
			instructions = instructions[3:]
		default:
			asm.ReportError(fmt.Errorf("unsupported: %v", insn))
			return
		}
	}
}

func (asm *Assembler) byte(b byte) {
	asm.Buffer = append(asm.Buffer, b)
}

func (asm *Assembler) int16(i uint16) {
	asm.Buffer = append(asm.Buffer, byte(i&0xFF), byte(i>>8))
}

func (asm *Assembler) int32(i uint32) {
	asm.Buffer = append(asm.Buffer,
		byte(i&0xFF),
		byte(i>>8),
		byte(i>>16),
		byte(i>>24))
}

func (asm *Assembler) int64(i uint64) {
	asm.Buffer = append(asm.Buffer,
		byte(i&0xFF),
		byte(i>>8),
		byte(i>>16),
		byte(i>>24),
		byte(i>>32),
		byte(i>>40),
		byte(i>>48),
		byte(i>>56))
}

func (asm *Assembler) rel32(addr uintptr) {
	off := uintptr(addr) - uintptr(unsafe.Pointer(&asm.Buffer[len(asm.Buffer)-1])) - 4
	if uintptr(int32(off)) != off {
		asm.ReportError(errors.New("call rel: target out of range"))
		return
	}
	asm.int32(uint32(off))
}

func (asm *Assembler) imm(imm Immediate) {
	switch imm.bits {
	case 8:
		asm.byte(byte(imm.val))
	case 16:
		asm.int16(uint16(imm.val))
	case 32:
		asm.int32(imm.val)
	}
}

func (asm *Assembler) MakeFunc(f interface{}) {
	pagesCount := (len(asm.Buffer) / PageSize) + 1
	executableMem, err := syscall.Mmap(
		-1,
		0,
		pagesCount*PageSize,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		asm.ReportError(err)
		return
	}
	copy(executableMem, asm.Buffer)
	typ := reflect2.TypeOf(f)
	ptr := unsafe.Pointer(&executableMem)
	typ.UnsafeSet(reflect2.PtrOf(f), unsafe.Pointer(&ptr))
}
