# Web services development on Go part 2, week 2 (2019/08/05)

## Interacting with databases
- Go provides an interface for interacting with SQL databases - `database/sql`, it allows to unify interactions with all common databases, just include driver for particular database for side effects (`init` function) with `import _ "github.com/go-sql-driver/mysql"` and then you can use `*sql.DB` instanse, created with `db = sql.Open("dbname", connOarams)`, then `db.SetMaxOpenConns(10)` to create a pool of reusable connections and then `db.Ping()` to test the connection, using `db.Querry` for `SELECT` and `db.Exec` for other nonreturning querries
- imstead of using SQL driver directly you can use `reflect`- or codegeneration-based libraries, for example `github.com/jinzhu/gorm`, it provides easier way to interact with database with just `db.Find` and `db.Create` with slice as an argument instead of writing SQL querry and using it inside `Querry` and `Exec` functions, you can also use `sql` and `gorm` struct tags to tune the GORM behaviour, moreover you can use numerous GORM hooks called before and after calling some db functionality
- `github.com/bradfitz/gomemcache/memcache` package allows to interact with key-value storage Memcached, it allows `GET`, `SET`, `DELETE`, `INCREMENT`, `EXPIRE`, to keep information up to date use version tags that are incremented each time information updated, if tag values in cached data are smaller then current tag values - cache need to be rebuilt
- to set value of one pointer to value of another pointer use code below
```go
inVal := reflect.ValueOf(in)
resultVal := reflect.ValueOf(result)
rv := reflect.Indirect(inVal)
rvpresult := reflect.Indirect(resultVal) rv.Set(rvpresult) // *in = *result
```
- `github.com/garyburd/redigo/redis` package could be used to interact with Redis, use `redis.DialURL` to set up connection, function `Do` allows to send any supported command, reseived value should be converted to the right type using for example `redis.Bytes`
- RabbitMQ could be used to send some heavily-cpu-time-consuming operations on another server and keep them in a queue, it is based on publisher-subscriber pattern
- `gopkg.in/mgo.v2` package could be used toninteract with MongoDB, data is stored as BSON (binarry JSON), struct tags `bson` could be used to tune field marshalling

## Sources
- [Разработка веб-сервисов на Go, часть 2. Лекция 2](golang-6.pdf)
- [The ultimate guide to building database-driven apps with Go](Database-Driven_Apps.pdf)
- [SQL database drivers](https://github.com/golang/go/wiki/SQLDrivers)
- [SQL Interface](https://github.com/golang/go/wiki/SQLInterface)
- [Go database/sql tutorial](http://go-database-sql.org)
- [Configuring sql.DB for Better Performance](https://www.alexedwards.net/blog/configuring-sqldb)
- [How To Build Microservice With MongoDB In Golang](https://goinbigdata.com/how-to-build-microservice-with-mongodb-in-golang/)
- [Communicating Go Applications through Redis Pub/Sub Messaging Paradigm](https://hackernoon.com/communicating-go-applications-through-redis-pub-sub-messaging-paradigm-df7317897b13)
- [Как работает реляционная БД](https://habr.com/ru/company/mailru/blog/266811/)
