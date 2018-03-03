package amd64

import (
	"syscall"
	"unsafe"
	"github.com/modern-go/reflect2"
)

const PageSize = 4096

type Assembler struct {
	Buffer []byte
	Error  error
}

func (assembler *Assembler) ReportError(err error) {
	if assembler.Error == nil {
		assembler.Error = err
	}
}

func (assembler *Assembler) Write(instructions ...interface{}) {
}

func (assembler *Assembler) Assemble(f interface{}) {
	pagesCount := (len(assembler.Buffer) / PageSize) + 1
	executableMem, err := syscall.Mmap(
		-1,
		0,
		pagesCount*PageSize,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		assembler.ReportError(err)
		return
	}
	copy(executableMem, assembler.Buffer)
	typ := reflect2.TypeOf(f)
	ptr := unsafe.Pointer(&executableMem)
	typ.UnsafeSet(reflect2.PtrOf(f), unsafe.Pointer(&ptr))
}

type Register string

var RAX = Register("RAX")
var RSP = Register("RSP")

type Indirect struct {
	Base   Register
	Offset int
}

func QWORD(base Register, offset int) Indirect {
	return Indirect{
		Base:   base,
		Offset: offset,
	}
}

var MOV = ""
var RET = ""
