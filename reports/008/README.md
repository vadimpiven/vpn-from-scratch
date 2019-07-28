# Web services development on Go - language basics, week 3 (2019/07/24)

## JSON
- packing and unoacking JSON in Go is called marshalling and unmarshalling, Go supports this out of the box, importsnt thing is that `json.Unmarshal(data, emptyStructObject)` function works with `[]byte`, not `string`, `json.Marshal(structObject)` function returns `[]byte`
- all structure fields that would be filled by `Unmarshal` must be public (their names must begin from capital letters), to specify JSON field name corresponding to structure field metainformetion could be given after type as in example below, `json` says that given metadata is intended for JSON, first field specifies the name of 
```go
type User struct {
    ID       int    `json:"user_id,string"`
    Username string
    Address  string `json:",omitempty"`
    Comnpany string `json:"-"`
}
```
- 

## Sources
- [Введение в Golang. Лекция 3](golang-3.pdf)
