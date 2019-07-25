# Web services development on Go - language basics, week 2 (2019/07/23)

## Go parallelism
- for parallel computation purposes Go implements goroutines - lightweight threads, `go func()` will execute function as new goroutine
- goroutines are executed concurrently, to make sure that all goroutines will have processor time `runtime.Goshed()` added at the end of function would call the task planer
- channels are used for synchronisation and data exchange between goroutines, they could be buffered and unbaffered, `range chan` is used to reseive new data from channel while `close(chan)` is not called
- read and write operations with channels could be performed inside `select` with multiple `case` actions (if no `default` implemented there is a risk of deadlock)
- when channel is passed to function left or right arrow could restrict using channel inside function to only read or write operations (checked at compile time)
- `chan struct{}` could be used for signalling as empty struct doesn't take place in memory
- when main goroutine execution is finished program exits without waiting while all other goroutines will complete, `wg := &sync.WaitGroup{}` could be used to wait for all goroutines (`wg.Add` will add new goroutines, `wg.Wait()` will block caller while all goroutines will not call `defer wg.Done()`)
- `time.NewTimer` (or shortage `time.After`) could be used as a channel `<-timer.C` inside goroutines in `select` to set timeout for operations
- `ticker := time.NewTicker` is used to do some job after specified time interval (`for tickTime := range ticker.C`) while `ticker.Stop()` is not called (`range time.Tick` could be used for shortage if calling stop function is not expected)
- `timer := time.AfterFunc` is used to call some function once after fixed timespan (`timer.Stop()` will stop the timer and function will never be called)
- `ctx, finish := context.WithCancel(context.Background())` allows to wait only for the first several goroutines to finish execution, when `finish()` is called, `select case <-ctx.Done()` is executed inside goroutine (if second argument is passed `ctx, finish := context.WithTimeout(context.Background(), workTime)` function `finish()` would be called automatically with `workTime` delay
- concurrent read and write operations cause races inside the program, `race` argument for build utility could be used to detect them
- `mu := &sync.Mutex{}` could be used to get read of races, lines `mu.Lock()` and `mu.Unlock()` should be used before and after concurrent read and write operations (mutex should be passed to function as a pointer `mu *sync.Mutex`)
- mutexes are rather bulk ans mostly used for complex io operations, they are built on top of atomic package, direct use of atomic (which provide functions for interacting with individual primitive variables) could be more effective, for example `atomic.AddInt32(&totalOperations, 1)` will atomically increment `totalOperations` variable

## Sources
- [Введение в Golang. Лекция 2](golang-2.pdf)
- [Как устроены каналы в Go](https://habr.com/ru/post/308070/)
- [Go Slices: usage and internals](https://blog.golang.org/go-slices-usage-and-internals)
- [Macro View of Map Internals In Go](https://www.ardanlabs.com/blog/2013/12/macro-view-of-map-internals-in-go.html)
- [Go Data Structures: Interfaces](https://research.swtch.com/interfaces)
- [Package atomic](https://golang.org/pkg/sync/atomic/)