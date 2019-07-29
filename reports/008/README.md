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
- memory and cpu usage profiling could be done with `go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_test.go` and then analised with `pprof` tool
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

## Sources
- [Введение в Golang. Лекция 3](golang-3.pdf)
- [Генерация кода в Go](https://habr.com/ru/post/269887/)
