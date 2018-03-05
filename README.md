# amd64

amd64 instruction assembler

Initial implementation copied from https://github.com/nelhage/gojit

```go
import . "github.com/modern-go/amd64"

asm := &Assembler{}
asm.Assemble(
    MOV, RAX, QWORD(RSP, 0x08),
    MOV, QWORD(RSP, 0x10), RAX,
    RET,
)
// ident func does nothing
// return identical value out
var ident func(i int) int
asm.MakeFunc(&ident)
fmt.Println(ident(100)) // will print 100
```
