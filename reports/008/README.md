# Web services development on Go - language basics, week 3 (2019/07/24)

## JSON
- packing and unoacking JSON in Go is called marshalling and unmarshalling, Go supports this out of the box, importsnt thing is that `json.Unmarshal(data, emptyStructObject)` function works with `[]byte`, not `string`, `json.Marshal(structObject)` function returns `[]byte`
- all structure fields that would be filled by `Unmarshal` must be public (their names must begin from capital letters), to specify JSON field name corresponding to structure field metainformetion (structure tags) could be given after type as in example below, `json` says that given metadata is intended for JSON, first field specifies the name of JSON field (leave `,` if it's the same as structure field name), second field specifies JSON field type (`omitempty` says that if field value is default - corresponding field of JSON structure should be omitted), if `-` goes after `json` - the field shouldn't be serialised or deserialised, if no structure tag is given - field will be serealised and deserialised with same name ans type
```go
type User struct {
    ID       int    `json:"user_id,string"`
    Username string
    Address  string `json:",omitempty"`
    Comnpany string `json:"-"`
}
```
- if JSON structure is unknown, JSON string could still be Unmarshalled into `interface{}`, which could be interpreted as `map[string]interface{}`
- if we don't know the structure type we can deal with it in runtime using the package `reflect` as in the example below
```go
func PrintReflect(u interface{}) error {
    val := reflect.ValueOf(u).Elem()
    fmt.Printf("%T have %d fields:\n", u, val.NumField())
    for i := 0; i < val.NumField(); i++ {
        valueField := val.Field(i)
        typeField := val.Type().Field(i)
        fmt.Printf("\tname=%v, type=%v, value=%v, tag=‘%v‘\n", typeField.Name,
                typeField.Type.Kind(),
                valueField,
                typeField.Tag,
    ) }
    return nil
}
```
- code generation in Go is possible as Go compiler is now written in Go and standard library contains packages for all things compiler need to do, to generate some code comment `//go:generate cmd` should plased and running `go generate` will perform all given `cmd`; examples are:
    - `easyjson` generating JSON marshalling and unmarshalling functions
    - `stringer` generating `String()` method for series of numbers
    - `unicode` generating Unicode tables
    - `encoding/gob` providing methods for effective encoding and decoding of arrays
    - `time` generating timezones
    - `yacc` generating `.go` files from `.y` syntax descryption
    - `protobufs` generating `.pb.go` files from protocol buffer definition `.proto`
    - `html` embedding HTML files into go sourcecode
    - `bindata` translating binary files such as JPEGs into byte arrays in Go source
- benchmarks in Go are placed inside `test` files, function names begin with `Benchmark`, argument is `b *testing.B`, inside should be a loop `for i := 0; i < b.N; i++ {...}`, command `go test -bench` is used to execute benchmarks, to display memory usage as well as time consumption command `go test -bench . -benchmem` should be performed
- memory and cpu usage profiling could be done with `go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_test.go` and then analised with `go tool pprof` (output types are `--text`, `--web`, `--list`)
- `sync.Pool` could be used to reduce number of new memory allocations inside program when some typical objects are created and then removed, to use it you should first create new pool as shown below, then get memory from pool with `data := dataPool.Get().(*bytes.Buffer)`, use it any way you want, reset to default values with `data.Reset()` and return to pool with `dataPool.Put(data)`
```go
var dataPool = sync.Pool{
    New: func() interface{} {
        return bytes.NewBuffer(make([]byte, 0, 64))
    },
}
```
- tests could be started not only in one thread but also in different goroutines, it could be done with the following code
```go
func BenchmarkAllocNew(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            ...
        }
    })
}
```
- test coverage persantage could be calculated with `go test -v -cover`, to get more reach output you can run `go test -coverprofile=cover.out` and then `go tool cover -html=cover.out -o cover.html` to visualise covered (green-colored) and uncovered (red-colored) parts of code
- XML could be decoded with package `xml`, important thing is that while decoding JSON or XML all data is simultanously loaded into memory, if data is large - we can bump into lack of memory, to avoid this input data should be processed in sycle, reader `input := bytes.NewReader(xmlData)` and decoder `decoder := xml.NewDecoder(input)` sould be created, then we read new portion of data `tok, tokenErr := decoder.Token()` while `tokenErr != io.EOF` and process it with `err := decoder.DecodeElement(&login, &tok)`
- `GOGC=off` before `go test`, `go run`, etc. allows to disable gabbage collector, doing this can speed up small scripts
- linux default profiler `perf` could also be used with Go programs (`sudo perf top -p $(pidof systemtap)`), it shows linux kernal stack trace in addition to Go profiler output
- test helpers could be used by multiple test cases, function name should start with `test`, first line of function should be `t.Helper()` to hide output of this function (it must never return an error)
- `func TestMain(m *testing.M)` could be used to make some setup before and after running tests with `m.Run()`
- `init()` functions can be used within a package block and regardless of how many times that package is imported, the `init()` function will only be called once, it will be executed before `func main()` and before `func TestMain(m *testing.M)`, moreover you can have two separate `init()` functions inside each `.go` file if you need it, they will execute following the order you write them
- goroutines in `syscall` state consume an OS thread, other goroutines do not (except for goroutines that called `runtime.LockOSThread`)
- `//go:noinline` comment just above the function name restricts Go compiler to optimise function into inline one
- inside test files example functions could go, their name should be `ExampleFuncname` where `Funcname` is a name of function the use of which example demonstrates, good example is shown below, this example functions are also performed inside tests to make sure that api haven't changed (the output result is automatically captured and compared with the result following `Output: ` comment)
```go
package stringutil_test

import (
    "fmt"

    "github.com/golang/example/stringutil"
)

func ExampleReverse() {
    fmt.Println(stringutil.Reverse("hello"))
    // Output: olleh
}
```
- test function output should follow the got-expected pattern, not reverse one

## Sources
- [Введение в Golang. Лекция 3](golang-3.pdf)
- [Генерация кода в Go](https://habr.com/ru/post/269887/)
- [Go Reflection: Creating Objects from Types -- Part I (Primitive Types)](https://medium.com/kokster/go-reflection-creating-objects-from-types-part-i-primitive-types-6119e3737f5d)
- [Go Reflection: Creating Objects from Types -- Part II (Composite Types)](https://medium.com/kokster/go-reflection-creating-objects-from-types-part-ii-composite-types-69a0e8134f20y)
- [Профилирование и оптимизация программ на Go](https://habr.com/ru/company/badoo/blog/301990/)
- [Профилирование и оптимизация веб-приложений на Go](https://habr.com/ru/company/badoo/blog/324682/)
- [pprof user interface](https://rakyll.org/pprof-ui/)
- [Инструменты для разработчика Go: знакомимся с лейблами профайлера](https://habr.com/ru/company/badoo/blog/332636/)
- [Five things that make Go fast](https://dave.cheney.net/2014/06/07/five-things-that-make-go-fast)
- [TestMain—What is it Good For?](http://cs-guy.com/blog/2015/01/test-main/)
- [The Go init Function](https://tutorialedge.net/golang/the-go-init-function/)
- [Advanced Testing in Go](https://about.sourcegraph.com/go/advanced-testing-in-go)
- [Debugging performance issues in Go programs](https://github.com/golang/go/wiki/Performance)
- [Testable Examples in Go](https://blog.golang.org/examples)
