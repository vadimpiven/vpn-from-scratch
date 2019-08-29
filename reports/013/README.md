# Web services development on Go part 2, week 4 (2019/08/26)

## Configs & Monitoring
- `flag` package could be used to parse command line flags and arguments, first defint all required flags with `flag.Bool("comments", false, "Enable comments after post")` where first argument is flag name, second is default value and third is a usage and then `flag.Parse()` must be called, to define flag of nonstandard type methods `String() string` and `Set(in string) error` should be provided
- it is possible to set value of uninitiolised global string variables using linker, to use this pass `-ldflags="-X 'pkg.Var=val'` to `go built` or `go run`
- to update configure without reloading program you could store configuration inside HashiCorp Consul, the only problem is that such config is recieved as `map[string]string` and so you should always remember which keys exist inside configuration
- service monitoring consists of three major parts: metrics (small, could be stored for long, mostly include timings and resource usage), logs (medium, could be stored for shorter period, include information about incoming and outgoing data, users requesting data and full error output), traces (large, could be stored only for short, could not be gathered for regular basis for each request, include full information about what is happening inside program while proceedeing the request), the most common tool for gathering metrics is Prometheus, logs are written by program itself and processed by ELK stack, traces are made with OpenTracing standard API which is best implemented for Go by Jaeger
- the default method for Go to produse metrics id `expvar` package, you could start using it with `expvar.NewMap` and adding some values to this map to indecate the resource usage, the statistics could be accessed at `/debug/var`, such statistics could be gathered by Graphite or any other tool
- to use Prometheus it's enough to import `github.com/prometheus/client_golang/prometheus/promhttp` and then register `promhttp.Handler()` providing the registered addres to Prometheus server to make it gather all statistics by itself, the resulting statistics is recommended to be visualised with Grafana, aside from standard metrics you could register your own params using the code below
```go
timings = prometheus.NewSummaryVec(
        prometheus.SummaryOpts{
                Name: "method_timing",
                Help: "Per method timing",
        },
        []string{"method"},
)
counter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
                Name: "method_counter",
                Help: "Per method counter",
        []string{"method"},
)
prometheus.MustRegister(timings)
prometheus.MustRegister(counter)

func timeTrackingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        // r.URL.Path приходит от юзера! не делайте так в проде!
        timings.
            WithLabelValues(r.URL.Path).
            Observe(float64(time.Since(start).Seconds()))
        counter.
            WithLabelValues(r.URL.Path).
            Inc()
    }
}
```
- package `unsafe` allows to use pointers as you would do it in C at a cost of no forward and backward compatibility guaranteed, moreover unsafe code will not work in Google IaaS (infrastructure as anservice), the important thing is that slices (including string) and maps are complex objects, the pointer you get for them will point not to the data but to the `reflect.SliceHeader` or other header which includes pointer to the data, its length and possibly some other fields, also there is a broblem with gabbage collector as it could clear the data when it's still used by unsafe pointer or do not collect it when it's no longer used, to solve this you could use `runtime.KeepAlive` on the original object in place when it's no longer used by unsafe pointer (it also works with passing Go objects to the C methods)
- to integrate Go with C you can either put C code as a block comment just above `import "C"` or store it inside independant `.c` file with `#include "_cgo_export.h"` at the first line, this way you steel have to plase C method signatures above `import "C"` inside Go code, important thing is that passing values to C code requires convertion to `C.int` or other type you need and the return result must be converted back to the Go type, there is also a problem with garbage collector mentioned above, to secure slices ans strings you have to call `C.Cstring` to make a C-like copy of string, call C method and then manually call `C.free` (`defer` usage is recommended), moreover the problem with C code is that while common go-routines can reuse single OS thread, when calling C code it will automatically call `runtime.LockOSThread()` and use OS thread exclusively (this security measure can seriously slow down the execution process, the only way to bypass ot is to use assembler instead of C)
- aside from C, you can use assembler code with Go, most of cryptography inside standard library is written on this assembler for performance optimisation, this assembler should be stored in `.s` files, it's syntax is taken from Plan9 OS assembler so it differs from NASM, intel ASM, etc., interesting thing is that while this is called assembler, in reality it is a virtual set of commands that will be then processed and optimised by compiler, so the Go assembler could be crosscompiled, the problem is that different architectures use different set of registers so you need to reimplement some methods for different architectures inside `#ifdef GOARCH_amd64` (or other architecture name) block, for all platform-dependent stdlib functions Go provides `#include "go_asm.h"`, to call assembler functions from Go code you have to define function signatures in Go code
-  to preserve clean and well-documented code stile and fix some potential errors it's recommended to use `go vet` and [github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint) toolkit

## Sources
- [Разработка веб-сервисов на Go, часть 2. Лекция 3](golang-6.pdf)
- [Кросс-компиляция в Go](https://habr.com/ru/post/249449/)
- [С-вызовы в Go: принцип работы и производительность](https://habr.com/ru/company/intel/blog/275709/)
- [Essential Go](https://www.programming-books.io/essential/go/)
- [Quick intro to go assembly](https://blog.hackercat.ninja/post/quick_intro_to_go_assembly/)
