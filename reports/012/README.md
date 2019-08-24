# Web services development on Go part 2, week 3 (2019/08/12)

## Remote Procedure Call (RPC)
-  microservices are a new trend idiology of programming software development, yousing them you could easyer scale your application, deploy small parts instead everything at once, change some parts with more ease, but microservices also have some disadvantages - you have to write more code, process of deployment becomes more sophisticated, more hardware resources are used, and sometimes you couldn't benefit at all from using microservices for some particular tasks
- interaction between microservices is mostly based on RPC, in Go the most simple implementation of RPS is provided by package `net/rpc`, this implementation uses binary data exchange in GOB format, it's simple and fast though the problem is that GOB is not widespread and you will have problems trying to find it's implementation for languages other then Go, to use it you should call `rpc.Register` and then `rpc.HandleHTTP()`, function signature should be of particular look: first argument is a structure of merhod parameters, second argument is an address of place in memory where the result of function should be placed, the only returning value is error, the intrerface defining the microservice functionality must be defined in both caller and microservice, the microservice must implement this interface
- `net/http/jsonrpc` is good alternative of GOB, JSONRPC standard is more wide spread, it's transmiting data as JSON which makes it human-readable and allows to use some authorisation based on transmited data without fully unpacking it, from the other side JSON format requires more channel space to transmit and more space to unpack, inside JSON first field should be `jsonrpc` containing the version of protocole used, second field should be `id` containing serial number of request, third field is the name of requested method and forth argument is array of method params, usage includes creating structure with `ServeHTTP` method and `net/rpc` field, connections will be recieved by `net/http` server, to use `jsonrpc` see the code below
```go
serverCodec := jsonrpc.NewServerCodec(&HttpConn{
    in:  r.Body,
    out: w,
})
w.Header().Set("Content-type", "application/json")
err := h.rpcServer.ServeRequest(serverCodec)
```
- protobuf protogol is another alternative to GOB, it's also binary but more wide spread (also developed by Google), it uses `Marshal` method as JSON, it's easily extandable as fields in transmited data are encoded as their numbers, `protoc` applet will generate the interface code from description for all popular programming languages including Go, it uses gRPC for data transmittion, its feature is that it recieves context as a first argument in all methods and it allows creating middleware
- gRPC allows implementing spdata streams, they will keep connection open until `io.EOF` error, to do this just add `stream` keyword before incoming or outcaming argument
- microservices are mostly used in clasters when one microservice is represented in maltiple copies, Consul by Hashicorp allows to solve a lot of problems with microservice clasters including service discovery and load balansing, package `github.com/hashicorp/consul/api` allows to interacct with Consul
- grpc-gateway allows to call gRPC by sending http request with JSON payload to specially generated proxy which then convert JSON to protobuf and call the gRPC as usual, moreover it's possible to generate Swagger documentation
- versioning of application could be easily implemented with the use of package [https://github.com/TrueFurby/version](https://github.com/TrueFurby/version)

## Sources
- [Разработка веб-сервисов на Go, часть 2ъ Лекция 3](golang-6.pdf)
- [gRPC in Production](https://about.sourcegraph.com/go/grpc-in-production-alan-shreve)
- [Getting Started with Microservices using Go, gRPC and Kubernetes](https://outcrawl.com/getting-started-microservices-go-grpc-kubernetes)
- [gRPC-Web: Moving past REST+JSON towards type-safe Web APIs](https://improbable.io/blog/grpc-web-moving-past-restjson-towards-type-safe-web-apis)
- [Go kit. A toolkit for microservices](https://gokit.io)
- [OpenAPI and gRPC Side-by-Side](https://medium.com/apis-and-digital-transformation/openapi-and-grpc-side-by-side-b6afb08f75ed)
- [How we use gRPC to build a client/server system in Go](https://medium.com/pantomath/how-we-use-grpc-to-build-a-client-server-system-in-go-dd20045fa1c2)
- [Microservices in Golang - Part 1](https://ewanvalentine.io/microservices-in-golang-part-1/)
- [Write a Kubernetes-ready service from zero step-by-step](https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/)
- [Command go](https://golang.org/cmd/go/)
- [Command link](https://golang.org/cmd/link/)
- [Implementing UDP vs TCP in Golang](http://www.minaandrawos.com/2016/05/14/udp-vs-tcp-in-golang/)
- [THE X-FILES: CONTROLLING THROUGHPUT WITH RATE.LIMITER](https://rodaine.com/2017/05/x-files-time-rate-golang/)
- [THE X-FILES: AVOIDING CONCURRENCY BOILERPLATE WITH GOLANG.ORG/X/SYNC](https://rodaine.com/2018/08/x-files-sync-golang/)
- [Introduction to modern network load balancing and proxying](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236)
