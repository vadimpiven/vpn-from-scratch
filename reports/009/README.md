# Web services development on Go - language basics, week 4 (2019/07/29)

## HTTP communication
- `net` and `net/http` packages contain most of standard library functionality used for web communications, first thing to do for creating TCP server is `listner, err := net.Listen("tcp", ":8080")` to listen for incoming messages, then `conn, err := listner.Accept()` to interect with new client on some unused port and finally `go handleConnection(conn)` to work with client in separate thread, example of `handleConnection` function you can see below, to test this server you should use `telnet 127.0.0.1 8080`
```go
func handleConnection(conn net.Conn) {
    name := conn.RemoteAddr().String()
    fmt.Printf("%+v connected\n", name)
    conn.Write([]byte("Hello, " + name + "\n\r"))
    defer conn.Close()
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        text := scanner.Text()
        if text == "Exit" {
            conn.Write([]byte("Bye\n\r"))
            fmt.Println(name, "disconnected")
            break
        } else if text != "" {
            fmt.Println(name, "enters", text)
            conn.Write([]byte("You enter " + text + "\n\r"))
        }
    }
}
```go
- HTTP server is more useful then TCP server, to create it you should register new URL with `http.HandleFunc("/", handler)` and start the server `http.ListenAndServe(":8080", nil)`, hendler should have signature `func handler(w http.ResponseWriter, r *http.Request)` where `w` implements `io.Writer`  interface, `http.HandleFunc("/page",...)` will proceed only `/page` URL, `http.HandleFunc("/pages/",...)` will response to every URL begining with `/pages/`, `http.HandleFunc("/", ...)` will serve all not specified URL addresses
- it is possible to attach HTTP request hendler to `struct`, itks enough to implement `func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)` for the structure (`Hendler` in this case) and then use `http.Handle("/", rootHandler)` (where `rootHandler` is an instance of `Handler`)
- to start more then one server from the same program multiplexer (MUX) could be used as in the example below
```go
mux := http.NewServeMux()
mux.HandleFunc("/", handler)
server := http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
server.ListenAndServe()
```
- to get query params use `r.URL.Query().Get("param")`, to get value from form use `r.FormValue("key")` (both will return empty string if no corresponding value found in request), to get value from cookie use `session, err := r.Cookie("session_id")` and check `err != http.ErrNoCookie`, to set cookie use the code below (you can prolong or expire the cookie lifetime inside any handler), to get header use `r.Header.Get("Accept")` (to set header use `w.Header().Set("Content-Type", "text/html")`)
```go
expiration := time.Now().Add(10 * time.Hour)
cookie := http.Cookie{
    Name:    "session_id",
    Value:   "rvasily",
    Expires: expiration,
}
http.SetCookie(w, &cookie)
```
- to redirect between pages use `http.Redirect(w, r, "/", http.StatusFound)`
- to get user agent use `r.UserAgent()`
- tunserve static files use static handler as in example below (in this example files should be placed in folder `static` near the program binary)
```go
staticHandler := http.StripPrefix(
    "/data/",
    http.FileServer(http.Dir("./static")),
)
http.Handle("/data/", staticHandler)
```
- to get file from form you should at first parse the form with `r.ParseMultipartForm(5 * 1024 * 1025)` (parse first 5 Mb, to parse entire request use `body, err := ioutil.ReadAll(r.Body)` and do not forget `defer r.Body.Close()`) and then read the contents with `file, handler, err := r.FormFile("my_file")` (handler allows to get `handler.Filename` and `handler.Header` - MIME type)
- to send a GET request use `resp, err := http.Get(url)` and do not forget `defer resp.Body.Close()`, example below shows how to add some headers and query params into such request
```go
req := &http.Request{
    Method: http.MethodGet,
    Header: http.Header{
        "User-Agent": {"coursera/golang"},
    },
}
req.URL, _ = url.Parse("http://127.0.0.1:8080/?id=42")
req.URL.Query().Set("user", "rvasily")
resp, err := http.DefaultClient.Do(req)
```
- to use even more parameters (such as timeouts) you should use transport as in the example below
```go
transport := &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}
client := &http.Client{
    Timeout:   time.Second * 10,
    Transport: transport,
}
data := ‘{"id": 42, "user": "rvasily"}‘
body := bytes.NewBufferString(data)
url := "http://127.0.0.1:8080/raw_body"
req, _ := http.NewRequest(http.MethodPost, url, body)
req.Header.Add("Content-Type", "application/json")
req.Header.Add("Content-Length", strconv.Itoa(len(data)))
resp, err := client.Do(req)
```
- to test a handler you can use package `httptest` which can simulate `http.ResponseWriter` with `httptest.NewRecorder()` and `*http.Request` with `httptest.NewRequest`, but better way of testing is not to sumulate but really implement the server sending requests using `net.Conn` or anything else from mentioned above, to test external service you can create dummy handler with `httptest.NewServer(http.HandlerFunc(CheckoutDummy))` and interact with it
- to generate HTML pages you could use package `text/template`, first you should specify the template text with `{{.Browser}}` where name after dot is a field of template parameters structure, then compile template with `tmpl := template.New(‘example‘)` and following `tmpl, _ = tmpl.Parse(EXAMPLE)` (`EXAMPLE` is a variable holding our templates text) and finally fill it with the right values using `tmpl.Execute(w, params)` (this will also automatically write the result into `w http.ResponseWriter`)
- it's recommended to use `html/template` instead of `text/template` as it escapes the parameters (this allows to protect for XSS and other similar attacks), for example write `tmpl := template.Must(template.ParseFiles("users.html"))`, moreover with `html/template` you can use structure methods inside templates as well as a fields, you can even use functions receiving data structure as in the example below (function `IsUserOdd`, usage `{{OddUser .}}` without leading dot)
```go
tmplFuncs := template.FuncMap{
    "OddUser": IsUserOdd,
}
tmpl, err := template.
        New("").
        Funcs(tmplFuncs).
        ParseFiles("func.html")
```
- to benchmark web application you could use default util `ab` (Apache benchmark)
- to use `pprof` in web app just use `import _ "net/http/pprof"` and then you can dump profile
    - `curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt` for memory usage
    - `curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt` for cpu usage
    - `curl http://localhost:8080/debug/pprof/goroutine?debug=2 -o goroutines.txt` for information on goroutines (`goroutines.txt` is a plain text file unlike others)
    - `curl http://localhost:8080/debug/pprof/trace?seconds=10 -o trace.out` for stack trace

## Sources
- [Введение в Golang. Лекция 4](golang-4.pdf)
- [Writing Web Applications](https://golang.org/doc/articles/wiki/#tmp_7)
- [Build Web Application with Golang](https://astaxie.gitbooks.io/build-web-application-with-golang/en/)
- [Webapps in Go, the anti textbook](antitextbookGo.pdf)
- [Network Programming with Go by Jan Newmarch](http://tumregels.github.io/Network-Programming-with-Go/)
- [The complete guide to Go net/http timeouts](https://new.blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
