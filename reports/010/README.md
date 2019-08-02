# Web services development on Go part 2, week 1 (2019/07/30)

## Web framework
- middleware could be used to not repeat yourself implementing some functionality essential for all incoming requests, such middleware could receive ahandler or even the whole mux, see the example below
```go
func panicMiddleware(next http.Handler)http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("panicMiddleware", r.URL.Path)
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("recovered", err)
                http.Error(w, "Internal server error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```
- `context.Context` is useful to store values which are essential during the whole request lifetime, no more (do not store database connection in context), no less, context has copy-on-write nature so assigning new value to existing `r.Context()` will just add new layer without overwriting previous values, new `context.WithValue` has a `string` key and `interface{}` value, so you have to cast the type each time you use the value
- errors in Go are structures implementing method `func (err MyError) Error() string {...}` so you can make custom errors with all data useful for tracing and use `swith err.(type) {...}` to determine the right error type, moreover it is recommended to use package `github.com/pkg/errors` to wrap your error with some text label, for example `errors.Wrap(err, "resource error")`, and then use it for enchanted tracing
- aside from default `net/http` router (`http.NewServeMux()`) there are several other packages which are slower but more functional: `github.com/gorilla/mux` (`mux.NewRouter()`) could parse complex request parameters - method, header, etc., `github.com/julienschmidt/httprouter` (`httprouter.New()`) could parse only URL parameter so it works significantly faster then gorilla implementation
- if you are not satisfied with prodactivity of default Go web server you can use `github.com/valyala/fasthttp` package providing faster implementation, it has different function signatures so you have to use different routers with it, for example `github.com/buaazp/fasthttprouter`, the problem is that it reuses data structures passed as arguments to request handler so after request is finished data inside context is invalid and couldn't be processed inside some goroutine started after the request is finished
- to parse URL query and validate it you could use packages `github.com/gorilla/schema` and `github.com/asaskevich/govalidator`
- the easiest way to write a web server is to use some framework, the most popular one in Go is Beego (it could generate skeleton using `bee` util, build documentation using Swagger and so on), another good framework is Gin - very fast and lightweight, recommended for use in a small projects
- writing log is essential for every reliable program, to do this in Go you could use standard library packaage `log`, another good logger is ZAP made by Google (you should always specify the type of argument which allows to make less memory allocation), and finally there is Logrus which can output log as JSON, but problem is that it's significantly slower then ZAP
- in Go there is no standard implementation of WebSockets, but you can use for example `github.com/gorilla/websocket`
- standard `html/template` implements everything you'll ever need but it's problem is that it uses reflection in runtime, to optimise it you can use Hero template package instead - it generates Go code from template and you do not have runtime overhead

## Sources
- [Разработка веб-сервисов на Go, часть 2. Лекция 1](golang-5.pdf)
- [The Go language guide. Web application secure coding practices](go-webapp-scp.pdf)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
