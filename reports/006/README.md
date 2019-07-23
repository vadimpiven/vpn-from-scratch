# Web services development on Go - language basics (2019/07/22)

## Go for new gophers
- garbage collector
- camelCase naming (`gofmt` preserves all other styling conventions)
- code generation with `go generate` and spesially designed comments `//go:generate`
- to use slice as a list of arguments use `...`
- `++i` not exists, only `i++`
- when data is appended to slice its capasity is doubled with allocating new underlaying array, but data is not transfered into this new array immidiately, small parts are moved each time some operations over this slice are performed, inside the code there is a number marking begining from which index data is not copied and should be taken from old array
- strings fully support UTF-8, so the char type is a 4-byte rune, `for i, c := range str` will place in `c` rune typed copies of characters from `str`, but `str[i]` will return one byte, not one symbol (important fact - strings could not be modified with `str[0] = 'a'` because this could brake UTF-8 symbol, to make this change you should convert string into slice of bytes with `var b = []byte(str)`
- `const a = 354` will leave `a` unassigned to type which makes possible to store numbers larger then `int64` (`int` without suffics is platform dependent and could be replaced with `int32` or `int64`)
- all type convertions musy be done manually, even when assigning `nil` to some link-type variable
- use of `unsafe.Pointer` makes program platform-dependent and leaves it without backwards-compartibility
- `iota` will be replased with serial numbers inside the block (begining from 0), to skip some number blank identifier `_` should be placed
- is standard library there are types `complex64` and `complex128`, functions for them are stored in `math/cmplx`
- you should write `fallthrough` inside `switch` to imitate default C language behaviour
- `for i,v := range map` will always give different order of keys because internally the begining index is not `0` but `rand()`
- second returning value of type `error` (all error types should implement error interface `func (e MyError) Error() string`) is used in GO insted try-catch blocks, in some situations `panic("error")` and `defer recover()` could be used for this purpose
- private and public keywords are not used, instead private fields and variables names begin from low-case letters while public field and global variable names - from capital letters
- new unassigned variables contain default value (0, false, etc.) instead of some rubbish as in C
- encapsulation is replased with embedding (all fields and functions of embedded structure will be merged into new one), operator overloading is not possible (all operations should be implemented through the interfaces of structure methods or, in case of new, with `pkg.New()` function)
- methods with `func (p Pers) funcname()` could not modify structure fields, only `func (p *Pers) funcname()` could
- begining from early August 2019 only officially supported package manager is go modules
- test files should have suffix `test` in their names, test function names should begin with `Test`, their signature should be `func TestFunc(t *testing.T)` and `go test -v` should be used to perform tests
- to make possible generating beautiful documentation with `godoc` comments for functions and packages should be written above function names and package name, they should begin with correspongind name and end with dot (in order to look as a finished sentence)
- no `;` are needed in GO, so curly bracket should be placed only in same line with control sequence name (`if true {`), not in the new line

## Sources
- [Введение в Golang. Лекция 1](golang-1.pdf)
- [Разбираемся в Go: пакет io ](https://habr.com/ru/post/306914/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [The Go Programming Language Specification](https://golang.org/ref/spec)
- [Godoc: documenting Go code](https://blog.golang.org/godoc-documenting-go-code)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Generating code](https://blog.golang.org/generate)
- [Go 1.11 Modules](https://github.com/golang/go/wiki/modules)
- [Package unsafe](https://golang.org/pkg/unsafe/)